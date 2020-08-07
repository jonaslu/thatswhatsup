nasm -f elf32 loader.s
nasm -f elf32 io.s

gcc -m32 -nostdlib -nostdinc -fno-builtin -fno-stack-protector -nostartfiles -nodefaultlibs -Wall -Wextra -Werror -c -I. kmain.c -o kmain.o
gcc -m32 -nostdlib -nostdinc -fno-builtin -fno-stack-protector -nostartfiles -nodefaultlibs -Wall -Wextra -Werror -c -I. framebuffer.c -o framebuffer.o
gcc -m32 -nostdlib -nostdinc -fno-builtin -fno-stack-protector -nostartfiles -nodefaultlibs -Wall -Wextra -Werror -c -I. serial.c -o serial.o

ld -T link.ld -melf_i386 io.o loader.o framebuffer.o serial.o kmain.o -o kernel.elf

mkdir -p iso/boot/grub
cp stage2_eltorito iso/boot/grub
cp kernel.elf iso/boot
cp menu.lst iso/boot/grub

genisoimage -R \
  -b boot/grub/stage2_eltorito \
  -no-emul-boot \
  -boot-load-size 4 \
  -A os \
  -input-charset utf8 \
  -quiet \
  -boot-info-table \
  -o os.iso \
  iso
