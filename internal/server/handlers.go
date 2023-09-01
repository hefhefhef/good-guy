package server

import (
	"github.com/gin-gonic/gin"
	ctmsHttp "github.com/openuniland/good-guy/external/ctms/delivery"
	facebookHttp "github.com/openuniland/good-guy/external/facebook/delivery"
	fithouHttp "github.com/openuniland/good-guy/external/fithou/delivery"
	articleHttp "github.com/openuniland/good-guy/internal/articles/delivery"
	examSchedulesHttp "github.com/openuniland/good-guy/internal/exam_schedules/delivery"
	userHttp "github.com/openuniland/good-guy/internal/users/delivery"
	"github.com/openuniland/good-guy/jobs"

	ctmsUS "github.com/openuniland/good-guy/external/ctms/usecase"
	facebookUS "github.com/openuniland/good-guy/external/facebook/usecase"
	fithouUS "github.com/openuniland/good-guy/external/fithou/usecase"
	articleUS "github.com/openuniland/good-guy/internal/articles/usecase"
	examSchedulesUS "github.com/openuniland/good-guy/internal/exam_schedules/usecase"
	userUS "github.com/openuniland/good-guy/internal/users/usecase"

	articleRepo "github.com/openuniland/good-guy/internal/articles/repository"
	examSchedulesRepo "github.com/openuniland/good-guy/internal/exam_schedules/repository"
	userRepo "github.com/openuniland/good-guy/internal/users/repository"
	cors "github.com/rs/cors/wrapper/gin"
)

func (server *Server) MapHandlers() {
	router := gin.Default()
	router.Use(cors.AllowAll())

	// Init repositories
	articleRepo := articleRepo.NewArticleRepository(server.configs, server.mongo)
	userRepo := userRepo.NewUserRepository(server.configs, server.mongo)
	examSchedulesRepo := examSchedulesRepo.NewExamSchedulesRepository(server.configs, server.mongo)

	// Init useCases
	fithouUS := fithouUS.NewFithouUseCase(server.configs)
	articleUS := articleUS.NewArticleUseCase(server.configs, articleRepo, fithouUS)
	facebookUS := facebookUS.NewFacebookUseCase(server.configs)
	userUS := userUS.NewUserUseCase(server.configs, userRepo)
	examSchedulesUS := examSchedulesUS.NewExamSchedulesUseCase(server.configs, examSchedulesRepo)
	ctmsUC := ctmsUS.NewCtmsUseCase(server.configs, examSchedulesUS, facebookUS)

	// Init handlers
	authHandlers := ctmsHttp.NewCtmsHandlers(server.configs, ctmsUC)
	articleHandlers := articleHttp.NewArticleHandlers(server.configs, articleUS)
	fithouHandlers := fithouHttp.NewFithouHandlers(server.configs, fithouUS)
	facebookHandlers := facebookHttp.NewFacebookHandlers(server.configs, facebookUS)
	userHandlers := userHttp.NewArticleHandlers(server.configs, userUS)
	examSchedulesHandlers := examSchedulesHttp.NewExamSchedulesHandlers(server.configs, examSchedulesUS)

	// Jobs
	jobs := jobs.NewJobs(server.configs, articleUS, userUS, facebookUS, ctmsUC)
	jobs.Run()

	// Init routes
	v1 := router.Group("/v1")

	health := v1.Group("/health")
	ctms := v1.Group("/ctms")
	articles := v1.Group("/articles")
	fithou := v1.Group("/fithou")
	facebook := v1.Group("/facebook")
	users := v1.Group("/users")
	examSchedules := v1.Group("/exam-schedules")

	health.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Map routes
	ctmsHttp.MapCtmsRoutes(ctms, authHandlers)
	articleHttp.MapArticleRoutes(articles, articleHandlers)
	fithouHttp.MapFithouRoutes(fithou, fithouHandlers)
	facebookHttp.MapFacebookRoutes(facebook, facebookHandlers)
	userHttp.MapUserRoutes(users, userHandlers)
	examSchedulesHttp.MapExamSchedulesRoutes(examSchedules, examSchedulesHandlers)
	server.router = router
}
