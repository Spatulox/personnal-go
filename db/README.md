# Database
This package provides a simple and efficient way to interact with MySQL databases in Go. It includes functions for connecting to a database, selecting data, inserting data, updating data, and deleting data.

## Features

- Easy database connection setup
- Flexible SELECT queries with support for JOINs and conditions
- Simple INSERT operations
- UPDATE operations with condition support
- DELETE operations with condition support
- Debug mode for query visualization

## Usage

### The Db struct
```go
type Db struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}
```

### Creating a db connection
```go
db := NewDb("localhost", 3306, "username", "password", "database_name")
```
Creates a new Db instance with the provided connection details

### Selecting data
```go
columns := []string{"id", "name", "email"}
condition := "age > 18"
results, err := db.SelectDB("users", columns, nil, &condition)
if err != nil {
    // Handle error
}
// Process results
```
Executes a SELECT query on the specified table. Supports JOINs and WHERE conditions.

### Inserting Data
```go
columns := []string{"name", "email", "age"}
values := []string{"John Doe", "john@example.com", "30"}
db.InsertDB("users", columns, values)
```
Inserts data into the specified table.

### Updating data
```go
columns := []string{"name", "email"}
values := []string{"Jane Doe", "jane@example.com"}
condition := "id = 1"
db.UpdateDB("users", columns, values, &condition)
```
Updates data in the specified table. Requires a condition for the update operation.

### Deleting data
```go
condition := "id = 1"
db.DeleteDB("users", &condition)
```
Deletes data from the specified table based on the given condition.

### Debug mode
To see the generated SQL query, add true as the last parameter:
```go
db.SelectDB("users", columns, nil, &condition, true)
```
Works with every query type


## Error Handling
All functions include error checking and logging. Make sure to handle returned errors appropriately in your application.