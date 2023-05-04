package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/retpolanne/whiterabbit/pkg"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func callTrack(cmd *cobra.Command, args []string) {
	reason := ""
	if len(args) > 0 {
		reason = strings.Join(args, " ")
	}
	err := pkg.Track(cmd.Use, reason)
	if err != nil {
		os.Exit(1)
	}
}

func printTimesheet(timesheet *pkg.Timesheet) {
	tbl := table.New("Weekdays", "S", "M", "T", "W", "T", "F", "S")
	tbl.AddRow(
		"Worked hours",
		timesheet.WorkedHours[0],
		timesheet.WorkedHours[1],
		timesheet.WorkedHours[2],
		timesheet.WorkedHours[3],
		timesheet.WorkedHours[4],
		timesheet.WorkedHours[5],
		timesheet.WorkedHours[6],
	)
	tbl.AddRow(
		"Breaks",
		timesheet.Breaks[0],
		timesheet.Breaks[1],
		timesheet.Breaks[2],
		timesheet.Breaks[3],
		timesheet.Breaks[4],
		timesheet.Breaks[5],
		timesheet.Breaks[6],
	)
	tbl.Print()
}

func callCalc(cmd *cobra.Command, args []string) {
	today, _ := cmd.Flags().GetBool("today")
	yesterday, _ := cmd.Flags().GetBool("yesterday")
	timesheet, _ := cmd.Flags().GetBool("timesheet")
	if today {
		diff, err := pkg.Calculate(true, false, time.Now(), "")
		if err != nil {
			log.Fatalf("Got the following error trying to calculate worked hours: %v", err)
		}
		fmt.Printf("Today you have worked %v – with 1 hour of lunchbreak\n", *diff)
	}
	if yesterday {
		diff, err := pkg.Calculate(false, true, time.Now(), "")
		if err != nil {
			log.Fatalf("Got the following error trying to calculate worked hours: %v", err)
		}
		fmt.Printf("Yesterday you have worked %v – with 1 hour of lunchbreak\n", *diff)
	}
	if timesheet {
		timesheet, err := pkg.CalculateTimesheet(time.Now(), "")
		if err != nil {
			log.Fatalf("Got the following error trying to calculate worked hours: %v", err)
		}
		firstDay, lastDay := pkg.CalculateFirstLastDayOfWeek(time.Now())
		fmt.Printf("week start: %v\n", firstDay)
		fmt.Printf("week end: %v\n", lastDay)
		printTimesheet(timesheet)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "whiterabbit",
	Short: "A simple time tracker",
	Long:  `This app helps you track the time you spend working`,
}

// Commands

var goodmorningCmd = &cobra.Command{
	Use:     "goodmorning",
	Aliases: []string{"gm"},
	Short:   "tracks when you start your day",
	Run:     callTrack,
}

var lunchbreakCmd = &cobra.Command{
	Use:   "lunchbreak",
	Short: "tracks when you stop for lunch",
	Run:   callTrack,
}

var brbCmd = &cobra.Command{
	Use:   "brb",
	Short: "tracks when you go out for an appointment or for commuting",
	Run:   callTrack,
}

var backCmd = &cobra.Command{
	Use:   "back",
	Short: "tracks when you are back from an appointment or commuting",
	Run:   callTrack,
}

var goodnightCmd = &cobra.Command{
	Use:     "goodnight",
	Aliases: []string{"gn"},
	Short:   "tracks when you end your day",
	Run:     callTrack,
}

var calculateCmd = &cobra.Command{
	Use:     "calculate",
	Aliases: []string{"calc"},
	Short:   "calculates time tracked",
	Run:     callCalc,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.whiterabbit.yaml)")

	// Cobra also supports local flags, which will only Run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(goodmorningCmd)
	rootCmd.AddCommand(lunchbreakCmd)
	rootCmd.AddCommand(brbCmd)
	rootCmd.AddCommand(backCmd)
	rootCmd.AddCommand(goodnightCmd)
	calculateCmd.Flags().BoolP("today", "t", false, "calculate worked hours today")
	calculateCmd.Flags().BoolP("yesterday", "y", false, "calculate worked hours yesterday")
	calculateCmd.Flags().BoolP("timesheet", "s", false, "calculate worked hours throughout the week in a timesheet form")
	rootCmd.AddCommand(calculateCmd)
}
