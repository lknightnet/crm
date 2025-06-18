package controller

import (
	"github.com/gin-gonic/gin"
	"project-service/internal/controller/lists"
	"project-service/internal/controller/project"
	"project-service/internal/controller/task"
	"project-service/internal/controller/timer"
	"project-service/internal/service"
)

func RouteAPI(route *gin.Engine, services *service.Service) {
	route.Use(gin.Logger())

	route.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := route.Group("/api")
	//api.Use(LoggerMiddleware())

	projectController := project.NewProjectController(services.ProjectService)
	projectApi := api.Group("/project")
	projectApi.Use(AuthMiddleware())

	projectApi.POST("/create", projectController.CreateProject)
	projectApi.GET("/get/list", projectController.GetProjectsByToken)
	projectApi.GET("/get/:id", projectController.GetProjectByID)
	projectApi.POST("/update", projectController.EditProject)
	//projectApi.GET("/like/:name", projectController.GetProjectsByName)

	taskController := task.NewTaskController(services.TaskService)
	taskApi := api.Group("/task")
	taskApi.Use(AuthMiddleware())

	taskApi.POST("/create", taskController.CreateTask)
	taskApi.GET("/get/project/:id", taskController.GetTasksByProjectID)
	taskApi.GET("/get/user", taskController.GetTasksByToken)           //uses
	taskApi.GET("/get/expired", taskController.GetTaskExpired)         //uses
	taskApi.GET("/get/today", taskController.GetTaskToday)             //uses
	taskApi.GET("/get/thisweek", taskController.GetTaskThisWeek)       //uses
	taskApi.GET("/get/nextweek", taskController.GetTaskNextWeek)       //uses
	taskApi.GET("/get/notdeadline", taskController.GetTaskNotDeadline) //uses
	taskApi.GET("/get/single/:id", taskController.GetTaskByID)
	taskApi.POST("/update", taskController.EditTask)

	timerController := timer.NewTimerController(services.TimerService)
	timerApi := api.Group("/timer")
	timerApi.Use(AuthMiddleware())

	timerApi.POST("/start", timerController.StartTimer)
	timerApi.POST("/stop", timerController.StopTimer)
	timerApi.GET("/resume/:id", timerController.ResumeTimer)
	//timerApi.GET("/get/task/:id", timerController.GetTimersByTaskID)
	//timerApi.GET("/get/user", timerController.GetTimersByUserID)

	listController := lists.NewInformationListController(services.InformationListService)
	listApi := api.Group("/lists")
	listApi.Use(AuthMiddleware())
	listApi.POST("/create", listController.CreateInformationList)
	listApi.GET("/get/project/:id", listController.GetListByProjectID)
	listApi.POST("/update", listController.UpdateInformationList)
	listApi.GET("/delete/:id", listController.DeleteInformationList)

}
