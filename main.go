package main

import (
	"main/pkg/models/postgre"
	"flag"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres dialect
	"log"
	"net/http"
	"os"
)

// Создаем структуру `application` для хранения зависимостей всего веб-приложения.
// Пока, что мы добавим поля только для двух логгеров, но
// мы будем расширять данную структуру по мере усложнения приложения.

// Добавляем поле snippets в структуру application. Это позволит
// сделать объект SnippetModel доступным для наших обработчиков.
type application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	users      *postgre.UserModel
	notes      *postgre.NoteModel
	subscribes *postgre.SubModel
	comments   *postgre.CommentModel
}

func main() {

	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Создаем логгер для записи сообщений об ошибках таким же образом, но используем stderr как
	// место для записи и используем флаг log.Lshortfile для включения в лог
	// названия файла и номера строки где обнаружилась ошибка.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Чтобы функция main() была более компактной, мы поместили код для создания
	// пула соединений в отдельную функцию openDB(). Мы передаем в нее полученный
	// источник данных (DSN) из флага командной строки.

	db, err := ConnectDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog:   errorLog,
		infoLog:    infoLog,
		users:      &postgre.UserModel{DB: db},
		notes:      &postgre.NoteModel{DB: db},
		subscribes: &postgre.SubModel{DB: db},
		comments:   &postgre.CommentModel{DB: db},
	}


	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	// Поскольку переменная `err` уже объявлена в приведенном выше коде, нужно
	// использовать оператор присваивания =
	// вместо оператора := (объявить и присвоить)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
