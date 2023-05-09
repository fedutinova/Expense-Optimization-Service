package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	wallet "wallet-app/pkg/models"
	errs "wallet-app/pkg/models/errors"
)

type WalletPostgres struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

func NewWalletPostgres(db *sqlx.DB, logger *logrus.Logger) *WalletPostgres {
	return &WalletPostgres{
		db:     db,
		logger: logger,
	}
}

func (r *WalletPostgres) Create(userId int, wallet wallet.Wallet) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Errorf("failed run db: %v", err)
		return 0, errs.InternalServer
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			r.logger.Errorf("failed create wallet: %v", err)
			fmt.Println(err)
		} else {
			err = tx.Commit()
			fmt.Println(err)
		}
	}()

	var id int
	createWalletQuery := fmt.Sprintf("INSERT INTO %s (isFamily, balance, description) VALUES ($1, $2, $3) RETURNING id", walletsTable)
	row := tx.QueryRow(createWalletQuery, wallet.IsFamily, wallet.Balance, wallet.Description)
	if err = row.Scan(&id); err != nil {
		r.logger.Errorf("failed insert wallet table: %v", err)
		return 0, errs.InternalServer
	}

	createUsersWalletQuery := fmt.Sprintf("INSERT INTO %s (user_id, wallet_id, is_holder) VALUES ($1, $2, $3)", usersWalletsTable)
	if _, err = tx.Exec(createUsersWalletQuery, userId, id, true); err != nil {
		r.logger.Errorf("failed insert users_wallet table: %v", err)
		return 0, errs.InternalServer
	}

	return id, nil
}

func (r *WalletPostgres) GetAll(userId int) ([]wallet.Wallet, error) {
	var wallets []wallet.Wallet

	query := fmt.Sprintf("SELECT w.id, w.isFamily, w.balance, w.description FROM %s w INNER JOIN %s uw on w.id = uw.wallet_id WHERE uw.user_id = $1",
		walletsTable, usersWalletsTable)
	if err := r.db.Select(&wallets, query, userId); err != nil {
		r.logger.Errorf("failed select wallets: %v", err)
		return nil, errs.InternalServer
	}

	return wallets, nil
}

func (r *WalletPostgres) GetById(userId, walletId int) (wallet.Wallet, error) {
	var wallet wallet.Wallet

	query := fmt.Sprintf(`SELECT tl.id, tl.isFamily, tl.balance, tl.description FROM %s tl
								INNER JOIN %s ul on tl.id = ul.wallet_id WHERE ul.user_id = $1 AND ul.wallet_id = $2`,
		walletsTable, usersWalletsTable)
	if err := r.db.Get(&wallet, query, userId, walletId); err != nil {
		r.logger.Errorf("failed select wallets: %v", err)
		return wallet, errs.InternalServer
	}

	return wallet, nil
}

func (r *WalletPostgres) Delete(userId, walletId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.wallet_id AND ul.user_id=$1 AND ul.wallet_id=$2",
		walletsTable, usersWalletsTable)
	if _, err := r.db.Exec(query, userId, walletId); err != nil {
		r.logger.Errorf("failed delete wallet: %v", err)
		return errs.InternalServer
	}

	return nil
}

func (r *WalletPostgres) AddMember(userId int, newUser wallet.Member) error {
	AddMemberQuery := fmt.Sprintf("INSERT INTO %s (user_id, wallet_id, is_holder) VALUES ($1, $2, $3)", usersWalletsTable)
	if _, err := r.db.Exec(AddMemberQuery, newUser.Id, newUser.WalletId, false); err != nil {
		r.logger.Errorf("failed to add user to wallet: %v", err)
		return errs.InternalServer
	}

	return nil

}

func (r *WalletPostgres) DeleteMember(userId int, userToDelete wallet.Member) error {
	query := fmt.Sprintf("DELETE FROM %s uw WHERE uw.user_id=$1 AND uw.wallet_id=$2",
		usersWalletsTable)
	if _, err := r.db.Exec(query, userToDelete.Id, userToDelete.WalletId); err != nil {
		r.logger.Errorf("failed delete wallet: %v", err)
		return errs.InternalServer
	}

	return nil

}

func (r *WalletPostgres) GetOwnerIdQuery(userId int, newMember wallet.Member) error {
	var ownerId int
	if err := r.db.QueryRow("SELECT uw.user_id FROM users_wallets uw INNER JOIN wallets w on uw.wallet_id = w.id WHERE uw.wallet_id = $1 AND uw.is_holder = true",
		newMember.WalletId).Scan(&ownerId); err != nil {
		r.logger.Errorf("failed check wallet owner: %v", err)
		return errs.InternalServer
	}
	if ownerId != userId {
		r.logger.Error("user is not owner")
		return fmt.Errorf("failed: not enough permissions")
	}

	return nil
}
