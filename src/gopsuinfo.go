/*A gopsutil-based command to display customizable system usage info in a single line
  Copyright (c) 2020 Piotr Miller
  e-mail: nwg.piotr@gmail.com
  Project: https://github.com/nwg-piotr/gopsuinfo
  License: GPL3
  gopsutil Copyright (c) 2014, WAKAYAMA Shirou, https://github.com/shirou/gopsutil
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

func cpuAvSpeed(as_icon bool, g glyphs, delay *string) string {
	duration, _ := time.ParseDuration(*delay)
	avSpeed, _ := cpu.Percent(duration, false)
	avs := fmt.Sprintf("%.0f", math.Round(float64(avSpeed[0])))
	if len(avs) < 2 {
		avs = " " + avs
	}
	if as_icon {
		return fmt.Sprintf("%s%%", avs)
	}
	return fmt.Sprintf("%s%s%%", g.glyphCPU, avs)
}

func temperatures(as_icon bool, g glyphs) string {
	output := ""
	if !as_icon {
		output += g.glyphTemp
	}
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

func uptime(as_icon bool, g glyphs) string {
	output := ""
	if !as_icon {
		output += g.glyphUptime
	}
	if t, e := host.Uptime(); e == nil {
		hh := t / 3600
		mm := t % 3600 / 60
		output += fmt.Sprintf("%02d:%02d", hh, mm)
	} else {
		output += "??:??"
	}
	return output
}

func memory(as_icon bool, g glyphs) string {
	output := ""
	if !as_icon {
		output += g.glyphMem + " "
	}
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
		output += fmt.Sprintf("%s:%.1f/%.0f", path, used, total) + " "
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
		if arg == "-h" {
			fmt.Println("Use gopsuinfo list_mountpoints to see available mount points.")
		}
	}
	// Glyphs below may be replaced, e.g. "MEM:" instead of ""
	g := glyphs{graphCPU: []rune("_▁▂▃▄▅▆▇███"), glyphCPU: "", glyphMem: "", glyphTemp: "", glyphUptime: " "}

	componentsPtr := flag.String("c", "gatmnu",
		`Output (c)omponents: (a)vg CPU load, (g)rahical CPU bar, disk usage by mou(n)tpoints, (t)emperatures, (m)emory, (u)ptime`)
	iconPtr := flag.String("i", "", "Returns (i)con path and a single component (a, n, t, m, u) value")
	cpuDelayPtr := flag.String("d", "900ms", "CPU measurement delay [timeout]")
	pathsPtr := flag.String("p", "/", "Quotation-delimited, space-separated list of mou(n)tpoints")
	setPtr := flag.Bool("dark", false, "Use dark icon set")

	flag.Parse()
	path := "/home/piotr/Obrazy/gopsuinfo/icons_light"
	if *setPtr {
		path = "/home/piotr/Obrazy/gopsuinfo/icons_dark"
	}

	output := ""

	if *iconPtr != "" {
		if *iconPtr == "a" {
			output += path + "/cpu.svg\n"
			output += cpuAvSpeed(true, g, cpuDelayPtr)
		} else if *iconPtr == "t" {
			output += path + "/temp.svg\n"
			output += temperatures(true, g)
		} else if *iconPtr == "n" {
			output += path + "/hdd.svg\n"
			output += diskUsage(pathsPtr)
		} else if *iconPtr == "m" {
			output += path + "/mem.svg\n"
			output += diskUsage(pathsPtr)
		} else if *iconPtr == "u" {
			output += path + "/up.svg\n"
			output += uptime(true, g)
		}
	} else {
		for _, char := range *componentsPtr {
			if string(char) == "g" {
				output += cpuGraph(g, cpuDelayPtr) + " "
			}
			if string(char) == "t" {
				output += cpuAvSpeed(false, g, cpuDelayPtr) + " "
			}
			if string(char) == "t" {
				output += temperatures(false, g) + " "
			}
			if string(char) == "u" {
				output += uptime(false, g) + " "
			}
			if string(char) == "m" {
				output += memory(false, g) + " "
			}
			if string(char) == "n" {
				output += diskUsage(pathsPtr) + " "
			}
		}
	}

	fmt.Println(strings.TrimSpace(output))
}
