package hardware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MikMuellerDev/smarthome/core/database"
	"github.com/MikMuellerDev/smarthome/core/event"
)

type HardwareRequest struct {
	Switch string `json:"switch"`
	Power  bool   `json:"power"`
}

// Checks if a node is online and updates the database entry accordingly
func checkNodeOnlineRequest(node database.HardwareNode) error {
	// Client has timeout of a second too
	client := http.Client{Timeout: time.Second}
	res, err := client.Get(fmt.Sprintf("%s/health", node.Url))
	if err != nil {
		log.Error("Hardware node checking request failed: ", err.Error())
		return err
	}
	if res.StatusCode != 200 {
		log.Error("Hardware node checking request failed: non 200 status code")
		return errors.New("checking node failed: non 200 status code")
	}
	return nil
}

// Runs the check request and updated the database entry accordingly
func checkNodeOnline(node database.HardwareNode) error {
	if err := checkNodeOnlineRequest(node); err != nil {
		if node.Online {
			log.Warn(fmt.Sprintf("Node `%s` failed to respond and is now offline", node.Name))
			go event.Error("Node Offline",
				fmt.Sprintf("Node %s went offline. Users will have to deal with increased wait times. It is advised to address this issue as soon as possible", node.Name))
		}
		if errDB := database.SetNodeOnline(node.Url, false); errDB != nil {
			log.Error("Failed to update power state of node: ", errDB.Error())
			return errDB
		}
		return nil
	}
	if !node.Online {
		log.Info(fmt.Sprintf("Node `%s` is now back online", node.Name))
		go event.Info("Node Online", fmt.Sprintf("Node %s is back online.", node.Name))
	}
	if errDB := database.SetNodeOnline(node.Url, true); errDB != nil {
		log.Error("Failed to update power state of node: ", errDB.Error())
		return errDB
	}
	return nil
}

// Delivers a power job to a given hardware node
// Returns an error if the job fails to execute on the hardware
// However, the preferred method of communication is by using the API `SetPower()` this way, priorities and interrupts are scheduled automatically
// TODO: add field in the node table for marking a node as unavailable (unavailable nodes will be excluded from the normal power request)
// A check if  a node is online again can be still executed afterwards
func sendPowerRequest(node database.HardwareNode, switchName string, powerOn bool) error {
	requestBody, err := json.Marshal(HardwareRequest{
		Switch: switchName,
		Power:  powerOn,
	})
	if err != nil {
		log.Error("Could not parse node request: ", err.Error())
		return err
	}
	// Create a client with a more realistic timeout of 1 second
	client := http.Client{Timeout: time.Second}
	res, err := client.Post(fmt.Sprintf("%s/power?token=%s", node.Url, node.Token), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Error("Hardware node request failed: ", err.Error())
		return err
	}
	if res.StatusCode != 200 {
		// TODO: check firmware version of the power nodes, analyze errors
		log.Error(fmt.Sprintf("Received non 200 status code: %d", res.StatusCode))
		return errors.New("set power failed: non 200 status code")
	}
	defer res.Body.Close()
	return nil
}

// More user-friendly API to directly address all hardware nodes
// However, the preferred method of communication is by using the API `ExecuteJob()` this way, priorities and interrupts are scheduled automatically
// This method is internally used by `ExecuteJob`
// Makes a database request at the beginning in order to obtain information about the available nodes
// Updates the power state in the database after the jobs have been sent to the hardware nodes
func setPowerOnAllNodes(switchName string, powerOn bool) error {
	var err error = nil
	// Retrieves available hardware nodes from the database
	nodes, err := database.GetHardwareNodes()
	if err != nil {
		log.Error("Failed to process power request: could not get nodes from database: ", err.Error())
		return err
	}
	for _, node := range nodes {
		errTemp := sendPowerRequest(node, switchName, powerOn)
		if errTemp != nil {
			// If the request failed, check the node and mark it as offline
			go checkNodeOnline(node)
			err = errTemp
		} else {
			if !node.Online {
				// If the node was previously offline and is now online
				go checkNodeOnline(node)
			}
			log.Debug("Successfully sent power request to: ", node.Name)
		}
	}
	if _, err := database.SetPowerState(switchName, powerOn); err != nil {
		log.Error("Failed to set power after addressing all nodes: updating database entry failed: ", err.Error())
		return err
	}
	return err
}
