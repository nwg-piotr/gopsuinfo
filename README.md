# gopsuinfo

This project is a Go version of my [psuinfo](https://github.com/nwg-piotr/psuinfo) python script, written for educational purposes, and also for better performance.
The code only implements these `psuinfo` features which I actually use.

The `gopsuinfo` command prints some customisable system usage information in a single line, and is intended to use with text-based panels, like [Waybar](https://github.com/Alexays/Waybar) or [polybar](https://github.com/polybar/polybar).

The program uses the [gopsutil](https://github.com/shirou/gopsutil) package, Copyright (c) 2014, WAKAYAMA Shirou.
