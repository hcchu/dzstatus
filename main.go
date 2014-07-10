package main

import (
    "bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	SPADDING     = " "
	SCREEN_WIDTH = 1920
	FONT_FAMILY  = "DejaVu Sans"
	FONT_SIZE    = "11"
)

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

/*
    dzen_cmd := exec.Command("dzen2", "-dock", "-ta", "l", "-title-name", "panel",
                        "-fn", dzen_font, "-fg", COLOR_FOREGROUND, "-bg", COLOR_BACKGROUND)
    dzen, err :=  dzen_cmd.StdinPipe()
    if err != nil {
        fmt.Println("failed to run dzen2")
    }

    if err := dzen_cmd.Start(); err != nil {
        fmt.Println("failed to run dzen2")
    }
*/

    dzen := startDzen()

    dzen_writer := bufio.NewWriter(dzen)

	var title string
	var windows string
	var clock string
	var battery string
	var weather string
    var cpu string

	wmessages := make(chan string)
	tmessages := make(chan string)
	cmessages := make(chan string)
	bmessages := make(chan string)
	weathermsgs := make(chan string)
	cpuchan := make(chan string)

	go showWindows(wmessages)
	go getTitle(tmessages)
	go getTime(cmessages)
	go getBattery(bmessages)
	go getWeather("KMMU", weathermsgs)
    go getCpu(cpuchan)

	for {

		select {

		// title
		case tmsg := <-tmessages:
			title = tmsg

		// windows
		case wmsg := <-wmessages:
			windows = wmsg

		// time
		case cmsg := <-cmessages:
			clock = cmsg

		// battery
		case bmsg := <-bmessages:
			battery = bmsg

		// weather
		case weathermsg := <-weathermsgs:
			weather = weathermsg

        case cpumsg := <-cpuchan:
            cpu = cpumsg
		}


		lwidth := stringWidth(windows)
		cwidth := stringWidth(title)

		// adding some extra padding on the right
		rwidth := 20

		sstatus := ""

		right_contents := []string{cpu, weather, battery, clock}
		//right_contents := []string{clock, battery, weather}

		/*
		   for _, x := range right_contents {
		       rwidth += stringWidth(x) + 4 * len(SPADDING)
		       sstatus += fmt.Sprintf("^pa(%d)%s%s%s", SCREEN_WIDTH - rwidth, SPADDING, x, SPADDING)
		   }
		*/

		for x := len(right_contents) - 1; x >= 0; x-- {
			rwidth += stringWidth(right_contents[x]) + 30*len(SPADDING)
			sstatus += fmt.Sprintf("^pa(%d)%s%s%s", SCREEN_WIDTH-rwidth, SPADDING, right_contents[x], SPADDING)
		}

		var cindent int
		var max_lrw int

		available_center := SCREEN_WIDTH - (lwidth + rwidth)
		if available_center < cwidth {
			cindent = lwidth
		} else {
			if lwidth < rwidth {
				max_lrw = lwidth
			} else {
				max_lrw = rwidth
			}
			if (2 * (max_lrw + cwidth)) > SCREEN_WIDTH {
				cindent = lwidth + (available_center-cwidth)/2
			} else {
				cindent = (SCREEN_WIDTH - cwidth) / 2
			}
		}

        dzstatus := fmt.Sprintf("^pa(0)%s^pa(%d)%s%s\n", windows, cindent, title, sstatus)
        _, err := dzen_writer.WriteString(dzstatus)
        dzen_writer.Flush()
        if err != nil {
            fmt.Println("failed to send status to dzen2")
        }

	}

}
