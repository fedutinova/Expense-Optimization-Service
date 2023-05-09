# Expense Optimization Service

This is a backend service written in Go for optimizing individual and family expenses. Users can create individual or family wallets, with family wallet owners being able to add or remove members. The service can calculate the total income and expenses for a month in general and by categories, as well as determine the current balance.

## Features
- Create individual or family wallets
- Family wallet owners can add or remove members
- Calculate total income and expenses for a month in general and by categories
- Determine the current balance
- Secure authentication and authorization mechanisms

## Benefits
- Helps individuals and families optimize their expenses
- Provides clear insights into income and expenses
- Enables better financial planning and decision-making
- Offers a flexible and scalable solution for managing finances

## Requirements
- Go 1.20 or higher
- PostgreSQL database
- JWT authentication

## Installation
1. Clone the repository
2. Install dependencies using `make build && make run`
3. Run the server using `go run cmd/main.go`
4. If the application is running for the first time `make migrate`

## API Endpoints

- `POST /sign-up` - register a new user
- `POST /sign-in` - login to the service and obtain a JWT token
- `POST /api/wallets` - create a new wallet
- `GET /api/wallets` - retrieve all wallets for the current user
- `GET /api/wallets/:id` - retrieve a specific wallet by ID
- `PUT /api/wallets` - update an existing wallet by ID
- `DELETE /api/wallets/:id` - delete a wallet by ID
- `DELETE /api/wallets` - delete a member from a wallet by ID
- `POST /api/wallets/:id/transactions` - create a new transaction for a specific wallet
- `GET /api/wallets/:id/transactions` - retrieve all transactions for a specific wallet
- `GET /api/wallets/:id/transactions/incomes` - retrieve all income transactions for a specific wallet
- `GET /api/wallets/:id/transactions/incomes/category` - retrieve income transactions for a specific category in a specific wallet
- `GET /api/wallets/:id/transactions/expenses` - retrieve all expense transactions for a specific wallet
- `GET /api/wallets/:id/transactions/expenses/category` - retrieve expense transactions for a specific category in a specific wallet
- `PUT /api/transactions/:id` - update an existing transaction by ID
- `DELETE /api/transactions/:id `- delete a transaction by ID

## Conclusion
The Expense Optimization Service is a powerful backend solution for individuals and families looking to optimize their expenses and gain insights into their financial situation. With secure authentication and authorization mechanisms and RESTful APIs, the service is flexible and scalable, and can easily integrate with other systems.