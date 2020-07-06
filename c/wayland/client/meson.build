project('client',
  'c',
  default_options: [
    'c_std=c11',
    'warning_level=2',
    'werror=true',
  ])

src = ['client.c', 'globalobj.c', 'surface.c', 'xdg-shell-protocol.c', 'text.c']

build_args = ['-D_POSIX_C_SOURCE=200112L', '-D_DEFAULT_SOURCE' , '-pedantic', '-Wno-unused-variable']

pixman=dependency('pixman-1')
wlroots=dependency('wlroots')
wayland_client=dependency('wayland-client')
# pango=dependency('pango')
# pangocairo=dependency('pangocairo')
cairo=dependency('cairo')

project_deps=[pixman, wlroots, wayland_client, cairo]

executable('client', src, dependencies: project_deps, c_args: build_args)
