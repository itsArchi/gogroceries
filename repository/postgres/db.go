package postgres

import (
	"fmt"
	"gogroceries/config"
	"log"
	// "gogroceries/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta", 	
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err !=nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection successful")

// 	err = db.AutoMigrate(
// 	&domain.User{},
// 	&domain.Toko{},       
// 	&domain.Alamat{},     
// 	&domain.Category{},
// 	&domain.Produk{},     
// 	&domain.FotoProduk{}, 
// 	&domain.Trx{},        
// 	&domain.DetailTrx{}, 
// 	&domain.LogProduk{},  
// )
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration successful")

	return db
}
