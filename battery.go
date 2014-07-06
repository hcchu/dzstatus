package main

import (
    "io/ioutil"
    "strconv"
    "strings"
)

const (
    BAT_FULL                  = "/sys/class/power_supply/BAT1/energy_full"
    BAT_NOW                   = "/sys/class/power_supply/BAT1/energy_now"
)

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

