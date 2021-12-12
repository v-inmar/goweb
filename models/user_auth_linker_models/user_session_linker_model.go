package user_auth_linker_models

import (
	"database/sql"
	"time"
)

type UserSessionLinkerModel struct {
	ID int64
	UserID int64
	SessionID int64
	DateCreated time.Time
}


func (m *UserSessionLinkerModel) Create(dbSession *sql.Tx, userID int64, sessionID int64) error {
	dt := time.Now().UTC()

	model, err := dbSession.Exec("insert into user_session_linker_model (user_id, session_id, date_created) values (?,?,?)", userID, sessionID, dt)
	if err != nil{
		return err
	}

	insertedID, err := model.LastInsertId()
	if err != nil {
		return err
	}

	m.ID = insertedID
	m.UserID = userID
	m.SessionID = sessionID
	m.DateCreated = dt
	return nil
}

func (m *UserSessionLinkerModel) ReadByID(db *sql.DB, id int64) error {
	row := db.QueryRow("select * from user_session_linker_model where id=?", id)
	if err := row.Scan(&m.ID, &m.UserID, &m.SessionID, &m.DateCreated); err != nil{
		if err != sql.ErrNoRows{ // No row or no match is handled in the caller's scope
			return err
		}
	}
	return nil
}

func (m *UserSessionLinkerModel) ReadBySessionID(db *sql.DB, sessionID int64) error {
	row := db.QueryRow("select * from user_session_linker_model where seesion_id=?", sessionID)
	if err := row.Scan(&m.ID, &m.UserID, &m.SessionID, &m.DateCreated); err != nil{
		if err != sql.ErrNoRows{ // No row or no match is handled in the caller's scope
			return err
		}
	}
	return nil
}

func (m *UserSessionLinkerModel) ReadByUserID(db *sql.DB, userID int64) error {
	row := db.QueryRow("select * from user_session_linker_model where user_id=?", userID)
	if err := row.Scan(&m.ID, &m.UserID, &m.SessionID, &m.DateCreated); err != nil{
		if err != sql.ErrNoRows{ // No row or no match is handled in the caller's scope
			return err
		}
	}
	return nil
}

func (m *UserSessionLinkerModel) Delete(dbSession *sql.Tx) error {
	_, err := dbSession.Exec("delete from user_session_linker_model where id=?", m.ID)
	if err != nil{
		return err
	}
	return nil
}