.PHONY: all clean install uninstall

all: uinetd

PREFIX=/usr/local

clean:
	rm -f ./uinetd


install: all
    install -Dm0755 uinetd "$(PREFIX)/bin/uinetd"
    $(MAKE) -C systemd install

uninstall:
	rm -f "$(PREFIX)/bin/uinetd"
	$(MAKE) -C systemd uninstall

uinetd: main.go
	go build -o uinetd
