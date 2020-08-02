; https://wiki.osdev.org/Text_Mode_Cursor

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
  mov ebx, 0x000B8000
  mov ecx, 0

.clear_screen:
  mov word [ebx], 0x0F20
  add ebx, 2
  inc ecx
  cmp ecx, 80*25
  jne .clear_screen

  ; Print A
  mov word [0x000B8000], 0x2841

  ; Show cursor
  mov dx, 0x3D4
  mov al, 0x0A
  out dx, al

  mov dx, 0x3D5
  mov al, 0x00
  out dx, al

  ; Move cursor
  mov dx, 0x3D4
  mov al, 14
  out dx, al

  mov dx, 0x3D5
  mov al, 0
  out dx, al

  mov dx, 0x3D4
  mov al, 15
  out dx, al

  mov dx, 0x3D5
  mov al, 80
  out dx, al
.loop:
  jmp .loop
