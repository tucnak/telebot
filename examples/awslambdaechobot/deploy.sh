#!/bin/bash
set -e 
set -o pipefail


terraform -chdir=./tf apply -input=false -var=bot_token=$TOKEN 

export FUNCTION_URL=$(terraform -chdir=./tf output function_url)  

curl -H "Content-Type: application/json"  -X POST "https://api.telegram.org/bot$TOKEN/setWebhook" -d "{
     \"url\": $FUNCTION_URL
     }"