### Run

if you are using docker-compose:
 - docker-compose up --build -d

if you are not using docker-compose:
 - set up config.yaml to your own local config

Then:
 - run 
 ```
 go build -o bin/ordent
 ```
 - run 
 ```
 bin/ordent
 ```

### Flow

There's 3 group API:
- user
- product
- transaction

### User

User's API including:
- login
- register
- logout
- check wallet: checking the amount of money the user had
- add wallet: add money to wallet of the user

Product's API including:
- insert product (admin)
- update product (admin)
- delete product (admin)
- get product
- get all products with pagination
- get all products by tag with pagination
- search products with pagination

Transaction's API including:
- Create transaction: will mutate product stock and sold
- Get all transactions by user
- Get all transactions (admin)
