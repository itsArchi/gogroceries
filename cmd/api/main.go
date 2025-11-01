package main

import (
	"fmt"
	"gogroceries/config"
	"gogroceries/delivery/http"
	"gogroceries/internal/helper"
	"gogroceries/domain"
	"gogroceries/repository/postgres"
	"gogroceries/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func Main() {
	
}

func main() {
	config.LoadConfig()
	cfg := config.AppConfig

	db := postgres.ConnectDatabase(cfg)

	log.Println("Migrating database...")
	errMigrate := db.AutoMigrate(
		&domain.User{},
		&domain.Toko{},
		&domain.Alamat{},
		&domain.Category{},
		&domain.Produk{},
		&domain.FotoProduk{},
		&domain.Trx{},
		&domain.DetailTrx{},
		&domain.LogProduk{},
	)
	if errMigrate != nil {
		log.Fatalf("Failed to migrate database: %v", errMigrate)
	}
	log.Println("Database migration successful!")

	jwtAuth := helper.NewJWTHelper(cfg)

	userRepo := postgres.NewPostgresUserRepository(db)
	tokoRepo := postgres.NewPostgresTokoRepository(db)
	produkRepo := postgres.NewPostgresProdukRepository(db)
	categoryRepo := postgres.NewPostgresCategoryRepository(db)
	trxRepo := postgres.NewPostgresTrxRepository(db)
	alamatRepo := postgres.NewPostgresAlamatRepository(db)

	authUC := usecase.NewAuthUsecase(userRepo, tokoRepo, jwtAuth)
	userUC := usecase.NewUserUsecase(userRepo)
	tokoUC := usecase.NewTokoUsecase(tokoRepo)
	produkUC := usecase.NewProdukUsecase(produkRepo, tokoRepo, categoryRepo)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo)
	trxUC := usecase.NewTrxUsecase(trxRepo, produkRepo, alamatRepo, categoryRepo, tokoRepo)
	alamatUC := usecase.NewAlamatUsecase(alamatRepo)

	engine := gin.Default()

	http.SetupRouter(
		engine,
		cfg,
		authUC,
		userUC,
		tokoUC,
		produkUC,
		categoryUC,
		trxUC,
		alamatUC,
		jwtAuth,
	)

	serverAddress := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server berjalan di http://localhost%s", serverAddress)
	err := engine.Run(serverAddress)
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}