package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"wallet-app/pkg/service"
)

type Handler struct {
	services *service.Service
	logger   *logrus.Logger
}

func NewHandler(services *service.Service, logger *logrus.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		wallets := api.Group("/wallets")
		{
			wallets.POST("/", h.createWallet)
			wallets.GET("/", h.getAllWallets)
			wallets.GET("/:id", h.getWalletById)
			wallets.PUT("/", h.addMember)
			wallets.DELETE("/:id", h.deleteWallet)
			wallets.DELETE("/", h.deleteMember)

			transactions := wallets.Group(":id/transactions")
			{
				transactions.POST("/", h.createTransaction)
				transactions.GET("/", h.getAllTransactions)
				transactions.GET("/incomes/", h.getAllIncomes)
				transactions.GET("/incomes/category", h.getByCategoryIncome)
				transactions.GET("/expenses", h.getAllExpenses)
				transactions.GET("/expenses/category", h.getByCategoryExpenses)
			}
		}
		//
		transactions := api.Group("/transactions")
		{
			//	items.GET("/:id", h.getTransactionsById)
			transactions.PUT("/:id", h.updateTransaction)
			transactions.DELETE("/:id", h.deleteTransaction)
		}
	}

	return router
}
