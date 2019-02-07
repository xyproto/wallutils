.PHONY: all clean install

DESTDIR ?=
PREFIX ?= /usr

all:
	go build
	(cd cmd/getdpi; go build)
	(cd cmd/lsmon; go build)
	(cd cmd/setrandom; go build)
	(cd cmd/setwallpaper; go build)
	(cd cmd/wayinfo; go build)
	(cd cmd/xinfo; go build)

install:
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/getdpi
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/lsmon
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/setrandom
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/setwallpaper
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/wayinfo
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/xinfo

clean:
	(cd cmd/getdpi; go clean)
	(cd cmd/lsmon; go clean)
	(cd cmd/setrandom; go clean)
	(cd cmd/setwallpaper; go clean)
	(cd cmd/wayinfo; go clean)
	(cd cmd/xinfo; go clean)
	go clean
