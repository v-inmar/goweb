package todo_models

import (
	"database/sql"
	"time"
)

type TitleModel struct {
	ID int64
	Value string
	DateCreated time.Time
}

// Create a new row in the title model for the given value
func (m *TitleModel) Create(dbSession *sql.Tx, value string) error {
	dt := time.Now().UTC()

	model, err := dbSession.Exec("insert into title_model (value, date_created) values (?, ?)", value, dt)
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

// Read row depending on the id passed in
func (m *TitleModel) ReadById(db *sql.DB, id int64)error{
	if err := db.QueryRow("select * from title_model where id=?", id).Scan(&m.ID, &m.Value, &m.DateCreated); err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}

// Read row depending on the value passed in
func (m *TitleModel) ReadByValue(db *sql.DB, value string)error{
	if err := db.QueryRow("select * from title_model where value=?", value).Scan(&m.ID, &m.Value, &m.DateCreated); err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}