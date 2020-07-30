global loader                   ; the entry symbol for ELF
extern kmain

MAGIC_NUMBER        equ 0x1BADB002
FLAGS               equ 0x0           ; multiboot flags
CHECKSUM            equ -MAGIC_NUMBER ; (magic number + checksum + flags should equal 0)
KERNEL_STACK_SIZE   equ 4096

section .bss
align 4
kernel_stack:                         ; label points to beginning of memory
  resb KERNEL_STACK_SIZE              ; reserve stack for the kernel

section .text                   ; write all of the constants to the data area - this has to do with the GRUB loader spec (methinks)
align 4
  dd MAGIC_NUMBER
  dd FLAGS
  dd CHECKSUM

loader:
  mov esp, kernel_stack + KERNEL_STACK_SIZE   ; point esb to the start of the stack (end of memory area)
  call kmain
