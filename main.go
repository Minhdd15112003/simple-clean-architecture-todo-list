package main

import (
	"log"
	"net/http"
	"os"
	"social-todo-list/components/tokenprovider/jwt"
	"social-todo-list/middleware"
	"social-todo-list/module/item/handler"
	ginitem "social-todo-list/module/item/handler/gin_item"
	"social-todo-list/module/item/storage"
	usecase "social-todo-list/module/item/use_case"
	"social-todo-list/module/upload"
	handlerUser "social-todo-list/module/user/handler"
	"social-todo-list/module/user/handler/ginuser"
	userStorage "social-todo-list/module/user/storage"
	userUseCase "social-todo-list/module/user/use_case"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	gin.SetMode(gin.DebugMode) // Set mode trước
	DB_CONN := os.Getenv("DB_CONN")
	SECRET_KEY := os.Getenv("SECRET_KEY")

	db, err := gorm.Open(mysql.Open(DB_CONN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln("DataBaseError: ", err)
	}
	log.Println("DB connection: ", db)

	router := gin.Default()
	router.Static("/static", "./static")

	v1 := router.Group("/v1")
	authStorage := userStorage.NewSqlStore(db)
	tokenprovider := jwt.NewTokenJWTProvider("jwt", SECRET_KEY)
	middlewareAuth := middleware.RequiredAuth(authStorage, tokenprovider)
	{
		v1.PUT("/upload", upload.Upload(db))
		users := v1.Group("auth")
		{
			useCase := userUseCase.NewAuthUseCase(authStorage, tokenprovider, 60*60*24*30)
			service := handlerUser.NewUserService(useCase)
			handler := ginuser.NewGinUserHandler(service)
			users.POST("/register", handler.Register)
			users.POST("/login", handler.Login)
			users.GET("/profile", middlewareAuth, handler.Profile)
		}

		items := v1.Group("items", middlewareAuth)
		{
			// Dependency Injection: Storage -> UseCase -> Service -> Handler
			storage := storage.NewSqlStore(db)            // layer Data Access/Storage: tương tác với cơ sở dữ liệu
			useCase := usecase.NewItemUseCase(storage)    // layer Business Logic/UseCase: thực hiện các nghiệp vụ chính của ứng dụng
			service := handler.NewItemService(useCase)    // layer Interface Adapter/Service: chuyển đổi dữ liệu giữa UseCase và Handler
			handler := ginitem.NewGinItemHandler(service) //layer Frameworks/Drivers/Transport: thực hiện giao tiếp với bên ngoài (HTTP, gRPC, ...)

			items.GET("", handler.GetItems)
			items.GET("/:id", handler.GetItem)
			items.POST("", handler.CreateItem)
			items.PATCH("/:id", handler.UpdateItem)
			items.DELETE("/:id", handler.DeleteItem)
		}
	}
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	if err := router.Run(":8000"); err != nil {
		log.Fatalln(err.Error())
	} // listens on 0.0.0.0:8080 by default
}
