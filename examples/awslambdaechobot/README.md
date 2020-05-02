This example shows how to write an echo telebot for AWS Lambda and how to launch it using Terraform.

This bot is different from a typical bot in two ways:

1. It is configured with `Settings.Synchronous = true`. This disables asynchronous handlers to let Lambda wait for their completion:

    ```go
    b, _ := tb.NewBot(tb.Settings{Token: token, Synchronous: true})
    ```

2. Instead of `Settings.Poller` and `bot.Start` it calls `bot.ProcessUpdate` inside `lambda.Start`:

    ```go
    lambda.Start(func(req events.APIGatewayProxyRequest) (err error) {
        var u tb.Update
        if err = json.Unmarshal([]byte(req.Body), &u); err == nil {
            b.ProcessUpdate(u)
        }
        return
    })
    ```

To launch the bot [install Terraform](https://www.terraform.io/downloads.html), run [`./init.sh`](init.sh) and then [`./deploy.sh`](deploy.sh). To tear down the cloud infrastructure run `terraform destroy`.
