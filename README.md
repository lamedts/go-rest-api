# Simple Go Restful API

## Introduction

1. A RESTful HTTP API listening to port `8080`
2. 3 endpoints
    - create an order
    - take an order
    - list orders
3. Using direction api from google api
4. Mysql as DB
5. using `start.sh` to do all the initialisation and installation

## To run
```sh
# 0. go to project folder
# 1. install dependencies
$ dep init
$ dep ensure
# 2.1 run in debug mode
$ mkdir devConfig
$ cp config/config.yaml ./devConfig
$ APP_ENV=DEBUG go run ../cmd/yay/main.go
# 2.2 run in prod mode
# add config yaml in docker folder
$ bash start.sh
```

## Api interface example

#### Place order

  - Method: `POST`
  - URL path: `/order`
  - Request body:

    ```
    {
        "origin": ["START_LATITUDE", "START_LONGTITUDE"],
        "destination": ["END_LATITUDE", "END_LONGTITUDE"]
    }
    ```

  - Response:

    Header: `HTTP 200`
    Body:
      ```
      {
          "id": <order_id>,
          "distance": <total_distance>,
          "status": "UNASSIGN"
      }
      ```
    or 
    
    Header: `HTTP 500`
    Body:
      ```json
      {
          "error": "ERROR_DESCRIPTION"
      }
      ```

#### Take order

  - Method: `PUT`
  - URL path: `/order/:id`
  - Request body:
    ```
    {
        "status":"taken"
    }
    ```
  - Response:
    Header: `HTTP 200`
    Body:
      ```
      {
          "status": "SUCCESS"
      }
      ```
    or
    
    Header: `HTTP 409`
    Body:
      ```
      {
          "error": "ORDER_ALREADY_BEEN_TAKEN"
      }
      ```

#### Order list

  - Method: `GET`
  - Url path: `/orders?page=:page&limit=:limit`
  - Response:

    ```
    [
        {
            "id": <order_id>,
            "distance": <total_distance>,
            "status": <ORDER_STATUS>
        },
        ...
    ]
    ```