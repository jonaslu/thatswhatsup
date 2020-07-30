ENTRY(loader) /* Name of the entry label in the asm file */

SECTIONS {
  . = 0x00100000; /* Load the code at 1MB (since GRUB, BIOS and mmap:ed I/O is below) */

  .text ALIGN (0x1000) : /* Align at 4 KB */
  {
    *(.text)
  }

  .rodata ALIGN (0x1000) :
  {
    *(.rodata*)
  }

  .data ALIGN (0x1000) :
  {
    *(.data)
  }

  .bss ALIGN (0x1000) :
  {
    *(COMMON)
    *(.bss)
  }
}
