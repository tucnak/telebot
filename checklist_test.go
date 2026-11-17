package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendAndEditChecklist(t *testing.T) {
	b, srv := newFakeBot(t)
	reply := json.RawMessage(`{"message_id":1,"chat":{"id":7},"checklist":{"title":"Trip","tasks":[{"id":1,"text":"Pack"}]}}`)
	srv.Reply("sendChecklist", reply)
	srv.Reply("editMessageChecklist", reply)

	list := &Checklist{
		Title:     "<b>Trip</b>",
		ParseMode: ModeHTML,
		Tasks: []ChecklistTask{
			{ID: 1, Text: "Pack"},
			{ID: 2, Text: "Go", TextEntities: []MessageEntity{{Type: EntityBold, Length: 2}}},
		},
		OthersCanAddTasks: true,
	}
	msg, err := b.Send(ChatID(7), list, &SendOptions{BusinessConnectionID: "bc", ParseMode: ModeMarkdown})
	require.NoError(t, err)
	assert.Equal(t, "Pack", msg.Checklist.Tasks[0].Text)

	req := srv.Last()
	assert.Equal(t, "sendChecklist", req.Method)
	assert.Equal(t, "bc", req.Params["business_connection_id"])
	assert.NotContains(t, req.Params, "parse_mode")
	assert.JSONEq(t, `{"title":"<b>Trip</b>","parse_mode":"HTML","others_can_add_tasks":true,"tasks":[
		{"id":1,"text":"Pack"},{"id":2,"text":"Go","text_entities":[{"type":"bold","offset":0,"length":2}]}]}`, req.Params["checklist"])

	_, err = b.Edit(msg, &Checklist{Title: "Trip", Tasks: []ChecklistTask{{ID: 1, Text: "Pack"}}}, &SendOptions{BusinessConnectionID: "bc"})
	require.NoError(t, err)
	req = srv.Last()
	assert.Equal(t, "editMessageChecklist", req.Method)
	assert.Equal(t, map[string]string{"business_connection_id": "bc", "chat_id": "7", "message_id": "1",
		"checklist": `{"tasks":[{"id":1,"text":"Pack"}],"title":"Trip"}`}, req.Params)
}

func TestChecklistServiceMessages(t *testing.T) {
	b, srv := newFakeBot(t)

	// A checklist is not media: OnMedia handlers expect Message.Media().
	var got []string
	for _, ev := range []string{OnMedia, OnChecklist, OnChecklistTasksDone, OnChecklistTasksAdded, OnDirectMessagePriceChanged} {
		ev := ev
		b.Handle(ev, func(c Context) error { got = append(got, ev); return nil })
	}

	for _, m := range []string{
		`"checklist":{"title":"T","tasks":[{"id":1,"text":"a","completed_by_user":{"id":5},"completion_date":9}],"others_can_mark_tasks_as_done":true}`,
		`"checklist_tasks_done":{"marked_as_done_task_ids":[1],"marked_as_not_done_task_ids":[2]}`,
		`"checklist_tasks_added":{"checklist_message":{"message_id":1,"date":1,"chat":{"id":7}},"tasks":[{"id":3,"text":"c"}]}`,
		`"direct_message_price_changed":{"are_direct_messages_enabled":true,"direct_message_star_count":10}`,
	} {
		var u Update
		require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":2,"date":1,"chat":{"id":7},`+m+`}}`), &u))
		b.ProcessUpdate(u)
	}
	assert.Equal(t, []string{OnChecklist, OnChecklistTasksDone, OnChecklistTasksAdded, OnDirectMessagePriceChanged}, got)

	delete(b.handlers, OnChecklist)
	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":2,"date":1,"chat":{"id":7},"checklist":{"title":"T","tasks":[]}}}`), &u))
	b.ProcessUpdate(u)
	assert.Len(t, got, 4, "an unhandled checklist does not fall back to OnMedia")

	srv.Reply("getMyStarBalance", json.RawMessage(`{"amount":12,"nanostar_amount":500}`))
	bal, err := b.StarBalance()
	require.NoError(t, err)
	assert.Equal(t, &StarAmount{Amount: 12, NanostarAmount: 500}, bal)
}

func TestReplyToChecklistTask(t *testing.T) {
	var msg Message
	require.NoError(t, json.Unmarshal([]byte(`{"message_id":2,"date":1,"chat":{"id":7},"reply_to_checklist_task_id":3}`), &msg))
	assert.Equal(t, 3, msg.ReplyToChecklistTaskID)

	data, err := json.Marshal(ReplyParams{MessageID: 1, ChecklistTaskID: 3})
	require.NoError(t, err)
	assert.JSONEq(t, `{"message_id":1,"checklist_task_id":3}`, string(data))
}

func TestStarBalanceLive(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	bal, err := b.StarBalance()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, bal.Amount, int64(0))
}
