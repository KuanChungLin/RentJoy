package main

import (
	"log"
	"net/http"
	"rentjoy/database"
	"rentjoy/internal/controllers"
	"rentjoy/internal/middleware"
	"rentjoy/internal/services"
	"rentjoy/pkg/redis"
	templateManager "rentjoy/pkg/template"
)

func main() {
	// 初始化模板
	tmplManager := templateManager.NewManager()
	if err := tmplManager.InitTemplates(); err != nil {
		log.Fatal("Failed to initialize templates: ", err)
	}

	// 連接至 DB
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to init database: ", err)
	}

	// 初始化 Redis 客戶端
	redisClient := redis.NewRedisClient()

	// 初始化 services
	homeService := services.NewHomeService(db)
	accountService := services.NewAccountService(db)
	searchService := services.NewSearchService(db)
	venueService := services.NewVenueService(db)
	orderService := services.NewOrderService(db)
	manageService := services.NewManageService(db)

	// 初始化排程
	scheduleService := services.NewScheduleService(db)
	// 啟動排程任務
	scheduleService.OrderSchedule()

	// 初始化 controllers
	homepageController := controllers.NewHomePageController(
		homeService,
		tmplManager.GetTemplates(),
	)
	accountController := controllers.NewAccountController(
		accountService,
		tmplManager.GetTemplates(),
	)
	searchpageController := controllers.NewSearchPageController(
		searchService,
		tmplManager.GetTemplates(),
	)
	venuepageController := controllers.NewVenuePageController(
		venueService,
		tmplManager.GetTemplates(),
		redisClient,
	)
	orderController := controllers.NewOrderController(
		orderService,
		tmplManager.GetTemplates(),
		redisClient,
	)
	ecpayController := controllers.NewEcpayController(
		tmplManager.GetTemplates(),
		redisClient,
	)
	manageController := controllers.NewManangeController(
		manageService,
		tmplManager.GetTemplates(),
	)

	// 配置靜態文件路徑
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../public"))))

	// 路由處理
	http.HandleFunc("/", homepageController.HomePage)
	http.HandleFunc("/Policy", homepageController.Privacy)
	http.HandleFunc("/Faq", homepageController.Faq)
	http.HandleFunc("/UserRules", homepageController.UserRules)
	http.HandleFunc("/OrderRules", homepageController.OrderRules)
	http.HandleFunc("/BecomeVenueOwner", homepageController.BecomeVenueOwner)
	http.HandleFunc("/SignUp", accountController.SignUp)
	http.HandleFunc("/Login", accountController.Login)
	http.HandleFunc("/Auth/Facebook", accountController.FacebookLogin)
	http.HandleFunc("/Auth/Facebook/Callback", accountController.FacebookCallback)
	http.HandleFunc("/Auth/Google", accountController.GoogleLogin)
	http.HandleFunc("/Auth/Google/Callback", accountController.GoogleCallback)
	http.HandleFunc("/Logout", accountController.Logout)
	http.HandleFunc("/SignUpDetail", accountController.SignUpDetail)
	http.HandleFunc("/Account", middleware.AuthMiddleware(accountController.Account))
	http.HandleFunc("/UpdateName", middleware.AuthMiddleware(accountController.UpdateName))
	http.HandleFunc("/UpdateEmail", middleware.AuthMiddleware(accountController.UpdateEmail))
	http.HandleFunc("/UpdatePhone", middleware.AuthMiddleware(accountController.UpdatePhone))
	http.HandleFunc("/UpdatePassword", middleware.AuthMiddleware(accountController.UpdatePassword))
	http.HandleFunc("/SearchPage", searchpageController.SearchPage)
	http.HandleFunc("/SearchPageLoading", searchpageController.SearchPageLoading)
	http.HandleFunc("/Venue/VenuePage", venuepageController.VenuePage)
	http.HandleFunc("/Venue/GetAvailableTime", venuepageController.GetAvailableTime)
	http.HandleFunc("/Venue/ReservedPage", venuepageController.ReservedPage)
	http.HandleFunc("/Venue/OrderPending", middleware.AuthMiddleware(venuepageController.OrderPending))
	http.HandleFunc("/Order/CreateOrder", middleware.AuthMiddleware(orderController.CreateOrder))
	http.HandleFunc("/Order/OrderReserved", middleware.AuthMiddleware(orderController.OrderReserved))
	http.HandleFunc("/Order/OrderProcessing", middleware.AuthMiddleware(orderController.OrderProcessing))
	http.HandleFunc("/Order/OrderCancel", middleware.AuthMiddleware(orderController.OrderCancel))
	http.HandleFunc("/Order/OrderFinished", middleware.AuthMiddleware(orderController.OrderFinished))
	http.HandleFunc("/Order/CancelReservation", middleware.AuthMiddleware(orderController.CancelReservation))
	http.HandleFunc("/Order/SaveEvaluate", middleware.AuthMiddleware(orderController.SaveEvaluate))
	http.HandleFunc("/Ecpay/Process", middleware.AuthMiddleware(ecpayController.Process))
	http.HandleFunc("/Ecpay/ReceivePaymentResult", ecpayController.ReceivePaymentResult)
	http.HandleFunc("/Manage/ReservedManagement", middleware.AuthMiddleware(manageController.ReservedManagement))
	http.HandleFunc("/Manage/VenueManagement", middleware.AuthMiddleware(manageController.VenueManagement))
	http.HandleFunc("/Manage/ReservedAccept", middleware.AuthMiddleware(manageController.ReservedAccept))
	http.HandleFunc("/Manage/ReservedReject", middleware.AuthMiddleware(manageController.ReservedReject))
	http.HandleFunc("/Manage/DelistVenue", middleware.AuthMiddleware(manageController.DelistVenue))
	http.HandleFunc("/Manage/DeleteVenue", middleware.AuthMiddleware(manageController.DeleteVenue))

	log.Println("伺服器運行中：https://localhost:8080")
	log.Fatal(http.ListenAndServeTLS(":8080", "../../cert.pem", "../../key.pem", nil))
}
