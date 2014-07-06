package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
    "regexp"
	"strings"
	//"time"
)

const (
    BAT_FULL                  = "/sys/class/power_supply/BAT1/energy_full"
    BAT_NOW                   = "/sys/class/power_supply/BAT1/energy_now"
    SCREEN_WIDTH              = 1920
	PADDING                   = " "
	FONT_FAMILY               = "DejaVu Sans"
	FONT_SIZE                 = "11"
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

var dzen_string string

func stringWidth(s string) int {
    re := regexp.MustCompile("\\^[a-z]+\\([^)]*\\)")
    s = re.ReplaceAllString(s, "")
    out, err := exec.Command("txtw", "-f", FONT_FAMILY, "-s", FONT_SIZE, s).Output()
    if err != nil {
        fmt.Println(err)
        return 2
    }
    w, _ := strconv.Atoi(strings.Replace(string(out), "\n", "", 1))
    return w
}

func getBattery() string {
    full, _ := ioutil.ReadFile(BAT_FULL)
    bat_full, _ := strconv.Atoi(strings.Replace(string(full), "\n", "", 1))
    fbat_full := float32(bat_full)
    now, _ := ioutil.ReadFile(BAT_NOW)
    bat_now, _ := strconv.Atoi(strings.Replace(string(now), "\n", "", 1))
    fbat_now := float32(bat_now)
    bat_pct := strconv.FormatFloat(float64(((fbat_now/fbat_full) * 100)), 'f', 0, 32)
    return bat_pct
}

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

func showWindows(messages chan string) {
	cmd := exec.Command("bspc", "control", "--subscribe")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error")
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("error")
	}
    wmstring := ""
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		s := scanner.Text()
        wmstring = ""
		for _, i := range strings.Split(s, ":") {
            wmatch, _ := regexp.MatchString("^[OoFfUu]\\d", i)
            lmatch, _ := regexp.MatchString("^L", i)
            if wmatch {
                name := string(i[1])
                var fg string
                var bg string
                switch {
                case strings.HasPrefix(i, "O"):
                    fg = COLOR_FOCUSED_OCCUPIED_FG
                    bg = COLOR_FOCUSED_OCCUPIED_BG
                case strings.HasPrefix(i, "F"):
                    fg = COLOR_FOCUSED_FREE_FG
                    bg = COLOR_FOCUSED_FREE_BG
                case strings.HasPrefix(i, "U"):
                    fg = COLOR_FOCUSED_URGENT_FG
                    bg = COLOR_FOCUSED_URGENT_BG
                case strings.HasPrefix(i, "o"):
                    fg = COLOR_OCCUPIED_FG
                    bg = COLOR_OCCUPIED_BG
                case strings.HasPrefix(i, "f"):
                    fg = COLOR_FREE_FG
                    bg = COLOR_FREE_BG
                case strings.HasPrefix(i, "u"):
                    fg = COLOR_URGENT_FG
                    bg = COLOR_URGENT_BG
                }
                wmstring = fmt.Sprintf("%s^fg(%s)^bg(%s)^ca(1, bspc desktop -f %s)^ca(2, bspc window -d %s)%s%s%s^ca()^ca()", wmstring, fg, bg, name, name, PADDING, name, PADDING)
            } else if lmatch {
                name := string(i[1])
                layout := strings.ToUpper(name)
                wmstring = fmt.Sprintf("%s^fg()^bg()%s%s^fg(%s)^bg(%s)^ca(1, bspc desktop -l next)%s%s%s^ca()", wmstring, PADDING, PADDING, COLOR_LAYOUT_FG, COLOR_LAYOUT_BG, PADDING, layout, PADDING)
            }
        }
        messages <- wmstring
    }
}

func getTitle(tmessages chan string) {
	cmd := exec.Command("xtitle", "-s")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error")
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("error")
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
        s := scanner.Text()
        fmt.Println(s)
        tmessages <-s
    }
}

func main() {
	//swidth := screen_width()
	//fmt.Println(swidth)
    fmt.Println(getBattery())

    wmessages := make(chan string)
    tmessages := make(chan string)

    go showWindows(wmessages)
    go getTitle(tmessages)

    /*
    go func() {
        t := time.Now().Local()
		fmt.Printf("S%s    %s\n", getWeather("KMMU"), t.Format("Jan 2 15:04"))
		time.Sleep(10 * time.Second)
    }
    */

    tmsg := <-tmessages
    fmt.Println(tmsg)

    wmsg := <-wmessages
    fmt.Println(wmsg)
        //fmt.Println(stringWidth(tmsg))
}
