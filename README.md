# gopsuinfo

This application is a part of the [nwg-shell](https://nwg-piotr.github.io/nwg-shell) project.

**Contributing:** please read the [general contributing rules for the nwg-shell project](https://nwg-piotr.github.io/nwg-shell/contribution).

This project is a Go version of my [psuinfo](https://github.com/nwg-piotr/psuinfo) python script, written for educational purposes, and also for better performance. The code only implements these `psuinfo` features that I actually use.

[![Packaging status](https://repology.org/badge/vertical-allrepos/gopsuinfo.svg)](https://repology.org/project/gopsuinfo/versions)

The program uses the [gopsutil](https://github.com/shirou/gopsutil) module, Copyright (c) 2014, WAKAYAMA Shirou.

For use with bars like [Waybar](https://github.com/Alexays/Waybar) or [polybar](https://github.com/polybar/polybar), the `gopsuinfo -c <components_string>` is suitable. It prints system information in a single line:

`gopsuinfo -c gatmnu`

![image](https://user-images.githubusercontent.com/20579136/171514998-3423165f-5628-4d49-8dde-06801d817993.png)

For panels capable of displaying an icon and text, like [tint2](https://gitlab.com/o9000/tint2)
or [nwg-panel](https://github.com/nwg-piotr/nwg-panel), you need to define an executor for each component.
For instance `gopsuinfo -i m` will return a path to the memory icon, and the current memory usage:

```text
/usr/share/gopsuinfo/icons_light/mem.svg
2103/22008 MiB
```

Add all the components you need to this way. Sample output with monochrome icons in nwg-panel:

![image](https://user-images.githubusercontent.com/20579136/171515322-f469d580-72e7-4950-9857-28746e380d6a.png)

```
$ gopsuinfo -h
Use gopsuinfo list_mountpoints to see available mount points.
Usage of bin/gopsuinfo:
  -c string
    	Output (c)omponents: (a)vg CPU load, (g)rahical CPU bar,
    			disk usage by mou(n)tpoints, (t)emperatures,
    			networ(k) traffic, (m)emory, (u)ptime (default "gatmnu")
  -d string
    	CPU measurement (d)elay [timeout] (default "900ms")
  -dark
    	use (dark) icon set
  -i string
    	returns (i)con path and a single component (a, n, t, m, u) value
  -p string
    	quotation-delimited, space-separated list of mount(p)oints (default "/")
  -t	Just (t)ext, no glyphs
  -v	display (v)ersion information
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
