#!/usr/bin/env bash
gcloud functions deploy TwilioHandler --runtime go111 --trigger-http
