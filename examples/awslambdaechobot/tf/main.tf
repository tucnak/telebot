terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"

}

provider "aws" {
  region = "eu-west-2"
}

//https://github.com/terraform-aws-modules/terraform-aws-lambda/tree/v6.0.0
module "bot_handler" {
  source = "terraform-aws-modules/lambda/aws"

  function_name              = "echo-bot"
  description                = "This is a lambda funtion that will respont to every message you send to the bot with the same content"
  runtime                    = "provided.al2"
  create_lambda_function_url = true
  architectures              = ["x86_64"]
  memory_size                = "128"
  timeout                    = "1"
  handler                    = "bootstrap"
  source_path = [{
    path     = "../"
    commands = ["GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o package/bootstrap", ":zip ../package/bootstrap"]
  }]

  environment_variables = {
    TOKEN : var.bot_token,
  }
  cloudwatch_logs_retention_in_days = 1
}


variable "bot_token" {
  description = "telegram bot token"
  type        = string
  sensitive   = true
}

output "function_url" {
  value = module.bot_handler.lambda_function_url
}
