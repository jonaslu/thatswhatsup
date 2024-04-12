global start
section .text
start:
    mov ah, 0x0e ; 0x0e is BIOS text mode

    mov al, 'H'
    int 0x10
    mov al, 'e'
    int 0x10
    mov al, 'l'
    int 0x10
    mov al, 'l'
    int 0x10
    mov al, 'o'
    int 0x10

    mov al, 'W'
    int 0x10
    mov al, 'o'
    int 0x10
    mov al, 'r'
    int 0x10
    mov al, 'l'
    int 0x10
    mov al, 'd'
    int 0x10

    jmp $

;; Fills up the rest of the binary with 0:s making
;; it 512 bytes, and the last two bytes are set to 0xaa55
;; signalling it's bootable to bootloaders
times 510-($-$$) db 0
dw 0xaa55
