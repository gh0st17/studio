package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	bt "studio/basic_types"
	"studio/studio"
	"sync"
)

const (
	TEMPLATES_PATH string = "web/html/"
	dateFormat     string = "02.01.2006 15:04:05"
)

type Web struct {
	st *studio.Studio

	sessionStore map[string]bt.Entity
	sessionMutex sync.RWMutex
}

func New(dbPath string) (web *Web, err error) {
	web = &Web{}
	if web.st, err = studio.New(dbPath); err != nil {
		return nil, err
	}
	web.sessionStore = make(map[string]bt.Entity)

	return web, nil
}

func (w *Web) Run() error {
	http.HandleFunc("/", w.mainHandler)
	http.HandleFunc("/login", w.loginHandler)
	http.HandleFunc("/do_login", w.doLoginHandler)
	http.HandleFunc("/register", w.registerHandler)
	http.HandleFunc("/orders", w.ordersHandler)
	http.HandleFunc("/order-items", w.orderItemsHandler)
	http.HandleFunc("/model", w.viewModelHandler)
	http.HandleFunc("/create-order", w.createOrderHandler)

	http.Handle("/styles/",
		http.StripPrefix(
			"/styles/",
			http.FileServer(http.Dir(TEMPLATES_PATH+"styles")),
		),
	)

	http.Handle("/scripts/",
		http.StripPrefix(
			"/scripts/",
			http.FileServer(http.Dir(TEMPLATES_PATH+"scripts")),
		),
	)

	fmt.Println("Запуск веб-интерфейса...")
	return http.ListenAndServe(":8080", nil)
}

func inc(a int) int {
	return a + 1
}

func (*Web) execTemplate(path string, w http.ResponseWriter, data interface{}) {
	// Создаем новый шаблон и регистрируем функцию инкремента
	tmpl := template.New(path).Funcs(template.FuncMap{
		"inc": inc,
	})

	// Парсим файлы шаблона
	tmpl, err := tmpl.ParseFiles(TEMPLATES_PATH + path)
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		log.Println("template error:", err)
		return
	}

	// Выполняем шаблон с данными
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка рендера шаблона", http.StatusInternalServerError)
		log.Println("execute error:", err)
	}
}

func (web *Web) isSessionExists(sessionID string) bool {
	web.sessionMutex.RLock()
	_, ok := web.sessionStore[sessionID]
	web.sessionMutex.RUnlock()

	return ok
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
