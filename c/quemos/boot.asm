BITS 32

section .text
global start
extern kmain

section .multiboot
    align 4
    dd 0x1BADB002               ; magic number
    dd 0x00                     ; flags
    dd - (0x1BADB002 + 0x00)    ; checksum

start:
    cli             ; clear any interupts

    mov ax, 0       ; clear all registers to avoid old values when kmain takes over
    mov ds, ax
    mov es, ax
    mov fs, ax
    mov gs, ax
    mov ss, ax

    mov esp, 0x90000

    call kmain

    hlt

section .bss
    resb 8192       ; stack size of 8192 bytes
