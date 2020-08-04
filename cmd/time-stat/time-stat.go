package main

import (
	"fmt"
	"net/http"
	"server/pkg/timesheet"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//GET /sleep?date=2020-06-18
func sleepGet(c *gin.Context) {
	fmt.Println("sleepGet")
	dateQuery := c.Query("date")
	date, _ := time.Parse("2006-01-02", dateQuery)
	sleep := timesheet.GetSleepTime(date)
	wakeTime := sleep.WakeTime.Format("03:04")
	fmt.Println(wakeTime)

	c.JSON(http.StatusOK, gin.H{
		"time": sleep.WakeTime.Format("03:04"),
	})
}

//GET /sleeptime?date=2020-08-05?last=8
func getSleepTime(c *gin.Context) {
	fmt.Println("getSleepTimes")

	lastQuery := c.DefaultQuery("last", "1")
	lastQueryInt, err := strconv.Atoi(lastQuery)
	if err != nil {
		lastQueryInt = 1
	}

	dateQuery := c.Query("date")
	date, _ := time.Parse("2006-01-02", dateQuery)
	var dates []time.Time

	for i := 0; i < lastQueryInt; i++ {
		dates = append(dates, date.AddDate(0, 0, -i))
	}

	sleepTimes := timesheet.GetSleepTimes(dates)
	sleepTimesJson := timesheet.SleepToJson(sleepTimes, "2006-01-02", "15:04")
	c.JSON(http.StatusOK, sleepTimesJson)
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	v1 := r.Group("/v1.0")
	{
		v1.GET("/sleep", sleepGet)
		v1.GET("/sleeptime", getSleepTime)
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
