package pkg

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type structuredTimes struct {
	goodmorning time.Time
	goodnight   time.Time
	lunchbreak  time.Time
	lunchback   time.Time
	brb         []time.Time
	back        []time.Time
}

func openFile(flags int, path string) (*os.File, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Got the following error trying to get the user home directory: %s\n", err)
		return nil, err
	}
	var filename string
	if path != "" {
		filename = filepath.Join(path, "whiterabbit.csv")
	} else {
		filename = filepath.Join(homedir, "whiterabbit.csv")
	}
	fd, err := os.OpenFile(filename, flags, 0644)
	return fd, err
}

func Track(command string, reason string) error {
	fmt.Printf("Received the command %s with reason %s\n", command, reason)
	fd, err := openFile(os.O_APPEND|os.O_CREATE|os.O_RDWR, "")
	defer fd.Close()
	if err != nil {
		fmt.Printf("Got the following error trying to open/create file: %s\n", err)
		return err
	}
	fmt.Println("Created/updated file $HOME/whiterabbit.csv")

	record := [][]string{
		{time.Now().Format(time.RFC1123), command, reason},
	}

	err = csv.NewWriter(fd).WriteAll(record)
	if err != nil {
		fmt.Printf("Got the following error writing csv: %s\n", err)
		return err
	}
	fmt.Printf("Added record:\n %s\n", strings.Join(record[0], ", "))
	return nil
}

func calculateDay(records [][]string, y int, m time.Month, d int) (diff *time.Duration, err error) {
	// TODO return total breaks
	diffRet, err := time.ParseDuration("0h0m0s")
	todayTime := &structuredTimes{}
	for i := range records {
		datetime, err := time.Parse(time.RFC1123, records[i][0])
		if err != nil {
			fmt.Printf("Got the following error trying to parse the date of record index %d: %s\n", i, err)
			return nil, err
		}
		yr, mr, dr := datetime.Date()
		if y == yr && m == mr && d == dr {
			switch records[i][1] {
			case "goodmorning":
				todayTime.goodmorning = datetime
			case "goodnight":
				todayTime.goodnight = datetime
			case "lunchbreak":
				todayTime.lunchbreak = datetime
			case "lunchback":
				todayTime.lunchback = datetime
			case "brb":
				todayTime.brb = append(todayTime.brb, datetime)
			case "back":
				todayTime.back = append(todayTime.back, datetime)
			}
		}
	}
	// TODO: parametrise lunchbreak
	if !todayTime.goodmorning.IsZero() && !todayTime.goodnight.IsZero() {
		if !todayTime.lunchbreak.IsZero() {
			diffRet = todayTime.goodnight.Sub(todayTime.goodmorning) - todayTime.lunchback.Sub(todayTime.lunchbreak)
		} else {
			diffRet = todayTime.goodnight.Sub(todayTime.goodmorning) - time.Hour
		}
	}
	var brbDiffs []time.Duration
	for i := range todayTime.brb {
		brbDiffs = append(brbDiffs, todayTime.back[i].Sub(todayTime.brb[i]))
	}
	for i := range brbDiffs {
		diffRet = diffRet - brbDiffs[i]
	}
	return &diffRet, nil
}

func Calculate(today, yesterday bool, now time.Time, filepath string) (diff *time.Duration, err error) {
	fd, err := openFile(os.O_RDONLY, filepath)
	defer fd.Close()
	if err != nil {
		fmt.Printf("Got the following error trying to open file: %s\n", err)
		return nil, err
	}
	fmt.Println("Opened file $HOME/whiterabbit.csv")

	records, err := csv.NewReader(fd).ReadAll()
	if err != nil {
		fmt.Printf("Got the following error reading the csv file: %s\n", err)
		return nil, err
	}

	// TODO: calculate longer/short lunchbreaks
	if today {
		yt, mt, dt := now.Date()
		diff, err = calculateDay(records, yt, mt, dt)
		return diff, err
	}
	if yesterday {
		yt, mt, dt := now.AddDate(0, 0, -1).Date()
		diff, err = calculateDay(records, yt, mt, dt)
		return diff, err
	}
	return nil, err
}

type Timesheet struct {
	WorkedHours []time.Duration
	Breaks      []time.Duration
}

func CalculateFirstLastDayOfWeek(now time.Time) (firstDay, lastDay time.Time) {
	firstDayOfWeekOffset := now.Weekday() - time.Sunday
	lastDayOfWeekOffset := now.Weekday() - time.Saturday
	firstDayOfWeek := now.AddDate(0, 0, -int(firstDayOfWeekOffset))
	lastDayOfWeek := now.AddDate(0, 0, -int(lastDayOfWeekOffset))
	return firstDayOfWeek, lastDayOfWeek
}

func CalculateTimesheet(now time.Time, filepath string) (*Timesheet, error) {
	var timesheet Timesheet
	fd, err := openFile(os.O_RDONLY, filepath)
	defer fd.Close()
	if err != nil {
		fmt.Printf("Got the following error trying to open file: %s\n", err)
		return nil, err
	}
	fmt.Println("Opened file $HOME/whiterabbit.csv")

	records, err := csv.NewReader(fd).ReadAll()
	if err != nil {
		fmt.Printf("Got the following error reading the csv file: %s\n", err)
		return nil, err
	}

	firstDay, lastDay := CalculateFirstLastDayOfWeek(now)
	dayIndex := 0
	for i := firstDay; i.After(lastDay) == false; i = i.AddDate(0, 0, 1) {
		yt, mt, dt := i.Date()
		diff, err := calculateDay(records, yt, mt, dt)
		fmt.Printf("Diff for day %v - %v\n", i, diff)
		if err != nil {
			fmt.Printf("Got the following error calculating day: %s\n", err)
			return nil, err
		}
		timesheet.WorkedHours = append(timesheet.WorkedHours, *diff)
		dayIndex++
	}
	return &timesheet, nil
}
