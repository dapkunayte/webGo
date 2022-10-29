package main

import (
	"main/pkg/models"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	//"math/rand"
	"net/http"
	"strconv"
	"time"
  "strings"
)

type ViewData struct {
	Text  string
	Check bool
}

type templateData struct {
	Note             *models.Note
	Notes            []*models.Note
	IsAuth           bool
	Username         string
	OtherUsername    string
	Subscribes_count int
	Follows_count    int
	Subscribes       []*models.Subscribe
	Sub_fact         bool
}
var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

//домашняя страница
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	switch session.Values["authenticated"].(type) {
	case nil:
		session.Values["authenticated"] = false
	}
	switch session.Values["name"].(type) {
	case nil:
		session.Values["name"] = ""
	}

	if r.URL.Path != "/" {
		app.notFound(w) // Использование помощника notFound()
		return
	}
<<<<<<< HEAD
	n, err := app.notes.Latest()
=======
  switch session.Values["authenticated"].(type)
  {
    case nil:
     session.Values["authenticated"] = false
  }
  switch session.Values["name"].(type)
  {
    case nil:
     session.Values["name"] = ""
  }
  n, err := app.notes.Latest()
>>>>>>> origin/main
	if err != nil {
		app.serverError(w, err)
		return
	}

	for i := range n {
		n[i].Created = strings.Replace(n[i].Created, "T", " ", -1)[0:19]
		t, _ := time.Parse("2006-01-02 15:04:05", n[i].Created)
		n[i].Created = t.Format("02.01.2006, 15:04")
	}

	//cockie_data := &CockieData{isAuth: session.Values["authenticated"].(bool), Username: session.Values["name"].(string)}
	data := &templateData{Notes: n, IsAuth: session.Values["authenticated"].(bool), Username: session.Values["name"].(string)}
	files := []string{
		"./ui/html/main_page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
		"./ui/html/header.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Использование помощника serverError()
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err) // Использование помощника serverError()
	}

}

//регистрация
func (app *application) singing(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/footer.partial.html",
	}
	users := &models.User{}
	if r.Method == http.MethodGet {
		data := ViewData{
			Text: "",
		}
		// Используем функцию template.ParseFiles() для чтения файла шаблона.
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err) // Использование помощника serverError()
			return
		}

		err = ts.Execute(w, data)
		if err != nil {
			app.serverError(w, err) // Использование помощника serverError()
		}
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		users.Password = r.Form["password"][0]
		users.Username = r.Form["username"][0]
		users.Email = r.Form["email"][0]
		checkPass := r.Form["checkPass"][0]
		checkMail := ValidEmail(users.Email)
		checkPassword, strPass := ValidPass(users.Password)
		checkLogin, strLog := ValidLogin(users.Username)
		checkUser, strUser := app.users.CheckUsers(*users)

		switch {
		case checkLogin == false:
			data := ViewData{
				Text: strLog,
			}
			tmpl, _ := template.ParseFiles(files...)
			tmpl.Execute(w, data)
		case checkPassword == false:
			data := ViewData{
				Text: strPass,
			}
			tmpl, _ := template.ParseFiles(files...)
			tmpl.Execute(w, data)
		case checkMail == false:
			data := ViewData{
				Text: "Некорректная почта",
			}
			tmpl, _ := template.ParseFiles(files...)
			tmpl.Execute(w, data)
		case users.Password != checkPass:
			data := ViewData{
				Text: "Пароли не совпадают",
			}
			tmpl, _ := template.ParseFiles(files...)
			tmpl.Execute(w, data)
		case users.Password == "" || users.Username == "" || users.Email == "":
			data := ViewData{
				Text: "Поля не могут быть пустыми",
			}
			tmpl, _ := template.ParseFiles(files...)
			tmpl.Execute(w, data)
		case checkUser == true:
			{
				data := ViewData{
					Text: strUser,
				}
				tmpl, _ := template.ParseFiles(files...)
				tmpl.Execute(w, data)
			}
		default:
			app.users.Singin(*users, checkUser)
			session, _ := store.Get(r, "cookie-name")
			session.Values["authenticated"] = true
			session.Values["name"] = users.Username
			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
	
			//447595loH!
		}
	}
}

