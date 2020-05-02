module gopkg.in/tucnak/telebot.v2/examples/awslambdaechobot

go 1.14

require (
	github.com/aws/aws-lambda-go v1.16.0
	github.com/yi-jiayu/terraform-provider-telegram v0.1.1 // indirect
	gopkg.in/tucnak/telebot.v2 v2.0.0-00010101000000-000000000000
)

replace gopkg.in/tucnak/telebot.v2 => ../..
