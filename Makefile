.PHONY: all clean install uninstall

all: install

PREFIX=/usr/local

clean:
	rm -f ./uinetd

build:
	go get -u -v github.com/evsio0n/log
	go build -o uinetd

install:build
	install -Dm0755 uinetd "$(PREFIX)/bin/uinetd"
	$(MAKE) -C systemd install

uninstall:
	rm -f "$(PREFIX)/bin/uinetd"
	$(MAKE) -C systemd uninstall
	rm -f ./uinetd

uinetd: main.go
	go get -v github.com/evsio0n/log
	go build -o uinetd
