build:
	go build -o catsel

install: build
	sudo mv catsel /usr/local/bin/

reinstall: build install

uninstall:
	sudo rm -f /usr/local/bin/catsel

