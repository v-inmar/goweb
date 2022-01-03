package todo_models

import (
	"database/sql"
	"time"
)

type TPIDModel struct {
	ID int64
	Value string
	DateCreated time.Time
}

func (m *TPIDModel) Create(dbSession sql.Tx, value string)error{
	dt := time.Now().UTC()
	model, err := dbSession.Exec("insert into tpid_model (value, date_created) values (?, ?)", value, dt)
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


func (m *TPIDModel) ReadByValue(db *sql.DB, value string) error {
	err := db.QueryRow("select * from tpid_model where value=?", value).Scan(&m.ID, &m.Value, &m.DateCreated)
	if err != nil {
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}

func (m *TPIDModel) ReadById(db *sql.DB, id int64) error {
	err := db.QueryRow("select * from tpid_model where id=?", id).Scan(&m.ID, &m.Value, &m.DateCreated)
	if err != nil {
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}