package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	BAT_FULL = "/sys/class/power_supply/BAT1/energy_full"
	BAT_FULL2 = "/sys/class/power_supply/BAT2/energy_full"
	BAT_NOW  = "/sys/class/power_supply/BAT1/energy_now"
	BAT_NOW2  = "/sys/class/power_supply/BAT2/energy_now"
)

func getBattery(bmessages chan string) {
	for {
		full, _ := ioutil.ReadFile(BAT_FULL)
		full2, _ := ioutil.ReadFile(BAT_FULL2)
		bat_full, _ := strconv.Atoi(strings.Replace(string(full), "\n", "", 1))
		bat_full2, _ := strconv.Atoi(strings.Replace(string(full2), "\n", "", 1))
		fbat_full := float32(bat_full) + float32(bat_full2)
		now, _ := ioutil.ReadFile(BAT_NOW)
		now2, _ := ioutil.ReadFile(BAT_NOW2)
		bat_now, _ := strconv.Atoi(strings.Replace(string(now), "\n", "", 1))
		bat_now2, _ := strconv.Atoi(strings.Replace(string(now2), "\n", "", 1))
		fbat_now := float32(bat_now) + float32(bat_now2)
		bat_pct := strconv.FormatFloat(float64(((fbat_now / fbat_full) * 100)), 'f', 0, 32)
		bmessages <- fmt.Sprintf("%s %s%% ", "^i(/home/hcchu/repos/dzstatus/dzen-icons/bat_full_02.xbm)", string(bat_pct))
		time.Sleep(60 * time.Second)
	}
}
