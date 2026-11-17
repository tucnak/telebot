package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v4/internal/spec"
)

// TestEveryUpdateIsRouted checks that each field of the Bot API Update
// reaches a handler and is listed in AllowedUpdates.
func TestEveryUpdateIsRouted(t *testing.T) {
	msg := `{"message_id":1,"date":1,"chat":{"id":1,"type":"private"},"from":{"id":2},"text":"hi"}`
	routes := map[string]struct {
		event string
		value string
	}{
		"message":                    {OnText, msg},
		"edited_message":             {OnEdited, msg},
		"channel_post":               {OnChannelPost, msg},
		"edited_channel_post":        {OnEditedChannelPost, msg},
		"business_connection":        {OnBusinessConnection, `{"id":"bc","user":{"id":2},"user_chat_id":2,"date":1,"is_enabled":true}`},
		"business_message":           {OnBusinessMessage, msg},
		"edited_business_message":    {OnEditedBusinessMessage, msg},
		"deleted_business_messages":  {OnDeletedBusinessMessages, `{"business_connection_id":"bc","chat":{"id":1},"message_ids":[1]}`},
		"guest_message":              {OnGuestMessage, msg},
		"message_reaction":           {OnReaction, `{"chat":{"id":1},"message_id":1,"user":{"id":2},"date":1,"old_reaction":[],"new_reaction":[{"type":"emoji","emoji":"👍"}]}`},
		"message_reaction_count":     {OnReactionCount, `{"chat":{"id":1},"message_id":1,"date":1,"reactions":[]}`},
		"inline_query":               {OnQuery, `{"id":"q","from":{"id":2},"query":"x","offset":""}`},
		"chosen_inline_result":       {OnInlineResult, `{"result_id":"r","from":{"id":2},"query":"x"}`},
		"callback_query":             {OnCallback, `{"id":"c","from":{"id":2},"chat_instance":"i","data":"d"}`},
		"widget_query":               {OnWidgetQuery, `{"id":"w","from":{"id":2},"chat_instance":"i","data":"d","type":{"type":"data"}}`},
		"shared_widget_query":        {OnSharedWidgetQuery, `{"id":"w","from":{"id":2},"data":"d","widget_id":"x","type":{"type":"data"}}`},
		"shipping_query":             {OnShipping, `{"id":"s","from":{"id":2},"invoice_payload":"p","shipping_address":{}}`},
		"pre_checkout_query":         {OnCheckout, `{"id":"p","from":{"id":2},"currency":"XTR","total_amount":1,"invoice_payload":"p"}`},
		"purchased_paid_media":       {OnPurchasedPaidMedia, `{"from":{"id":2},"paid_media_payload":"p"}`},
		"poll":                       {OnPoll, `{"id":"p","question":"q","options":[],"total_voter_count":0}`},
		"poll_answer":                {OnPollAnswer, `{"poll_id":"p","user":{"id":2},"option_ids":[0]}`},
		"my_chat_member":             {OnMyChatMember, `{"chat":{"id":1},"from":{"id":2},"date":1,"old_chat_member":{"status":"member","user":{"id":3}},"new_chat_member":{"status":"left","user":{"id":3}}}`},
		"chat_member":                {OnChatMember, `{"chat":{"id":1},"from":{"id":2},"date":1,"old_chat_member":{"status":"member","user":{"id":3}},"new_chat_member":{"status":"left","user":{"id":3}}}`},
		"chat_join_request":          {OnChatJoinRequest, `{"chat":{"id":1},"from":{"id":2},"user_chat_id":2,"date":1}`},
		"chat_boost":                 {OnBoost, `{"chat":{"id":1},"boost":{"boost_id":"b","add_date":1,"expiration_date":2,"source":{"source":"premium","user":{"id":2}}}}`},
		"removed_chat_boost":         {OnBoostRemoved, `{"chat":{"id":1},"boost_id":"b","remove_date":1,"source":{"source":"premium","user":{"id":2}}}`},
		"managed_bot":                {OnManagedBot, `{"user":{"id":2},"bot":{"id":3}}`},
		"subscription":               {OnSubscription, `{"user":{"id":2},"invoice_payload":"p","state":"active"}`},
		"stopped_message_generation": {OnGenerationStopped, `{"chat":{"id":1},"draft_id":1}`},
	}

	api, err := spec.Load()
	require.NoError(t, err)
	for _, ty := range api.Types {
		if ty.Name != "Update" {
			continue
		}
		for _, f := range ty.Fields {
			if f.Name == "update_id" {
				continue
			}
			route, ok := routes[f.Name]
			if !assert.True(t, ok, "Update.%s has no route in this test", f.Name) {
				continue
			}
			assert.Contains(t, AllowedUpdates, f.Name)

			b, _ := newFakeBot(t)
			fired := false
			b.Handle(route.event, func(c Context) error { fired = true; return nil })

			var u Update
			require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"`+f.Name+`":`+route.value+`}`), &u), f.Name)
			b.ProcessUpdate(u)
			assert.True(t, fired, "Update.%s did not reach its handler", f.Name)
		}
	}
}

func TestReactionContext(t *testing.T) {
	b, _ := newFakeBot(t)

	var got *MessageReaction
	var sender *User
	var chat *Chat
	b.Handle(OnReaction, func(c Context) error {
		got, sender, chat = c.MessageReaction(), c.Sender(), c.Chat()
		return nil
	})

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message_reaction":{"chat":{"id":-100},"message_id":5,"user":{"id":7},"date":1,
		"old_reaction":[],"new_reaction":[{"type":"emoji","emoji":"🔥"}]}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, got)
	assert.Equal(t, "🔥", got.NewReaction[0].Emoji)
	assert.Equal(t, int64(7), sender.ID)
	assert.Equal(t, int64(-100), chat.ID)
}
