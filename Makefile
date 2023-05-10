.PHONY: all clean install

DESTDIR ?=
PREFIX ?= /usr
UNAME_R ?= $(shell uname -r)

ifneq (,$(findstring arch,$(UNAME_R)))
# Arch Linux
LDFLAGS ?= -Wl,-O2,--sort-common,--as-needed,-z,relro,-z,now
BUILDFLAGS ?= -mod=vendor -buildmode=pie -trimpath -ldflags "-s -w -linkmode=external -extldflags $(LDFLAGS)"
else
# Default settings
BUILDFLAGS ?= -mod=vendor -trimpath
endif

# build all utilities, but allow heic2stw to fail if the wrong version of libheif is installed
all:
	go build ${BUILDFLAGS}
	(cd cmd/getdpi; go build ${BUILDFLAGS})
	-(cd cmd/heic2stw; go build ${BUILDFLAGS})
	(cd cmd/lscollection; go build ${BUILDFLAGS})
	(cd cmd/lsmon; go build ${BUILDFLAGS})
	(cd cmd/lstimed; go build ${BUILDFLAGS})
	(cd cmd/lswallpaper; go build ${BUILDFLAGS})
	(cd cmd/setcollection; go build ${BUILDFLAGS})
	(cd cmd/setrandom; go build ${BUILDFLAGS})
	(cd cmd/settimed; go build ${BUILDFLAGS})
	(cd cmd/setwallpaper; go build ${BUILDFLAGS})
	(cd cmd/timedinfo; go build ${BUILDFLAGS})
	(cd cmd/wayinfo; go build ${BUILDFLAGS})
	(cd cmd/xinfo; go build ${BUILDFLAGS})
	(cd cmd/xml2stw; go build ${BUILDFLAGS})
	(cd cmd/vram; go build ${BUILDFLAGS})

# utilities that depend on dynamic libraries are commented out
static:
	CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a
	@#(cd cmd/getdpi; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	@#(cd cmd/heic2stw; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	(cd cmd/lscollection; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	@#(cd cmd/lsmon; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	(cd cmd/lstimed; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	(cd cmd/lswallpaper; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	@#(cd cmd/setcollection; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	(cd cmd/setrandom; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	(cd cmd/settimed; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	(cd cmd/setwallpaper; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	(cd cmd/timedinfo; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	@#(cd cmd/wayinfo; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	@#(cd cmd/xinfo; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	(cd cmd/xml2stw; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)
	(cd cmd/vram; CGO_ENABLED=0 go build ${BUILDFLAGS} -ldflags "-s" -a)

# install all utilities, but let the ones that depend on dynamic libraries be optional
install:
	-install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/getdpi/getdpi
	-install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/heic2stw/heic2stw && \
	  install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" scripts/heic-install
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lscollection/lscollection
	-install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lsmon/lsmon
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lstimed/lstimed
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/lswallpaper/lswallpaper
	-install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/setcollection/setcollection
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/setrandom/setrandom
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/settimed/settimed
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/setwallpaper/setwallpaper
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/timedinfo/timedinfo
	-install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/wayinfo/wayinfo
	-install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/xinfo/xinfo
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/xml2stw/xml2stw
	install -Dm755 -t "$(DESTDIR)$(PREFIX)/bin" cmd/vram/vram

clean:
	(cd cmd/getdpi; go clean)
	(cd cmd/heic2stw; go clean)
	(cd cmd/lscollection; go clean)
	(cd cmd/lsmon; go clean)
	(cd cmd/lstimed; go clean)
	(cd cmd/lswallpaper; go clean)
	(cd cmd/setcollection; go clean)
	(cd cmd/setrandom; go clean)
	(cd cmd/settimed; go clean)
	(cd cmd/setwallpaper; go clean)
	(cd cmd/timedinfo; go clean)
	(cd cmd/wayinfo; go clean)
	(cd cmd/xinfo; go clean)
	(cd cmd/xml2stw; go clean)
	(cd cmd/vram; go clean)
	go clean
