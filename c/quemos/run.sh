nasm -f bin -o boot.bin loader.asm
qemu-system-x86_64 -drive format=raw,file=boot.bin
