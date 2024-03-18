# Tutorial

A simple (client > server) project in Golang.

The Server provides a  proxy that calls [this third party API](https://economia.awesomeapi.com.br/json/last/USD-BRL).
This URL supplies an exchange rate from US dollar to Brazilian real.

Server also inserts the retrieved information from the external source into a SQLite DB.

The client makes a call to the server and inserts the retrieved value into a .txt file.

## Before you start

Make sure that:
- You have GCC installed
  ```bash
  choco install mingw
   ```
- Sets CGO_ENABLED to 1
    ```bash
    go env -w CGO_ENABLED=1
   ```

## Execution

## If this is your first time running, you could execute the migrations:

1. Execute the following command in the terminal, go to server directory and run:

   ```bash
    go run server.go -migratedb=true
   ```
   The flag `-migratedb=true` ensure that the migrations will be executed
2. Open another Terminal, go to client directory and run:
  ```bash
      go run client.go
  ```

## If not:

1. Just go to server directory and run:
    ```bash
        go run server.go
    ```

2. Open another Terminal, go to client directory and run:
    ```bash
        go run client.go
    ```

