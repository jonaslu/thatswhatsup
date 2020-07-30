global loader                   ; the entry symbol for ELF

MAGIC_NUMBER  equ 0x1BADB002
FLAGS         equ 0x0           ; multiboot flags
CHECKSUM      equ -MAGIC_NUMBER ; (magic number + checksum + flags should equal 0)

section .text                   ; write all of the constants to the data area - this has to do with the GRUB loader spec (methinks)
align 4
  dd MAGIC_NUMBER
  dd FLAGS
  dd CHECKSUM

loader:
  mov eax, 0xCAFEBABE
.loop:
  jmp .loop
