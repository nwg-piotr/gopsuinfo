get:
	go get github.com/shirou/gopsutil

build:
	go build -o bin/gopsuinfo gopsuinfo.go

install:
	mkdir -p /usr/share/gopsuinfo
	cp -R icons_light /usr/share/gopsuinfo
	cp -R icons_dark /usr/share/gopsuinfo
	cp bin/gopsuinfo /usr/bin

uninstall:
	rm -r /usr/share/gopsuinfo
	rm /usr/bin/gopsuinfo
	rm /usr/share/applications/nwgocc.desktop
	rm /usr/share/pixmaps/nwgocc.svg

run:
	go run gopsuinfo.go
