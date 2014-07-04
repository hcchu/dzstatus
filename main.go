package main

import (
    "bufio"
	"fmt"
    "net/http"
    "io/ioutil"
	"os/exec"
//	"strconv"
	"strings"
    "time"
)

const (
	FONT_FAMILY               = "DejaVu Sans"
	FONT_SIZE                 = 11
	COLOR_FOREGROUND          = "#A3A6AB"
	COLOR_BACKGROUND          = "#34322E"
	COLOR_FOCUSED_OCCUPIED_FG = "#F6F9FF"
	COLOR_FOCUSED_OCCUPIED_BG = "#5C5955"
	COLOR_FOCUSED_FREE_FG     = "#F6F9FF"
	COLOR_FOCUSED_FREE_BG     = "#6D561C"
	COLOR_FOCUSED_URGENT_FG   = "#34322E"
	COLOR_FOCUSED_URGENT_BG   = "#F9A299"
	COLOR_OCCUPIED_FG         = "#A3A6AB"
	COLOR_OCCUPIED_BG         = "#34322E"
	COLOR_FREE_FG             = "#6F7277"
	COLOR_FREE_BG             = "#34322E"
	COLOR_URGENT_FG           = "#F9A299"
	COLOR_URGENT_BG           = "#34322E"
	COLOR_LAYOUT_FG           = "#A3A6AB"
	COLOR_LAYOUT_BG           = "#34322E"
	COLOR_TITLE_FG            = "#A3A6AB"
	COLOR_TITLE_BG            = "#34322E"
	COLOR_STATUS_FG           = "#A3A6AB"
	COLOR_STATUS_BG           = "#34322E"
)

/*
func screen_width() int {
	out, err := exec.Command("sres", "-W").Output()
	if err != nil {
		fmt.Println(err)
		return 2
	}

	swidth, err := strconv.Atoi(strings.Replace(string(out), "\n", "", 1))
	if err != nil {
		fmt.Println(err)
		return 2
	}
	return swidth
}
*/

func getWeather(station string) string {

    var icon string
    var temp string
    sunny := "☀"
    cloudy := "☁"
    rain := "☔"
    snow := "☃"

    url := "http://weather.noaa.gov/pub/data/observations/metar/decoded/" + station + ".TXT"
    res, err := http.Get(url)
    if err != nil {
        return "error"
    }

    body, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        return "error"
    }

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

    return weather

}

    

func main() {
        //swidth := screen_width()
        //fmt.Println(swidth)

    cmd := exec.Command("bspc", "control", "--subscribe")
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        fmt.Println("error")
    }
    if err := cmd.Start(); err != nil {
        fmt.Println("error")
    }
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }

        
    for true {
        t := time.Now().Local()
        fmt.Printf("S%s    %s\n", getWeather("KMMU"), t.Format("Jan 2 15:04"))
        time.Sleep(10 * time.Second)
    }
}
