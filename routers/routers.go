package routers

import (
	"ResumeMamageSystem/controller"
	"ResumeMamageSystem/setting"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if setting.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionId", store))
	// 告诉gin框架模板文件引用的静态文件去哪里找
	//r.Static("/static", "static")
	//// 告诉gin框架去哪里找模板文件
	//r.LoadHTMLGlob("templates/*")

	//user
	userRoutes := r.Group("/api")
	{

		userRoutes.POST("/user", controller.Register)
		userRoutes.GET("/user", controller.Me)
		userRoutes.PUT("/user", controller.UpdateUser)
		userRoutes.GET("/logout", controller.Logout)
		userRoutes.POST("/login", controller.Login)
		userRoutes.POST("/reset-password", controller.ResetPassword)
		//userRoutes.POST("/new-validation", controller.NewValidation)

	}
	rewardRoutes := r.Group("/api")
	{

		rewardRoutes.POST("/reward", controller.CreateReward)
		rewardRoutes.POST("/upload-reward", controller.UploadReward)
		rewardRoutes.DELETE("/reward", controller.DeleteReward)
		rewardRoutes.GET("/reward", controller.GetReward)
		rewardRoutes.GET("/my-reward", controller.GetMyReward)
		rewardRoutes.PUT("/reward", controller.UpdateReward)
		rewardRoutes.GET("/download-reward", controller.DownloadReward)
	}
	resumeRoutes := r.Group("/api")
	{

		resumeRoutes.POST("/resume", controller.CreateResume)
		resumeRoutes.POST("/upload-resume", controller.UploadResume)
		resumeRoutes.DELETE("/resume", controller.DeleteResume)
		resumeRoutes.GET("/resume", controller.GetResume)
		resumeRoutes.GET("/my-resume", controller.GetMyResume)
		resumeRoutes.PUT("/resume", controller.UpdateResume)
		resumeRoutes.GET("/download-resume", controller.DownloadResume)
	}
	return r
}
