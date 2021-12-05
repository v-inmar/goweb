package user_models

import (
	"database/sql"
	"time"
)

type FirstnameModel struct {
	ID int64
	Value string
	DateCreated time.Time
}

func (m *FirstnameModel) Create(dbSession *sql.Tx, value string) error {

	dt := time.Now().UTC()

	model, err := dbSession.Exec("insert into firstname_model (value, date_created) values (?, ?)", value, dt)
	if err != nil{
		return err
	}

	insertedID, err := model.LastInsertId()
	if err != nil{
		return err
	}

	m.ID = insertedID
	m.Value = value
	m.DateCreated = dt
	return nil
}

func (m *FirstnameModel) ReadByValue(db *sql.DB, value string) error {
	err := db.QueryRow("select * from firstname_model where value=?", value).Scan(&m.ID, &m.Value, &m.DateCreated)
	if err != nil {
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}

func (m *FirstnameModel) ReadById(db *sql.DB, id int64) error {
	err := db.QueryRow("select * from firstname_model where id=?", id).Scan(&m.ID, &m.Value, &m.DateCreated)
	if err != nil {
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}