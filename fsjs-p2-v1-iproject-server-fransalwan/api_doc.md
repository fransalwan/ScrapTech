#Endpoints
List of Avaliable Endpoints:

# Customer

- POST /register
- POST /login
- POST /OAuthLogin
- GET /materials
- GET /categories
- GET /order
- POST /order/:id
- DELETE /order/:id
- PATCH /payment
- POST /generate-midtrans-token

# POST / register

## Description

      -  Register Customer

## Response

200 - OK - Body
{
"id": 5,
"email": "frans13@gmail.com"
}

400 - Bad Request

{
"message": string
}

# POST / login

## Description

      -  Login Customer

## Response

200 - OK - Body
{
"username": "username"
"access_token": "access_token"
}

401 - Not Authorized

{
"message": string
}

# GET / materials

## Description

      -  Get All Material data

## Response

200 - OK - Body

# POST /categories

## Description

      -  Add New Categories data

## Response

200 - OK - Body
{
"message": {
"id": 3,
"name": "baru",
"updatedAt": "2023-07-14T03:23:20.791Z",
"createdAt": "2023-07-14T03:23:20.791Z"
}
}

# GET /order

## Description

      - get all order by customer

## Response

200 - OK

404 - Not Found
Body
{
"message": "Error not found"
}

# POST /order/:id

## Description

      - add order by material id

## Response

200 - OK
Body

404 - Not Found
Body
{
"message": "Error not found"
}

# DELETE /order/:id

## Description

      - delete order by material id

## Response

200 - OK
Body

404 - Not Found
Body
{
"message": "Error not found"
}

# PATCH /payment

## Description

      - Edit status materials and status order

## Response

200 - OK
Body

404 - Not Found
Body
{
"message": "Error not found"
}

- POST /generate-midtrans-token

# POST /generate-midtrans-token

## Description

      -       - generate midtrans token

## Response

200 - OK
Body

404 - Not Found
Body
{
"message": "Error not found"
}

# Global Error

500 - Internal Server Error

- Body
  {
  "message": "Internal Server Error"
  }

// coba ulik Open API biar lebih rapih doc nya
