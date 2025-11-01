package http

import (
	"gogroceries/config"
	"gogroceries/delivery/middleware"
	"gogroceries/domain"
	"gogroceries/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	engine *gin.Engine,
	cfg config.Config,
	authUC domain.AuthUsecase,
	userUC domain.UserUsecase,
	tokoUC domain.TokoUsecase, 
	produkUC domain.ProdukUsecase, 
	categoryUC domain.CategoryUsecase, 
	trxUC domain.TrxUsecase, 
	alamatUC domain.AlamatUsecase,
	jwtAuth helper.JWTInterface,
) {
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})

	apiV1 := engine.Group("/api/v1")

	NewAuthHandler(apiV1, authUC) 

	userRoutes := apiV1.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware(jwtAuth)) 
	{
		userHandler := NewUserHandler(userUC, jwtAuth)
		alamatHandler := NewAlamatHandler(alamatUC, jwtAuth)

		userRoutes.GET("", userHandler.GetMyProfile) 
		userRoutes.PUT("", userHandler.UpdateProfile) 

		alamatRoutes := userRoutes.Group("/alamat")
		{
			alamatRoutes.GET("", alamatHandler.GetAllAlamatUser)    
			alamatRoutes.POST("", alamatHandler.CreateAlamat)   
			alamatRoutes.GET("/:id", alamatHandler.GetAlamatByID) 
			alamatRoutes.PUT("/:id", alamatHandler.UpdateAlamat) 
			alamatRoutes.DELETE("/:id", alamatHandler.DeleteAlamat) 
		}
	}
	tokoRoutes := apiV1.Group("/toko")
	tokoHandler := NewTokoHandler(tokoUC, jwtAuth)
	{
		tokoRoutes.GET("/my", middleware.AuthMiddleware(jwtAuth), tokoHandler.GetMyToko)
		tokoRoutes.PUT("/:id_toko", middleware.AuthMiddleware(jwtAuth), tokoHandler.UpdateToko) 
		tokoRoutes.GET("", tokoHandler.GetAllToko) 
		tokoRoutes.GET("/:id_toko", tokoHandler.GetTokoByID) 
	}

	productRoutes := apiV1.Group("/product")
	produkHandler := NewProdukHandler(produkUC, jwtAuth)
	{
		productRoutes.POST("", middleware.AuthMiddleware(jwtAuth), produkHandler.CreateProduk)   
		productRoutes.PUT("/:id", middleware.AuthMiddleware(jwtAuth), produkHandler.UpdateProduk) 
		productRoutes.DELETE("/:id", middleware.AuthMiddleware(jwtAuth), produkHandler.DeleteProduk) 
		productRoutes.GET("", produkHandler.GetAllProduk)    
		productRoutes.GET("/:id", produkHandler.GetProdukByID) 
	}

	categoryRoutes := apiV1.Group("/category")
	categoryHandler := NewCategoryHandler(categoryUC)
	{
		categoryRoutes.GET("", categoryHandler.GetAllCategories)    
		categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID) 

		adminCategoryRoutes := categoryRoutes.Use(middleware.AuthMiddleware(jwtAuth), middleware.AdminMiddleware())
		{
			adminCategoryRoutes.POST("", categoryHandler.CreateCategory)  
			adminCategoryRoutes.PUT("/:id", categoryHandler.UpdateCategory) 
			adminCategoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory) 
		}
	}

	trxRoutes := apiV1.Group("/trx")
	trxRoutes.Use(middleware.AuthMiddleware(jwtAuth)) 
	{
		trxHandler := NewTrxHandler(trxUC, jwtAuth) 
		trxRoutes.POST("", trxHandler.CreateTransaksi)  
		trxRoutes.GET("", trxHandler.GetAllTransaksiUser)    
		trxRoutes.GET("/:id", trxHandler.GetTransaksiByID) 
	}
}