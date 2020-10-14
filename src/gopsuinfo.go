/*
Fed up with learning just by reading the manual, I'll try to code my old psuinfo python script from scratch in Go.
*/
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func cpuGraph(g glyphs, delay *string) string {
	bar := ""
	duration, _ := time.ParseDuration(*delay)
	speeds, _ := cpu.Percent(duration, true)
	for _, speed := range speeds {
		bar += string(g.graphCPU[int8(math.Round(speed/10))])
	}
	return bar
}

func cpuAvSpeed(g glyphs, delay *string) string {
	duration, _ := time.ParseDuration(*delay)
	avSpeed, _ := cpu.Percent(duration, false)
	avs := fmt.Sprintf("%.2f", avSpeed[0])
	if len(avs) < 5 {
		avs = " " + avs
	}
	return fmt.Sprintf("%s%s%%", g.glyphCPU, avs)
}

func temperatures(g glyphs) string {
	output := g.glyphTemp
	vals := make(map[string]int)
	temps, _ := host.SensorsTemperatures()
	for _, temp := range temps {
		if temp.SensorKey == "acpitz_input" {
			vals["acpitz"] = int(temp.Temperature)
		}
		if temp.SensorKey == "k10temp_tdie_input" {
			vals["k10temp"] = int(temp.Temperature)
		}
		if temp.SensorKey == "amdgpu_mem_input" {
			vals["amdgpu"] = int(temp.Temperature)
		}
	}
	if v, ok := vals["k10temp"]; ok {
		output += fmt.Sprint(v)
	} else {
		if v, ok := vals["acpitz"]; ok {
			output += fmt.Sprint(v)
		}
	}
	if v, ok := vals["amdgpu"]; ok {
		output += "|" + fmt.Sprint(v)
	}
	output += "℃"
	return output
}

func uptime(g glyphs) string {
	output := g.glyphUptime
	if t, e := host.Uptime(); e == nil {
		hh := t / 3600
		mm := t % 3600 / 60
		output += fmt.Sprintf("%02d:%02d", hh, mm)
	} else {
		output += "??:??"
	}
	return output
}

func memory(g glyphs) string {
	output := g.glyphMem + " "
	stats, _ := mem.VirtualMemory()
	used := math.Round(float64(stats.Used)) / 1048576
	total := math.Round(float64(stats.Total)) / 1048576
	output += fmt.Sprintf("%.0f", used) + "/" + fmt.Sprintf("%.0f", total) + "MiB"
	return output
}

func listMountpoints() {
	partitions, _ := disk.Partitions(true)
	fmt.Println("List in format: [device] mountpoint")
	for _, p := range partitions {
		fmt.Printf("[%s] %s\n", p.Device, p.Mountpoint)
	}
	os.Exit(0)
}

func diskUsage(paths *string) string {
	output := ""
	sliced := strings.Fields(*paths)
	for _, path := range sliced {
		usage, _ := disk.Usage(path)
		used := math.Round(float64(usage.Used)) / 1073741824
		total := math.Round(float64(usage.Total)) / 1073741824
		output += fmt.Sprintf("%s:%.1f/%.0f", path, used, total)
	}
	output += "GiB"
	return output
}

// Settings for now will store glyphs only
type glyphs struct {
	graphCPU    []rune
	glyphCPU    string
	glyphTemp   string
	glyphMem    string
	glyphUptime string
}

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "list_mountpoints" {
			listMountpoints()
		}
	}
	// Glyphs below may be replaced, e.g. "MEM:" instead of ""
	g := glyphs{graphCPU: []rune("_▁▂▃▄▅▆▇███"), glyphCPU: "", glyphMem: "", glyphTemp: "", glyphUptime: " "}

	componentsPtr := flag.String("c", "gatmnu",
		`Output (c)omponents: (a)vg CPU load, (f)an speed, (g)rahical bar, (t)emperatures,
		(m)emory, (u)ptime`)
	cpuDelayPtr := flag.String("d", "500ms", "CPU measurement delay [timeout]")
	pathsPtr := flag.String("p", "/", "Space-separated list of mou(n)tpoints")
	flag.Parse()

	output := ""

	for _, char := range *componentsPtr {
		if string(char) == "g" {
			output += cpuGraph(g, cpuDelayPtr) + " "
		}
		if string(char) == "a" {
			output += cpuAvSpeed(g, cpuDelayPtr) + " "
		}
		if string(char) == "t" {
			output += temperatures(g) + " "
		}
		if string(char) == "u" {
			output += uptime(g) + " "
		}
		if string(char) == "m" {
			output += memory(g) + " "
		}
		if string(char) == "n" {
			output += diskUsage(pathsPtr) + " "
		}
	}

	fmt.Println(strings.TrimSpace(output))
}
