#!/bin/bash
# basic script to run some basic api testing using curl

# set the base url
base_url="http://localhost:8081"

# generate a random valid uuid

# Create uuid to use for PaymentID
uuidgen=$(uuidgen)

# Create a payment
curl -s -L -m 5 -X POST "$base_url/payments" \
-H "Content-Type: application/json" \
-d "{
    \"payment\": {
        \"id\": \"$uuidgen\",
        \"price\": \"100.00\",
        \"orderID\": \"3fa85f64-5717-4562-b3fc-2c963f66afa6\"
    }
}" | jq

# Get the payment
curl -s -L -m 5 -X GET "$base_url/payments" \
-H "Content-Type: application/json" \
-d "{
    \"payment_id\": \"$uuidgen\"
}" | jq

# Update the payment its not easily tested because the main usage is in the background task using redis