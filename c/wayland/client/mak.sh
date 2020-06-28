XDG_SHELL_XML=/usr/share/wayland-protocols/stable/xdg-shell/xdg-shell.xml

wayland-scanner client-header $XDG_SHELL_XML xdg-shell-client-protocol.h
wayland-scanner private-code $XDG_SHELL_XML xdg-shell-protocol.c

gcc -Werror -O2 -std=c99 -pedantic -Wno-unused-variable -lwayland-client -o connect connect.c xdg-shell-protocol.c
