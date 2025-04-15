package main

import (
	"github.com/DecodeWorms/BorderBlitz/config"
	"github.com/DecodeWorms/BorderBlitz/repository"
	"github.com/DecodeWorms/BorderBlitz/utils"
)

var c config.Config

func main() {

	//Call database
	c = *config.NewConfig()
	db, err := utils.SetUpDatabase(&c)
	if err != nil {
		panic("error connecting to DB..")
	}
	wallet := repository.NewWalletRepository(db)
	coin := repository.NewStablecoinRepository(db)
	user := repository.NewUsers(db)
	tx := repository.NewTransactionRepository(db)

	// Call services
	walletService := utils.SetUpWalletService(wallet, coin, user)
	usService := utils.SetUpUserService(user)
	txService := utils.SetUpTransactionService(db, tx, wallet, coin)
	swapService := utils.SetUpSwapService(db, wallet, coin, tx)
	_ = utils.SetUpFxService(coin)

	// Call handlers
	explorerHandler := utils.SetUpExplorerHandler(txService, walletService)
	swapHandler := utils.SetUpSwapHandler(swapService)
	transactionHandler := utils.SetUpTransactionHandler(txService)
	walletHandler := utils.SetUpWalletHandler(walletService)
	userHandler := utils.SetUpUserHandler(usService)

	// Call router
	router := utils.SetUpRouter(explorerHandler, transactionHandler, swapHandler, walletHandler, userHandler, db)

	utils.StartServer(router)

}