//вход в систему
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		files := []string{
			"./ui/html/login.html",
			"./ui/html/footer.partial.html",
		}
		// Используем функцию template.ParseFiles() для чтения файла шаблона.
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err) // Использование помощника serverError()
		}

		err = ts.Execute(w, nil)
		if err != nil {
			app.serverError(w, err) // Использование помощника serverError()
		}
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		login := r.Form["username"][0]
		password := r.Form["password"][0]

		checkLogin, strLogin := app.users.Login(login, password)
		if checkLogin == false {
			//app.users.Singin(*users, false)
			session, _ := store.Get(r, "cookie-name")
			session.Values["authenticated"] = true
			session.Values["name"] = login
			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			files := []string{
				"./ui/html/login.html",
				"./ui/html/footer.partial.html",
			}
			// Используем функцию template.ParseFiles() для чтения файла шаблона.
			data := ViewData{Text: strLogin}
			tmpl, _ := template.ParseFiles(files...)
			tmpl.Execute(w, data)
		}
	}
}

//отображение заметки - не реализован шаблон
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	switch session.Values["authenticated"].(type) {
	case nil:
		session.Values["authenticated"] = false
	}

	switch session.Values["name"].(type) {
	case nil:
		session.Values["name"] = ""
	}

	// Извлекаем значение параметра id из URL и попытаемся
	// конвертировать строку в integer используя функцию strconv.Atoi(). Если его нельзя
	// конвертировать в integer, или значение меньше 1, возвращаем ответ
	// 404 - страница не найдена!
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Использование помощника notFound()
	}

	// Вызываем метода Get из модели Snipping для извлечения данных для
	// конкретной записи на основе её ID. Если подходящей записи не найдено,
	// то возвращается ответ 404 Not Found (Страница не найдена).
	s, err := app.notes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	s.Created = strings.Replace(s.Created, "T", " ", -1)[0:19]
	t, _ := time.Parse("2006-01-02 15:04:05", s.Created)
	s.Created = t.Format("02.01.2006, 15:04")

	data := &templateData{Note: s, IsAuth: session.Values["authenticated"].(bool), Username: session.Values["name"].(string)}
	files := []string{
		"./ui/html/note.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
		"./ui/html/header.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Использование помощника serverError()
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err) // Использование помощника serverError()
	}

	// Отображаем весь вывод на странице.
}

// Обработчик для создания новой заметки. 
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Используем r.Method для проверки, использует ли запрос метод POST или нет. Обратите внимание,
	// что http.MethodPost является строкой и содержит текст "POST".
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if r.Method == http.MethodGet {
		files := []string{
			"./ui/html/add_note.html",
			"./ui/html/footer.partial.html",
			"./ui/html/header.partial.html",
			"./ui/html/base.layout.html",
		}
		// Используем функцию template.ParseFiles() для чтения файла шаблона.
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err) // Использование помощника serverError()
		}
		data := &templateData{IsAuth: session.Values["authenticated"].(bool), Username: session.Values["name"].(string)}
		err = ts.Execute(w, data)
		if err != nil {
			app.serverError(w, err) // Использование помощника serverError()
		}
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		title := r.Form["title"][0]
		content := r.Form["content"][0]
		str := fmt.Sprintf("%v", session.Values["name"])
		username := str
		//fmt.Println(title, content, username)
		// Передаем данные в метод SnippetModel.Insert(), получая обратно
		// ID только что созданной записи в базу данных.

		id, err := app.notes.Insert(title, content, username)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// Перенаправляем пользователя на соответствующую страницу заметки.
		http.Redirect(w, r, fmt.Sprintf("/notes?id=%d", id), http.StatusSeeOther)
	}
	// Создаем несколько переменных, содержащих тестовые данные. Мы удалим их позже.

}

//функция выхода (очистка куки)
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Values["name"] = ""
  session.Values["sub_fact"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//возможность пользователя просматрирвать информацию о своем аккаунте
func (app *application) userInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, _ := store.Get(r, "cookie-name")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		s, err := app.notes.GetUsersNotes(session.Values["name"].(string))
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		for i := range s {
			s[i].Created = strings.Replace(s[i].Created, "T", " ", -1)[0:19]
			t, _ := time.Parse("2006-01-02 15:04:05", s[i].Created)
			s[i].Created = t.Format("02.01.2006, 15:04")
		}

		// получение списка подписчиков
		subs, err := app.subscribes.GetUsersSub(session.Values["name"].(string))
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		//получение списка подписок
		folls, err := app.subscribes.GetUsersFolls(session.Values["name"].(string))
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		// Отображаем весь вывод на странице.
		//fmt.Fprintf(w, "%v", u, s, len(s))
		files := []string{
			"./ui/html/account.html",
			"./ui/html/footer.partial.html",
			"./ui/html/header.partial.html",
			"./ui/html/base.layout.html",
		}
		// Используем функцию template.ParseFiles() для чтения файла шаблона.
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err) // Использование помощника serverError()
		}
		data := &templateData{Notes: s, IsAuth: session.Values["authenticated"].(bool), Username: session.Values["name"].(string), Subscribes_count: len(subs), Follows_count: len(folls)}

		err = ts.Execute(w, data)
		if err != nil {
			app.serverError(w, err) // Использование помощника serverError()
		}
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

