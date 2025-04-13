package web

import (
	"context"
	"html/template"
	"log"
	"net/http"
	bt "studio/basic_types"
	"studio/studio"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const TEMPLATES_PATH string = "web/html/"

type Web struct {
	st  *studio.Studio
	rdb *redis.Client
	ctx context.Context
}

type User struct {
	Id       uint
	Login    string
	FullName string
	AccLevel bt.AccessLevel
}

func (u *User) FirstLastName() string       { return u.FullName }
func (u *User) AccessLevel() bt.AccessLevel { return u.AccLevel }
func (u *User) GetId() uint                 { return u.Id }
func (u *User) GetLogin() string            { return u.Login }

func New(dbPath string) (web *Web, err error) {
	web = &Web{}
	if web.st, err = studio.New(dbPath); err != nil {
		return nil, err
	}
	web.ctx = context.Background()

	web.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := web.rdb.Ping(web.ctx).Result()
	if err != nil {
		panic(err)
	}
	log.Println("REDIS", pong)

	return web, nil
}

func (web *Web) Run() error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Загрузка шаблонов
	r.FuncMap = template.FuncMap{
		"inc": inc,
	}
	r.LoadHTMLGlob(TEMPLATES_PATH + "*.html")

	r.Use(web.checkCookies)

	// Маршруты
	r.GET("/", web.mainHandler)
	r.POST("/", web.mainHandler)
	r.GET("/login", web.loginHandler)
	r.POST("/do_login", web.doLoginHandler)
	r.GET("/register", web.registerHandler)
	r.POST("/register", web.registerHandler)
	r.GET("/orders", web.ordersHandler)
	r.POST("/orders", web.ordersHandler)
	r.GET("/order-items", web.orderItemsHandler)
	r.GET("/model", web.viewModelHandler)
	r.GET("/create-order", web.createOrderHandler)
	r.POST("/create-order", web.createOrderHandler)

	// Обработка статических файлов
	r.Static("/styles", TEMPLATES_PATH+"styles")
	r.Static("/scripts", TEMPLATES_PATH+"scripts")

	log.Println("Запуск веб-интерфейса...")
	return r.Run(":8080")
}

func (web *Web) checkCookies(c *gin.Context) {
	if !web.allCookiesExists(c) &&
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
			c.Request.URL.Path == "/register") {
		c.Redirect(http.StatusSeeOther, "/")
		c.Abort()
		return
	}

	c.Next()
}

func inc(a int) int {
	return a + 1
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
