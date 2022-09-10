package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// Используется функция http.NewServeMux() для инициализации нового рутера, затем
	// функцию "home" регистрируется как обработчик для URL-шаблона "/".
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
  mux.HandleFunc("/singup", app.singing)
	//mux.HandleFunc("email", app.singing)
	//mux.HandleFunc("/email", app.checkEmail)
	mux.HandleFunc("/auth", app.login)
	mux.HandleFunc("/notes", app.showSnippet)
	mux.HandleFunc("/notes/create", app.createSnippet)
	mux.HandleFunc("/logout", app.logout)
	mux.HandleFunc("/account", app.userInfo)
	mux.HandleFunc("/user", app.userPage)

	// Инициализируем FileServer, он будет обрабатывать
	// HTTP-запросы к статическим файлам из папки "./ui/static".
	// Обратите внимание, что переданный в функцию http.Dir путь
	// является относительным корневой папке проекта
	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Используем функцию mux.Handle() для регистрации обработчика для
	// всех запросов, которые начинаются с "/static/". Мы убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
