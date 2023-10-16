package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/krishnamdhanush/web-dawn/handlers"
)

func Run() (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Recovery())
	albumRouters := router.Group("/albums")
	{
		albumRouters.GET("", handlers.AlbumCache(), handlers.GetAlbums)
		albumRouters.POST("", handlers.PostAlbum)
		albumRouters.GET("/:id", handlers.AlbumCacheById(), handlers.GetAlbumById)
	}
	err := router.Run(getPath())
	return router, err
}

func getPath() string {
	err := godotenv.Load()
	if err != nil {
		return "localhost:8000"
	}

	return os.Getenv("BASE_PATH")
}
