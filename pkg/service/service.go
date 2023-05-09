package service

import (
	"github.com/sirupsen/logrus"
	"wallet-app/pkg/models"
	"wallet-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Wallet interface {
	Create(userId int, wallet models.Wallet) (int, error)
	GetAll(userId int) ([]models.Wallet, error)
	GetById(userId, walletId int) (models.Wallet, error)
	Delete(userId, walletId int) error
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

type Service struct {
	Authorization Authorization
	Wallet        Wallet
	Transaction   Transaction
}

func NewService(repos *repository.Repository, logger *logrus.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, logger),
		Wallet:        NewWalletsService(repos.Wallets, logger),
		Transaction:   NewTransactionService(repos.Transaction, repos.Wallets, logger),
	}
}
