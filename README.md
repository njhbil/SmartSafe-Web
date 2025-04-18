# SmartSafe Project

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
  - [Backend (Go)](#backend-go)
  - [Frontend (Svelte)](#frontend-svelte)
  - [Redis Setup](#redis-setup)
  - [PostgreSQL Setup](#postgresql-setup)
    - [Create Tables](#create-tables)
- [Contributing](#contributing)
- [License](#license)

## Features

- Secure storage of sensitive information
- User authentication and authorization
- Responsive and intuitive UI built with Svelte
- High performance backend powered by Go

## Requirements

- [Go](https://golang.org/doc/install) (version 1.16 or higher)
- [Bun](https://bun.sh/) (version 0.1.10 or higher)
- [Redis](https://redis.io/) (version 6.2 or higher)
- [PostgreSQL](https://www.postgresql.org/) (version 13 or higher)

## Installation

### Backend (Go)

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/SmartSafe.git
    cd SmartSafe/server
    ```

2. Install dependencies:
    ```sh
    go mod download
    ```

3. Create a `.env` file based on the provided template:
    ```sh
    cp .envExample .env
    ```

4. Update the `.env` file with your configuration details.

    Example `.env` file:
    ```
    DATABASE_URL=postgres://user:password@localhost:5432/smartsafe
    REDIS_URL=redis://localhost:6379
    ```

5. Run the backend server:
    ```sh
    go run main.go
    ```

### Frontend (Svelte)

1. Navigate to the frontend directory:
    ```sh
    cd ../SmartSafe
    ```

2. Install dependencies:
    ```sh
    bun install
    ```

3. Run the frontend development server:
    ```sh
    bun run dev
    ```

### Redis Setup

1. Install Redis:
    - On macOS:
        ```sh
        brew install redis
        ```
    - On Linux:
        ```sh
        sudo apt update
        sudo apt install redis-server
        ```
    - On Windows: Download Redis from the [official website](https://redis.io/docs/getting-started/installation/).

2. Start the Redis server:
    ```sh
    redis-server
    ```

3. Test the Redis connection:
    ```sh
    redis-cli ping
    ```

    You should see:
    ```
    PONG
    ```

### PostgreSQL Setup

1. Install PostgreSQL:
    - On macOS:
        ```sh
        brew install postgresql
        ```
    - On Linux:
        ```sh
        sudo apt update
        sudo apt install postgresql postgresql-contrib
        ```
    - On Windows: Download PostgreSQL from the [official website](https://www.postgresql.org/download/).

2. Start the PostgreSQL service:
    ```sh
    sudo service postgresql start
    ```

3. Access the PostgreSQL shell:
    ```sh
    psql -U postgres
    ```

4. Create a new database for the project:
    ```sql
    CREATE DATABASE smartsafe;
    ```

5. Create a new user and grant privileges:
    ```sql
    CREATE USER smartsafe_user WITH PASSWORD 'yourpassword';
    GRANT ALL PRIVILEGES ON DATABASE smartsafe TO smartsafe_user;
    ```

#### Create Tables

After setting up your PostgreSQL database, create the required tables for the project.

1. Switch to the `smartsafe` database:
    ```sh
    \c smartsafe
    ```

2. Create the `accounts` table:
    ```sql
    CREATE TABLE accounts (
        username VARCHAR(100),
        email VARCHAR(255),
        password VARCHAR(255),
    );
    ```

3. Create the `verifytokens` table:
    ```sql
    CREATE TABLE verifytokens (
        token VARCHAR(255),
    );
    ```

4. Verify the tables were created:
    ```sql
    \dt
    ```

You should see:
```
          List of relations
 Schema |   Name       | Type  |  Owner
--------+--------------+-------+----------
 public | accounts     | table | postgres
 public | verifytokens | table | postgres
(2 rows)
```

## Contributing

We welcome contributions! Please follow the standard [GitHub workflow](https://guides.github.com/introduction/flow/) for submitting pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
