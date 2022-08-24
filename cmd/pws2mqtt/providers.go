package main

import (
	"encoding/json"
)

type Pwsdata interface {
	ToJSON() string
}

type Ecowittdata struct {
	Passkey        string `mapstructure:"PASSKEY"`
	Stationtype    string `mapstructure:"stationtype"`
	Dateutc        string `mapstructure:"dateutc"`
	Tempinf        string `mapstructure:"tempinf"`    //Inside Temp
	Humidityin     string `mapstructure:"humidityin"` // Inside Humidity
	Baromrelin     string `mapstructure:"baromrelin"` // Barometric relative
	Baromabsin     string `mapstructure:"baromabsin"` // Barometric Absolute
	Tempf          string `mapstructure:"tempf"`      //outside temp
	Humidity       string `mapstructure:"humidity"`
	Winddir        string `mapstructure:"winddir"`
	Windspeedmph   string `mapstructure:"windspeedmph"`
	Windgustmph    string `mapstructure:"windgustmph"`
	Maxdailygust   string `mapstructure:"maxdailygust"`
	Rainratein     string `mapstructure:"rainratein"`
	Eventrainin    string `mapstructure:"eventrainin"`
	Hourlyrainin   string `mapstructure:"hourlyrainin"`
	Dailyrainin    string `mapstructure:"dailyrainin"`
	Weeklyrainin   string `mapstructure:"weeklyrainin"`
	Monthlyrainin  string `mapstructure:"monthlyrainin"`
	Yearlyrainin   string `mapstructure:"yearlyrainin"`
	Totalrainin    string `mapstructure:"totalrainin"`
	Solarradiation string `mapstructure:"solarradiation"`
	Uv             string `mapstructure:"uv"`
	Wh65batt       string `mapstructure:"wh65batt"`
	Freq           string `mapstructure:"freq"`
	Model          string `mapstructure:"model"`
}

func (wd *Ecowittdata) ToJSON() string {
	b, err := json.Marshal(wd)
	if err != nil {
		return ""
	}
	return string(b)
}

type Wundergrounddata struct {
	Id             string `mapstructure:"id"`
	Password       string `mapstructure:"PASSWORD"`
	Indoortempf    string `mapstructure:"indoortempf"`
	Tempf          string `mapstructure:"tempf"`
	Dewptf         string `mapstructure:"dewptf"`
	Windchillf     string `mapstructure:"windchillf"`
	Indoorhumidity string `mapstructure:"indoorhumidity"`
	Humidity       string `mapstructure:"humidity"`
	Windspeedmph   string `mapstructure:"windspeedmph"`
	Windgustmph    string `mapstructure:"windgustmph"`
	Winddir        string `mapstructure:"winddir"`
	Absbaromin     string `mapstructure:"absbaromin"`
	Baromin        string `mapstructure:"baromin"`
	Rainin         string `mapstructure:"rainin"`
	Dailyrainin    string `mapstructure:"dailyrainin"`
	Weeklyrainin   string `mapstructure:"weeklyrainin"`
	Monthlyrainin  string `mapstructure:"monthlyrainin"`
	Solarradiation string `mapstructure:"solarradiation"`
	Uv             string `mapstructure:"UV"`
	Dateutc        string `mapstructure:"dateutc"`
	Action         string `mapstructure:"action"`
	Realtime       string `mapstructure:"realtime"`
	Rtfreq         string `mapstructure:"rtfreq"`
}

func (wd *Wundergrounddata) ToJSON() string {
	b, err := json.Marshal(wd)
	if err != nil {
		return ""
	}
	return string(b)
}
