package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	wallet "wallet-app/pkg/models"
)

type getTransactionsByCategory struct {
	Data []wallet.TransactionsByCategory `json:"data"`
}

func (h *Handler) createTransaction(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input wallet.Transaction
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id int
	id, err = h.services.Transaction.Create(userId, walletId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllTransactions(c *gin.Context) {
	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid wallet id param")
		return
	}

	var input wallet.TransactionDate
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	transactions, err := h.services.Transaction.GetAll(walletId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *Handler) getAllIncomes(c *gin.Context) {
	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid wallet id param")
		return
	}

	var input wallet.TransactionDate
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	incomes, err := h.services.Transaction.GetAllIncomes(walletId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, incomes)
}

func (h *Handler) getAllExpenses(c *gin.Context) {
	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid wallet id param")
		return
	}

	var input wallet.TransactionDate
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	expenses, err := h.services.Transaction.GetAllExpenses(walletId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func (h *Handler) getByCategoryIncome(c *gin.Context) {
	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid wallet id param")
		return
	}

	var input wallet.TransactionDate
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	transactions, err := h.services.Transaction.GetByCategoryIncome(walletId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getTransactionsByCategory{
		Data: transactions,
	})
}

func (h *Handler) getByCategoryExpenses(c *gin.Context) {
	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid wallet id param")
		return
	}

	var input wallet.TransactionDate
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	transactions, err := h.services.Transaction.GetByCategoryExpenses(walletId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getTransactionsByCategory{
		Data: transactions,
	})
}

func (h *Handler) updateTransaction(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	transactionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid transaction id param")
		return
	}

	var input wallet.UpdateTransactionInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Transaction.Update(userId, transactionId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteTransaction(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	transactionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid wallet id param")
		return
	}

	err = h.services.Transaction.Delete(userId, transactionId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
