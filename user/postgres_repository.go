package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var errUserNotFound = errors.New("user not found")

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepositoryFromEnv() (*PostgresUserRepository, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		if host == "" || port == "" || user == "" || password == "" || dbname == "" {
			return nil, fmt.Errorf("missing DB config: set DATABASE_URL or DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME")
		}

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	}

	return NewPostgresUserRepository(connStr)
}

func NewPostgresUserRepository(connStr string) (*PostgresUserRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := connectWithRetry(ctx, "pgx", connStr)
	if err != nil {
		return nil, err
	}

	repo := &PostgresUserRepository{db: db}
	if err := repo.initSchema(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return repo, nil
}

func connectWithRetry(ctx context.Context, driver, connStr string) (*sql.DB, error) {
	var lastErr error
	for {
		if ctx.Err() != nil {
			if lastErr != nil {
				return nil, fmt.Errorf("db connect timeout: %w", lastErr)
			}
			return nil, ctx.Err()
		}

		db, err := sql.Open(driver, connStr)
		if err != nil {
			lastErr = err
			select {
			case <-time.After(1 * time.Second):
				continue
			case <-ctx.Done():
				continue
			}
		}

		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		err = db.PingContext(pingCtx)
		cancel()
		if err == nil {
			return db, nil
		}

		lastErr = err
		_ = db.Close()

		select {
		case <-time.After(1 * time.Second):
			continue
		case <-ctx.Done():
			continue
		}
	}
}

func (r *PostgresUserRepository) initSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		);
	`)
	return err
}

func (r *PostgresUserRepository) GetAll() ([]User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PostgresUserRepository) GetByID(id int) (*User, error) {
	var u User
	err := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) Create(u User) (*User, error) {
	err := r.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) Update(id int, u User) (*User, error) {
	res, err := r.db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
	if err != nil {
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errUserNotFound
	}

	u.ID = id
	return &u, nil
}

func (r *PostgresUserRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errUserNotFound
	}
	return nil
}
