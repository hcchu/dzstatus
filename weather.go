package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func getWeather(station string, weathermsgs chan string) {

    // sleep to let networking start
//    time.Sleep(60 * time.Second)

	var icon string
	var temp string
	sunny := "☀"
	cloudy := "☁"
	rain := "☔"
	snow := "☃"

	for {
		url := "http://weather.noaa.gov/pub/data/observations/metar/decoded/" + station + ".TXT"
		res, err := http.Get(url)
        if err != nil {
            time.Sleep(10 * time.Second)
            continue
        }

		body, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()

		lines := strings.Split(string(body), "\n")
		for _, line := range lines {
			switch {
			case strings.HasPrefix(line, "Sky"):
				switch {
				case strings.Contains(line, "sun") || strings.Contains(line, "clear"):
					icon = sunny
				case strings.Contains(line, "cloud") || strings.Contains(line, "overcast"):
					icon = cloudy
				case strings.Contains(line, "rain"):
					icon = rain
				case strings.Contains(line, "snow"):
					icon = snow
				}
			case strings.HasPrefix(line, "Weather"):
				switch {
				case strings.Contains(line, "rain"):
					icon = rain
				case strings.Contains(line, "snow"):
					icon = snow
				}
			case strings.HasPrefix(line, "Temp"):
				temp = strings.Fields(line)[1] + "F"
			}
		}

		weather := icon + " " + temp

		weathermsgs <- weather
		time.Sleep(60 * time.Second)
	}

}
