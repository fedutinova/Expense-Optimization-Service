package service

import (
	"github.com/sirupsen/logrus"
	wallet "wallet-app/pkg/models"
	"wallet-app/pkg/repository"
)

type WalletsService struct {
	repo   repository.Wallets
	logger *logrus.Logger
}

func NewWalletsService(repo repository.Wallets, logger *logrus.Logger) *WalletsService {
	return &WalletsService{
		repo:   repo,
		logger: logger,
	}
}

func (s *WalletsService) Create(userId int, wallet wallet.Wallet) (int, error) {
	return s.repo.Create(userId, wallet)
}

func (s *WalletsService) GetAll(userId int) ([]wallet.Wallet, error) {
	return s.repo.GetAll(userId)
}

func (s *WalletsService) GetById(userId, walletId int) (wallet.Wallet, error) {
	return s.repo.GetById(userId, walletId)
}

func (s *WalletsService) Delete(userId, walletId int) error {
	return s.repo.Delete(userId, walletId)
}

func (s *WalletsService) AddMember(userId int, newUser wallet.Member) error {
	return s.repo.AddMember(userId, newUser)
}

func (s *WalletsService) DeleteMember(userId int, userToDelete wallet.Member) error {
	return s.repo.DeleteMember(userId, userToDelete)
}

func (s *WalletsService) GetOwnerIdQuery(userId int, newMember wallet.Member) error {
	return s.repo.GetOwnerIdQuery(userId, newMember)
}
