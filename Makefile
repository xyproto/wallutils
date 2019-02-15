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
	(cd cmd/lscollection; go build)
	(cd cmd/timedinfo; go build)
	(cd cmd/lswallpaper; go build)
	(cd cmd/setcollection; go build)
	(cd cmd/lstimed; go build)
	(cd cmd/settimed; go build)
	(cd cmd/xml2stw; go build)

install:
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/getdpi/getdpi
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lsmon/lsmon
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/setrandom/setrandom
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/setwallpaper/setwallpaper
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/wayinfo/wayinfo
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/xinfo/xinfo
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lscollection/lscollection
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/timedinfo/timedinfo
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/lswallpaper
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/setcollection/setcollection
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/settimed/settimed
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lstimed/lstimed
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/xml2stw/xml2stw

clean:
	(cd cmd/getdpi; go clean)
	(cd cmd/lsmon; go clean)
	(cd cmd/setrandom; go clean)
	(cd cmd/setwallpaper; go clean)
	(cd cmd/wayinfo; go clean)
	(cd cmd/xinfo; go clean)
	(cd cmd/lscollection; go clean)
	(cd cmd/timedinfo; go clean)
	(cd cmd/lswallpaper; go clean)
	(cd cmd/setcollection; go clean)
	(cd cmd/lstimed; go clean)
	(cd cmd/settimed; go clean)
	(cd cmd/xml2stw; go clean)
	go clean
