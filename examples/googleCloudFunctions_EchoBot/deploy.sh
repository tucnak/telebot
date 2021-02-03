#!/bin/bash
gcloud functions deploy telebot \
       --entry-point EchoBot \
       --runtime go113 \
       --trigger-http \
       --allow-unauthenticated \
       --env-vars-file env.yaml \
       --memory 128MB \
       --timeout 180s \
       --max-instances 1
