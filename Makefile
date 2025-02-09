PREFIX=$(HOME)/.local
BINDIR=$(PREFIX)/bin

all: build

build:
	go build -o spotigo main.go

install: build
	install -Dm755 spotigo $(BINDIR)/spotigo

uninstall:
	rm -f $(BINDIR)/spotigo

clean:
	rm -f spotigo
