get:
	go get github.com/shirou/gopsutil

build:
	go build -o bin/gopsuinfo gopsuinfo.go

install:
	mkdir -p "${DESTDIR}/usr/share/gopsuinfo" "${DESTDIR}/usr/bin"
	cp -R icons_light "${DESTDIR}/usr/share/gopsuinfo"
	cp -R icons_dark "${DESTDIR}/usr/share/gopsuinfo"
	cp bin/gopsuinfo "${DESTDIR}/usr/bin/"

uninstall:
	rm -r "${DESTDIR}/usr/share/gopsuinfo"
	rm "${DESTDIR}/usr/bin/gopsuinfo"

run:
	go run gopsuinfo.go
