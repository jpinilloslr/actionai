BINARY   := actionai
CMD_PATH := ./cmd/cli
GNOME_EXT_UUID := actionai@jpinillos.dev.com
GNOME_EXT_SRC := ./extras/gnome-extension/
GNOME_EXT_DEST = $(HOME)/.local/share/gnome-shell/extensions/$(GNOME_EXT_UUID)

.PHONY: all build run clean install-gnome-ext test-gnome-ext uninstall-gnome-ext

all: build

build:
	go vet ./...
	go build -o bin/$(BINARY) $(CMD_PATH)
	cp ./actions.json ./bin/

install-gnome-ext:
	rm -rf $(GNOME_EXT_DEST)
	mkdir -p $(GNOME_EXT_DEST)
	cp -r $(GNOME_EXT_SRC)/* $(GNOME_EXT_DEST)/

test-gnome-ext:
	env MUTTER_DEBUG_DUMMY_MODE_SPECS=1920x1200 \
	dbus-run-session -- gnome-shell --nested --wayland

uninstall-gnome-ext:
	gnome-extensions disable $(GNOME_EXT_UUID)
	rm -rf $(GNOME_EXT_DEST)

run:
	go run $(CMD_PATH)

clean:
	rm -rf bin/


