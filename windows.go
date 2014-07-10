package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const (
	PADDING                   = " "
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
