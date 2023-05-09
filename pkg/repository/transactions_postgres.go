package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"strings"
	wallet "wallet-app/pkg/models"
	errs "wallet-app/pkg/models/errors"
)

type TransactionPostgres struct {
	logger *logrus.Logger
	db     *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB, logger *logrus.Logger) *TransactionPostgres {
	return &TransactionPostgres{
		db:     db,
		logger: logger,
	}
}

func (r *TransactionPostgres) Create(userId, walletId int, transaction wallet.Transaction) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Errorf("failed to run db: %v", err)
		return 0, errs.InternalServer
	}
	defer tx.Rollback()

	var transactionId int

	if transaction.Amount == decimal.Zero {
		return 0, fmt.Errorf("fail: amount must not be zero")
	}

	createTransactionsQuery := fmt.Sprintf("INSERT INTO %s (user_id, date, category, amount, description) values ($1, $2, $3, $4, $5) RETURNING id",
		transactionsTable)
	if err = tx.QueryRow(createTransactionsQuery, userId, transaction.Date, transaction.Category, transaction.Amount,
		transaction.Description).Scan(&transactionId); err != nil {
		r.logger.Errorf("failed query transaction table: %v", err)
		return 0, errs.InternalServer
	}

	createWalletTransactionsQuery := fmt.Sprintf("INSERT INTO %s (wallet_id, transaction_id) VALUES($1, $2)",
		walletsTransactionsTable)
	if _, err = tx.Exec(createWalletTransactionsQuery, walletId, transactionId); err != nil {
		r.logger.Errorf("failed query wallets_transactions table: %v", err)
		return 0, errs.InternalServer
	}

	var currentBalance, newBalance decimal.Decimal

	if err = r.db.QueryRow("SELECT w.balance FROM wallets w WHERE w.id = $1",
		walletId).Scan(&currentBalance); err != nil {
		r.logger.Errorf("failed scan current balance: %v", err)
		return 0, errs.InternalServer
	}

	newBalance = currentBalance.Add(transaction.Amount)
	if err != nil {
		r.logger.Errorf("failed calculate new balance: %v", err)
		return 0, errs.InternalServer
	}

	if _, err = r.db.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", newBalance.String(), walletId); err != nil {
		r.logger.Errorf("failed update balance: %v", err)
		return 0, errs.InternalServer
	}

	return transactionId, tx.Commit()
}

func (r *TransactionPostgres) GetAll(walletId int, date wallet.TransactionDate) ([]wallet.Transaction, error) {
	transactions := make([]wallet.Transaction, 0, 100)
	query := fmt.Sprintf(`SELECT t.id, t.user_id, t.date, t.category, t.amount, t.description
								FROM %s t 
								INNER JOIN %s wt on wt.transaction_id = t.id 
								WHERE wt.wallet_id = $1 
									AND EXTRACT(MONTH FROM CAST(t.date AS timestamp)) = $2
									AND EXTRACT(YEAR FROM CAST(t.date AS timestamp)) = $3`,
		transactionsTable, walletsTransactionsTable)
	if err := r.db.Select(&transactions, query, walletId, date.Month, date.Year); err != nil {
		r.logger.Errorf("failed query transactions: %v", err)
		return nil, errs.InternalServer
	}

	return transactions, nil
}

func (r *TransactionPostgres) GetAllIncomes(walletId int, date wallet.TransactionDate) ([]wallet.Transaction, error) {
	transactions := make([]wallet.Transaction, 0, 100)
	query := fmt.Sprintf(`SELECT t.id, t.user_id, t.date, t.category, t.amount, t.description
								FROM %s t 
								INNER JOIN %s wt on wt.transaction_id = t.id 
								WHERE wt.wallet_id = $1 
									AND t.amount > 0
									AND EXTRACT(MONTH FROM CAST(t.date AS timestamp)) = $2
									AND EXTRACT(YEAR FROM CAST(t.date AS timestamp)) = $3`,
		transactionsTable, walletsTransactionsTable)
	if err := r.db.Select(&transactions, query, walletId, date.Month, date.Year); err != nil {
		r.logger.Errorf("failed query incomes: %v", err)
		return nil, errs.InternalServer
	}

	return transactions, nil
}

func (r *TransactionPostgres) GetAllExpenses(walletId int, date wallet.TransactionDate) ([]wallet.Transaction, error) {
	transactions := make([]wallet.Transaction, 0, 100)
	query := fmt.Sprintf(`SELECT t.id, t.user_id, t.date, t.category, t.amount, t.description
								FROM %s t 
                                INNER JOIN %s wt on wt.transaction_id = t.id 
                                WHERE wt.wallet_id = $1 
                                    AND t.amount < 0
                                    AND EXTRACT(MONTH FROM CAST(t.date AS timestamp)) = $2
									AND EXTRACT(YEAR FROM CAST(t.date AS timestamp)) = $3`,
		transactionsTable, walletsTransactionsTable)
	if err := r.db.Select(&transactions, query, walletId, date.Month, date.Year); err != nil {
		r.logger.Errorf("failed query expenses: %v", err)
		return nil, errs.InternalServer
	}

	return transactions, nil
}

