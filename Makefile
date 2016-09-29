SSH = tcache
VERSION = latest

build:
	git tag --force --annotate $(VERSION) -m "$(VERSION)"
	git push --force origin $(VERSION)

restart:
	ssh $(SSH) 'systemctl --user restart tcache.service'

.DEFAULT_GOAL := build

.PHONY: build restart
