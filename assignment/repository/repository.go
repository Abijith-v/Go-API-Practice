package repository

import (
	"assignment/misc"
	_ "assignment/misc"
	"database/sql"
	_ "database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"log"
)

type accessLevel string
type gender string

const HOST = "root:@(localhost:3306)/assignment?charset=utf8mb4&parseTime=True"
const (
	ADMIN  accessLevel = "admin"
	USER   accessLevel = "user"
	MALE   gender      = "male"
	FEMALE gender      = "female"
)

func (al *accessLevel) Scan(value interface{}) error {
	*al = accessLevel(value.([]byte))
	return nil
}

func (al accessLevel) Value() (driver.Value, error) {
	return string(al), nil
}

func (g *gender) Scan(value interface{}) error {
	*g = gender(value.([]byte))
	return nil
}

func (g gender) Value() (driver.Value, error) {
	return string(g), nil
}

func (UsersModel) TableName() string {
	return "users"
}

func (DetailsModel) TableName() string {
	return "details"
}

type UsersModel struct {
	EntityId    int         `gorm:"primary_key;auto_increment"`
	Username    string      `gorm:"size:55"`
	Password    string      `gorm:"size:55"`
	AccessLevel accessLevel `gorm:"type:enum('ADMIN', 'USER');column:access_level"`
}

type DetailsModel struct {
	UsersModel   UsersModel `gorm:"foreignkey:UserId;references:EntityId"`
	UserId       int
	FirstName    string `gorm:"size:55"`
	LastName     string `gorm:"size:55"`
	Age          int
	Gender       gender `gorm:"type:enum('MALE', 'FEMALE');column:gender"`
	Address      string `gorm:"size:256"`
	Email        string `gorm:"size:100"`
	MobileNumber string `gorm:"size:100"`
}

func Create() {
	db, err := gorm.Open(mysql.Open(HOST))
	if err != nil {
		fmt.Println("Failed to open mysql")
	}
	fmt.Println("Connection ready")
	dbErr1 := db.AutoMigrate(&UsersModel{})
	dbErr2 := db.AutoMigrate(&DetailsModel{})
	if dbErr1 != nil || dbErr2 != nil {
		fmt.Println("Error while creating tables")
	}
}

func GetUserFromDB(username string) (UsersModel, error) {
	query := "SELECT * FROM users WHERE username = ?"

	db, err := sql.Open("mysql", HOST)
	defer db.Close()
	if err != nil {
		log.Fatal("Failed to open DB")
	}
	row := db.QueryRow(query, username)

	var user UsersModel
	err = row.Scan(&user.EntityId, &user.Username, &user.Password, &user.AccessLevel)
	if err != nil {
		return UsersModel{}, err
	}

	return user, nil
}

func GetUserDetailsFromDB(userIdF float64) (DetailsModel, error) {
	query := "SELECT first_name, last_name, age, gender, address, email, mobile_number FROM details WHERE user_id = ?"
	userId := int(userIdF)
	db, err := sql.Open("mysql", HOST)
	defer db.Close()
	if err != nil {
		log.Fatal("Failed to open DB")
	}
	row := db.QueryRow(query, userId)

	var details DetailsModel
	err = row.Scan(
		&details.FirstName,
		&details.LastName,
		&details.Age,
		&details.Gender,
		&details.Address,
		&details.Email,
		&details.MobileNumber,
	)

	if err != nil {
		fmt.Println("error here")
		return DetailsModel{}, err
	}

	return details, nil
}

func AddPerson(userRequest misc.UserDetailsRequest) error {

	db, err := sql.Open("mysql", HOST)
	defer db.Close()

	if err != nil {
		return err
	}

	// Start transaction
	tnx, err := db.Begin()
	if err != nil {
		return err
	}

	userQuery := "INSERT INTO users (username, password, access_level) VALUES (?, ?, ?)"
	randomPass, err := misc.GenerateString()
	if err != nil {
		log.Fatal("Error while generating password")
		tnx.Rollback()
		return err
	}
	res, err := tnx.Exec(userQuery, userRequest.UserName, randomPass, accessLevel(userRequest.AccessLevel))

	if err != nil {
		return err
	}

	entityId, err := res.LastInsertId()
	// Rollback tnx in case of failure
	if err != nil {
		tnx.Rollback()
		return err
	}

	detailsQuery := "INSERT INTO details " +
		"(user_id, first_name, last_name, age, gender, address, email, mobile_number)" +
		" VALUES " +
		"(?, ?, ?, ?, ?, ?, ?, ?)"

	address, err := json.Marshal(userRequest.Address)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Exec query
	_, err = tnx.Exec(detailsQuery,
		entityId,
		userRequest.FirstName,
		userRequest.LastName,
		userRequest.Age,
		userRequest.Gender,
		string(address),
		userRequest.Email,
		userRequest.MobileNumber)

	if err != nil {
		tnx.Rollback()
		return err
	}
	// commut changes
	err = tnx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("New entry created")
	return nil
}
