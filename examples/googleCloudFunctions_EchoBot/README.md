
<p align="center">
  <img src="https://seeklogo.com/images/G/google-cloud-functions-logo-AECD57BFA2-seeklogo.com.png" alt="Google Cloud Functions"></a>
</p>
<p align="center">
    <em>Deploy a Telegram Bot at Google Cloud Functions</em>
</p>

---

This is a simple example how to deploy a Telegram Bot at [Google Cloud Functions](https://cloud.google.com/functions/ "Google Cloud Functions"). We assume you have a Google Cloud account active, and you have succeeded deploy the [hello world example](https://cloud.google.com/functions/docs/first-go "hello world example").



## Requirements
- A Google Cloud account activated
- [gcloud cli](https://cloud.google.com/sdk/gcloud "gcloud cli") installed (optional, but useful)
- Bot Token (created with https://t.me/BotFather)

## First Steps

- Create or clone this project and download this library module `go get -u gopkg.in/tucnak/telebot.v2
  `
- Add your token at `env.yaml` file
- Deploy executing `deploy.sh` file if you have [gcloud cli](https://cloud.google.com/sdk/gcloud "gcloud cli ") installed in your system. If you haven't gcloud cli, you can [deploy using web browser console](https://cloud.google.com/functions/docs/deploying/console "deploy using web browser console") following these instructions.
- Set your webhook follow the instructions bellow.

## Webhook

To Telegram delivers new messages to your bot, you need inform where to send those messages. This is where the webhook comes to.

Webhook is simple an url to receive new messages. Everytime this endpoint is called, it will execute the bot and processes the update.

If you have deployed using gcloud cli, the deploy output command should be something like this:

```yaml
httpsTrigger:
  securityLevel: SECURE_OPTIONAL
  url: https://us-miami-your-project-name.cloudfunctions.net/telebot
```

And to execute your function, just append the function name at the end. For example:
`https://us-miami-your-project-name.cloudfunctions.net/telebot/EchoBot`

This is your public url to execute your function. But to pass to Telegram as webhook, you must
`https://api.telegram.org/bot<your-bot-token>/setWebhook?url=<your-cloud-functions-public-url>`

There is also a nice [guide how to setup a webhook here](https://panjeh.medium.com/telegram-bot-get-webhook-updates-send-message-49156ac02375 "guide how to setup a webhook here"). Check it in case of any issue.

For example:
https://api.telegram.org/bot1233456789:AeAeaAEAeaAEaEAeAeA/setWebhook?url=https://us-miami-your-project-name.cloudfunctions.net/telebot/EchoBot

Just copy and paste this url in your browser, and you should receive the following response:

```json
{
  "ok": true,
  "result": true,
  "description": "Webhook was set"
}
```

This means your bot is ready to receive and response new messages! Go ahead and start talking to the bot.

I also will let a bot deploy as example available, if you would like to check the final result: https://t.me/TelebotEchoBot

## Troubleshooting
If you have not succeeded, try first understand how Cloud Functions works. The [documentation](https://cloud.google.com/functions/docs "documentation") is very well written and is available in many languages. As far as you [deploy a simple hello world example](https://cloud.google.com/functions/docs/first-go "deploy a simple hello world example"), you be able to deploy this guide too.

## Motivation
There is a lot of material about AWS products, as AWS Lambda for example, which works pretty similar Google Cloud Functions. But about Google Cloud Platform there is not.

One day I was a beginner, and I looked for a guide like that. I did this hoping it will help someone who needs it.