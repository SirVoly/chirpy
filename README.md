# Chirp 

An RSS feed aggregator in Go, called "Gator"

This project was build with the intention of hands-on practice in GO, HTTP Servers, RESTful API and a secure authentication/authorization system.

## Technologies Used

*   **GO**: The primary programming language.
*   **PSQL**: The primary database.
*   **Goose**: A database migration tool written in Go

## Setup

To run this project, follow these steps:

    Copy the repo on your pc
    Install [GO](https://go.dev/doc/install)
        webi go@1.25.4
    Install PostGres:
        sudo apt install postgresql postgresql-contrib

    Setup Database:
        sudo service postgresql start
        su postgres
        psql
        CREATE DATABASE chirpy;
        \c chirpy
    
    Install Goose:
        go install github.com/pressly/goose/v3/cmd/goose@latest
        cd sql/schema && goose postgres postgres://postgres:postgres@localhost:5432/chirpy up && cd ../..

    Setup your environment file:
        .env file in the root of the project
        DB_URL= Your database string
        PLATFORM= What you use it for, example: user
        JWTSECRET= A long random secret for your tokens
        POLKA_KEY= A key to verify when a request comes from Polka

## Usage

This server application is build to let people post their chirps.
Most of the interactions are mainly supported from an API standpoint, not a browser standpoint.

All endpoints can be found in server.go, with their method and needs.

# Credit
This project was completed as part of a guided course on [Boot.dev](https://www.boot.dev).
It was build following along with the [Learn HTTP Servers in Go](https://www.boot.dev/courses/learn-http-servers-golang) course.