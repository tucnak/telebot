package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:       os.Getenv("TELEBOT_SECRET"),
		Synchronous: true,
	})
	if err != nil {
		panic(err)
	}

	b.Handle(tb.OnText, func(m *tb.Message) { b.Send(m.Chat, m.Text) })

	lambda.Start(func(req events.APIGatewayProxyRequest) (err error) {
		var u tb.Update
		if err = json.Unmarshal([]byte(req.Body), &u); err == nil {
			b.ProcessUpdate(u)
		}
		return
	})
}
