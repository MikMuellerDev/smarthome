package automation

import (
	"errors"
	"fmt"
	"time"

	"github.com/nathan-osman/go-sunrise"

	"github.com/MikMuellerDev/smarthome/core/database"
)

// Utils for determining the times for sunrise and sunset
// Will be used if the automation's mode is set to either 'sunset' or 'sunrise'

type SunTime struct {
	Hour   uint
	Minute uint
}

// Returns (sunrise, sunset) based on the provided coordinates which are stored in the server configuration
func CalculateSunRiseSet(lat float32, lon float32) (SunTime, SunTime) {
	sunRise, sunSet := sunrise.SunriseSunset(
		float64(lat), float64(lon),
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
	)
	return SunTime{
			uint(sunRise.Local().Hour()), uint(sunRise.Local().Minute()),
		}, SunTime{
			uint(sunSet.Local().Hour()), uint(sunSet.Local().Minute()),
		}
}

// Given a jobId and whether sunrise or sunset should is activated, the next execution time is modified
func updateJobTime(id uint, useSunRise bool) error {
	// Obtain the server's configuration in order to determine the latitude and longitude
	config, found, err := database.GetServerConfiguration()
	if err != nil || !found {
		log.Error("Failed to update job launch time: could not obtain the server's configuration")
		return errors.New("could not update launch time: failed to obtain server config")
	}
	// Retrieve the current job in order to get its current cron-expression (for the days)
	job, found, err := database.GetAutomationById(id)
	if err != nil || !found {
		return errors.New("could not update launch time: invalid id supplied")
	}
	// Calculate both the sunrise and sunset time
	sunRise, sunSet := CalculateSunRiseSet(config.Latitude, config.Longitude)
	// Select the time which is desired
	var finalTime SunTime
	if useSunRise {
		finalTime = sunRise
	} else {
		finalTime = sunSet
	}
	// Extract the days from the cron-expression
	days, err := GetDaysFromCronExpression(job.CronExpression)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to extract days from cron-expression '%s': Error: %s", job.CronExpression, err))
		return err
	}
	cronExpression, err := GenerateCronExpression(uint8(finalTime.Hour), uint8(finalTime.Minute), days)
	if err != nil {
		return err
	}
	if err := ModifyAutomationById(id, database.AutomationWithoutIdAndUsername{
		Name:           job.Name,
		Description:    job.Description,
		CronExpression: cronExpression,
		HomescriptId:   job.HomescriptId,
		Enabled:        job.Enabled,
		TimingMode:     job.TimingMode,
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to update next execution time of automation '%d': could not modify automation: %s", id, err.Error()))
		return err
	}
	log.Trace(fmt.Sprintf("Successfully updated the next execution time of automation '%d'", id))
	return nil
}
