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
│   │   ├── routes.go                    # Routes config
│   │   ├── server.go                    # Server config
│   ├── database/
│   │   └── database.go                  # MongoDB connection
│   └── features/
│       ├── users/
│       │   ├── commands/
│       │   │   └── command.go
│       │   ├── queries/
│       │   │   └── queries.go
│       │   ├── handler.go
│       │   │   └── handler.go
│       │   ├── models/
│       │   │   └── models.go
│       │   └── routes/
│       │       └── routes.go
│       ├── customers/
│       │   ├── commands/
│       │   │   └── command.go
│       │   ├── queries/
│       │   │   └── queries.go
│       │   ├── handler.go
│       │   │   └── handler.go
│       │   ├── models/
│       │   │   └── models.go
│       │   └── routes/
│       │       └── routes.go
│       ├── invoices/
│       │   ├── commands/
│       │   │   └── command.go
│       │   ├── queries/
│       │   │   └── queries.go
│       │   ├── handler.go
│       │   │   └── handler.go
│       │   ├── models/
│       │   │   └── models.go
│       │   └── routes/
│       │       └── routes.go
│       └── revenue/
│           ├── commands/
│           │   └── command.go
│           ├── queries/
│           │   └── queries.go
│           ├── handler.go
│           │   └── handler.go
│           ├── models/
│           │   └── models.go
│           └── routes/
│               └── routes.go
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

## Key Features

✅ **Complete Vertical Slice Architecture** - Each feature is fully self-contained  
✅ **CQRS Pattern** - Clear separation of Commands and Queries  
✅ **Full CRUD Operations** - All entities have Create, Read, Update, Delete  
✅ **MongoDB Aggregation** - Latest invoices query joins with customers  
✅ **Production Ready** - Docker, Makefile, proper error handling  
✅ **Clean Code** - Well-organized, maintainable structure

All files are now properly separated into their respective packages following Go best practices and Vertical Slice Architecture!
