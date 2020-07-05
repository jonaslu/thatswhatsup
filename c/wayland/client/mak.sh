#!/bin/bash
XDG_SHELL_XML=/usr/share/wayland-protocols/stable/xdg-shell/xdg-shell.xml

wayland-scanner client-header $XDG_SHELL_XML xdg-shell-client-protocol.h
wayland-scanner private-code $XDG_SHELL_XML xdg-shell-protocol.c

CAIRO_FLAGS=$(pkg-config pango -libs -cflags)
PANGO_FLAGS=$(pkg-config cairo -libs -cflags)

gcc -Werror -std=c99 -D_POSIX_C_SOURCE=200112L -D_DEFAULT_SOURCE -pedantic -Wno-unused-variable -lwayland-client $CAIRO_FLAGS $PANGO_FLAGS -o client client.c globalobj.c surface.c xdg-shell-protocol.c
