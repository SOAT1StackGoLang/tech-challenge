# Test Script README

This script is designed to test the functionality of an API that manages users, categories, products, and payments. 
But only testing flow funcionality , not to be confused with automated testing or test coverage tool.

## The script performs the following actions:

1. Checks if the `jq` command-line tool is installed.
2. Gets an admin user.
   - This tests the API's ability to retrieve an existing user by their ID.
3. Creates a new user.
   - This tests the API's ability to create a new user with the given data.
4. Gets a user by document.
   - This tests the API's ability to retrieve an existing user by their document number.
5. Creates a new category.
   - This tests the API's ability to create a new category with the given data.
6. Gets a category by ID.
   - This tests the API's ability to retrieve an existing category by its ID.
7. Deletes a category by ID.
   - This tests the API's ability to delete an existing category by its ID.
8. Creates a new product.
   - This tests the API's ability to create a new product with the given data.
9. Edits an existing product.
   - This tests the API's ability to edit an existing product with the given data.
10. Gets a product by ID.
    - This tests the API's ability to retrieve an existing product by its ID.
11. Creates a second product.
    - This tests the API's ability to create a new product with the given data.
12. Lists products by category.
    - This tests the API's ability to retrieve a list of products by category ID.
13. Deletes a product by ID.
    - This tests the API's ability to delete an existing product by its ID.
14. Tests payment functionality.
    - This tests the API's ability to process payments for products.

To run the script, simply execute the `autotest.sh` file in a Bash terminal. The script will output the results of each test to the console, including any errors or failures.

Note: The script assumes that the API is running on `localhost:8000`. If the API is running on a different host or port, you will need to modify the script accordingly.