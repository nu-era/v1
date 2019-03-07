package handlers

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

const get = "select * from alerts where id = ?"
const getDev = "select * from alerts where deviceID = ?"
const getUN = "select * from alerts where user_name = ?"
const ins = "insert into alerts(email,pass_hash,user_name,first_name,last_name,photo_url) values(?,?,?,?,?,?)"
const upd = "update alerts set first_name = ?, last_name = ? where id = ?"
const del = "delete from alerts where id = ?"
const getUSERS = "select id, user_name, first_name, last_name from alerts"

type MySqlStore struct {
	db *sql.DB
}

func NewMySqlStore(db *sql.DB) *MySqlStore {
	if db == nil {
		return nil
	}
	return &MySqlStore{
		db: db,
	}
}

//GetByID returns the Alert with the given ID
func (mysql *MySqlStore) GetByID(id int64) (*Alert, error) {
	var row *sql.Row
	row = mysql.db.QueryRow(get, id)
	alert := Alert{}

	// scan row values into alert struct
	if err := row.Scan(&alert.ID, &alert.Message, &alert.DeviceID, &alert.Status,
		&alert.CreatedAt, &alert.EditedAt); err != nil {
		return nil, errors.New("Alert not found")
	}
	return &alert, nil
}

//Insert inserts the alert into the database, and returns
//the newly-inserted alert, complete with the DBMS-assigned ID
func (mysql *MySqlStore) Insert(alt *Alert) (*Alert, error) {
	tx, err := mysql.db.Begin() // begin transaction
	if err != nil {
		return nil, err
	}

	res, err := tx.Exec(ins, alt.Message, alt.DeviceID, alt.Status, alt.CreatedAt, alt.EditedAt)
	if err != nil {
		tx.Rollback() // rollback transaction
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback() // rollback transaction
		return nil, err
	}
	alt.ID = id
	tx.Commit()
	return alt, nil
}

func (mysql *MySqlStore) Update(id int64, upd *AlertUpdates) (*Alert, error) {
	return nil, nil
}
