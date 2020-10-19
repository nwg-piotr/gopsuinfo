# gopsuinfo

This project is a Go version of my [psuinfo](https://github.com/nwg-piotr/psuinfo) python script, written for educational purposes, and also for better performance.
The code only implements these `psuinfo` features which I actually use.

The `gopsuinfo` command prints some customisable system usage information in a single line, and is intended to use with text-based panels, like [Waybar](https://github.com/Alexays/Waybar) or [polybar](https://github.com/polybar/polybar).

The program uses the [gopsutil](https://github.com/shirou/gopsutil) package, Copyright (c) 2014, WAKAYAMA Shirou.

```
$ gopsuinfo -h
Use gopsuinfo list_mountpoints to see available mount points.
Usage of gopsuinfo:
  -c string
    	Output (c)omponents: (a)vg CPU load, (g)rahical CPU bar, disk usage by mou(n)tpoints, (t)emperatures, (m)emory, (u)ptime (default "gatmnu")
  -d string
    	CPU measurement delay [timeout] (default "450ms")
  -p string
    	Quotation-delimited, space-separated list of mou(n)tpoints (default "/")
```

Sample usage:

![gopsuinfo-waybar](http://nwg.pl/Lychee/uploads/big/29020400991f20e10272b4c3c65d37c1.png)

Waybar custom modules:

```json
"custom/cpubar": {
    "interval": 1,
    "return-type": "text",
    "exec": "gopsuinfo -c g",
    "escape": true
},
"custom/cpuavs": {
    "interval": 1,
    "return-type": "text",
    "exec": "gopsuinfo -c a",
    "escape": true
,
"custom/swayinfo": {
    "interval": 5,
    "return-type": "text",
    "exec": "gopsuinfo -c tmnu",
    "escape": true
}
```
