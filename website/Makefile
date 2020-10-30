# cerebrum make

.PHONY: list # show all make targets
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' #| sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

.PHONY: help # show all make targets with descriptions
help:
	@echo "-----------------------"
	@echo "|         help        |"
	@echo "-----------------------"
	@echo "| TARGET: DESCRIPTION |"
	@echo "-----------------------"
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1: \2/'

.PHONY: install # runs npm install
install:
	npm install

.PHONY: update # run npm update, to make major dependency updates install `npm install -g npm-check-updates` then run `ncu -u`
update:
	npm update

.PHONY: serve # run npm start, to serve the UI
serve:
	npm start

.PHONY: lint # run es lint
lint:
	./node_modules/.bin/eslint --fix src
