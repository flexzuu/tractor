.PHONY: build setup clean clobber dev versions studio kill qtalk

build: clean local/bin/tractor-agent local/bin/tractor studio/plugins/*/lib

setup: local/workspace local/bin qtalk studio
	make build

dev:
	./local/bin/tractor-agent --dev

kill:
	@killall node || true
	@killall tractor-agent || true

clean:
	rm -rf local/bin/tractor
	rm -rf local/bin/tractor-agent
	rm -rf studio/plugins/*/lib

clobber: clean
	rm -rf studio/node_modules	
	rm -rf studio/shell/lib
	rm -rf studio/shell/src-gen
	rm -rf studio/shell/node_modules
	rm -rf studio/shell/webpack.config.js
	rm -rf studio/extension/lib
	rm -rf studio/extension/node_modules

versions:
	@go version
	@echo "node $(shell node --version)"
	@git --version
	@echo "yarn $(shell yarn --version)"
	@echo "typescript $(shell tsc --version)"

qtalk:
	git submodule update --init --recursive
	make -C qtalk link

local/bin:
	mkdir -p local/bin

local/bin/tractor-agent: local/bin
	go build -o ./local/bin/tractor-agent ./cmd/tractor-agent

local/bin/tractor: local/bin
	go build -o ./local/bin/tractor ./cmd/tractor

local/workspace:
	mkdir -p local
	cp -r data/workspace local/workspace
	mv local/workspace/tractor.go.data local/workspace/tractor.go
	mkdir -p ~/.tractor/workspaces
	rm ~/.tractor/workspaces/dev || true
	ln -fs $(PWD)/local/workspace ~/.tractor/workspaces/dev

studio: studio/node_modules studio/extension/lib studio/shell/src-gen studio/plugins/*/lib

studio/node_modules:
	cd studio && yarn install
	cd studio && yarn link qmux qrpc

studio/extension/lib:
	cd studio/extension && yarn build

studio/plugins/%/lib:
	tsc -p $@/..

studio/shell/src-gen: studio/extension/lib
	cd studio/shell && yarn build
