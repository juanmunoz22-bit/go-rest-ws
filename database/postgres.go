package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/juanmunoz22-bit/go-rest-ws/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO users (email, password, id) VALUES ($1, $2, $3)",
		user.Email,
		user.Password,
		user.Id,
	)
	return err
}

func (repo *PostgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, email FROM users WHERE id = $1",
		id,
	)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal()
		}
	}()
	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row, err := repo.db.QueryContext(
		ctx,
		"SELECT id, email, password FROM users WHERE email = $1",
		email,
	)
	defer func() {
		err = row.Close()
		if err != nil {
			log.Fatal()
		}
	}()
	var user = models.User{}
	for row.Next() {
		if err = row.Scan(
			&user.Id,
			&user.Email,
			&user.Password,
		); err == nil {
			return &user, nil
		}
	}
	if err = row.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)",
		post.Id,
		post.PostContent,
		post.UserId,
	)
	return err
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
