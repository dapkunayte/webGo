package postgre

import (
	"main/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Get(username string) (*models.User, error) {
	stmt := `SELECT username, email FROM users WHERE username = $1`

	// Используем метод QueryRow() для выполнения SQL запроса,
	// передавая ненадежную переменную id в качестве значения для плейсхолдера
	// Возвращается указатель на объект sql.Row, который содержит данные записи.
	row, _ := m.DB.Query(stmt, username)
	defer row.Close()

	// Инициализируем указатель на новую структуру Snippet.
	s := &models.User{}

	// Используйте row.Scan(), чтобы скопировать значения из каждого поля от sql.Row в
	// соответствующее поле в структуре Snippet. Обратите внимание, что аргументы
	// для row.Scan - это указатели на место, куда требуется скопировать данные
	// и количество аргументов должно быть точно таким же, как количество
	// столбцов в таблице базы данных.
	err := row.Scan(&s.Username, &s.Email)
	if err != nil {
		// Специально для этого случая, мы проверим при помощи функции errors.Is()
		// если запрос был выполнен с ошибкой. Если ошибка обнаружена, то
		// возвращаем нашу ошибку из модели models.ErrNoRecord.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Если все хорошо, возвращается объект Snippet.
	return s, nil
}

func (m *UserModel) CheckUsers(users models.User) (bool, string) {
	checkForUsers := "SELECT * FROM users"
	rows, err := m.DB.Query(checkForUsers)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var checkedUsers []*models.User

	for rows.Next() {
		cu := &models.User{}
		err = rows.Scan(&cu.Username, &cu.Password, &cu.Email)
		if err != nil {
			panic(err)
		}
		checkedUsers = append(checkedUsers, cu)
	}

	for i := range checkedUsers {
		if users.Username == checkedUsers[i].Username {
			return true, "Такой пользователь уже есть"
			fmt.Println("Такой пользователь уже есть")
			break
		}
	}

	return false, ""
}

func (m *UserModel) Login(login string, password string) (bool, string) {
	checkForUsers := "SELECT * FROM users"
	rows, err := m.DB.Query(checkForUsers)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var checkedUsers []*models.User

	for rows.Next() {
		cu := &models.User{}
		err = rows.Scan(&cu.Username, &cu.Password, &cu.Email)
		if err != nil {
			panic(err)
		}
		checkedUsers = append(checkedUsers, cu)
	}
	//hashedPasswordFromForm, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	var matchedUserPassword string
	var matchedUser bool = false
	for i := range checkedUsers {
		if login == checkedUsers[i].Username {
			matchedUserPassword = checkedUsers[i].Password
			matchedUser = true
		}
	}

	if matchedUser == true && bcrypt.CompareHashAndPassword([]byte(matchedUserPassword), []byte(password)) != nil {
		return true, "Неверный логин или пароль"
	} else if matchedUser == false {
		return true, "Такого пользователя не существует"
	}

	return false, ""
}

func (m *UserModel) Singin(users models.User, checkBool bool) {
	//users := &models.User{}

	if checkBool == false {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), 8)
		stmt := "INSERT INTO users VALUES ($1, $2, $3)"

		row, err := m.DB.Query(stmt, users.Username, string(hashedPassword), users.Email)
		defer row.Close()
		if err != nil {
			panic(err)
		} else {
			fmt.Println("User has added")
		}
	}
}
