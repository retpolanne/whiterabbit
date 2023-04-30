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
	filename := filepath.Join(homedir, "whiterabbit.csv")
	fd, err := os.OpenFile(filename, flags, 0644)
	return fd, err
}

func Track(command string, reason string) error {
	fmt.Printf("Received the command %s with reason %s\n", command, reason)
	fd, err := openFile(os.O_APPEND | os.O_CREATE | os.O_RDWR)
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

func Calculate() error {
	fd, err := openFile(os.O_RDONLY)
	defer fd.Close()
	if err != nil {
		fmt.Printf("Got the following error trying to open file: %s\n", err)
		return err
	}
	fmt.Println("Opened file $HOME/whiterabbit.csv")

	records, err := csv.NewReader(fd).ReadAll()
	if err != nil {
		fmt.Printf("Got the following error reading the csv file: %s\n", err)
		return err
	}

	// Iterate through records for today
	// TODO: add yesterday and workweek
	// TODO: add yesterday and workweek flags to cmd
	// TODO: calculate longer/short lunchbreaks

	todayTime := &structuredTimes{}
	for i := range records {
		datetime, err := time.Parse(time.RFC1123, records[i][0])
		if err != nil {
			fmt.Printf("Got the following error trying to parse the date of record index %d: %s\n", i, err)
			return err
		}
		// Today's date
		yt, mt, dt := time.Now().Date()
		yr, mr, dr := datetime.Date()
		if yt == yr && mt == mr && dt == dr {
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
	// TODO: calculate time brb/back and add tests
	diff := todayTime.goodnight.Sub(todayTime.goodmorning) - 1
	fmt.Printf("Today you worked %.f hours and %.f minutes – with 1h of lunchbreak\n", diff.Hours(), diff.Minutes())
	return nil
}
