package auth

import (
	"database/sql"
	"fmt"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	db *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Store) RegisterUser(user User) error {
	// Проверка, существует ли уже пользователь с таким логином
	var existingUser User
	err := s.db.QueryRow("SELECT id, username FROM users WHERE username = ?", user.Username).Scan(&existingUser.ID, &existingUser.Username)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing user: %v", err)
	}
	if existingUser.Username != "" {
		return errors.New("user already exists")
	}

	// Создание нового пользователя
	_, err = s.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}

	return nil
}

func NewStore(path string) (*Store, error) {
	// Проверим, существует ли база
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("База данных не найдена по пути '%s'. Будет создан новый файл.", path)
	} else {
		log.Printf("База данных найдена по пути '%s'.", path)
	}

	// Открываем подключение
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	// Проверим подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Создание таблицы users
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
	}
	log.Println("Таблица users создана или уже существует.")

	// Для отладки: вывод структуры базы
	rows, err := db.Query("PRAGMA table_info(users);")
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения структуры таблицы users: %w", err)
	}
	defer rows.Close()

	log.Println("Структура таблицы users:")
	for rows.Next() {
		var cid int
		var name string
		var ctype string
		var notnull int
		var dfltValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			return nil, fmt.Errorf("ошибка сканирования структуры таблицы: %w", err)
		}
		log.Printf("  - %s %s", name, ctype)
	}

	return &Store{db: db}, nil
}

func (s *Store) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *Store) ValidateUser(username, password string) error {
	var hashedPassword string
	err := s.db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}
