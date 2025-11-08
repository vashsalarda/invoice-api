# Invoice API - Go + Fiber + MongoDB

A RESTful API with Vertical Slice Architecture and CQRS pattern.

---

## Project Structure

```
invoice-api/
├── cmd/api/
│   └── main.go                          # Application entry point
├── internal/
│   ├── server/
│   │   ├── route.go                    # Routes config
│   │   ├── server.go                    # Server config
│   ├── database/
│   │   └── database.go                  # MongoDB connection
│   └── features/
│       ├── user/
│       │   ├── command/
│       │   │   └── command.go
│       │   ├── query/
│       │   │   └── query.go
│       │   ├── controller/
│       │   │   └── controller.go
│       │   ├── model/
│       │   │   └── model.go
│       │   └── route/
│       │       └── route.go
│       ├── customer/
│       │   ├── command/
│       │   │   └── command.go
│       │   ├── query/
│       │   │   └── query.go
│       │   ├── controller/
│       │   │   └── controller.go
│       │   ├── model/
│       │   │   └── model.go
│       │   └── route/
│       │       └── route.go
│       ├── invoice/
│       │   ├── command/
│       │   │   └── command.go
│       │   ├── query/
│       │   │   └── query.go
│       │   ├── controller/
│       │   │   └── controller.go
│       │   ├── model/
│       │   │   └── model.go
│       │   └── route/
│       │       └── route.go
│       └── revenue/
│           ├── command/
│           │   └── command.go
│           ├── query/
│           │   └── query.go
│           ├── controller/
│           │   └── controller.go
│           ├── model/
│           │   └── model.go
│           └── route/
│               └── route.go
├── .env                            # Environment varible
├── .gitignore                      # Git ignore
├── Dockerfile                      # Docker image
├── docker-compose.yml              # Docker orchestration
├── go.mod                          # Go dependencies
├── Makefile                        # Build commands
└── README.md                       # Project documentation
```

## Setup

1. Manage dependencies:

```bash
go mod tidy
```

2. Set up MongoDB:

```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

3. Configure environment variables in `.env`

4. Run the application:

```bash
go run cmd/api/main.go
```

or

```bash
go run ./cmd/api
```

## API Endpoints

### Health Check
- `GET /` - API health check

### Users
- `POST /api/users` - Create user
- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Customers
- `POST /api/customers` - Create customer
- `GET /api/customers` - Get all customers
- `GET /api/customers/:id` - Get customer by ID
- `PUT /api/customers/:id` - Update customer
- `DELETE /api/customers/:id` - Delete customer

### Invoices
- `POST /api/invoices` - Create invoice
- `GET /api/invoices` - Get all invoices
- `GET /api/invoices/latest` - Get latest 5 invoices with customer details
- `GET /api/invoices/:id` - Get invoice by ID
- `PUT /api/invoices/:id` - Update invoice
- `DELETE /api/invoices/:id` - Delete invoice

### Revenue
- `POST /api/revenue` - Create revenue record
- `GET /api/revenue` - Get all revenue
- `GET /api/revenue/:month` - Get revenue by month
- `PUT /api/revenue/:month` - Update revenue
- `DELETE /api/revenue/:month` - Delete revenue

---

## Example API Calls

### Create User
```bash
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Wick",
    "middleName": "D",
    "email": "john.wick@example.com",
    "password": "password123"
  }'
```

### Create Customer
```bash
curl -X POST http://localhost:3000/api/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Simple Corporation",
    "email": "contact.simple.corporation@gmail.com",
    "imageUrl": "https://ui-avatars.com/api/?name=AC&background=8763C5&color=fff&size=150"
  }'
```

### Create Invoice
```bash
curl -X POST http://localhost:3000/api/invoices \
  -H "Content-Type: application/json" \
  -d '{
    "customerId": "customer-id-string",
    "amount": 15000,
    "date": "2024-11-06",
    "status": "pending"
  }'
```

### Get Latest Invoices
```bash
curl http://localhost:3000/api/invoices/latest
```

---

## Architecture Overview

### Vertical Slice Architecture
Each feature (users, customers, invoices, revenue) is self-contained with:
- **command.go** - Write operations (Create, Update, Delete)
- **query.go** - Read operations (Get(GetItem), GetAll(GetItems))
- **controller.go** - HTTP request handlers
- **route.go** - Route registration

### CQRS Pattern
- **Command**: Modify state, return success/error
- **Query**: Read state, return data
- Clear separation of concerns

### No Repository Pattern
Handlers directly interact with MongoDB collections for simplicity.

---

## Key Features

✅ **Complete Vertical Slice Architecture** - Each feature is fully self-contained  
✅ **CQRS Pattern** - Clear separation of Commands and Queries  
✅ **Full CRUD Operations** - All entities have Create, Read, Update, Delete  
✅ **MongoDB Aggregation** - Latest invoices query joins with customers  
✅ **Production Ready** - Docker, Makefile, proper error handling  
✅ **Clean Code** - Well-organized, maintainable structure

All files are now properly separated into their respective packages following Go best practices and Vertical Slice Architecture!
