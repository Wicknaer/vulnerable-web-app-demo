package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "./demo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
				"id" integer PRIMARY KEY AUTOINCREMENT,
				"username" TEXT,
				"password" TEXT );`
	db.Exec(createTableSQL)

	db.Exec("INSERT INTO users (username, password) VALUES ('admin', 'gizlisifre123')")

	http.HandleFunc("/", loginPage)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/users", listUsersPage) // Veritabanındaki herkesi listeleyen sayga. Bkz: 60. satır

	fmt.Println("Sunucu çalışıyor: https://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func loginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	query := fmt.Sprintf("SELECT id FROM users WHERE username='%s' AND password='%s'", username, password)
	fmt.Println("SQL Sorgusu:", query)

	var id int
	err := db.QueryRow(query).Scan(&id)

	if err != nil {
		fmt.Fprintf(w, "<h1>Giriş Başarısız! %s</h1>", username)
		return
	}
	fmt.Fprintf(w, "<h1>Giriş Başarılı! Hoşgeldin Admin (ID: %d)</h1>", id)
}

// Veritabanındaki tüm kullanıcıları listeleyerek veritabanı kurulumunun başarılı olduğunu kontrol etmek için...
func listUsersPage(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username, password FROM users")
	if err != nil {
		http.Error(w, "Sorgu hatası: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	fmt.Fprintln(w, "<h1>Veritabanı Kayıtları</h1><ul>")
	for rows.Next() {
		var id int
		var user, pass string
		rows.Scan(&id, &user, &pass)
		fmt.Fprintf(w, "<li>ID: %d | User: %s | Pass: %s</li>", id, user, pass)
	}
	fmt.Fprintln(w, "</ul>")
}
