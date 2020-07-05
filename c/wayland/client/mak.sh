XDG_SHELL_XML=/usr/share/wayland-protocols/stable/xdg-shell/xdg-shell.xml

wayland-scanner client-header $XDG_SHELL_XML xdg-shell-client-protocol.h
wayland-scanner private-code $XDG_SHELL_XML xdg-shell-protocol.c

gcc -Werror -O2 -std=c99 -D_POSIX_C_SOURCE=200112L -D_DEFAULT_SOURCE -pedantic -Wno-unused-variable -lwayland-client -o client client.c globalobj.c surface.c xdg-shell-protocol.c
