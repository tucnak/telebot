## How to deploy

- Export an environment variable `TOKEN` with your bot token obtained from @BotFather
- Export your AWS credentials `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` so terraform can create resources.
- Run `./deploy.sh` (it will apply terraform and make a call to telegram server to set the webhook)
- Enjoy!

Once you are done, you can run `./destory.sh` to remove all infrastructure created.