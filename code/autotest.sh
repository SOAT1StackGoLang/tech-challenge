#!/bin/bash

# Define ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Check if jq is installed
if ! command -v jq &> /dev/null
then
    echo -e "${RED}jq could not be found. Please install jq to parse JSON output.${NC}"
    exit
fi

# Get admin user
echo -e "${YELLOW}Getting admin user${NC}"
curl -s --location 'localhost:8000/v1/users/validate/97580053080' | jq
sleep 2

# Create user
echo -e "\n-----"
echo -e "${YELLOW}Creating User${NC}"
curl -s --location 'localhost:8000/v1/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Name": "Test User",
    "Document": "97580053081",
    "Email": "test@user.com"
}' | jq
sleep 2

# Get user by document
echo -e "\n-----"
echo -e "${YELLOW}Getting user by document${NC}"
curl -s --location 'localhost:8000/v1/users/validate/97580053081' | jq
sleep 2

# Category
echo -e "\n-----"
echo -e "${YELLOW}Category${NC}"

# Create category
echo -e "\n-----"
echo -e "${YELLOW}Creating Category${NC}"
CATEGORY=$(curl -s --location 'localhost:8000/v1/categories' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Drink",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
}')
echo -e "$CATEGORY" | jq
CATEGORY_ID=$(echo $CATEGORY | jq -r '.id')
sleep 2

# Get category
echo -e "\n-----"
echo -e "${YELLOW}Getting Category${NC}"
CATEGORY=$(curl -s --location localhost:8000/v1/categories/$CATEGORY_ID)
if command -v jq &> /dev/null; then
    echo $CATEGORY | jq
else
    echo $CATEGORY
fi

# Delete category
echo -e "\n-----"
echo -e "${YELLOW}Deleting Category${NC}"
curl -s --location --request DELETE 'localhost:8000/v1/categories' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "'$CATEGORY_ID'",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
}'
sleep 2

# Get category after delete
echo -e "\n-----"
echo -e "${YELLOW}Getting Category after delete${NC}"
CATEGORY=$(curl -s --location localhost:8000/v1/categories/$CATEGORY_ID)
echo $CATEGORY
sleep 2
echo -e "\n-----"

# Product
echo -e "\n-----"
echo -e "${YELLOW}Product${NC}"

# Create category
echo -e "\n-----"
echo -e "${YELLOW}Creating Category${NC}"
CATEGORY=$(curl -s --location 'localhost:8000/v1/categories' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Drink",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
}')
echo -e "$CATEGORY" | jq
CATEGORY_ID=$(echo $CATEGORY | jq -r '.id')
sleep 2

# Create product
echo -e "\n-----"
echo -e "${YELLOW}Creating Product${NC}"
PRODUCT=$(curl -s --location 'localhost:8000/v1/products' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Coca Cola",
    "description": "Coca Cola 350ml",
    "price": "5.00",
    "category_id": "'$CATEGORY_ID'",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
}')
echo -e "$PRODUCT" | jq
PRODUCT_ID=$(echo $PRODUCT | jq -r '.id')
sleep 2

# Edit product
echo -e "\n-----"
echo -e "${YELLOW}Editing Product${NC}"
curl -s --location --request PUT 'localhost:8000/v1/products' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "'$PRODUCT_ID'",
    "name": "Coca Cola 350ml",
    "description": "Coca Cola 350ml Lata",
    "price": "5.00",
    "category_id": "'$CATEGORY_ID'",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
}' | jq .
sleep 2

# Get product
echo -e "\n-----"
echo -e "${YELLOW}Getting Product${NC}"
PRODUCT=$(curl -s --location localhost:8000/v1/products/$PRODUCT_ID)
if command -v jq &> /dev/null; then
    echo $PRODUCT | jq
else
    echo $PRODUCT
fi

# Create product 2
echo -e "\n-----"
echo -e "${YELLOW}Creating Product${NC}"
PRODUCT2=$(curl -s --location 'localhost:8000/v1/products' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Coca Cola 2L",
    "description": "Coca Cola 2L",
    "price": "5.00",
    "category_id": "'$CATEGORY_ID'",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
}')
echo -e "$PRODUCT2" | jq
PRODUCT_ID2=$(echo $PRODUCT2 | jq -r '.id')
sleep 2

# List products by category
echo -e "\n-----"
echo -e "${YELLOW}Listing Products by Category${NC}"
PRODUCTS=$(curl -s --location localhost:8000/v1/products/$CATEGORY_ID/10/0)
if command -v jq &> /dev/null; then
    echo $PRODUCTS | jq
    echo $PRODUCTS
else
    echo $PRODUCTS
fi



# Delete product
echo -e "\n-----"
echo -e "${YELLOW}Deleting Product${NC}"
curl -s --location --request DELETE 'localhost:8000/v1/products' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "'$PRODUCT_ID'",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
}'
sleep 2

# Get product after delete
echo -e "\n-----"
echo -e "${YELLOW}Getting Product after delete${NC}"
PRODUCT=$(curl -s --location localhost:8000/v1/products/$PRODUCT_ID)
echo $PRODUCT
sleep 2
echo -e "\n-----"