package web

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	bt "studio/basic_types"
	"studio/studio"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const TEMPLATES_PATH string = "web/html/"

type Web struct {
	st  *studio.Studio
	srv *http.Server
	rdb *redis.Client
	ctx context.Context
}

type User struct {
	Id           uint
	Login        string
	UserFullName string
	AccLevel     bt.AccessLevel
}

func (u *User) FullName() string            { return u.UserFullName }
func (u *User) AccessLevel() bt.AccessLevel { return u.AccLevel }
func (u *User) GetId() uint                 { return u.Id }
func (u *User) GetLogin() string            { return u.Login }

func New(pgSqlSocket, redisSocket, httpSocket string) (web *Web, err error) {
	web = &Web{}
	if web.st, err = studio.New(pgSqlSocket); err != nil {
		return nil, err
	}
	web.ctx = context.Background()

	web.rdb = redis.NewClient(&redis.Options{
		Addr:     redisSocket,
		Password: "",
		DB:       0,
	})

	pong, err := web.rdb.Ping(web.ctx).Result()
	if err != nil {
		panic(err)
	}
	log.Println("REDIS", pong)

	web.srv = web.initHttp(httpSocket)

	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)

	return web, nil
}

func (web *Web) Run() error {
	log.Println("Запуск веб-интерфейса...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Прерываю...")
		if err := web.srv.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	if err := web.srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Веб-сервер закрыт")
		} else {
			return err
		}
	}

	return nil
}

func (web *Web) initHttp(webSocket string) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Загрузка шаблонов
	router.FuncMap = template.FuncMap{
		"inc": inc,
		"eq":  eq,
	}
	router.LoadHTMLGlob(TEMPLATES_PATH + "*.html")

	router.Use(
		web.checkCookies,
		web.checkSession,
		metricsMiddleware(),
	)

	// Маршруты
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/", web.mainHandler)
	router.POST("/", web.mainHandler)
	router.GET("/login", web.loginHandler)
	router.POST("/do_login", web.doLoginHandler)
	router.GET("/register", web.registerHandler)
	router.POST("/register", web.registerHandler)
	router.GET("/orders", web.ordersHandler)
	router.POST("/orders", web.ordersHandler)
	router.GET("/order-items", web.orderItemsHandler)
	router.GET("/model", web.viewModelHandler)
	router.GET("/create-order", web.createOrderHandler)
	router.POST("/create-order", web.createOrderHandler)

	// Обработка статических файлов
	router.Static("/styles", TEMPLATES_PATH+"styles")
	router.Static("/scripts", TEMPLATES_PATH+"scripts")

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound,
			"alert.html",
			gin.H{
				"Msg": "Страница не найдена",
			},
		)
	})

	return &http.Server{
		Addr:    webSocket,
		Handler: router,
	}
}

func (web *Web) checkCookies(c *gin.Context) {
	if !web.allCookiesExists(c) &&
		c.Request.URL.Path != "/metrics" &&
		c.Request.URL.Path != "/login" &&
		c.Request.URL.Path != "/do_login" &&
		c.Request.URL.Path != "/register" &&
		c.Request.URL.Path != "/styles/style.css" {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	if web.allCookiesExists(c) &&
		(c.Request.URL.Path == "/login" ||
			c.Request.URL.Path == "/do_login" ||
			c.Request.URL.Path == "/register") {
		c.Redirect(http.StatusSeeOther, "/")
		c.Abort()
		return
	}

	c.Next()
}

func (web *Web) checkSession(c *gin.Context) {
	// Список защищённых путей
	protectedPaths := map[string]bool{
		"/":             true,
		"/orders":       true,
		"/order-items":  true,
		"/create-order": true,
	}

	// Проверяем, защищён ли путь
	if protectedPaths[c.Request.URL.Path] {
		entity := web.entityFromSession(c)
		if entity == nil {
			c.SetCookie("login", "", 0, "/", "", false, true)
			c.SetCookie("session_id", "", 0, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
	}

	c.Next()
}

func inc(a int) int {
	return a + 1
}

func eq(a bt.OrderStatus, b uint) bool {
	return uint(a) == b
}

func customerOptions() []string {
	return []string{
		"Создать заказ",
		"Просмотреть заказы",
		"Выход",
	}
}

func operatorOptions() []string {
	return []string{
		"Просмотреть заказы",
		"Выход",
	}
}
