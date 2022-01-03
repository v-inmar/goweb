package todo_models

import (
	"database/sql"
	"time"
)

type BodyModel struct {
	ID int64
	Value string
	DateCreated time.Time
}

func (m *BodyModel) Create(dbSession *sql.Tx, value string)error{
	dt := time.Now().UTC()

	model, err := dbSession.Exec("insert into body_model (value, date_created) values (?, ?)", value, dt)
	if err != nil{
		return err
	}

	id, err := model.LastInsertId()
	if err != nil{
		return err
	}

	m.ID = id
	m.Value = value
	m.DateCreated = dt
	return nil
}

func (m *BodyModel) ReadById(db *sql.DB, id int64)error{
	if err := db.QueryRow("select * from body_model where id=?", id).Scan(&m.ID, &m.Value, &m.DateCreated); err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}