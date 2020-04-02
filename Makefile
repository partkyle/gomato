all: package ;

build: build/gomato

build/gomato: *.go
	./bin/build

package: build/gomato.app ;

build/gomato.app: icons build plist
	./bin/createapp
plist: resources/Info.plist ;

icons: build/icon.icns
build/icon.icns: resources/mainicon.png
	./bin/mkicons resources/mainicon.png

clean:
	rm -rf build