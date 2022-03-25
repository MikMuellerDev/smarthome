package user

import (
	"github.com/MikMuellerDev/smarthome/core/database"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var log *logrus.Logger

func InitLogger(logger *logrus.Logger) {
	log = logger
}

// Will return <true / false> based on authentication validity
// <true> means valid authentication
// Can return an error if the database fails to return a valid result, meaning service downtime
func ValidateCredentials(username string, password string) (bool, error) {
	userExists, err := database.DoesUserExist(username)
	if err != nil {
		log.Error("Failed to validate password: could not check if user exists: ", err.Error())
		return false, err
	}
	if !userExists {
		log.Trace("Credentials invalid: user does not exist")
		return false, nil
	}
	hash, err := database.GetUserPasswordHash(username)
	if err != nil {
		log.Error("Failed to validate password: database failure")
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true, nil
	}
	if err.Error() != "crypto/bcrypt: hashedPassword is not the hash of the given password" {
		log.Error("failed to check password: ", err.Error())
		return false, err
	}
	log.Trace("password check using bcrypt failed: passwords do not match")
	return false, nil
}

// Removes a user, also removes everything that depends on the user (permissions, switchPermissions)
func DeleteUser(username string) error {
	if err := RemoveAvatar(username); err != nil {
		log.Error("Failed to delete user: removing avatar failed: ", err.Error())
		return err
	}
	if err := database.DeleteUser(username); err != nil {
		log.Error("Failed to delete user: fatabase error: ", err.Error())
		return err
	}
	return nil
}
