variable "name" {
  type        = string
  default     = "awslambdaechobot" # basename(path.root)
  description = "Name of the binary (basename of this directory)"
}
variable "token" {
  type        = string
  description = "Telegram bot token"
}
variable "region" {
  type        = string
  description = "AWS deployment region"
}

locals {
  exe      = var.name
  zip      = "${local.exe}.zip"
  func     = "${var.name}Function"
  api      = "${var.name}API"
  role     = "${var.name}Role"
  policy   = "${var.name}Policy"
  perm     = "${var.name}Permission"
  log      = "/aws/lambda/${local.func}"
  envToken = "TELEBOT_SECRET"
}

# Run ./init.sh to install github.com/yi-jiayu/terraform-provider-telegram
provider "telegram" {
  bot_token = var.token
}
resource "telegram_bot_webhook" "a" {
  url             = aws_apigatewayv2_api.a.api_endpoint
  max_connections = 100
}

provider "archive" {
  version = "~> 1.3"
}
data "archive_file" "a" {
  type        = "zip"
  source_file = local.exe
  output_path = local.zip
}

provider "aws" {
  version = "~> 2.59"
  region  = var.region
}

resource "aws_apigatewayv2_api" "a" {
  name          = local.api
  protocol_type = "HTTP"
  target        = aws_lambda_function.a.arn
  route_key     = "POST /"
}
resource "aws_lambda_permission" "a" {
  statement_id  = local.perm
  function_name = aws_lambda_function.a.function_name
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.a.execution_arn}/*/*/*" # Any stage, method, resource.
}

resource "aws_lambda_function" "a" {
  function_name = local.func
  runtime       = "go1.x"
  handler       = local.exe
  memory_size   = 128 # MB, 128 + 64*x
  timeout       = 60  # seconds
  role          = aws_iam_role.a.arn

  filename         = data.archive_file.a.output_path
  source_code_hash = data.archive_file.a.output_base64sha256

  environment {
    variables = {
      (local.envToken) = var.token
    }
  }
}
resource "aws_cloudwatch_log_group" "a" {
  # AWS Lambda automatically logs to the group with this name.
  name              = local.log
  retention_in_days = 1
}

resource "aws_iam_role_policy_attachment" "a" {
  role       = aws_iam_role.a.name
  policy_arn = aws_iam_policy.a.arn
}
resource "aws_iam_role" "a" {
  name               = local.role
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
}
data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}
resource "aws_iam_policy" "a" {
  name   = local.policy
  policy = data.aws_iam_policy_document.a.json
}
data "aws_iam_policy_document" "a" {
  # Based on AWSLambdaBasicExecutionRole.
  statement {
    actions   = ["logs:CreateLogStream", "logs:PutLogEvents"]
    resources = [aws_cloudwatch_log_group.a.arn]
  }
}
