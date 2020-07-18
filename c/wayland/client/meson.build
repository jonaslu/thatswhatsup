project('client',
  'c',
  default_options: [
    'c_std=c11',
    'warning_level=2',
    'werror=true',
  ])

src = [
  'client.c',
  'globalobj.c',
  'surface.c',
  'xdg-shell-protocol.c',
  'input.c',
  'pseudoterm.c',
]

build_args = ['-D_POSIX_C_SOURCE=200112L', '-D_DEFAULT_SOURCE' , '-pedantic', '-Wno-unused-variable']

pixman=dependency('pixman-1')
wlroots=dependency('wlroots')
wayland_client=dependency('wayland-client')
pango=dependency('pango')
cairo=dependency('cairo')
pangocairo=dependency('pangocairo')
xkbcommon=dependency('xkbcommon')

compiler=meson.get_compiler('c')
util=compiler.find_library('util')

project_deps=[pixman, wlroots, wayland_client, cairo, xkbcommon, util, pango, pangocairo]

executable('client', src, dependencies: project_deps, c_args: build_args)
