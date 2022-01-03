package todo_models

import (
	"database/sql"
	"time"
)

type TodoModel struct {
	ID int64
	DateCreated time.Time
	DateDeleted sql.NullTime
}

func (m *TodoModel) Create(dbSession *sql.Tx)error{
	dt := time.Now().UTC()

	model, err := dbSession.Exec("insert into todo_model (date_created) values (?)", dt)
	if err != nil{
		return err
	}

	id, err := model.LastInsertId()
	if err != nil{
		return err
	}

	m.ID = id
	m.DateCreated = dt
	return nil
}

func (m *TodoModel) ReadById(db *sql.DB, id int64)error{
	if err := db.QueryRow("select * from todo_model where id=?", id).Scan(&m.ID, &m.DateCreated, &m.DateDeleted); err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}