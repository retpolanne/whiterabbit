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
			case "brb":
				todayTime.brb = append(todayTime.brb, datetime)
			case "back":
				todayTime.back = append(todayTime.brb, datetime)
			}
		}
	}
	// TODO: parametrise lunchbreak
	diffRet := todayTime.goodnight.Sub(todayTime.goodmorning) - time.Hour
	return &diffRet, nil
}

func Calculate(today, yesterday, weekdays bool, now time.Time, filepath string) (diff *time.Duration, err error) {
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

	// Iterate through records for today
	// TODO: add yesterday and workweek
	// TODO: add yesterday and workweek flags to cmd
	// TODO: calculate longer/short lunchbreaks
	if today {
		yt, mt, dt := now.Date()
		diff, err = calculateDay(records, yt, mt, dt)
		return diff, nil
	}
	return nil, err
}
