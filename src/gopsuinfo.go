/*
Fed up with learning just by reading the manual, I'll try to code my old psuinfo python script from scratch in Go.
*/
package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func cpuGraph(delay *string) string {
	bar := ""
	graph := []rune("_▁▂▃▄▅▆▇███")
	duration, _ := time.ParseDuration(*delay)
	speeds, _ := cpu.Percent(duration, true)
	for _, speed := range speeds {
		bar += string(graph[int8(math.Round(speed/10))])
	}
	return bar
}

func cpuAvSpeed(delay *string) string {
	duration, _ := time.ParseDuration(*delay)
	avSpeed, _ := cpu.Percent(duration, false)
	avs := fmt.Sprintf("%.2f", avSpeed[0])
	if len(avs) < 5 {
		avs = " " + avs
	}
	return fmt.Sprintf("%s%%", avs)

}

func main() {
	componentsPtr := flag.String("C", "gatmd", `Components: (a)vg CPU load, (f)an speed, (g)rahical bar, (t)emperatures,
	(m)emory, (u)ptime`)
	cpuDelayPtr := flag.String("delay", "400ms", "CPU measurement delay [timeout]")
	flag.Parse()

	output := ""

	for _, char := range *componentsPtr {
		if string(char) == "g" {
			output += cpuGraph(cpuDelayPtr) + " "
		}
		if string(char) == "a" {
			output += cpuAvSpeed(cpuDelayPtr) + " "
		}
	}

	fmt.Println(output)
}
