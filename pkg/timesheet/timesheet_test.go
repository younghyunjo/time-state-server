package timesheet

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Given_Date_When_GetSleep_Then_GetRightTime(t *testing.T) {
	// given
	date := time.Date(2020, 5, 10, 0, 0, 0, 0, time.UTC)

	//when
	sleepTime := GetWakeTime(date)
	expected, _ := time.Parse("3:04", "5:00")

	assert.Equal(t, expected, sleepTime)
}

func Test_Given_Date_When_GetSleepTime_Then_GetRightTime(t *testing.T) {
	// given
	date := time.Date(2020, 5, 10, 0, 0, 0, 0, time.UTC)

	sleepTime := GetSleepTime(date)
	expectedWakeTime, _ := time.Parse("3:04", "5:00")
	expectedBedTime, _ := time.Parse("15:04", "21:30")

	assert.Equal(t, date, sleepTime.Date)
	assert.Equal(t, expectedWakeTime, sleepTime.WakeTime)
	assert.Equal(t, expectedBedTime, sleepTime.BedTime)
}

func Test_Given_Dates_When_GetSleepTimes_Then_GetRightTime(t *testing.T) {
	//given
	var dates []time.Time
	date0 := time.Date(2020, 5, 10, 0, 0, 0, 0, time.UTC)
	dates = append(dates, date0)
	date1 := time.Date(2020, 5, 11, 0, 0, 0, 0, time.UTC)
	dates = append(dates, date1)

	sleepTimes := GetSleepTimes(dates)

	assert.Equal(t, date0, sleepTimes[0].Date)
	assert.Equal(t, date1, sleepTimes[1].Date)
}

func Test_Json(t *testing.T) {
	SleepToJson()
}
