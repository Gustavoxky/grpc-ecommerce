package server

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
    "golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func InitDB() {
    _ = godotenv.Load(".env") // Tenta carregar variáveis do .env

    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        user, password, host, port, dbname,
    )

    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Erro ao conectar no PostgreSQL: %v", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("Erro ao pingar DB: %v", err)
    }

    log.Println("Conectado ao PostgreSQL com sucesso!")
    createTables()
}

func createTables() {
    orders := `
        CREATE TABLE IF NOT EXISTS orders (
            id TEXT PRIMARY KEY,
            customer_name TEXT,
            status TEXT
        );`

    users := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL
        );`

    _, err := db.Exec(orders)
    if err != nil {
        log.Fatalf("Erro ao criar tabela orders: %v", err)
    }

    _, err = db.Exec(users)
    if err != nil {
        log.Fatalf("Erro ao criar tabela users: %v", err)
    }

    fmt.Println("Tabelas criadas ou já existentes.")
}

func InsertOrder(id, customerName, status string) error {
    _, err := db.Exec("INSERT INTO orders (id, customer_name, status) VALUES ($1, $2, $3)", id, customerName, status)
    return err
}

func GetOrderById(id string) (string, string, string, error) {
    var orderID, customerName, status string
    err := db.QueryRow("SELECT id, customer_name, status FROM orders WHERE id = $1", id).
        Scan(&orderID, &customerName, &status)
    if err != nil {
        return "", "", "", err
    }
    return orderID, customerName, status, nil
}

func UpdateOrderStatus(id, status string) error {
    _, err := db.Exec("UPDATE orders SET status = $1 WHERE id = $2", status, id)
    return err
}

func DeleteOrder(id string) error {
    _, err := db.Exec("DELETE FROM orders WHERE id = $1", id)
    return err
}

func ListOrders() ([]map[string]string, error) {
    rows, err := db.Query("SELECT id, customer_name, status FROM orders")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []map[string]string
    for rows.Next() {
        var id, customerName, status string
        if err := rows.Scan(&id, &customerName, &status); err != nil {
            return nil, err
        }
        orders = append(orders, map[string]string{
            "id":            id,
            "customer_name": customerName,
            "status":        status,
        })
    }
    return orders, nil
}

func RegisterUser(username, password string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    _, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, string(hash))
    return err
}

func ValidateUser(username, password string) bool {
    var hash string
    err := db.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).Scan(&hash)
    if err != nil {
        return false
    }

    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func CloseDB() {
    db.Close()
}
