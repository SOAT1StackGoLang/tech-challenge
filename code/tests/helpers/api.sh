#!/bin/bash
## bash best pratices
#set -euo pipefail
IFS=$'\n\t'

## default values
API_URL="http://localhost:8000"

## Default app values
USER_ID="123e4567-e89b-12d3-a456-426614174000"

## colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

## This function will wrapper the curl command to make the requests to the API
## It will receive the method, the endpoint and the body as parameters
## It will return the response body
## will also validate the response code and exit if it is not 200

function make_request() {
    local method=$1
    local endpoint=$2
    local body=$3
    local response=$(curl -s -X "$method" --location "$API_URL$endpoint" \
        --header 'Content-Type: application/json' \
        --data-raw "$body")
    local response_code=$?
    if [ "$response_code" -ne 0 ]; then
        echo -e "${RED}Failed to connect to API after 5 retries.${NC}"
        exit 1
    fi
    echo "$response"
}


## This function will check if the API is working
## It will return the response body
function check_api() {
    echo -e "${YELLOW}Check if api is working${NC}"
    make_request "POST" "/v1/users/validate" '{"document": "97580053080"}'
}

## This function will get the admin user
## It will return the response body
function get_admin_user() {
    echo -e "${YELLOW}Getting admin user${NC}"
    make_request "POST" "/v1/users/validate" '{"document": "97580053080"}'
}

## This function will create a product under drink category with default values
## values are name, description, price, category_id and user_id
## It will return the response body
CATEGORY_ID_DRINK="a557b0c0-3bcf-11ee-be56-0242ac120002"

function create_product_drink() {
    echo -e "${YELLOW}Creating product${NC}"
    if get_product_drink; then
        echo -e "${YELLOW}Product already exists${NC}"
        echo -e "${YELLOW}Product id: $PRODUCT_ID${NC}"
        return
    fi
    make_request "POST" "/v1/products" "{\"name\": \"Coca Cola Test\", \"description\": \"Coca Cola 350ml\", \"price\": \"5.00\", \"category_id\": \"$CATEGORY_ID_DRINK\", \"user_id\": \"$USER_ID\"}"
}

## this function will get the list of products from drink category
## and will save the product_id of the product with the name of Coca Cola Test in a variable
## use the jq to search for the name and get the id

function get_product_drink() {
    #echo -e "${YELLOW}Getting product${NC}"
    PRODUCT=$(make_request "GET" "/v1/products?category-id=$CATEGORY_ID_DRINK&limit=10&offset=0")
    PRODUCT_ID=$(echo $PRODUCT | jq '.products[] | select(.name == "Coca Cola Test") | .id')
    #echo -e "${YELLOW}Product id: $PRODUCT_ID${NC}"
}

## Create default order with default values, will receive a list of product_ids as parameter
## It will return the response body

function create_order_default() {
    local product_ids=$1
    if [ -z "$product_ids" ]; then
        get_product_drink
        product_ids="$PRODUCT_ID"
    fi
    #echo -e "${YELLOW}Creating order${NC}"
    #echo -e "${YELLOW}Product ids: $product_ids${NC}"
    payload="{\"products_ids\": [$product_ids], \"user_id\": \"$USER_ID\"}"
    order=$(make_request "POST" "/v1/orders" "$payload")
    echo "$order"
}

## Create checkout will recevie the order_id or a list of order_id as parameter
## will do a for each to create the checkout for each order_id, as json order_id and user_id
## It will return the response body for each order_id with the return of the checkout

function create_checkout() {
    local order_id=$1
    if [ -z "$order_id" ]; then
        order_id=$(create_order_default | jq '.id')
    fi
    #echo -e "${YELLOW}Creating checkout${NC}"
    payload="{\"order_id\": $order_id, \"user_id\": \"$USER_ID\"}"
    checkout=$(make_request "POST" "/v1/orders/checkout" "$payload")
    echo "$checkout"
}

## this function will create 5 orders and checkout for each order
## will return a list with the responses of the checkouts generated

function create_orders_and_checkout() {
    local max=${1:-5}
    #echo -e "${YELLOW}Creating orders and checkout${NC}"
    order_ids=()
    for i in $(seq 1 "$max"); do
        order=$(create_order_default)
        order_id=$(echo "$order" | jq '.id')
        order_ids+=("$order_id")
    done
    checkout_responses=()
    for order_id in "${order_ids[@]}"; do
        checkout=$(create_checkout "$order_id")
        checkout_responses+=("$checkout")
    done
    echo "${checkout_responses[@]}"
}

## this function will pay the order using the webhook /v1/webhook/payment-notification
## will receive the order_id, payment_id and approved(boolean default true) as parameters
## will return the response body

function pay_order() {
    local order_id=$1
    local payment_id=$2
    local approved=${3:-true}
    #echo -e "${YELLOW}Paying order${NC}"
    payload="{\"approved\": $approved, \"payment_id\": \"$payment_id\", \"order_id\": \"$order_id\"}"
    payment=$(make_request "POST" "/v1/webhook/payment-notification" "$payload")
    echo "$payment"
}

## this function will move the order to the next status /v1/orders/status-update
## will receive the order_id and status as parameters and user_id as optional
## will return the response body

function move_order_to_next_status() {
    local order_id=$1
    local status=$2
    local user_id=${3:-$USER_ID}
    #echo -e "${YELLOW}Moving order to next status${NC}"
    payload="{\"order_id\": \"$order_id\", \"status\": \"$status\", \"user_id\": \"$user_id\"}"
    status_update=$(make_request "PUT" "/v1/orders/status-update" "$payload")
    echo "$status_update"
}