func (r *TransactionPostgres) GetByCategoryIncome(walletId int, date wallet.TransactionDate) ([]wallet.TransactionsByCategory, error) {
	transactions := make([]wallet.TransactionsByCategory, 0, 100)
	query := fmt.Sprintf(`SELECT t.category, SUM(t.amount) AS sum 
								FROM %s t 
								INNER JOIN %s wt ON wt.transaction_id = t.id 
								WHERE wt.wallet_id = $1 
								    AND t.amount > 0
									AND EXTRACT(MONTH FROM CAST(t.date AS timestamp)) = $2
									AND EXTRACT(YEAR FROM CAST(t.date AS timestamp)) = $3
								GROUP BY t.category
								ORDER BY SUM(t.amount) DESC`,
		transactionsTable, walletsTransactionsTable)
	if err := r.db.Select(&transactions, query, walletId, date.Month, date.Year); err != nil {
		r.logger.Errorf("failed query transactions: %v", err)
		return nil, errs.InternalServer
	}

	return transactions, nil
}

func (r *TransactionPostgres) GetByCategoryExpenses(walletId int, date wallet.TransactionDate) ([]wallet.TransactionsByCategory, error) {
	transactions := make([]wallet.TransactionsByCategory, 0, 100)
	query := fmt.Sprintf(`SELECT t.category, SUM(t.amount) AS sum 
								FROM %s t 
								INNER JOIN %s wt ON wt.transaction_id = t.id 
								WHERE wt.wallet_id = $1 
								    AND t.amount < 0
									AND EXTRACT(MONTH FROM CAST(t.date AS timestamp)) = $2
									AND EXTRACT(YEAR FROM CAST(t.date AS timestamp)) = $3
								GROUP BY t.category
								ORDER BY SUM(t.amount) DESC`,
		transactionsTable, walletsTransactionsTable)
	if err := r.db.Select(&transactions, query, walletId, date.Month, date.Year); err != nil {
		r.logger.Errorf("failed query transactions: %v", err)
		return nil, errs.InternalServer
	}

	return transactions, nil
}

func (r *TransactionPostgres) Update(userId, transactionId int, input wallet.UpdateTransactionInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	var (
		oldBalance, oldAmount, currentBalance, newBalance decimal.Decimal
	)

	if input.Date != nil {
		setValues = append(setValues, fmt.Sprintf("date=$%d", argId))
		args = append(args, *input.Date)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Amount != (decimal.Decimal{}) {
		setValues = append(setValues, fmt.Sprintf("amount=$%d", argId))
		args = append(args, input.Amount)
		argId++

		if err := r.db.QueryRow(`
    SELECT t.amount, w.balance
    FROM transactions t
    INNER JOIN wallets_transactions wt ON wt.transaction_id = t.id
    INNER JOIN wallets w ON wt.wallet_id = w.id
    WHERE t.id = $1 AND t.user_id = $2
`, transactionId, userId).Scan(&oldAmount, &oldBalance); err != nil {
			r.logger.Errorf("failed scan amount transaction: %v", err)
			return errs.InternalServer
		}
		currentBalance = oldBalance.Sub(oldAmount)
	}

	if input.Category != nil {
		setValues = append(setValues, fmt.Sprintf("category=$%d", argId))
		args = append(args, *input.Category)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s t SET %s FROM %s wt, %s uw
									WHERE t.id = wt.transaction_id AND wt.wallet_id = uw.wallet_id AND t.user_id = $%d AND t.id = $%d`,
		transactionsTable, setQuery, walletsTransactionsTable, usersWalletsTable, argId, argId+1)
	args = append(args, userId, transactionId)

	if _, err := r.db.Exec(query, args...); err != nil {
		r.logger.Errorf("failed to execute query: %v", err)
		return errs.InternalServer
	}

	if input.Amount != (decimal.Decimal{}) {
		newBalance = currentBalance.Add(input.Amount)

		if _, err := r.db.Exec("UPDATE wallets w SET balance = $1 FROM wallets_transactions wt WHERE w.id = wt.wallet_id AND wt.transaction_id = $2", newBalance.String(), transactionId); err != nil {
			r.logger.Errorf("failed update balance: %v", err)
			return errs.InternalServer
		}

	}

	return nil
}

func (r *TransactionPostgres) Delete(userId, transactionId int) error {
	query := fmt.Sprintf(`DELETE FROM %s t WHERE t.id = $1 AND t.user_id = $2`, transactionsTable)
	if _, err := r.db.Exec(query, transactionId, userId); err != nil {
		r.logger.Errorf("failed delete transaction: %v", err)
		return errs.InternalServer
	}

	return nil
}
