package main

type Config struct {
	DiscordToken             string `split_words:"true" required:"true"`
	DiscordChannel           string `split_words:"true" required:"true"`
	AlmanaxDailyReportHour   string `split_words:"true" default:"0"`
	AlmanaxDailyReportMinute string `split_words:"true" default:"1"`
}
