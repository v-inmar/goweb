package user_linker_models

import (
	"database/sql"
	"time"
)


type AuthLinkerModel struct{
	ID int64
	UserID int64
	AuthID int64
	DateCreated time.Time
	DateUpdated sql.NullTime
}

func (m *AuthLinkerModel) Create(dbSession *sql.Tx, userID, authID int64) error{
	dt := time.Now().UTC()
	model, err := dbSession.Exec("insert into user_auth_linker_model (user_id, auth_id, date_created) values (?,?,?)", userID, authID, dt)
	if err != nil{
		return err
	}

	insertedID, err := model.LastInsertId()
	if err != nil{
		return err
	}

	m.ID = insertedID
	m.UserID = userID
	m.AuthID = authID
	m.DateCreated = dt
	return nil
}

