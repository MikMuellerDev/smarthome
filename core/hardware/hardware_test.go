package hardware

import (
	"testing"

	"github.com/MikMuellerDev/smarthome/core/database"
	"github.com/MikMuellerDev/smarthome/core/event"
	"github.com/sirupsen/logrus"
)

func initDB(args ...bool) error {
	database.InitLogger(logrus.New())
	if err := database.Init(database.DatabaseConfig{
		Username: "smarthome",
		Password: "testing",
		Hostname: "localhost",
		Database: "smarthome",
		Port:     3330,
	}, "admin",
	); err != nil {
		return err
	}
	if len(args) > 0 {
		if err := database.DeleteTables(); err != nil {
			return err
		}
		initDB()
	}
	return nil
}

func TestPower(t *testing.T) {
	InitLogger(logrus.New())
	if err := initDB(true); err != nil {
		t.Error(err.Error())
	}
	table := []struct {
		Switch string
		Power  bool
	}{
		{"1", true},
		{"1", false},
		{"2", true},
		{"2", false},
		{"3", true},
		{"3", false},
		{"4", true},
		{"4", false},
		{"5", true},
		{"5", false},
		{"6", true},
		{"6", false},
	}
	if err := database.CreateRoom("testing", "testing", "testing"); err != nil {
		t.Error(err.Error())
		return
	}
	for _, item := range table {
		if err := database.CreateSwitch(item.Switch, item.Switch, "testing", 0); err != nil {
			t.Error(err.Error())
			return
		}
		if err := setPowerOnAllNodes(item.Switch, item.Power); err != nil {
			t.Error(err.Error())
			return
		}
		power, err := GetPowerState(item.Switch)
		if err != nil {
			t.Error(err.Error())
		}
		if power != item.Power {
			t.Errorf("Failed to set power: got: %t, want: %t", power, item.Power)
			return
		}
	}
}

// Requests
func TestCheckNodeOnline(t *testing.T) {
	InitLogger(logrus.New())
	if err := initDB(true); err != nil {
		t.Error(err.Error())
	}
	if err := checkNodeOnline(database.HardwareNode{
		Name:    "test",
		Online:  true,
		Enabled: true,
		Url:     "https://example.com",
		Token:   "",
	}); err != nil {
		t.Error("Node check failed: ", err.Error())
	}
}

func TestSendPowerRequest(t *testing.T) {
	log := logrus.New()
	InitLogger(log)
	event.InitLogger(log)
	if err := initDB(true); err != nil {
		t.Error(err.Error())
	}
	table := []struct {
		Node  database.HardwareNode
		Error string
	}{
		{
			Node: database.HardwareNode{
				Name:    "test1",
				Online:  true,
				Enabled: true,
				Url:     "http://localhost:1",
				Token:   "",
			},
			Error: `Post "http://localhost:1/power?token=": dial tcp [::1]:1: connect: connection refused`,
		},
		{
			Node: database.HardwareNode{
				Name:    "test2",
				Online:  false,
				Enabled: true,
				Url:     "http://localhost:2",
				Token:   "",
			},
			Error: `Post "http://localhost:2/power?token=": dial tcp [::1]:2: connect: connection refused`,
		},
		{
			Node: database.HardwareNode{
				Name:    "test3",
				Online:  true,
				Enabled: false,
				Url:     "http://localhost:3",
				Token:   "",
			},
			Error: ``,
		},
	}
	for _, item := range table {
		if got := sendPowerRequest(item.Node, "", false); got != nil {
			if item.Error == "" {
				t.Errorf("Node: %s Error is not expected: want: '', got %s", item.Node.Name, got.Error())
				continue
			}
			if item.Error != got.Error() {
				t.Errorf("Node: %s Error is not expected: want: %s, got %s", item.Node.Name, item.Error, got.Error())
				continue
			}
		} else {
			if item.Error != "" {
				t.Errorf("Node: %s Expected error which did not occur: %s", item.Node.Name, item.Error)
			}
		}
	}
}
