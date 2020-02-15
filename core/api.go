package core

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

////-----------Structs-------------------------
type Client struct {
	Id         int
	Login      string
	Password   string
	Name       string
	Surname    string
	SerialPass string
	Phone      int
	MainAcc    int
}

type Accounts struct {
	Id       int
	Name     string
	Number   int
	Money    int
	ClientId int
}

type Services struct {
	Id            int
	Name          string
	IdPayment     string
	AccountNumber int
}

type Terminals struct {
	Id      int
	Number  string
	Address string
}

//----------------------------------------------

func Init(db *sql.DB) (err error) {
	ddls := []string{managersDDL, managersInitialData, clientsDDL, accountsDDL, initAccountDDL, terminalsDDL, servicesDDL}
	for _, ddl := range ddls {
		_, err = db.Exec(ddl)
		if err != nil {
			return err
		}
	}
	return nil
}

func Login(login, password, nameTable string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string

	if nameTable == "managers" {
		loginSQL = loginManSQL
	}
	if nameTable == "clients" {
		loginSQL = loginCliSQL
	}
	err := db.QueryRow(loginSQL, login).Scan(&dbLogin, &dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	if dbPassword != password {
		return false, nil
	}

	return true, nil
}

func AddTerminals(db *sql.DB, number, address string) error {
	var idTerminal int
	err := db.QueryRow(exitsTerminalSQL, number).Scan(&idTerminal)
	if err == nil {
		return errors.New("Terminal already exits")
	} else if err == sql.ErrNoRows {
		_, err = db.Exec(insertTerminal, number, address)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

//----------------------------------------------------------------------------------------------------------------------

func AddServices(db *sql.DB, name string, idPayment string) error {
	name = strings.Trim(name, "\n")
	idPayment = strings.Trim(idPayment, "\n")

	var idService int
	err := db.QueryRow(exitsServicesSQL, name).Scan(&idService)
	if err == nil {
		return errors.New("Service already exits")
	} else if err == sql.ErrNoRows {
		var numberLastAcc int
		err := db.QueryRow(lastAccountId).Scan(&numberLastAcc)
		if err != nil {
			return err
		}
		numberLastAcc++

		_, err = db.Exec(insertAccountSQL, name, numberLastAcc, 0, 0)
		if err != nil {
			return err
		}

		_, err = db.Exec(insertService, name, idPayment, numberLastAcc)
		if err != nil {
			return err
		}

		return nil
	}
	return err
}

func AddNewClient(db *sql.DB, newClinet Client) (int, error) {
	var numberLastAcc, clientId int
	err := db.QueryRow(lastAccountId).Scan(&numberLastAcc)
	if err != nil {
		return 0, err
	}
	numberLastAcc++

	_, err = db.Exec(insertClientSQL, newClinet.Login, newClinet.Password, newClinet.Name, newClinet.Surname, newClinet.SerialPass, newClinet.Phone, numberLastAcc)
	if err != nil {
		return 0, err
	}
	err = db.QueryRow(lastClientId, newClinet.Login).Scan(&clientId)
	if err != nil {
		return 0, err
	}

	_, err = db.Exec(insertAccountSQL, "Alif", numberLastAcc, 0, clientId)
	if err != nil {
		return 0, err
	}

	return numberLastAcc, nil
}

func AddAccountByLogin(db *sql.DB, login, name string) (int, error) {
	var clientId int
	err := db.QueryRow(idClientByLogin, login).Scan(&clientId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Login was not detected!")
		}
		return 0, err
	}

	//-----------------------------------------------------
	var numberLastAcc int
	err = db.QueryRow(lastAccountId).Scan(&numberLastAcc)
	if err != nil {
		return 0, err
	}
	numberLastAcc++
	_, err = db.Exec(insertAccountSQL, name, numberLastAcc, 0, clientId)
	if err != nil {
		return 0, err
	}
	//-----------------------------------------------------

	return numberLastAcc, nil
}

func AddAccountByPhone(db *sql.DB, phone, name string) (int, error) {
	var clientId int
	err := db.QueryRow(idClientByPhone, phone).Scan(&clientId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Login was not detected!")
		}
		return 0, err
	}

	//-----------------------------------------------------
	var numberLastAcc int
	err = db.QueryRow(lastAccountId).Scan(&numberLastAcc)
	if err != nil {
		return 0, err
	}
	numberLastAcc++
	_, err = db.Exec(insertAccountSQL, name, numberLastAcc, 0, clientId)

	if err != nil {
		return 0, err
	}
	//-----------------------------------------------------

	return numberLastAcc, nil
}

//----------------------------------------------------------------------------------------------------------------------

func ShowAccounts(db *sql.DB, accounts *[]Accounts) error {
	rows, err := db.Query(allAccounts)
	if err != nil {
		return err
	}
	for rows.Next() {
		var res Accounts
		err = rows.Scan(&res.Id, &res.Name, &res.Number, &res.Money, &res.ClientId)
		if err != nil {
			return err
		}
		*accounts = append(*accounts, res)
	}
	return nil
}

func ShowTerminals(db *sql.DB, terminal *[]Terminals) error {
	rows, err := db.Query(allTerminalSQL)
	if err != nil {
		log.Println(err)
		return err
	}
	var number, address string
	var id int
	for rows.Next() {
		err = rows.Scan(&id, &number, &address)
		*terminal = append(*terminal, Terminals{id, number, address})
		if err != nil {
			return err
		}

	}
	return nil
}

func ShowClients(db *sql.DB, clients *[]Client) error {
	rows, err := db.Query(allClients)
	if err != nil {
		return err
	}

	for rows.Next() {
		var res Client
		err = rows.Scan(&res.Id, &res.Login, &res.Password, &res.Name, &res.Surname, &res.SerialPass, &res.Phone, &res.MainAcc)
		if err != nil {
			return err
		}
		*clients = append(*clients, res)
	}
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}

func UpdateAccounts(db *sql.DB, accounts *[]Accounts) error {
	//ddls := []string{managersDDL,managersInitialData,clientsDDL,accountsDDL,initAccountDDL,terminalsDDL,servicesDDL}
	_, err := db.Exec(accountsDDL)
	if err != nil {
		return err
	}
	for _, val := range *accounts {
		//-----------:id, :name, :number, :money, :client_id
		_, err := db.Exec(updateAccount, val.Id, val.Name, val.Number, val.Money, val.ClientId)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateTerminals(db *sql.DB, terminals *[]Terminals) error {
	//ddls := []string{managersDDL,managersInitialData,clientsDDL,accountsDDL,initAccountDDL,terminalsDDL,servicesDDL}
	_, err := db.Exec(terminalsDDL)
	if err != nil {
		return err
	}
	for _, val := range *terminals {
		//-----------(:id, :number, :address)
		_, err := db.Exec(updateTerminals, val.Id, val.Number, val.Address)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateClients(db *sql.DB, clients *[]Client) error {
	//ddls := []string{managersDDL,managersInitialData,clientsDDL,accountsDDL,initAccountDDL,terminalsDDL,servicesDDL}
	_, err := db.Exec(clientsDDL)
	if err != nil {
		return err
	}
	for _, val := range *clients {
		//-----------(:id, :login, :password, :name, :surname, :serialPass, :phone, :mainAcc)
		_, err := db.Exec(updateClients, val.Id, val.Login, val.Password, val.Name, val.Surname, val.SerialPass, val.Phone, val.MainAcc)
		if err != nil {
			return err
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------

func ShowServices(db *sql.DB, services *[]Services) error {
	rows, err := db.Query(allServices)
	if err != nil {
		return err
	}

	for rows.Next() {
		//--------------------Id,Name,Payment,AccountNumber
		var res Services
		err = rows.Scan(&res.Id, &res.Name, &res.IdPayment, &res.AccountNumber)
		if err != nil {
			return err
		}
		*services = append(*services, res)
	}
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}

func ShowAccountById(db *sql.DB, accounts *[]Accounts, mainAcc *int, id int) error {
	rows, err := db.Query(allAccountsById, id)
	if err != nil {
		return err
	}
	var mainAccount int
	err = db.QueryRow(mainAccountSqlById, id).Scan(&mainAccount)
	*mainAcc = mainAccount
	if err != nil {
		return err
	}

	for rows.Next() {
		var res Accounts
		err = rows.Scan(&res.Id, &res.Name, &res.Number, &res.Money, &res.ClientId)
		if err != nil {
			return err
		}
		*accounts = append(*accounts, res)
	}
	return nil
}

func IdClientByLogin(db *sql.DB, login string) (int, error) {
	var clientId int
	err := db.QueryRow(idClientByLogin, login).Scan(&clientId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Login was not detected!")
		}
		return 0, err
	}
	return clientId, nil
}

func IdClientByAccount(db *sql.DB, numberAcc int) (int, error) {
	var clientId int
	err := db.QueryRow(idClientByAccount, numberAcc).Scan(&clientId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Login was not detected!")
		}
		return 0, err
	}
	return clientId, nil
}

func NumberAccountByPhone(db *sql.DB, phone int) (int, error) {
	var numberAccount int
	err := db.QueryRow(numberAccountByPhone, phone).Scan(&numberAccount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Undefined phone")
		}
		return 0, errors.New("Server problem")
	}

	return numberAccount, nil
}

func CheckAccount(db *sql.DB, numberAcc int) error {
	var id int
	err := db.QueryRow(checkAccount, numberAcc).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Undefined number account")
		}
		return err
	}
	return nil
}

func AddMoneyAccountNumber(db *sql.DB, numberAcc int, amount int) error {
	err := CheckAccount(db, numberAcc)
	if err != nil {
		return err
	}
	_, err = db.Exec(updateBalanceAccount, amount, numberAcc)
	if err != nil {
		return err
	}
	return nil
}

func ChangeMainAcc(db *sql.DB, clientId, accountId int) error {
	_, err := db.Exec(updateMainAcc, accountId, clientId)
	if err != nil {
		return err
	}
	return nil
}
