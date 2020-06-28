# set -x
XDG_SHELL_XML=/usr/share/wayland-protocols/stable/xdg-shell/xdg-shell.xml

wayland-scanner server-header $XDG_SHELL_XML xdg-shell-protocol.h
wayland-scanner private-code $XDG_SHELL_XML xdg-shell-protocol.c

WLROOTS_FLAGS=$(pkg-config wlroots -libs -cflags)
WAYLAND_FLAGS=$(pkg-config wayland-server -libs -cflags)

gcc --std=c11 -Werror -Wall -D_POSIX_C_SOURCE=200112L -DWLR_USE_UNSTABLE -lrt $WLROOTS_FLAGS $WAYLAND_FLAGS -I. -o server xdg-shell-protocol.c main.c
# set +x