//просмотр информации о другом пользователе
func (app *application) userPage(w http.ResponseWriter, r *http.Request) {
	username := string(r.URL.Query().Get("id"))
	//url := r.URL.Path
	session, _ := store.Get(r, "cookie-name")
	session.Values["sub_name"] = username
	if (session.Values["sub_name"].(string) == session.Values["name"].(string)) && (session.Values["sub_name"] != nil) && (session.Values["name"] != nil) {
		http.Redirect(w, r, "/account", http.StatusSeeOther)
	}

	switch session.Values["sub_fact"].(type) {
	case nil:
		session.Values["sub_fact"] = false
	}

	if r.Method == http.MethodGet {

		//получение списка записей пользователя
		s, err := app.notes.GetUsersNotes(username)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		for i := range s {
			s[i].Created = strings.Replace(s[i].Created, "T", " ", -1)[0:19]
			t, _ := time.Parse("2006-01-02 15:04:05", s[i].Created)
			s[i].Created = t.Format("02.01.2006, 15:04")
		}

		// получение списка подписчиков
		subs, err := app.subscribes.GetUsersSub(username)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		//получение списка подписок
		folls, err := app.subscribes.GetUsersFolls(username)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		// проверка факта подписки на пользователя
		if len(folls) != 0 {
			for i := range folls {
				if string(folls[i].SubId) == session.Values["name"].(string) {
					session.Values["sub_fact"] = true
					session.Save(r, w)
					break
				} else {
					session.Values["sub_fact"] = false
					session.Save(r, w)
				}
			}
		}
		// Отображаем весь вывод на странице.
		//fmt.Fprintf(w, "%v", u, s, len(s))
		files := []string{
			"./ui/html/account.html",
			"./ui/html/footer.partial.html",
			"./ui/html/header.partial.html",
			"./ui/html/base.layout.html",
		}
		// Используем функцию template.ParseFiles() для чтения файла шаблона.
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err) // Использование помощника serverError()
		}

		data := &templateData{Notes: s, IsAuth: session.Values["authenticated"].(bool), OtherUsername: username, Username: session.Values["name"].(string), Subscribes_count: len(subs), Follows_count: len(folls), Sub_fact: session.Values["sub_fact"].(bool)}

		err = ts.Execute(w, data)
		if err != nil {
			app.serverError(w, err) // Использование помощника serverError()
		}

	} else {
		if session.Values["sub_fact"] == true {
			session.Save(r, w)
			http.Redirect(w, r, "/unfollow", http.StatusSeeOther)
		} else {
			session.Save(r, w)
			http.Redirect(w, r, "/follow", http.StatusSeeOther)
		}
		//session.Save(r, w)
		//http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//обработка подписки
func (app *application) unfollow(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	app.subscribes.Delete(session.Values["name"].(string), session.Values["sub_name"].(string))
	session.Values["sub_fact"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/user?id="+session.Values["sub_name"].(string), http.StatusSeeOther)
}

//обработка отписки
func (app *application) follow(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	app.subscribes.Insert(session.Values["name"].(string), session.Values["sub_name"].(string))
	session.Values["sub_fact"] = true
	session.Save(r, w)
	http.Redirect(w, r, "/user?id="+session.Values["sub_name"].(string), http.StatusSeeOther)
}

func (app *application) showSubList(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	subs, err := app.subscribes.GetUsersSub(session.Values["name"].(string))
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	files := []string{
		"./ui/html/subscribes.html",
		"./ui/html/footer.partial.html",
		"./ui/html/header.partial.html",
		"./ui/html/base.layout.html",
	}
	// Используем функцию template.ParseFiles() для чтения файла шаблона.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Использование помощника serverError()
	}

	data := &templateData{Subscribes: subs, IsAuth: session.Values["authenticated"].(bool), Username: session.Values["name"].(string)}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err) // Использование помощника serverError()
	}
}

func (app *application) showFollowList(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	folls, err := app.subscribes.GetUsersFolls(session.Values["name"].(string))
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	files := []string{
		"./ui/html/followers.html",
		"./ui/html/footer.partial.html",
		"./ui/html/header.partial.html",
		"./ui/html/base.layout.html",
	}
	// Используем функцию template.ParseFiles() для чтения файла шаблона.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Использование помощника serverError()
	}

	data := &templateData{Subscribes: folls, IsAuth: session.Values["authenticated"].(bool), Username: session.Values["name"].(string)}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err) // Использование помощника serverError()
	}
}