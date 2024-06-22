#!/bin/bash

gcloud functions deploy callback \
    --gen2 \
    --region asia-northeast1 \
    --entry-point Entrypoint \
    --runtime go122 \
    --trigger-http \
    --memory 128Mi \
    --cpu 0.083 \
    --timeout 30 \
    --max-instances 1 \
    --allow-unauthenticated
