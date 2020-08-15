global outb                     ; method to send data to I/O port
global inb
global lgdt
global lidt

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

; Read a byte from an I/O port
; [esp + 4] <- first argument is the port to read from
; ret the byte read
inb:
  mov dx, [esp + 4]
  in al, dx
  ret

lgdt:
  mov eax, [esp + 4]
  lgdt [eax]

  mov ax, 0x10                ; kernel data-segment selector
  mov ds, ax
  mov es, ax
  mov fs, ax
  mov gs, ax
  mov ss, ax
  jmp 0x08:.flush
.flush:
  ret

lidt:
  mov eax, [esp + 4]
  lidt [eax]
  ret
