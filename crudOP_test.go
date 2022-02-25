package crud

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var emp = Empl{
	id:    1,
	Name:  "sukant",
	Email: "sukant@zopsmart.com",
	role:  "sde",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestUpdateById(t *testing.T) {

	db, mock := NewMock()

	query := "update Employee_Details set Name=?, Email=?, role=? where id=?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("sukant", "sukant@zopsmart.com", "sde", 1).WillReturnResult(sqlmock.NewResult(0, 1))

	err := UpdateById(db, emp.id, emp.Name, emp.Email, emp.role)

	assert.NoError(t, err)
}

func TestUpdateByIdError(t *testing.T) {

	db, mock := NewMock()

	query := "update Employee_Detail set Name=?, Email=?, role=? where id=?"

	defer db.Close()

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("sukant", "sukant@zopsmart.com", "sde", 1).WillReturnError(errors.New("Id Not Present"))

	err := UpdateById(db, emp.id, emp.Name, emp.Email, emp.role)

	assert.Error(t, err)
}

func TestDeleteByIdError(t *testing.T) {

	db, mock := NewMock()
	query := "delete from Employee_Detail where id=?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(1).WillReturnError(errors.New("Id Not Present"))

	err := DeleteById(db, 1)

	assert.Error(t, err)

}

func TestDeleteById(t *testing.T) {

	db, mock := NewMock()
	query := "delete from Employee_Details where id=?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err := DeleteById(db, 1)

	assert.NoError(t, err)

}

func TestInsertData(t *testing.T) {
	db, mock := NewMock()

	query := "INSERT INTO Employee_Details (Name,Email,role) values (?,?,?);"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("sukant", "sukant@zopsmart.com", "sde").WillReturnResult(sqlmock.NewResult(0, 1))

	err := InsertData("Employee_Details", db, "sukant", "sukant@zopsmart.com", "sde")

	assert.NoError(t, err)

}

func TestInsertDataError(t *testing.T) {
	db, mock := NewMock()

	query := "INSERT INTO Employee_Detail (Name,Email,role) values (?,?,?);"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("sukant", "sukant@zopsmart.com", "sde").WillReturnError(errors.New("Inserting Error"))

	err := InsertData("Employee_Details", db, "sukant", "sukant@zopsmart.com", "sde")

	assert.Error(t, err)

}

func TestGetDetailsById(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Errorf(err.Error())
	}

	defer db.Close()

	query := "select * from Employee_Details where id=?"
	wquery := "select * from Employee_Detail where id=?"
	testCases := []struct {
		id          int
		emp         *Empl
		mockQuery   interface{}
		expectError error
	}{
		{
			id:          1,
			emp:         &Empl{id: 1, Name: "sukant", Email: "sukant@zopsmart.com", role: "sde"},
			mockQuery:   mock.ExpectPrepare(query).ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "Name", "Email", "role"}).AddRow(1, "sukant", "sukant@zopsmart.com", "sde")),
			expectError: nil,
		},
		{
			id:          2,
			emp:         &Empl{id: 2, Name: "Jane", Email: "j@j.com", role: "SDE-II"},
			mockQuery:   mock.ExpectPrepare(query).ExpectQuery().WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id", "Name", "Email", "role"}).AddRow(2, "Jane", "j@j.com", "SDE-II")),
			expectError: nil,
		},
		// Failure
		{
			id:          3,
			emp:         nil,
			mockQuery:   mock.ExpectPrepare(query).ExpectQuery().WithArgs(3).WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
		},
		{
			id:          4,
			emp:         nil,
			mockQuery:   mock.ExpectPrepare(wquery).WillReturnError(errors.New("wrong Table Name")),
			expectError: errors.New("wrong Table Name"),
		},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {

			emp, err := GetById(db, testCase.id)
			if err != nil && err.Error() != testCase.expectError.Error() {
				t.Errorf("expected error:%v, got:%v", testCase.expectError, err)
			} else {

				if !reflect.DeepEqual(testCase.emp, emp) {
					t.Errorf("expected users %v, got: %v", testCase.emp.Name, emp)
				}
			}

		})
	}
}
