package hardware

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/MikMuellerDev/smarthome/core/database"
)

type PowerJob struct {
	Id     int64  `json:"id"`
	Switch string `json:"switch"`
	Power  bool   `json:"power"`
}

type JobResult struct {
	Id    int64 `json:"id"`
	Error error `json:"error"`
}

var log *logrus.Logger

func InitLogger(logger *logrus.Logger) {
	log = logger
}

// Returns the power state of a given switch
// Checks if the switch exists beforehand
func GetPowerState(switchId string) (bool, error) {
	switchExists, err := database.DoesSwitchExist(switchId)
	if err != nil {
		log.Error("Failed to get power state: database failure: ", err.Error())
		return false, err
	}
	if !switchExists {
		return false, fmt.Errorf("can not get power state of switch '%s': switch does not exists", switchId)
	}
	powerState, err := database.GetPowerStateOfSwitch(switchId)
	if err != nil {
		log.Error("Failed to get power state: database failure: ", err.Error())
		return false, err
	}
	return powerState, nil
}

// Sets the powerstate of a specific switch
// Checks if the switch exists
// Checks if the user has all required permissions
func SetSwitchPowerAll(switchId string, powerOn bool, username string) error {
	switchExists, err := database.DoesSwitchExist(switchId)
	if err != nil {
		return err
	}
	if !switchExists {
		return fmt.Errorf("Failed to set power: switch '%s' does not exist", switchId)
	}
	userHasPowerPermission, err := database.UserHasPermission(username, database.PermissionSetPower)
	if err != nil {
		return fmt.Errorf("Failed to set power: could not check if user is allowed to interact with switches: %s", err.Error())
	}
	if !userHasPowerPermission {
		return errors.New("Failed to set power: user is not allowed to interact with switches")
	}
	if err := SetPower(switchId, powerOn); err != nil {
		return fmt.Errorf("Failed to set power: hardware error: %s", err.Error())
	}
	return nil
}
