package core

//---------------Manager
const managersDDL = `
CREATE TABLE IF NOT EXISTS managers
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL
);`
const managersInitialData = `INSERT INTO managers
VALUES (1, 'Akram', '900658008', '8008'),
       (2, 'Azam', '985658008', '8008')
       ON CONFLICT DO NOTHING;`

//------------Client-------------------------------

const clientsDDL = `
CREATE TABLE IF NOT EXISTS clients
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    login   TEXT   NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name    TEXT    NOT NULL,
    surname    TEXT    NOT NULL,
	serialPass TEXT    NOT NULL UNIQUE,
	phone INTEGER NOT NULL UNIQUE,
	mainAcc INTEGER UNIQUE
);`

const loginManSQL = `SELECT login, password FROM managers WHERE login = ? ;`
const loginCliSQL = `SELECT login, password FROM clients WHERE login = ? ;`
const idClientByLogin = `SELECT id FROM clients WHERE login = ? ;`
const idClientByAccount = `SELECT client_id FROM accounts WHERE number = ? ;`
const idClientByPhone = `SELECT id FROM clients WHERE phone = ? ;`
const numberAccountByPhone = `SELECT mainAcc FROM clients WHERE phone = ? ;`

var loginSQL string

const insertClientSQL = `INSERT INTO clients(login, password, name, surname, serialPass, phone, mainAcc)  VALUES (?, ?, ?, ?, ?, ?, ?);`
const insertAccountSQL = `INSERT INTO accounts (name, number, money, client_id)  VALUES (?, ?, ?, ?);`
const lastClientId = `SELECT id FROM clients WHERE login=:log;`

//---------------Account------------
const accountsDDL = `
CREATE TABLE IF NOT EXISTS accounts
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
	name   TEXT NOT NULL,
    number INTEGER NOT NULL,
    money INTEGER NOT NULL,
    client_id      TEXT NOT NULL

);`
const initAccountDDL = `INSERT INTO accounts
VALUES (1 , 'Initial', 100000000,0,0)
       ON CONFLICT DO NOTHING;`
const lastAccountId = `SELECT number FROM accounts ORDER BY id DESC;`
const allAccountsById = `SELECT * FROM accounts  WHERE client_id=? ORDER BY money DESC;`
const allAccounts = `SELECT * FROM accounts;`
const mainAccountSqlById = `SELECT mainAcc FROM clients  WHERE id=?;`
const updateBalanceAccount = `UPDATE accounts SET money = money + :money WHERE number=:numberAcc ;`
const checkAccount = `SELECT id FROM accounts WHERE number= :numberAcc ;`
const updateAccount = `INSERT INTO accounts VALUES (:id, :name, :number, :money, :client_id) ON CONFLICT DO NOTHING;`
const updateTerminals = `INSERT INTO terminals VALUES (:id, :number, :address) ON CONFLICT DO NOTHING;`
const updateClients = `INSERT INTO clients VALUES (:id, :login, :password, :name, :surname, :serialPass, :phone, :mainAcc) ON CONFLICT DO NOTHING;`

//-------------------------Bankomat-----------------
const terminalsDDL = `CREATE TABLE IF NOT EXISTS terminals
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
	number   INTEGER NOT NULL UNIQUE,
    address TEXT NOT NULL
);`
const exitsTerminalSQL = `SELECT id FROM terminals WHERE number=?;`
const allTerminalSQL = `SELECT * FROM terminals;`
const insertTerminal = `INSERT INTO terminals(number, address)  VALUES (?, ?);`
const updateMainAcc = `UPDATE clients SET mainAcc = :AccId WHERE id=:clID ;`

//---------------------------Uslugi

const servicesDDL = `CREATE TABLE IF NOT EXISTS services
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
	name   TEXT NOT NULL UNIQUE,
    idPayment TEXT NOT NULL,
    accountNumber INTEGER NOT NULL


);`
const exitsServicesSQL = `SELECT id FROM services WHERE name=?;`
const insertService = `INSERT INTO services (name, idPayment, accountNumber)  VALUES (?, ?, ?);`
const allServices = `SELECT * FROM services;`
const allClients = `SELECT * FROM clients;`
