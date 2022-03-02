package database

import "errors"

// Used during <Init> of the database, only called once
// Creates the table containing <users> if it doesn't already exist
// Can return an error if the database fails
func createUserTable() error {
	query := `
	CREATE TABLE
	IF NOT EXISTS
	user(
		Username VARCHAR(20) PRIMARY KEY,
		Password text
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Error("Failed to create user table: Executing query failed: ", err.Error())
		return err
	}
	return nil
}

// Lists users which are currently in the Database
// Returns an empty list with an error when failing
func ListUsers() ([]User, error) {
	query := `SELECT Username, Password FROM user`
	res, err := db.Query(query)
	if err != nil {
		log.Error("Could not list users. Failed to execute query: ", err.Error())
		return []User{}, err
	}
	var userList []User
	for res.Next() {
		var user User
		err := res.Scan(&user.Username, &user.Password)
		if err != nil {
			log.Error("Failed tp scan user values from database results: ", err.Error())
		}
		userList = append(userList, user)
	}
	return userList, nil
}

// Creates a new user based on a the supplied `User` struct
// Won't panic if user already exists, but will change password
func InsertUser(user User) error {
	query, err := db.Prepare("INSERT INTO user(Username, Password) VALUES(?,?) ON DUPLICATE KEY UPDATE Password=VALUES(Password)")
	if err != nil {
		log.Error("Could not create user. Failed to prepare query: ", err.Error())
		return err
	}
	_, err = query.Exec(user.Username, user.Password)
	if err != nil {
		log.Error("Could not create user. Failed to execute query: ", err.Error())
		return err
	}
	return nil
}

// Deletes a User based on a given Username, can return an error if the database fails
// The function does not validate the existence of this username itself, so additional checks should be done beforehand
func DeleteUser(Username string) error {
	query, err := db.Prepare(`
	DELETE FROM user WHERE Username=? 
	`)
	if err != nil {
		log.Error("Could not delete user. Failed to prepare query: ", err.Error())
		return err
	}
	_, err = query.Exec(Username)
	if err != nil {
		log.Error("Could not delete user. Failed to execute query: ", err.Error())
		return err
	}
	return nil
}

// Helper function to create a User which is given a set of basic permissions
// Will return an error if the database fails
// Does not check for duplicate users
func AddUser(user User) error {
	userExists, err := DoesUserExist(user.Username)
	if err != nil {
		return err
	}
	if userExists {
		return errors.New("could not add user: user already exists")
	}
	err = InsertUser(user)
	if err != nil {
		return err
	}
	err = AddUserPermission(user.Username, "authentication")
	if err != nil {
		return err
	}
	return nil
}

// Returns <true> if a provided user exists
// If the database fails, it returns an error
func DoesUserExist(username string) (bool, error) {
	userList, err := ListUsers()
	if err != nil {
		return false, err
	}
	for _, userItem := range userList {
		if userItem.Username == username {
			return true, nil
		}
	}
	return false, nil
}
