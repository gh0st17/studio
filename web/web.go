package web

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	bt "github.com/gh0st17/studio/basic_types"
	"github.com/gh0st17/studio/studio"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const TEMPLATES_PATH string = "web/html/"

type Web struct {
	st         *studio.Studio
	srv        *http.Server
	rdb        *redis.Client
	ctx        context.Context
	rdbPresent atomic.Bool
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
		log.Printf("REDIS: %v", err)
	} else {
		web.rdbPresent.Store(true)
		log.Println("REDIS", pong)
	}

	web.srv = web.initHttp(httpSocket)

	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(serverErrors)

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

	web.startRedisMonitor()

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
		"inc":     inc,
		"eq":      eq,
		"timeStr": timeToStr,
	}
	router.LoadHTMLGlob(TEMPLATES_PATH + "*.html")

	router.Use(
		web.checkCookies,
		requestsMetric,
		serverErrorsMetric,
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
		web.alert(c, http.StatusNotFound, notFound)
	})

	return &http.Server{
		Addr:    webSocket,
		Handler: router,
	}
}

func (web *Web) startRedisMonitor() {
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if _, err := web.rdb.Ping(web.ctx).Result(); err != nil {
					log.Printf("Redis is offline: %v", err)
					web.rdbPresent.Store(false)
				} else {
					web.rdbPresent.Store(true)
				}
			case <-web.ctx.Done():
				return
			}
		}
	}()
}

func (web *Web) checkCookies(c *gin.Context) {
	if !web.dataCookiesExists(c) &&
		c.Request.URL.Path != "/metrics" &&
		c.Request.URL.Path != "/login" &&
		c.Request.URL.Path != "/do_login" &&
		c.Request.URL.Path != "/register" &&
		c.Request.URL.Path != "/order-items" &&
		c.Request.URL.Path != "/styles/style.css" {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	if web.dataCookiesExists(c) &&
		(c.Request.URL.Path == "/login" ||
			c.Request.URL.Path == "/do_login" ||
			c.Request.URL.Path == "/register") {
		c.Redirect(http.StatusSeeOther, "/")
		c.Abort()
		return
	}

	c.Next()
}

func (web *Web) alert(c *gin.Context, code int, msg string) {
	c.HTML(code, "alert.html", gin.H{"Msg": msg})
}

func inc(a int) int {
	return a + 1
}

func eq(a bt.OrderStatus, b uint) bool {
	return uint(a) == b
}

func timeToStr(t time.Time) string {
	if t.Equal(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return "---"
	} else {
		return t.Format(bt.DateFormat)
	}
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
