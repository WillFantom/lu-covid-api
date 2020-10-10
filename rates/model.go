package rates

import "time"

type Rate struct {
	Date           time.Time `mapstructure:"date"`
	CampusStudents uint64    `mapstructure:"campus"`
	CityStudents   uint64    `mapstructure:"city"`
	Staff          uint64    `mapstructure:"staff"`
}

type ResponseContent struct {
	Title    string `json:"title"`
	Abstract string `json:"abstract`
	Content  string `json:"main"`
}

type WebResponse struct {
	Key     string            `json:"key"`
	Title   string            `json:"title"`
	Content []ResponseContent `json:"contentItems"`
}
