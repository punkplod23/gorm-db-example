# GORM DB Example

This project demonstrates various ways to interact with a database using GORM in Go. The main program can be executed with different flags to run specific examples.

## Performance Metrics

| Example            | Execution Time   |
|--------------------|------------------|
| Eager Loading      | 107.3493ms       |
| JSON Aggregate     | 66.8975ms        |
| Join               | 51.7037ms        |
| Lazy Loading       | 30.1419633s      |

## Flags

- `-eager`: Run the eager loading example.
- `-join`: Run the join example.
- `-lazy`: Run the lazy loading example.
- `-json`: Run the JSON aggregate example.
- `-all`: Run all examples and print metrics.

## How to Execute

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/gorm-db-example.git
    cd gorm-db-example
    ```

2. Ensure you have Go installed and set up your database connection in the `main.go` file.

3. Run the program with the desired flag:
    ```sh
    go run . -eager
    go run . -join
    go run . -lazy
    go run . -json
    go run . -all
    ```

4. Alternatively, you can use Docker Compose to set up the environment and run the full experiment:
    ```sh
    docker-compose up --build
    ```

## Examples

### Eager Loading
Run the eager loading example:
```sh
go run . -eager
```

### Join
Run the join example:
```sh
go run . -join
```

### Lazy Loading
Run the lazy loading example:
```sh
go run . -lazy
```

### JSON Aggregate
Run the JSON aggregate example:
```sh
go run . -json
```

### All Examples
Run all examples and print metrics:
```sh
go run . -all
```

## Load Database SQL
Use the schema in migration db-schema.sql to load the database

## Using a local database
Change the dsn in main.go

## License

This project is licensed under the MIT License.
