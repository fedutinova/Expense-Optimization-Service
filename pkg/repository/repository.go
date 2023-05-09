package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"wallet-app/pkg/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Wallets interface {
	Create(userId int, wallet models.Wallet) (int, error)
	GetAll(userId int) ([]models.Wallet, error)
	GetById(userId, listId int) (models.Wallet, error)
	Delete(userId, listId int) error
	AddMember(userId int, newUser models.Member) error
	DeleteMember(userId int, userToDelete models.Member) error
	GetOwnerIdQuery(userId int, newMember models.Member) error
}

type Transaction interface {
	Create(userId, walletId int, transaction models.Transaction) (int, error)
	Delete(userId, transactionId int) error
	GetAll(walletId int, date models.TransactionDate) ([]models.Transaction, error)
	GetAllIncomes(walletId int, date models.TransactionDate) ([]models.Transaction, error)
	GetAllExpenses(walletId int, date models.TransactionDate) ([]models.Transaction, error)
	GetByCategoryIncome(walletId int, date models.TransactionDate) ([]models.TransactionsByCategory, error)
	GetByCategoryExpenses(walletId int, date models.TransactionDate) ([]models.TransactionsByCategory, error)
	Update(userId, transactionId int, input models.UpdateTransactionInput) error
}

type Repository struct {
	Authorization Authorization
	Wallets       Wallets
	Transaction   Transaction
}

func NewRepository(db *sqlx.DB, logger *logrus.Logger) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db, logger),
		Wallets:       NewWalletPostgres(db, logger),
		Transaction:   NewTransactionPostgres(db, logger),
	}
}
