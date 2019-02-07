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
	(cd cmd/lscollections; go build)
	(cd cmd/lsgnomewallpaper; go build)
	(cd cmd/lswallpaper; go build)
	(cd cmd/setcollection; go build)

install:
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/getdpi
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/lsmon
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/setrandom
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/setwallpaper
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/wayinfo
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/xinfo
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/lscollections
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/lsgnomewallpaper
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/lswallpaper
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/setcollection

clean:
	(cd cmd/getdpi; go clean)
	(cd cmd/lsmon; go clean)
	(cd cmd/setrandom; go clean)
	(cd cmd/setwallpaper; go clean)
	(cd cmd/wayinfo; go clean)
	(cd cmd/xinfo; go clean)
	(cd cmd/lscollections; go clean)
	(cd cmd/lsgnomewallpaper; go clean)
	(cd cmd/lswallpaper; go clean)
	(cd cmd/setcollection; go clean)
	go clean
