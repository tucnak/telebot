package telebot

import "testing"

const TESTING_TOKEN = "107177593:AAHBJfF3nv3pZXVjXpoowVhv_KSGw56s8zo"

func TestCreate(t *testing.T) {
	_, err := Create(TESTING_TOKEN)
	if err != nil {
		t.Fatal(err)
	}
}
