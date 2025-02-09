PREFIX=$(HOME)/.local
BINDIR=$(PREFIX)/bin
UNITDIR=$(HOME)/.config/systemd/user

all: build

build:
	go build -o spotigo main.go

install: build
	install -Dm755 spotigo $(BINDIR)/spotigo
	install -Dm644 systemd-user/spotigo.service $(UNITDIR)/spotigo.service
	systemctl --user daemon-reload

uninstall:
	rm -f $(BINDIR)/spotigo
	rm -f $(UNITDIR)/spotigo.service
	systemctl --user daemon-reload

clean:
	rm -f spotigo
