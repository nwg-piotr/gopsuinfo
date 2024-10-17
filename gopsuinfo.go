/*
A gopsutil-based command to display customizable system usage info in a single line

	Copyright (c) 2020-2021 Piotr Miller
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
	"github.com/shirou/gopsutil/net"
)

const version = "0.1.7"

var g glyphs
var path string

func cpuGraph(delay *string) string {
	bar := ""
	duration, _ := time.ParseDuration(*delay)
	speeds, _ := cpu.Percent(duration, true)
	for _, speed := range speeds {
		bar += string(g.graphCPU[int8(math.Round(speed/10))])
	}
	return bar
}

func cpuAvSpeed(asIcon bool, delay *string) string {
	output := ""
	if !asIcon {
		output += g.glyphCPU
	}
	duration, _ := time.ParseDuration(*delay)
	avSpeed, _ := cpu.Percent(duration, false)
	avs := fmt.Sprintf("%.1f", avSpeed[0])
	output += fmt.Sprintf("%s%%", avs)
	return output
}

func temperatures(asIcon bool) string {
	output := ""
	if !asIcon {
		output += g.glyphTemp
	}
	vals := make(map[string]int)

	temps, _ := host.SensorsTemperatures()
	for _, temp := range temps {
		// Some machines may return multiple sensors of the same name. Let's accept the 1st non-zero temp value.
		if vals["acpitz"] == 0 && temp.SensorKey == "acpitz_input" {
			vals["acpitz"] = int(temp.Temperature)
		}
		if vals["coretemp"] == 0 && temp.SensorKey == "coretemp_packageid0_input" || temp.SensorKey == "coretemp_core0_input" {
			vals["coretemp"] = int(temp.Temperature)
		}
		if temp.SensorKey == "k10temp_tctl_input" || temp.SensorKey == "k10temp_tdie_input" {
			vals["k10temp"] = int(temp.Temperature)
		}
		fmt.Println(temp.SensorKey, temp.Temperature)
	}
	// in case we still have no temperature value, let's try calculating average per-core value
	//sum := 0
	//for i, t := range temps {
	//	key := fmt.Sprintf("coretemp_core%v_input", i)
	//	val :=
	//}

	if v, ok := vals["k10temp"]; ok {
		output += fmt.Sprint(v)
	} else {
		if v, ok := vals["coretemp"]; ok {
			output += fmt.Sprint(v)
		} else {
			if v, ok := vals["acpitz"]; ok {
				output += fmt.Sprint(v)
			}
		}
	}

	output += "℃"
	return output
}

func uptime(asIcon bool) string {
	output := ""
	if !asIcon {
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

func memory(asIcon bool) string {
	output := ""
	if !asIcon {
		output += g.glyphMem + " "
	}
	stats, _ := mem.VirtualMemory()
	used := math.Round(float64(stats.Used)) / 1048576
	total := math.Round(float64(stats.Total)) / 1048576
	output += fmt.Sprintf("%.0f", used) + "/" + fmt.Sprintf("%.0f", total) + " MiB"
	return output
}

func traffic(asIcon bool) string {
	t0, _ := net.IOCounters(false)
	time.Sleep(time.Second)
	t1, _ := net.IOCounters(false)
	ul := math.Round(float64(t1[0].BytesSent-t0[0].BytesSent)) / 1024
	dl := math.Round(float64(t1[0].BytesRecv-t0[0].BytesRecv)) / 1024
	if asIcon {
		icon := net_icon(ul, dl)
		text := fmt.Sprintf("%.2f", ul) + " " + fmt.Sprintf("%.2f", dl) + " kB/s"
		return fmt.Sprintf("%s\n%s", icon, text)
	}
	return fmt.Sprintf("%.2f", ul) + " / " + fmt.Sprintf("%.2f", dl) + " kB/s"
}

func net_icon(ul float64, dl float64) string {
	fname := ""
	if ul >= 0.01 && dl >= 0.01 {
		fname = "xfer-b.svg"
	} else if ul >= 0.01 {
		fname = "xfer-u.svg"
	} else if dl >= 0.01 {
		fname = "xfer-d.svg"
	} else {
		fname = "xfer.svg"
	}
	return fmt.Sprintf("%s/%s", path, fname)
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

	componentsPtr := flag.String("c", "gatmnu",
		`Output (c)omponents: (a)vg CPU load, (g)rahical CPU bar,
		disk usage by mou(n)tpoints, (t)emperatures,
		networ(k) traffic, (m)emory, (u)ptime`)
	iconPtr := flag.String("i", "", "returns (i)con path and a single component (a, n, t, m, u) value")
	cpuDelayPtr := flag.String("d", "900ms", "CPU measurement (d)elay [timeout]")
	pathsPtr := flag.String("p", "/", "quotation-delimited, space-separated list of mount(p)oints")
	setPtr := flag.Bool("dark", false, "use (dark) icon set")
	textPtr := flag.Bool("t", false, "Just (t)ext, no glyphs")
	displayVersion := flag.Bool("v", false, "display (v)ersion information")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("gopsuinfo version %s\n", version)
		os.Exit(0)
	}

	if *textPtr {
		g = glyphs{graphCPU: []rune("_▁▂▃▄▅▆▇███"), glyphCPU: "", glyphMem: "", glyphTemp: "", glyphUptime: ""}
	} else {
		// Glyphs below may be replaced, e.g. "MEM:" instead of ""
		g = glyphs{graphCPU: []rune("_▁▂▃▄▅▆▇███"), glyphCPU: "", glyphMem: "", glyphTemp: "", glyphUptime: " "}
	}

	path = "/usr/share/gopsuinfo/icons_light"
	if *setPtr {
		path = "/usr/share/gopsuinfo/icons_dark"
	}

	output := ""

	if *iconPtr != "" {
		if *iconPtr == "g" {
			output += cpuGraph(cpuDelayPtr)
		}
		if *iconPtr == "a" {
			output += path + "/cpu.svg\n"
			output += cpuAvSpeed(true, cpuDelayPtr)
		} else if *iconPtr == "t" {
			output += path + "/temp.svg\n"
			output += temperatures(true)
		} else if *iconPtr == "n" {
			output += path + "/hdd.svg\n"
			output += diskUsage(pathsPtr)
		} else if *iconPtr == "m" {
			output += path + "/mem.svg\n"
			output += memory(true)
		} else if *iconPtr == "u" {
			output += path + "/up.svg\n"
			output += uptime(true)
		} else if *iconPtr == "k" {
			output += traffic(true)
		}
	} else {
		for _, char := range *componentsPtr {
			if string(char) == "g" {
				output += cpuGraph(cpuDelayPtr) + " "
			}
			if string(char) == "a" {
				output += cpuAvSpeed(false, cpuDelayPtr) + " "
			}
			if string(char) == "t" {
				output += temperatures(false) + " "
			}
			if string(char) == "u" {
				output += uptime(false) + " "
			}
			if string(char) == "m" {
				output += memory(false) + " "
			}
			if string(char) == "n" {
				output += diskUsage(pathsPtr) + " "
			}
			if string(char) == "k" {
				output += traffic(false)
			}
		}
	}

	fmt.Println(strings.TrimSpace(output))
}
