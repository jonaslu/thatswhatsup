global outb                     ; method to send data to I/O port

; Write a byte to an I/O port
; Usage in C - outb(0x3DA, 14);
; [esp + 8] <- second argument is the value
; [esp + 4] <- first argument is the I/O port
; [esp] <- return address
outb:
  mov al, [esp + 8]
  mov dx, [esp + 4]
  out dx, al
  ret
