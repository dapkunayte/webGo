package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"runtime/debug"
	"unicode"
 // "github.com/eaigner/jet"
)

// Помощник serverError записывает сообщение об ошибке в errorLog и
// затем отправляет пользователю ответ 500 "Внутренняя ошибка сервера".
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Помощник clientError отправляет определенный код состояния и соответствующее описание
// пользователю. Мы будем использовать это в следующий уроках, чтобы отправлять ответы вроде 400 "Bad
// Request", когда есть проблема с пользовательским запросом.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Мы также реализуем помощник notFound. Это просто
// удобная оболочка вокруг clientError, которая отправляет пользователю ответ "404 Страница не найдена".
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func ConnectDB() (*sql.DB, error) {
  db, err := sql.Open("postgres", "host=abul.db.elephantsql.com port=5432 user=vfslaqjo dbname=vfslaqjo password=CKBQsUjB8sfEyCgYcG1kwI7cQkE0b2Kt sslmode=disable")

	if err != nil {
		return nil, err
	}
	return db, nil
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidPass(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Пароль должен состоять не менее чем из 8 символов"
	} else {
		permissions := []string{"!", "\"", "#", "$", "%", "&", "‘", "(", ")", "*", "+", ",", "–", ".", "/", ":", ";", "?", "@", "[", "]", "^", "_", "`"}
		isUpper := true
		isLower := true
		isSpecial := true
		//isNotPermiss := false
		for i := range password {
			if unicode.IsUpper(rune(password[i])) {
				isUpper = false
			}
			if unicode.IsLower(rune(password[i])) {
				isLower = false
			}
			for j := range permissions {
				if string(password[i]) == permissions[j] {
					isSpecial = false
				}
			}
		}
		if isUpper == true || isLower == true {
			return false, "Пароль должнен содержать хотя бы одну строчную и заглавную букву"
		}
		if isSpecial == true {
			return false, "Пароль должен содержать хотя бы один специальный символ"
		}
		/*if isNotPermiss == false{
			return false, "Пароль содержит запрещённые символы"
		}

		*/
	}
	return true, ""
}

func ValidLogin(password string) (bool, string) {
	if len(password) < 4 {
		return false, "Логин должен состоять не менее чем из 4 символов"
	} else {
		permissions := []string{"!", "\"", "#", "$", "%", "&", "‘", "(", ")", "*", "+", ",", ".", "/", ":", ";", "?", "@", "[", "]", "^", "`", " "}
		isSpecial := true
		for i := range password {
			for j := range permissions {
				if string(password[i]) == permissions[j] {
					isSpecial = false
				}
			}
		}
		if isSpecial == false {
			return false, "Логин содержит запрещённые символы"
		}
		/*if isNotPermiss == false{
			return false, "Пароль содержит запрещённые символы"
		}

		*/
	}
	return true, ""
}
