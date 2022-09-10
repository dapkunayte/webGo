package postgre

import (
	"main/pkg/models"
	"database/sql"
)

type SubModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *SubModel) Insert(sub string, follow string) error {
	// Ниже будет SQL запрос, который мы хотим выполнить. Мы разделили его на две строки
	// для удобства чтения (поэтому он окружен обратными кавычками
	// вместо обычных двойных кавычек).
	stmt := `INSERT INTO subscribes (date, sub_name, follow_name)
    VALUES(now(), $1, $2)`

	// Используем метод Exec() из встроенного пула подключений для выполнения
	// запроса. Первый параметр это сам SQL запрос, за которым следует
	// заголовок заметки, содержимое и срока жизни заметки. Этот
	// метод возвращает объект sql.Result, который содержит некоторые основные
	// данные о том, что произошло после выполнения запроса.
	result, _ := m.DB.Query(stmt, sub, follow)
	defer result.Close()
	// Возвращаемый ID имеет тип int64, поэтому мы конвертируем его в тип int
	// перед возвратом из метода.
	return nil
}

func (m *SubModel) GetUsersSub(username string) ([]*models.Subscribe, error) {
	// Пишем SQL запрос, который мы хотим выполнить.
	stmt := `SELECT id, follow_name FROM subscribes WHERE sub_name = $1
    ORDER BY date DESC`

	// Используем метод Query() для выполнения нашего SQL запроса.
	// В ответ мы получим sql.Rows, который содержит результат нашего запроса.
	rows, err := m.DB.Query(stmt, username)
	if err != nil {
		return nil, err
	}

	// Откладываем вызов rows.Close(), чтобы быть уверенным, что набор результатов из sql.Rows
	// правильно закроется перед вызовом метода Latest(). Этот оператор откладывания
	// должен выполнится *после* проверки на наличие ошибки в методе Query().
	// В противном случае, если Query() вернет ошибку, это приведет к панике
	// так как он попытается закрыть набор результатов у которого значение: nil.
	defer rows.Close()

	// Инициализируем пустой срез для хранения объектов models.Snippets.
	var subs []*models.Subscribe

	// Используем rows.Next() для перебора результата. Этот метод предоставляем
	// первый а затем каждую следующею запись из базы данных для обработки
	// методом rows.Scan().
	for rows.Next() {
		// Создаем указатель на новую структуру Snippet
		s := &models.Subscribe{}
		// Используем rows.Scan(), чтобы скопировать значения полей в структуру.
		// Опять же, аргументы предоставленные в row.Scan()
		// должны быть указателями на место, куда требуется скопировать данные и
		// количество аргументов должно быть точно таким же, как количество
		// столбцов из таблицы базы данных, возвращаемых вашим SQL запросом.
		err = rows.Scan(&s.ID, &s.FollowId)
		if err != nil {
			return nil, err
		}
		// Добавляем структуру в срез.
		subs = append(subs, s)
	}

	// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
	// если в ходе работы у нас не возникла какая либо ошибка.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Если все в порядке, возвращаем срез с данными.
	return subs, nil
}

func (m *SubModel) GetUsersFolls(username string) ([]*models.Subscribe, error) {
	// Пишем SQL запрос, который мы хотим выполнить.
	stmt := `SELECT id, sub_name FROM subscribes WHERE follow_name = $1
    ORDER BY date DESC`

	rows, err := m.DB.Query(stmt, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Инициализируем пустой срез для хранения объектов models.Snippets.
	var subs []*models.Subscribe

	for rows.Next() {
		// Создаем указатель на новую структуру Snippet
		s := &models.Subscribe{}
		err = rows.Scan(&s.ID, &s.SubId)
		if err != nil {
			return nil, err
		}
		// Добавляем структуру в срез.
		subs = append(subs, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Если все в порядке, возвращаем срез с данными.
	return subs, nil
}
func (m *SubModel) Delete(sub string, follow string) error {
	// Ниже будет SQL запрос, который мы хотим выполнить. Мы разделили его на две строки
	// для удобства чтения (поэтому он окружен обратными кавычками
	// вместо обычных двойных кавычек).
	stmt := "DELETE FROM subscribes WHERE sub_name=$1 AND follow_name=$2"

	// Используем метод Exec() из встроенного пула подключений для выполнения
	// запроса. Первый параметр это сам SQL запрос, за которым следует
	// заголовок заметки, содержимое и срока жизни заметки. Этот
	// метод возвращает объект sql.Result, который содержит некоторые основные
	// данные о том, что произошло после выполнения запроса.
	result, _ := m.DB.Query(stmt, sub, follow)
	defer result.Close()
	// Возвращаемый ID имеет тип int64, поэтому мы конвертируем его в тип int
	// перед возвратом из метода.
	return nil
}