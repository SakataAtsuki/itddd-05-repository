package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/SakataAtsuki/itddd-05-repository/domain/model/user"
)

func main() {
	uri := fmt.Sprintf("postgres://%s/%s?sslmode=disable&user=%s&password=%s&port=%s&timezone=Asia/Tokyo",
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"))
	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("successfully connected to database")

	err = CreateUser(db, "test-user", "test-user-id")
	if err != nil {
		log.Println(err)
	}
}

func CreateUser(db *sql.DB, name string, id string) (err error) {
	defer func() {
		if err != nil {
			err = &CreateUserError{Message: fmt.Sprintf("main.CreateUser err: %v", err), Err: err}
		}
	}()
	userName, err := user.NewUserName("username")
	if err != nil {
		return err
	}
	userId, err := user.NewUserId("userid")
	if err != nil {
		return err
	}
	newUser, err := user.NewUser(*userId, *userName)
	if err != nil {
		return err
	}

	userRepository, err := user.NewUserRepository(db)
	if err != nil {
		return err
	}

	userService, err := user.NewUserService(userRepository)
	if err != nil {
		return err
	}

	isExists, err := userService.Exists(newUser)
	if err != nil {
		return err
	}

	if isExists {
		return fmt.Errorf("the user %v is already exists", newUser)
	}

	if err := userRepository.Save(newUser); err != nil {
		return err
	}
	log.Println("test-user is successfully added in users table")
	return nil
}

type CreateUserError struct {
	Message string
	Err     error
}

func (err *CreateUserError) Error() string {
	return err.Message
}
