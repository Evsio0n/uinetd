.PHONY: install uninstall

install:
	install -Dm0644 uinetd.service /usr/lib/systemd/system/uinetd.service
	cp uinetd.conf /etc/uinetd/uinetd.conf
	chmod 0644 /etc/uinetd/uinetd.conf
	touch /var/log/uinetd.log
	chmod 0644 /var/log/uinetd.log
	systemctl daemon-reload || true

uninstall:
	rm -f /usr/lib/systemd/system/uinetd.service /etc/uinetd/uinetd.conf /var/log/uinetd.log
	systemctl daemon-reload || true