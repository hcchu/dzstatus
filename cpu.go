package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
    "time"
)

func cpuNow() (float64, float64) {
    f, err := os.Open("/proc/stat")
    if err != nil {
        fmt.Println("can't open /proc/stat")
    }

    r := bufio.NewReader(f)
    s, e := r.ReadString('\n')
    if e != nil {
        fmt.Println("error reading from /proc/stat")
    }
    f.Close()
    cpu := strings.Fields(s)
    tcpu := 0.0
    for _, i := range(cpu[1:]) {
        x,_ := strconv.ParseFloat(i,32)
        tcpu += x
    }
    idle,_ := strconv.ParseFloat(cpu[4],32)
    return tcpu, idle
}

func getCpu(cpuchan chan string) {

    for {

        ptotal, pidle := cpuNow()
        time.Sleep(1 * time.Second)
        total, idle := cpuNow()
        
        cpu := (((total-ptotal) - (idle-pidle)) / (total-ptotal)) * 100

        result := strconv.FormatFloat(cpu, 'f', 0, 32)
        cpumsg := fmt.Sprintf("%s %s%% ", "^i(/home/hcchu/repos/dzstatus/dzen-icons/cpu.xbm)", result)
        
        cpuchan <-cpumsg
        time.Sleep(2 * time.Second)

    }

}
