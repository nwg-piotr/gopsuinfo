/*
Fed up with learning just by reading the manual, I'll try to code my old psuinfo python script from scratch in Go.
*/
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func main() {
	var graph = []rune("_▁▂▃▄▅▆▇███")
	var bar = ""

	duration, _ := time.ParseDuration("400ms")
	speeds, _ := cpu.Percent(duration, true)
	for _, speed := range speeds {
		bar += string(graph[int8(math.Round(speed/10))])
	}
	avSpeed, _ := cpu.Percent(duration, false)
	avs := fmt.Sprintf("%.1f", avSpeed[0])
	if len(avs) < 4 {
		avs = " " + avs
	}
	var output = fmt.Sprintf("%s %s%%", bar, avs)
	fmt.Println(output)
}
