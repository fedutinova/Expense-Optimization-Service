package service

import (
	"github.com/sirupsen/logrus"
	wallet "wallet-app/pkg/models"
	"wallet-app/pkg/repository"
)

type TransactionService struct {
	repo       repository.Transaction
	walletRepo repository.Wallets
	logger     *logrus.Logger
}

func NewTransactionService(repo repository.Transaction, walletRepo repository.Wallets, logger *logrus.Logger) *TransactionService {
	return &TransactionService{
		repo:       repo,
		walletRepo: walletRepo,
		logger:     logger,
	}
}

func (s *TransactionService) Create(userId, walletId int, transaction wallet.Transaction) (int, error) {
	if err := transaction.Validate(); err != nil {
		return 0, err
	}
	_, err := s.walletRepo.GetById(userId, walletId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(userId, walletId, transaction)
}

func (s *TransactionService) Delete(userId, transactionId int) error {
	return s.repo.Delete(userId, transactionId)
}

func (s *TransactionService) GetAll(walletId int, date wallet.TransactionDate) ([]wallet.Transaction, error) {
	return s.repo.GetAll(walletId, date)
}

func (s *TransactionService) GetAllIncomes(walletId int, date wallet.TransactionDate) ([]wallet.Transaction, error) {
	return s.repo.GetAllIncomes(walletId, date)
}

func (s *TransactionService) GetAllExpenses(walletId int, date wallet.TransactionDate) ([]wallet.Transaction, error) {
	return s.repo.GetAllExpenses(walletId, date)
}

func (s *TransactionService) GetByCategoryIncome(walletId int, date wallet.TransactionDate) ([]wallet.TransactionsByCategory, error) {
	return s.repo.GetByCategoryIncome(walletId, date)
}

func (s *TransactionService) GetByCategoryExpenses(walletId int, date wallet.TransactionDate) ([]wallet.TransactionsByCategory, error) {
	return s.repo.GetByCategoryExpenses(walletId, date)
}

func (s *TransactionService) Update(userId, transactionId int, input wallet.UpdateTransactionInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, transactionId, input)
}
