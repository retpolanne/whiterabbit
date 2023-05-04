# White Rabbit

![white rabbit from Alice](rabbit.png "white rabbit from Alice")

_"Oh dear! Oh dear! I shall be too late!"_

White Rabbit is a time tracker. You can use it to save information about:

- when you start your day
- if you have a lunchbreak longer than 1h
- if you have any other breaks throughout the day
- when you finish your day

You can then calculate the worked hours for today, yesterday, and generate a handy timesheet for your week.

## Limitations 

Since this is based off Brazil's working laws, it calculates a default 1h break for lunch. 

## Usage

```
This app helps you track the time you spend working

Usage:
  whiterabbit [command]

Available Commands:
  back        tracks when you are back from an appointment or commuting
  brb         tracks when you go out for an appointment or for commuting
  calculate   calculates time tracked
  completion  Generate the autocompletion script for the specified shell
  goodmorning tracks when you start your day
  goodnight   tracks when you end your day
  help        Help about any command
  lunchbreak  tracks when you stop for lunch

Flags:
  -h, --help     help for whiterabbit
  -t, --toggle   Help message for toggle

Use "whiterabbit [command] --help" for more information about a command.
```

### Track times

`whiterabbit goodmorning`

`whiterabbit brb "reason"`

`whiterabbit back "reason"`

`whiterabbit lunchbreak`

`whiterabbit lunchback`

`whiterabbit goodnight`

### Calculate times

`whiterabbit calc --today`

`whiterabbit calc --yesterday`

`whiterabbit calc --timesheet`

## Contributing

Please open issues with bugs and feature requests. 

To contribute, fork this repo, install Go, run `make test` and make sure to add tests at least to pkg. 
