package cmd

import (
	"os"
  "strings"

	"github.com/spf13/cobra"
  "github.com/retpolanne/whiterabbit/pkg"
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

func callCalc(cmd *cobra.Command, args []string) {
  err := pkg.Calculate()
  if err != nil {
    os.Exit(1)
  }
} 

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "whiterabbit",
	Short: "A simple time tracker",
	Long: `This app helps you track the time you spend working`,
}

// Commands

var goodmorningCmd = &cobra.Command{
	Use:     "goodmorning",
  Aliases: []string{"gm"},
	Short:   "tracks when you start your day",
	Run: callTrack,
}

var lunchbreakCmd = &cobra.Command{
	Use:   "lunchbreak",
	Short: "tracks when you stop for lunch",
	Run: callTrack,
}

var brbCmd = &cobra.Command{
	Use:   "brb",
	Short: "tracks when you go out for an appointment or for commuting",
	Run: callTrack,
}

var backCmd = &cobra.Command{
	Use:   "back",
	Short: "tracks when you are back from an appointment or commuting",
	Run: callTrack,
}

var goodnightCmd = &cobra.Command{
	Use:     "goodnight",
  Aliases: []string{"gn"},
	Short:   "tracks when you end your day",
	Run: callTrack,
}

var calculateCmd = &cobra.Command{
	Use:     "calculate",
  Aliases: []string{"calc"},
	Short:   "calculates time tracked",
	Run: callCalc,
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
	rootCmd.AddCommand(calculateCmd)
}


