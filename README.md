# Schedule API

**⚠️ Work in Progress ⚠️**

This project is currently under active development. Features may be incomplete or subject to change.

## Description

This is a RESTful API for scheduling appointments. It allows users to create and manage their availability, and to book appointments with other users.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

*   [Go](https://golang.org/)
*   [Docker](https://www.docker.com/)
*   [Docker Compose](https://docs.docker.com/compose/)

### Installation

1.  Clone the repo
    ```sh
    git clone https://github.com/lopesmarcello/schedule-api.git
    ```
2.  Create a `.env` file from the `.env.example` file and update the environment variables.
    ```sh
    cp .env.example .env
    ```
3.  Build and run the application using Docker Compose.
    ```sh
    docker-compose up --build
    ```

The API will be available at `http://localhost:8080`.

## Usage

The API documentation is not yet available.

## Running the tests

To run the tests, use the following command:

```sh
go test ./...
```

## Built With

*   [Go](https://golang.org/) - The programming language used
*   [Gin](https://gin-gonic.com/) - The web framework used
*   [PostgreSQL](https://www.postgresql.org/) - The database used
*   [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
*   [sqlc](https://sqlc.dev/) - A tool for generating type-safe Go code from SQL
*   [Docker](https://www.docker.com/) - The containerization platform used
