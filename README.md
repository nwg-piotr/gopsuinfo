# gopsuinfo

This project is a Go version of my [psuinfo](https://github.com/nwg-piotr/psuinfo) python script, written for educational purposes, and also for better performance. The code only implements these `psuinfo` features which I actually use.

The program uses the [gopsutil](https://github.com/shirou/gopsutil) module, Copyright (c) 2014, WAKAYAMA Shirou.

For use with bars like [Waybar](https://github.com/Alexays/Waybar) or [polybar](https://github.com/polybar/polybar), the `gopsuinfo -c <components_string>` is suitable. It prints system information in a single line:

`gopsuinfo -c gatmnu`

![gopsuinfo-waybar](http://nwg.pl/Lychee/uploads/big/29020400991f20e10272b4c3c65d37c1.png)

For panels capable of displaying an icon and text, like [tint2](https://gitlab.com/o9000/tint2)
or[nwg-panel](https://github.com/nwg-piotr/nwg-panel), you need to define an executor for each component.
For instance `gopsuinfo -i m` will return a path to the memory icon, and the current memory usage:

```text
/usr/share/gopsuinfo/icons_light/mem.svg
2103/22008 MiB
```

Add all the components you need this way.

![gopsuinfo.png](https://scrot.cloud/images/2021/02/24/gopsuinfo.png)

```
$ gopsuinfo -h
Use gopsuinfo list_mountpoints to see available mount points.
Usage of gopsuinfo:
  -c string
    	Output (c)omponents: (a)vg CPU load, (g)rahical CPU bar, disk usage by mou(n)tpoints, (t)emperatures, (m)emory, (u)ptime (default "gatmnu")
  -d string
    	CPU measurement delay [timeout] (default "900ms")
  -dark
    	Use dark icon set
  -i string
    	Returns (i)con path and a single component (a, n, t, m, u) value
  -p string
    	Quotation-delimited, space-separated list of mou(n)tpoints (default "/")
```

## Installation

Clone the repository:

```
git clone https://github.com/nwg-piotr/gopsuinfo.git
cd gopsuinfo
```

Get the gopsutil library:

```
make get
```

Build binary:

```
make build
```

Install files:

```
sudo make install
```

## To uninstall

```
sudo make uninstall
```
