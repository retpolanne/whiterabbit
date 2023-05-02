package pkg

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[Error] T - got the following error trying to get current working directory: %s\n", err)
	}
	fd, err := openFile(os.O_APPEND|os.O_CREATE|os.O_RDWR, cwd)
	defer fd.Close()
	if err != nil {
		log.Fatalf("[Error] T - Got the following error trying to open/create file: %s\n", err)
	}
	log.Println("T - Created/updated file $CWD/whiterabbit.csv")
	// Adding records
	record := [][]string{
		{"Mon, 01 May 2023 10:00:00 -03", "goodmorning", ""},
		{"Mon, 01 May 2023 19:00:00 -03", "goodnight", ""},
		{"Tue, 02 May 2023 10:00:00 -03", "goodmorning", ""},
		{"Tue, 02 May 2023 19:00:00 -03", "goodnight", ""},
		{"Wed, 03 May 2023 10:00:00 -03", "goodmorning", ""},
		{"Wed, 03 May 2023 14:00:00 -03", "brb", "appointment"},
		{"Wed, 03 May 2023 15:00:00 -03", "back", "appointment"},
		{"Wed, 03 May 2023 19:00:00 -03", "goodnight", ""},
		{"Thu, 04 May 2023 10:00:00 -03", "goodmorning", ""},
		{"Thu, 04 May 2023 12:00:00 -03", "lunchbreak", ""},
		{"Thu, 04 May 2023 14:00:00 -03", "lunchback", ""},
		{"Thu, 04 May 2023 19:00:00 -03", "goodnight", ""},
	}

	err = csv.NewWriter(fd).WriteAll(record)
	if err != nil {
		log.Fatalf("[Error] T - Got the following error writing csv: %s\n", err)
	}
	log.Printf("Added record:\n %s\n", strings.Join(record[0], ", "))
}

func teardown() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[Error] T - got the following error trying to get current working directory: %s\n", err)
	}
	err = os.Remove(filepath.Join(cwd, "whiterabbit.csv"))
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestCalculateTodayNoBreaks(t *testing.T) {
	// Mock today as tuesday, 2 May 2023
	datetime, err := time.Parse(time.RFC1123, "Tue, 02 May 2023 10:00:00 -03")
	if err != nil {
		log.Fatalf("[Error] T - Got the following error trying to parse the mock date: %s\n", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[Error] T - got the following error trying to get current working directory: %s\n", err)
	}
	diff, err := Calculate(true, false, false, datetime, cwd)
	if err != nil {
		log.Fatalf("[Error] T - Got the following error trying calculate day: %s\n", err)
	}
	expectedDuration, err := time.ParseDuration("8h0m0s")
	if err != nil {
		log.Fatalf("[Error] T - Got the following error parsing duration: %s\n", err)
	}
	assert.Equal(t, expectedDuration, *diff)
}

func TestCalculateYesterdayNoBreaks(t *testing.T) {
	// Mock today as tuesday, 2 May 2023
	datetime, err := time.Parse(time.RFC1123, "Tue, 02 May 2023 10:00:00 -03")
	if err != nil {
		log.Fatalf("[Error] T - Got the following error trying to parse the mock date: %s\n", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[Error] T - got the following error trying to get current working directory: %s\n", err)
	}
	diff, err := Calculate(false, true, false, datetime, cwd)
	if err != nil {
		log.Fatalf("[Error] T - Got the following error trying calculate day: %s\n", err)
	}
	expectedDuration, err := time.ParseDuration("8h0m0s")
	if err != nil {
		log.Fatalf("[Error] T - Got the following error parsing duration: %s\n", err)
	}
	assert.Equal(t, expectedDuration, *diff)
}

func TestCalculateTodayWithBreaks(t *testing.T) {
	// Mock today as wednesday, 3 May 2023
	datetime, err := time.Parse(time.RFC1123, "Wed, 03 May 2023 10:00:00 -03")
	if err != nil {
		log.Fatalf("[Error] T - Got the following error trying to parse the mock date: %s\n", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[Error] T - got the following error trying to get current working directory: %s\n", err)
	}
	diff, err := Calculate(true, false, false, datetime, cwd)
	if err != nil {
		log.Fatalf("[Error] T - Got the following error trying calculate day: %s\n", err)
	}
	expectedDuration, err := time.ParseDuration("7h0m0s")
	if err != nil {
		log.Fatalf("[Error] T - Got the following error parsing duration: %s\n", err)
	}
	assert.Equal(t, expectedDuration, *diff)
}

func TestCalculateYesterdayWithBreaks(t *testing.T) {
	// Mock today as Thursday, 4 May 2023
	datetime, err := time.Parse(time.RFC1123, "Thu, 04 May 2023 10:00:00 -03")
	if err != nil {
		log.Fatalf("[Error] T - Got the following error trying to parse the mock date: %s\n", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[Error] T - got the following error trying to get current working directory: %s\n", err)
	}
	diff, err := Calculate(false, true, false, datetime, cwd)
	if err != nil {
		log.Fatalf("[Error] T - Got the following error trying calculate day: %s\n", err)
	}
	expectedDuration, err := time.ParseDuration("7h0m0s")
	if err != nil {
		log.Fatalf("[Error] T - Got the following error parsing duration: %s\n", err)
	}
	assert.Equal(t, expectedDuration, *diff)
}
