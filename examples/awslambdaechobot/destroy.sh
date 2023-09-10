#!/bin/bash
set -e 
set -o pipefail


terraform -chdir=./tf destroy -input=false -var=bot_token=$TOKEN 

curl -H "Content-Type: application/json"  -X POST "https://api.telegram.org/bot$TOKEN/setWebhook" -d "{
     \"url\": \"\" 
     }"