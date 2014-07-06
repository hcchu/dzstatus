package main

import (
	"fmt"
	"os/exec"
	"strconv"
    "regexp"
	"strings"
	//"time"
)

const (
    SCREEN_WIDTH              = 1920
	FONT_FAMILY               = "DejaVu Sans"
	FONT_SIZE                 = "11"
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

    for {

        select {

        case tmsg := <-tmessages:
            fmt.Println(tmsg)

        case wmsg := <-wmessages:
            fmt.Println(wmsg)
        //fmt.Println(stringWidth(tmsg))
    }
    }
}
