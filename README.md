# SmartSafe Project

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
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

## Installation

### Backend (Go)

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/SmartSafe.git
    cd SmartSafe/backend
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

## Contributing

We welcome contributions!

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
