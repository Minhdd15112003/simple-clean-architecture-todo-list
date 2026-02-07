package main

import (
	"log"
	"net/http"
	"os"
	"social-todo-list/common"
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
	likeItemHandler "social-todo-list/module/userlikeitem/handler"
	likeItemGin "social-todo-list/module/userlikeitem/handler/gin_item"
	likeItemRepo "social-todo-list/module/userlikeitem/storage"
	likeItemUsecase "social-todo-list/module/userlikeitem/use_case"
	"social-todo-list/pubsub"
	"social-todo-list/subscriber"

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
	router.Use(middleware.Recover())

	v1 := router.Group("/v1")

	//Setup AUTH
	authStorage := userStorage.NewSqlStore(db)
	tokenprovider := jwt.NewTokenJWTProvider(common.PluginJWT, SECRET_KEY)
	middlewareAuth := middleware.RequiredAuth(authStorage, tokenprovider)

	// Layer 1: Storage (Data Access)
	itemStore := storage.NewSqlStore(db)
	likeStore := likeItemRepo.NewSqlStore(db)

	// Setup PubSub Subscribers
	ps := pubsub.NewPubSub()
	subEngine := subscriber.NewEngine(ps)
	subEngine.Start(itemStore)

	// Layer 2: Use Cases (Business Logic)
	itemUseCase := usecase.NewItemUseCase(itemStore, likeStore)
	likeUseCase := likeItemUsecase.NewUserLikeItemUseCase(likeStore, ps)

	// Layer 3: Services (Application Logic)
	itemService := handler.NewItemService(itemUseCase)
	likeService := likeItemHandler.NewLikeService(likeUseCase)

	// Layer 4: Handlers (HTTP Transport)
	itemHandler := ginitem.NewGinItemHandler(itemService)
	likeHandler := likeItemGin.NewGinLikeHandler(likeService)
	likeRpc := likeItemGin.NewLikeRPC(likeStore)
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
			// Item routes
			items.GET("", itemHandler.GetItems)
			items.GET("/:id", itemHandler.GetItem)
			items.POST("", itemHandler.CreateItem)
			items.PATCH("/:id", itemHandler.UpdateItem)
			items.DELETE("/:id", itemHandler.DeleteItem)

			// Like routes
			items.GET("/:id/list-user", likeHandler.ListUserLikeItem)
			items.POST("/:id/like", likeHandler.LikeItem)
			items.DELETE("/:id/unlike", likeHandler.UnLikeItem)
		}
		rpc := v1.Group("rpc")
		{
			rpc.POST("/item-like", likeRpc.GetLikeItem)
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
