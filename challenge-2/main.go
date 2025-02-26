package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed static/* templates/*
var content embed.FS

// Flag to be extracted via SQL injection
var flag string

func main() {
	flag = os.Getenv("FLAG")
	if flag == "" {
		log.Fatal("Flag not set")
	}
	// Set up the database
	err := setupDatabase()
	if err != nil {
		log.Fatal("Failed to set up database:", err)
	}

	// Parse templates
	tmpl, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}

	// Set up routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})

	// Serve static files
	http.Handle("/static/", http.FileServer(http.FS(content)))

	// Set up the vulnerable login route
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Open database connection
		db, err := sql.Open("sqlite3", "./challenge.db")
		if err != nil {
			log.Println("Database connection error:", err)
			tmpl.ExecuteTemplate(w, "index.html", "Database error")
			return
		}
		defer db.Close()

		// Vulnerable SQL query (intentionally vulnerable to SQL injection)
		query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s' AND password = '%s'", username, password)

		// Log the query for debugging
		log.Println("Executing query:", query)

		rows, err := db.Query(query)
		if err != nil {
			log.Println("Query error:", err)
			tmpl.ExecuteTemplate(w, "index.html", fmt.Sprintf("Query error: %v", err))
			return
		}
		defer rows.Close()

		if rows.Next() {
			// User logged in successfully
			var id int
			var username, password string
			err := rows.Scan(&id, &username, &password)
			if err != nil {
				log.Println("Scan error:", err)
				tmpl.ExecuteTemplate(w, "index.html", fmt.Sprintf("Scan error: %v", err))
				return
			}

			tmpl.ExecuteTemplate(w, "success.html", username)
		} else {
			// Login failed
			tmpl.ExecuteTemplate(w, "index.html", "Invalid username or password")
		}
	})

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Starting server on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setupDatabase() error {
	// Create a new database if it doesn't exist
	db, err := sql.Open("sqlite3", "./challenge.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			username TEXT,
			password TEXT
		)
	`)
	if err != nil {
		return err
	}

	// Create flags table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS flags (
			id INTEGER PRIMARY KEY,
			flag TEXT
		)
	`)
	if err != nil {
		return err
	}

	// Check if we need to insert data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil || count == 0 {
		// Insert demo users
		_, err = db.Exec("INSERT INTO users (username, password) VALUES ('admin', 'supersecretpassword')")
		if err != nil {
			return err
		}
		_, err = db.Exec("INSERT INTO users (username, password) VALUES ('user', 'password123')")
		if err != nil {
			return err
		}

		// Insert the flag
		_, err = db.Exec("INSERT INTO flags (flag) VALUES (?)", flag)
		if err != nil {
			return err
		}
	}

	return nil
}
