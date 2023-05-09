package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	walletApp "wallet-app/pkg/models"
	errs "wallet-app/pkg/models/errors"
)

type AuthPostgres struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

func NewAuthPostgres(db *sqlx.DB, logger *logrus.Logger) *AuthPostgres {
	return &AuthPostgres{
		db:     db,
		logger: logger,
	}
}

func (r *AuthPostgres) CreateUser(user walletApp.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		r.logger.Errorf("failed create user id: %v", err)
		return 0, errs.InternalServer
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (walletApp.User, error) {
	var user walletApp.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	if err := r.db.Get(&user, query, username, password); err != nil {
		r.logger.Errorf("failed get user: %v", err)
		return user, fmt.Errorf("user does not exist")
	}

	return user, nil
}
