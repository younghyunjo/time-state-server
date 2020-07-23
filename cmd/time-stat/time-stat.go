package main

import (
	"fmt"
	"net/http"
	"server/pkg/timesheet"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// get : http://localhost:8080/v1.0/sleep?date=2020-06-18
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

// post {date: '2020-06-15', time:'04:30'}
func sleepPost(c *gin.Context) {
	fmt.Println("sleepPost")
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	v1 := r.Group("/v1.0")
	{
		v1.GET("/sleep", sleepGet)
		v1.POST("/sleep", sleepPost)
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
