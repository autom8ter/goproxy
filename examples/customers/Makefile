callback:
	gcloud functions deploy CustomerCallback --runtime go111 --trigger-http

proxy:
	gcloud functions deploy CreateCustomer --set-env-vars ACCOUNT=XXXXXX,KEY=XXXXXX--runtime go111 --trigger-http
