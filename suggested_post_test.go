package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuggestPost(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendMessage", json.RawMessage(`{"message_id":1,"chat":{"id":-200,"type":"supergroup","is_direct_messages":true},
		"direct_messages_topic":{"topic_id":77,"user":{"id":5}},"suggested_post_info":{"state":"pending","price":{"currency":"XTR","amount":100}}}`))

	msg, err := b.Send(ChatID(-200), "post me", &SendOptions{DirectMessagesTopicID: 77},
		&SuggestedPostParameters{Price: &SuggestedPostPrice{Currency: CurrencyStars, Amount: 100}, SendDate: 1700000000})
	require.NoError(t, err)

	req := srv.Last()
	assert.Equal(t, "77", req.Params["direct_messages_topic_id"])
	assert.JSONEq(t, `{"price":{"currency":"XTR","amount":100},"send_date":1700000000}`, req.Params["suggested_post_parameters"])

	assert.True(t, msg.Chat.IsDirectMessages)
	assert.Equal(t, int64(77), msg.DirectMessagesTopic.ID)
	assert.Equal(t, SuggestedPostPending, msg.SuggestedPostInfo.State)

	require.NoError(t, b.ApproveSuggestedPost(msg, 1700000500))
	assert.Equal(t, "approveSuggestedPost", srv.Last().Method)
	assert.Equal(t, map[string]string{"chat_id": "-200", "message_id": "1", "send_date": "1700000500"}, srv.Last().Params)

	require.NoError(t, b.DeclineSuggestedPost(msg, "no"))
	assert.Equal(t, "declineSuggestedPost", srv.Last().Method)
	assert.Equal(t, map[string]string{"chat_id": "-200", "message_id": "1", "comment": "no"}, srv.Last().Params)
}

func TestSuggestedPostServiceMessages(t *testing.T) {
	b, _ := newFakeBot(t)

	var got []string
	for _, ev := range []string{OnSuggestedPostApproved, OnSuggestedPostApprovalFailed, OnSuggestedPostDeclined, OnSuggestedPostPaid, OnSuggestedPostRefunded} {
		ev := ev
		b.Handle(ev, func(c Context) error { got = append(got, ev); return nil })
	}

	var paid *SuggestedPostPaid
	for _, m := range []string{
		`"suggested_post_approved":{"send_date":5,"price":{"currency":"TON","amount":1000000000}}`,
		`"suggested_post_approval_failed":{"price":{"currency":"XTR","amount":50}}`,
		`"suggested_post_declined":{"comment":"later"}`,
		`"is_paid_post":true,"suggested_post_paid":{"currency":"XTR","star_amount":{"amount":42}}`,
		`"suggested_post_refunded":{"reason":"post_deleted"}`,
	} {
		var u Update
		require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":2,"date":1,"chat":{"id":-200},`+m+`}}`), &u))
		b.ProcessUpdate(u)
		if u.Message.SuggestedPostPaid != nil {
			paid = u.Message.SuggestedPostPaid
			assert.True(t, u.Message.IsPaidPost)
		}
	}
	assert.Equal(t, []string{OnSuggestedPostApproved, OnSuggestedPostApprovalFailed, OnSuggestedPostDeclined, OnSuggestedPostPaid, OnSuggestedPostRefunded}, got)
	assert.Equal(t, int64(42), paid.StarAmount.Amount)
}

func TestDirectMessagesChat(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("getChat", json.RawMessage(`{"id":-200,"type":"supergroup","is_direct_messages":true,"parent_chat":{"id":-100,"type":"channel"}}`))

	chat, err := b.ChatByID(-200)
	require.NoError(t, err)
	assert.True(t, chat.IsDirectMessages)
	assert.Equal(t, int64(-100), chat.ParentChat.ID)

	require.NoError(t, b.Promote(&Chat{ID: -100}, &ChatMember{User: &User{ID: 7}, Rights: Rights{CanManageDirectMessages: true}}))
	assert.Equal(t, "true", srv.Last().Params["can_manage_direct_messages"])
}
