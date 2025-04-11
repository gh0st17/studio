package web

import (
	"fmt"
	"net/http"
	bt "studio/basic_types"

	"github.com/google/uuid"
)

// Получаем имя пользователя из сессии
func (web *Web) getEntityFromSession(r *http.Request) (bt.Entity, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, false
	}

	web.sessionMutex.RLock()
	entity, exists := web.sessionStore[cookie.Value]
	web.sessionMutex.RUnlock()

	return entity, exists
}

func (web *Web) homeHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := web.getEntityFromSession(r); !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/main", http.StatusSeeOther)
	}
}

func (web *Web) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<form method="POST" action="/do_login">
			<input type="text" name="login" placeholder="Введите логин" required />
			<button type="submit">Log in</button>
		</form>
	`)
}

func (web *Web) doLoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm() error", http.StatusBadRequest)
		return
	}
	login := r.FormValue("login")
	if login == "" {
		http.Error(w, "login is required", http.StatusBadRequest)
		return
	}

	sessionID := uuid.New().String()

	web.sessionMutex.Lock()
	entity, err := web.st.Login(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	web.sessionStore[sessionID] = entity
	web.sessionMutex.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
	})

	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

func (web *Web) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		_ = r.FormValue("login")
		http.Redirect(w, r, fmt.Sprintf("/alert?msg=Регистрация+успешна,+%s+%s!&next=/login", firstName, lastName), http.StatusSeeOther)
		return
	}

	fmt.Fprint(w, `
		<html>
		<body>
			<h2>Регистрация</h2>
			<form method="POST">
				<label>Имя: <input type="text" name="first_name"></label><br><br>
				<label>Фамилия: <input type="text" name="last_name"></label><br><br>
				<label>Логин: <input type="text" name="login"></label><br><br>
				<input type="submit" value="Зарегистрироваться">
			</form>
		</body>
		</html>
	`)
}

func (web *Web) mainHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")

	fmt.Fprintf(w, `
		<html>
		<body>
			<h2>Главное меню</h2>
			<p>Добро пожаловать, %s!</p>
			<ul>
				<li><a href="/orders">Посмотреть заказы</a></li>
				<li><a href="/create-order">Создать заказ</a></li>
				<li><a href="/order-items">Просмотреть содержимое заказа</a></li>
			</ul>
		</body>
		</html>
	`, web.sessionStore[cookie.Value].FirstLastName())
}

func (web *Web) ordersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
		<html>
		<body>
			<h2>Список заказов</h2>
			<p>Здесь будет таблица заказов.</p>
			<a href="/main">Назад в меню</a>
		</body>
		</html>
	`)
}

func (web *Web) orderItemsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
		<html>
		<body>
			<h2>Содержимое заказа</h2>
			<p>Здесь будет содержимое выбранного заказа.</p>
			<a href="/main">Назад в меню</a>
		</body>
		</html>
	`)
}

func (web *Web) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/alert?msg=Заказ+успешно+создан!&next=/orders", http.StatusSeeOther)
		return
	}

	fmt.Fprint(w, `
		<html>
		<body>
			<h2>Создать заказ</h2>
			<form method="POST">
				<label>Описание заказа: <input type="text" name="description"></label><br><br>
				<input type="submit" value="Создать">
			</form>
			<a href="/main">Назад в меню</a>
		</body>
		</html>
	`)
}

func (web *Web) alertHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	next := r.URL.Query().Get("next")

	fmt.Fprintf(w, `
		<html>
		<body>
			<script>
				alert("%s");
				window.location.href = "%s";
			</script>
		</body>
		</html>
	`, msg, next)
}
