## Overview
This project utilizes the popular 'echo' library for configuring HTTP REST endpoints as
well as the popular 'testify' library for assertion-driven testing.

## Build
Install external dependencies by running the following commands:
* `go get "github.com/labstack/echo"`
* `go get "github.com/dgrijalva/jwt-go"`
* `go get "github.com/stretchr/testify/assert"`

## Running unit tests
Run `go test`

## REST API Standards:
Create Address
POST /address

* Accepts JSON payload
* On success 201 - Created
* On error 500 - Internal Server Error

Get Address
GET /address/:id

* On success 200 - OK
* On error 404 - Not Found if address is not found otherwise 500 - Internal Server Error

List Addresses
GET /address

* On success 200 - OK

Modify Address
PUT /address/:id

* On success 200 - OK
* On error 404 - Not Found if address is not found otherwise 500 - Internal Server Error

Delete Address
DELETE /address/:id

* On success/error 200 - OK