#!/bin/bash

gcloud functions deploy $FUNCTION_NAME \
    --gen2 \
    --region $DEFAULT_REGION \
    --entry-point $FUNCTION_ENTRYPOINT \
    --runtime go122 \
    --trigger-http \
    --memory 128Mi \
    --cpu 0.083 \
    --timeout 30 \
    --max-instances 1 \
    --allow-unauthenticated \
    --run-service-account=$FUNCTION_RUNNER_SERVICE_ACCOUNT_EMAIL \
    --set-secrets=LINE_CHANNEL_TOKEN=$LINE_CHANNEL_TOKEN_SECRET_ID,LINE_CHANNEL_SECRET=$LINE_CHANNEL_SECRET_SECRET_ID \
    --set-env-vars=LOG_LEVEL=debug
