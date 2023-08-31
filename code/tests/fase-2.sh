#!/bin/bash

## import the api.sh file
source ./helpers/api.sh

## this Script will run the tests for the fase 2 of the challenge
## the following steps will be executed:
## 1. check if the api is working
## 2. create 6 orders to validate the status: Recebido|Preparacao|Pronto|Finalizado|Cancelado
## 3. will pay 5 orders
## 4. will refuse the last(6) order and move to Cancelado status
## 5. will change the status of the order 2 to "Preparacao"
## 6. will change the status of the order 3 to "Pronto"
## 7. will change the status of the order 4 to "Finalizado"
## 8. will change the status of the order 5 to "Preparacao"
## 9. will change the status of the order 5 to "Pronto"
## 10. will change the status of the order 5 to "Finalizado"    
## 11. will list all orders and filter to show only ID, status, payment_status, as a result of success

## this function will run all the tests for the challange fase 2

function run_tests_fase_2() {
    ## check if the api is working
    echo -e "\n-----"
    check_api

    # Create 6 orders
    echo -e "\n-----"
    echo -e "${YELLOW}Creating 6 orders and checkout - FASE 2${NC}"

    create_product_drink
    fullorders=$(create_orders_and_checkout 6)
    #echo "$fullorders" | jq -s '.'
    echo $fullorders | jq -s '.'

    # Pay 5 orders
    echo -e "\n-----"
    echo -e "${YELLOW}Paying 5 orders - FASE 2${NC}"

    oders_id_with_payment_id=$(echo "$fullorders" | jq -r '. | {order_id: .order.id, payment_id: .payment_info.payment_id}'| jq -s '.')
    echo $oders_id_with_payment_id | jq

    for order in $(echo "${oders_id_with_payment_id}" | jq -c '.[]'); do
        order_id=$(echo "$order" | jq -r '.order_id')
        payment_id=$(echo "$order" | jq -r '.payment_id')
        echo "pay_order \"$order_id\" \"$payment_id\""
        #pay_order "$order_id" "$payment_id" | jq
        #sleep 2
    done


}





