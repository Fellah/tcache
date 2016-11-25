SSH = leta@tcache
VERSION = latest

build:
	git tag --force --annotate $(VERSION) -m "$(VERSION)"
	git push --force origin $(VERSION)
	ssh $(SSH) 'docker pull fellah/tcache:latest'

restart:
	ssh $(SSH) 'systemctl --user restart tcache.service'

start:
	ssh $(SSH) 'systemctl --user start tcache.service'

stop:
	ssh $(SSH) 'systemctl --user stop tcache.service'

.DEFAULT_GOAL := build

.PHONY: build restart
