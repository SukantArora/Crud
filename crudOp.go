package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Empl struct {
	id    int
	Name  string
	Email string
	role  string
}

func DbConn(db_name string) (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "password"
	dbName := db_name
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// func CreateTable()
func CreateTable(db_name string, table_name string) {
	db := DbConn("Employee_Db")
	defer db.Close()
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v(id int PRIMARY KEY AUTO_INCREMENT, Name varchar(30) NOT NULL, Email varchar(30), role varchar(30));", table_name)
	fmt.Println(query)
	res, err := db.Exec(query)
	if err != nil {
		fmt.Println("here")
		log.Fatal(err)
	}
	fmt.Println(res.RowsAffected())
}

func InsertData(table_name string, db *sql.DB, name string, email string, role string) error {
	query := "INSERT INTO Employee_Details (Name,Email,role) values (?,?,?);"

	res, err := db.Prepare(query)

	if err != nil {
		return err
	}
	_, err = res.Exec(name, email, role)

	return err
}

func GetById(db *sql.DB, id int) (*Empl, error) {
	query := "select * from Employee_Details where id=?"
	res, err := db.Prepare(query)
	var e Empl

	if err != nil {
		return nil, errors.New("wrong Table Name")
	}

	err = res.QueryRow(id).Scan(&e.id, &e.Name, &e.Email, &e.role)

	if err != nil {
		return nil, err
	}
	return &e, nil

}

// // Update
func UpdateById(db *sql.DB, id int, Name string, Email string, role string) error {

	stmt, err := db.Prepare("update Employee_Details set Name=?, Email=?, role=? where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(Name, Email, role, id)
	return err

}

// // Delete
func DeleteById(db *sql.DB, id int) error {

	stmt, err := db.Prepare("delete from Employee_Details where id=?")

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)

	return err
}
