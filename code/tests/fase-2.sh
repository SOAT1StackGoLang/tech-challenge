#!/bin/bash

## import the api.sh file
source ./helpers/api.sh

## this Script will run the tests for the fase 2 of the challenge
## the following steps will be executed:
## 1. check if the api is working
## 2. create 6 orders to validate the status: Recebido|Preparacao|Pronto|Finalizado
## 3. will pay 5 orders
## 4. will change the status of the order 2 to "Recebido"
## 5. will change the status of the order 3 to "Preparacao"
## 6. will change the status of the order 4 to "Pronto"
## 7. will change the status of the order 5 to "Finalizado"
## 8. will change the status of the order 6 to "Preparacao"
## 9. will change the status of the order 6 to "Pronto"
## 10. will change the status of the order 6 to "Finalizado"    
## 11. will list all orders and filter to show only ID, status, payment_status, as a result of success

## this function will run all the tests for the challange fase 2

function run_fase_2() {
    ## check if the api is working
    echo -e "\n-----"
    check_api

    # Create 6 orders
    echo -e "\n-----"
    echo -e "${YELLOW}Creating 6 orders and checkout - FASE 2${NC}"

    create_product_drink > /dev/null
    fullorders=$(create_orders_and_checkout 6)
    echo "$fullorders" | jq -s '.'
    #echo $fullorders | jq -s '.'

    # Pay 5 orders
    echo -e "\n-----"
    echo -e "${YELLOW}Paying 5 orders - FASE 2${NC}"

    oders_id_with_payment_id=$(echo "$fullorders" | jq -r '. | {order_id: .order.id, payment_id: .payment_info.payment_id}'| jq -s '.')
    #echo $oders_id_with_payment_id | jq


    ## this function will pay the orders
    delivery_2 "$oders_id_with_payment_id" #> /dev/null

    # this function will wrapper all the move order to next status
    delivery_4 "$oders_id_with_payment_id" #> /dev/null

    # List all orders
    echo -e "\n-----"
    echo -e "${YELLOW}Listing all Valid Orders - FASE 2${NC}"
    list_all_orders | jq -r '.[] | {id: .id, status: .status}'
}

## this function will pay the orders

function delivery_2() {
    local list_orders_id_with_payment_id=$1

    for order in $(echo "${oders_id_with_payment_id}" | jq -c '.[]'); do
        order_id=$(echo "$order" | jq -r '.order_id')
        payment_id=$(echo "$order" | jq -r '.payment_id')
        #echo "pay_order \"$order_id\" \"$payment_id\""
        pay_order "$order_id" "$payment_id" | jq -r '. | {order_id: .order_id, status_do_pagamento: ."status do pagamento"}'
        sleep 1
    done

}

## this function will wrapper all the move order to next status
## no paramters
## no return

function delivery_4() {
    local list_orders_id_with_payment_id=$1
    # Change status of order 2 to "Preparacao"
    echo -e "\n---------------------------------------------------------------"
    echo -e "${YELLOW}Changing status of order 1 to \"Recebido\" - FASE 2${NC}"
    order_id_1=$(echo "$oders_id_with_payment_id" | jq -r '.[0].order_id')
    move_order_to_next_status "$order_id_1" "Recebido" | jq -r '. | {id: .id, status: .status}'
    sleep 2

    echo -e "${YELLOW}Changing status of order 2 to \"Preparacao\" - FASE 2${NC}"
    order_id_2=$(echo "$oders_id_with_payment_id" | jq -r '.[1].order_id')
    move_order_to_next_status "$order_id_2" "Preparacao" | jq -r '. | {id: .id, status: .status}'
    sleep 2

    echo -e "${YELLOW}Changing status of order 3 to \"Pronto\" - FASE 2${NC}"
    order_id_3=$(echo "$oders_id_with_payment_id" | jq -r '.[2].order_id')
    move_order_to_next_status "$order_id_3" "Pronto" | jq -r '. | {id: .id, status: .status}'
    sleep 2

    echo -e "${YELLOW}Changing status of order 4 to \"Finalizado\" - FASE 2${NC}"
    order_id_4=$(echo "$oders_id_with_payment_id" | jq -r '.[3].order_id')
    move_order_to_next_status "$order_id_4" "Finalizado" | jq -r '. | {id: .id, status: .status}'
    sleep 2

    echo -e "${YELLOW}Listing all Valid Orders${NC}"
    list_all_orders | jq -r '.[] | {id: .id, status: .status}'
    sleep 2
    
    echo -e "${YELLOW}Changing status of order 5 to \"Recebido\" - FASE 2${NC}"
    order_id_5=$(echo "$oders_id_with_payment_id" | jq -r '.[4].order_id')
    move_order_to_next_status "$order_id_5" "Recebido" | jq -r '. | {id: .id, status: .status}'

    echo -e "${YELLOW}Listing all Valid Orders${NC}"
    echo -e "${YELLOW} Moving Order 5 - ID: $order_id_5 - FASE 2${NC}"
    list_all_orders | jq -r '.[] | {id: .id, status: .status}'
    sleep 2

    echo -e "${YELLOW}Changing status of order 5 to \"Preparacao\" - FASE 2${NC}"
    order_id_5=$(echo "$oders_id_with_payment_id" | jq -r '.[4].order_id')
    move_order_to_next_status "$order_id_5" "Preparacao" | jq -r '. | {id: .id, status: .status}'

    echo -e "${YELLOW}Listing all Valid Orders${NC}"
    echo -e "${YELLOW} Moving Order 5 - ID: $order_id_5 - FASE 2${NC}"
    list_all_orders | jq -r '.[] | {id: .id, status: .status}'
    sleep 2

    echo -e "${YELLOW}Changing status of order 5 to \"Pronto\" - FASE 2${NC}"
    order_id_5=$(echo "$oders_id_with_payment_id" | jq -r '.[4].order_id')
    move_order_to_next_status "$order_id_5" "Pronto" | jq -r '. | {id: .id, status: .status}'

    echo -e "${YELLOW}Listing all Valid Orders${NC}"
    echo -e "${YELLOW} Moving Order 5 - ID: $order_id_5 - FASE 2${NC}"
    list_all_orders | jq -r '.[] | {id: .id, status: .status}'
    sleep 2

    echo -e "${YELLOW}Changing status of order 5 to \"Finalizado\" - FASE 2${NC}"
    order_id_5=$(echo "$oders_id_with_payment_id" | jq -r '.[4].order_id')
    move_order_to_next_status "$order_id_5" "Finalizado" | jq -r '. | {id: .id, status: .status}'
    sleep 2

    echo -e "${YELLOW}Listing all Valid Orders${NC}"
    echo -e "${YELLOW} Moving Order 5 - ID: $order_id_5 - FASE 2${NC}"
    list_all_orders | jq -r '.[] | {id: .id, status: .status}'
    sleep 2
}





