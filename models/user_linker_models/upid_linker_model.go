package user_linker_models

import (
	"database/sql"
	"time"
)

type UPIDLinkerModel struct {
	ID int64
	UserID int64
	UpidID int64
	DateCreated time.Time
}

func (m *UPIDLinkerModel) Create(dbSession *sql.Tx, userID, upidID int64) error {
	// dbSession, err := db.Begin()
	// if err != nil {
	// 	return err
	// }
	// defer dbSession.Rollback()

	dt := time.Now().UTC()

	model, err := dbSession.Exec("insert into user_upid_linker_model (user_id, upid_id, date_created) values (?,?,?)", userID, upidID, dt)
	if err != nil{
		return err
	}

	insertedID, err := model.LastInsertId()
	if err != nil{
		return err
	}

	// if err := dbSession.Commit(); err != nil{
	// 	return err
	// }

	m.ID = insertedID
	m.UserID = userID
	m.UpidID = upidID
	m.DateCreated = dt
	return nil
}


func (m *UPIDLinkerModel) ReadById(db *sql.DB, id int64) error {
	err := db.QueryRow("select * from user_upid_linker_model where id=?", id).Scan(&m.ID, &m.UserID, &m.UpidID, &m.DateCreated)
	if err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}

func (m *UPIDLinkerModel) ReadByUserId(db *sql.DB, userID int64) error {
	err := db.QueryRow("select * from user_upid_linker_model where user_id=?", userID).Scan(&m.ID, &m.UserID, &m.UpidID, &m.DateCreated)
	if err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}

func (m *UPIDLinkerModel) ReadByUPIDId(db *sql.DB, upidID int64) error {
	err := db.QueryRow("select * from user_upid_linker_model where upid_id=?", upidID).Scan(&m.ID, &m.UserID, &m.UpidID, &m.DateCreated)
	if err != nil{
		if err != sql.ErrNoRows{
			return err
		}
	}
	return nil
}