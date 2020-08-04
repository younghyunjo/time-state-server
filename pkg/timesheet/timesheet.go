package timesheet

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type configsJson struct {
	Credentials    string `json:"credentials"`
	Token          string `json:"token"`
	Spreadsheet_id string `json:"spreadsheet_id"`
	Sheet_name     string `json:"sheet_name"`
}

type sleepJson struct {
	Date     string `json:"date"`
	WakeTime string `json:"wakeTime"`
	BedTime  string `json:"bedTime"`
}

type Sleep struct {
	Date     time.Time
	WakeTime time.Time
	BedTime  time.Time
}

/*
type Work struct {
	date    time.Time
	spent   time.Time
	working time.Time
}
*/

var configs configsJson
var sleepSheet map[time.Time]Sleep
var sheetLock sync.RWMutex
var UpdateTicker *time.Ticker
var updaterChannel chan bool

func init() {
	configsFile, err := os.Open("configs/timesheets.json")
	if err != nil {
		log.Fatalln("config file open failed")
		return
	}
	defer configsFile.Close()

	bytes, err := ioutil.ReadAll(configsFile)
	if err != nil {
		log.Fatalln("ioutil.ReadAll error")
		return
	}
	err = json.Unmarshal(bytes, &configs)
	if err != nil {
		log.Fatalln("unmarshal failed")
	}

	downloadSheet()
	startSheetUpdater()
}

func startSheetUpdater() {
	UpdateTicker = time.NewTicker(1 * time.Hour)
	updaterChannel = make(chan bool)

	go func() {
		for {
			select {
			case <-updaterChannel:
				return
			case tm := <-UpdateTicker.C:
				fmt.Println("downloadSheet", tm)
				downloadSheet()
			}
		}
	}()
}

func stopSheetUpdater() {
	UpdateTicker.Stop()
	updaterChannel <- true
}

func downloadSheet() {
	client, err := getClient()
	if err != nil {
		return
	}

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(configs.Spreadsheet_id, configs.Sheet_name).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return
	}

	sheetLock.Lock()
	sleepSheet = make(map[time.Time]Sleep)
	for i, row := range resp.Values {
		if i == 0 {
			continue
		}
		sleep := Sleep{}
		sleep.Date, _ = time.Parse("2006. 01. 02.", row[0].(string))
		sleep.WakeTime, _ = time.Parse("15:04", row[1].(string))
		if len(row) > 2 {
			sleep.BedTime, _ = time.Parse("15:04", row[2].(string))
		}
		sleepSheet[sleep.Date] = sleep
	}
	sheetLock.Unlock()
}

func GetWakeTime(date time.Time) time.Time {
	sheetLock.RLock()
	wakeTime := sleepSheet[date].WakeTime
	sheetLock.RUnlock()
	return wakeTime
}

func GetSleepTime(date time.Time) Sleep {
	sheetLock.RLock()
	sleepTime := sleepSheet[date]
	sheetLock.RUnlock()
	return sleepTime
}

func GetSleepTimes(date []time.Time) []Sleep {
	var sleepTimes []Sleep

	sheetLock.RLock()
	for _, date := range date {
		if s, ok := sleepSheet[date]; ok {
			sleepTimes = append(sleepTimes, s)
		} else {
			sleepTimes = append(sleepTimes, Sleep{Date: date})
		}
	}
	sheetLock.RUnlock()

	return sleepTimes
}

func SleepToJson(sleepTimes []Sleep, dateLayout string, timeLayout string) interface{} {
	var sleepTimeJson []sleepJson
	for _, s := range sleepTimes {
		sleepTimeJson = append(sleepTimeJson, sleepJson{
			Date:     s.Date.Format(dateLayout),
			WakeTime: s.WakeTime.Format(timeLayout),
			BedTime:  s.BedTime.Format(timeLayout),
		})
	}

	return sleepTimeJson
}
