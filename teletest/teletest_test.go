package teletest_test

import (
	"fmt"
	"strings"
	"testing"

	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/teletest"
)

func newBot(t *testing.T, srv *teletest.Server) *tele.Bot {
	b, err := tele.NewBot(tele.Settings{URL: srv.URL, Token: srv.Token, Offline: true})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestServerRecordsJSON(t *testing.T) {
	srv := teletest.NewServer(t)
	srv.Reply("sendMessage", map[string]interface{}{
		"message_id": 7,
		"chat":       map[string]interface{}{"id": 42, "type": "private"},
	})

	msg, err := newBot(t, srv).Send(tele.ChatID(42), "hello", tele.Silent)
	if err != nil {
		t.Fatal(err)
	}
	if msg.ID != 7 || msg.Chat.ID != 42 {
		t.Fatalf("decoded %+v", msg)
	}

	req := srv.Last()
	if req.Method != "sendMessage" || req.Params["text"] != "hello" || req.Params["chat_id"] != "42" || req.Params["disable_notification"] != "true" {
		t.Fatalf("recorded %+v", req)
	}
}

func TestServerRecordsMultipart(t *testing.T) {
	srv := teletest.NewServer(t)
	srv.Reply("sendDocument", map[string]interface{}{"message_id": 1, "chat": map[string]interface{}{"id": 42}})

	doc := &tele.Document{File: tele.FromReader(strings.NewReader("data")), FileName: "a.txt", Caption: "c"}
	if _, err := newBot(t, srv).Send(tele.ChatID(42), doc); err != nil {
		t.Fatal(err)
	}

	req := srv.Last()
	if req.Method != "sendDocument" || req.Params["caption"] != "c" {
		t.Fatalf("recorded %+v", req)
	}
	if f := req.Files["document"]; f.Name != "a.txt" || string(f.Data) != "data" {
		t.Fatalf("file %+v", f)
	}
}

func TestServerErrorsAndQueue(t *testing.T) {
	srv := teletest.NewServer(t)
	srv.ReplyError("getChat", 400, "Bad Request: chat not found")
	srv.Reply("getChat", map[string]interface{}{"id": 1, "type": "private"})

	b := newBot(t, srv)
	if _, err := b.ChatByID(1); err != tele.ErrChatNotFound {
		t.Fatalf("first call: %v", err)
	}
	for i := 0; i < 2; i++ {
		if c, err := b.ChatByID(1); err != nil || c.ID != 1 {
			t.Fatalf("call %d: %v %v", i, c, err)
		}
	}
	if n := len(srv.Calls("getChat")); n != 3 {
		t.Fatalf("recorded %d calls", n)
	}

	srv.Reset()
	if len(srv.Requests()) != 0 {
		t.Fatal("Reset kept requests")
	}
}

func TestServerRejectsWrongToken(t *testing.T) {
	srv := teletest.NewServer(t)
	b, err := tele.NewBot(tele.Settings{URL: srv.URL, Token: "1:other", Offline: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.ChatByID(1); err != tele.ErrUnauthorized {
		t.Fatalf("got %v", err)
	}
	if strings.Contains(fmt.Sprint(err), "1:other") {
		t.Fatal("error carries the token")
	}
}
