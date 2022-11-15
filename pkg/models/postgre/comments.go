package postgre

import (
    "main/pkg/models"
	"database/sql"
)

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Get(noteId int) ([]*models.Comment, error) {
	// Пишем SQL запрос, который мы хотим выполнить.
	stmt := `SELECT id,username,content FROM comments WHERE noteId = $1
    ORDER BY date DESC`

	rows, err := m.DB.Query(stmt, noteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment

	for rows.Next() {
		// Создаем указатель на новую структуру Snippet
		c := &models.Comment{}
		err = rows.Scan(&c.ID, &c.Username, &c.Content)
		if err != nil {
			return nil, err
		}
		// Добавляем структуру в срез.
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Если все в порядке, возвращаем срез с данными.
	return comments, nil
}

func (m *CommentModel) Add(comment models.Comment) {
	//users := &models.User{}
	stmt := "INSERT INTO comments (date,content,username,noteId) VALUES (now(), $1, $2, $3)"

	row, err := m.DB.Query(stmt, comment.Content, comment.Username, comment.NoteId)
	defer row.Close()
	if err != nil {
		panic(err)
	}
}

func (m *CommentModel) Delete(commentId int) error {
	stmt := "DELETE FROM comments WHERE id = $1"
	row, err := m.DB.Query(stmt, commentId)
	defer row.Close()
	if err != nil {
		panic(err)
		return err
	}
	return nil
}
