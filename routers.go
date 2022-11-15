package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// Используется функция http.NewServeMux() для инициализации нового рутера, затем

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)          // обработчик для домашней страницы
	mux.HandleFunc("/singup", app.singing) // обработчик для домашней страницы
	//mux.HandleFunc("/email", app.checkEmail) // обработчик для домашней страницы
	mux.HandleFunc("/auth", app.login)                 // обработчик авторизации
	mux.HandleFunc("/notes", app.showSnippet)          // обработчик для отображения записей
	mux.HandleFunc("/notes/create", app.createSnippet) // обработчик для создания записей
	mux.HandleFunc("/logout", app.logout)              // обработчик для выхода из аккаунта
	mux.HandleFunc("/account", app.userInfo)           // обработчик для просмотра информации о своем аккаунте
	mux.HandleFunc("/user", app.userPage)              // обработчик для просмотра информации об аккаунте другого пользователя
	mux.HandleFunc("/follow", app.follow)              // обработчик для подписки на пользователя
	mux.HandleFunc("/unfollow", app.unfollow)          // обработчик отписки от пользователя
	mux.HandleFunc("/subscribes", app.showSubList)
	mux.HandleFunc("/followers", app.showFollowList)
    mux.HandleFunc("/comment", app.comment)
	mux.HandleFunc("/update", app.updateNote)
    mux.HandleFunc("/deleteComment", app.deleteComment)


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
