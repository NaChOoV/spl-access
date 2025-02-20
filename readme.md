# SPL Access API

This is a simple API to manage access control using Go.

## Prerequisites

- Go 1.19 or higher
- PostgreSQL

## Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/spl-access.git
cd spl-access
```

2. Install dependencies

```bash
go mod download
```

3. Environment Setup

- Copy the example environment file
```bash
cp .env.example .env
```
- Update the `.env` file with your configuration

## Development

For hot reloading during development, it's recommended to use [Air](https://github.com/cosmtrek/air)

```bash
# Install Air globally
go install github.com/cosmtrek/air@latest

# Run the application with Air
air
```

## Environment Variables

Make sure to set up the following environment variables in your `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=spl_access
DB_SSLMODE=disable

SERVER_PORT=8080
```

## Running the Application

```bash
go run main.go
```


#### **Pending**
- [ ] **Features**
  - [ ] Create a CronJob to remove garbage data from the database.
  - [ ] Calculate the mean time of each person base on historic access data.
  - [ ] Calculate the total people on a place base in time interval.
    - [ ] Create an endpoint to show the data, filtered by date range.

- [ ] **Testing**
  - [ ] Create unit test.
  - [ ] Configure CI.
