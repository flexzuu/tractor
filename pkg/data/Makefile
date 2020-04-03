.PHONY: icons editors

all: icons editors

icons:
	2goarray Available icons < icons/available.ico > icons/available.go
	2goarray Partially icons < icons/partially.ico > icons/partially.go
	2goarray Unavailable icons < icons/unavailable.ico > icons/unavailable.go
	2goarray Tractor icons < icons/tractor.ico > icons/tractor.go

editors:
	cd ../.. && go-bindata -pkg editors -o pkg/data/editors/bindata.go $(shell cd ../.. && find ./studio/editors -type d | xargs echo)