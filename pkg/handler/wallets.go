package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	wallet "wallet-app/pkg/models"
)

type getAllWalletsResponse struct {
	Data []wallet.Wallet `json:"data"`
}

func (h *Handler) createWallet(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input wallet.Wallet
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Wallet.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllWallets(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	wallets, err := h.services.Wallet.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllWalletsResponse{
		Data: wallets,
	})
}

func (h *Handler) getWalletById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	wallet, err := h.services.Wallet.GetById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *Handler) deleteWallet(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Wallet.Delete(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) addMember(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input wallet.Member
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Wallet.GetOwnerIdQuery(userId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err = h.services.Wallet.AddMember(userId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteMember(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input wallet.Member
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Wallet.GetOwnerIdQuery(userId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.services.Wallet.DeleteMember(userId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
