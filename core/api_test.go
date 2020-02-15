package core

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)



func TestLoginClient_LoginNotOkForInvalidPassword(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
  CREATE TABLE clients(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
  login TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL);`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO clients (id, login, password) VALUES (1, 'don', 'don');`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	 _, err = Login("don", "xer","clients", db)

	if err == nil {
		t.Errorf("Not ErrInvalidPass error for invalid pass: %v", err)
	}
}

func TestAddClientWithTheSameData(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(clientsDDL)
	if err != nil {
		t.Errorf("can't create table clients err: %v", err)
	}
	_,err = AddNewClient(db, Client{1, "user", "user", "user", "user","A65",0,0})
	if err == nil {
		t.Errorf("can't insert client user: %v", err)
	}
	_,err = AddNewClient (db, Client{1, "user", "user", "user", "user","A65",0,0})
	if err == nil {
		t.Errorf("error client  user already exist but new user added: %v", err)
	}
}

func TestAddAccountToClientWhitWrong_ID_AndNegative_Balance(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(accountsDDL)
	if err != nil {
		t.Errorf("can't create table account err: %v", err)
	}
	_, err = AddNewClient(db, Client{1, "user", "user", "user", "user","A85",0,0})
	if err ==nil{
		t.Errorf("Wrong id. Account wasn't added")
	}
	if err ==nil{
		t.Errorf("Error balance can't be negative")
	}
}