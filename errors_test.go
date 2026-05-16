package telebot

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// Regression for #770. The stdlib HTTP client embeds the full request
// URL in transport errors, and that URL contains /bot<TOKEN>/<method>.
// wrapError used to surface that token verbatim. It now redacts the
// token before wrapping.
func TestWrapErrorRedactsBotToken(t *testing.T) {
	token := "123456789:AAEhBP0av28DKrJB1ULTGYojCZ0HEYJBotk"
	rawErr := fmt.Errorf(`Post "https://api.telegram.org/bot%s/getFile": dial: timeout`, token)

	wrapped := wrapError(rawErr)
	if strings.Contains(wrapped.Error(), token) {
		t.Fatalf("wrapped error leaked bot token: %q", wrapped.Error())
	}
	if !strings.Contains(wrapped.Error(), "/bot<token>") {
		t.Fatalf("wrapped error missing /bot<token> sentinel: %q", wrapped.Error())
	}
	if !errors.Is(wrapped, rawErr) {
		t.Fatal("wrapped error broke errors.Is chain to the original error")
	}
}

func TestWrapErrorLeavesUnrelatedErrorsAlone(t *testing.T) {
	original := errors.New("something unrelated")
	wrapped := wrapError(original)
	if !errors.Is(wrapped, original) {
		t.Fatal("expected wrapped error to keep wrapping the original")
	}
	if want := "telebot: something unrelated"; wrapped.Error() != want {
		t.Fatalf("wrapped error = %q, want %q", wrapped.Error(), want)
	}
}
