OUTPUT_FORMAT("elf32-i386")
ENTRY(start)

SECTIONS {
    . = 0x100000;  /* code starts at 1MB */

    .text : {
        *(.multiboot)  
        *(.text)
        . = ALIGN(4096);
    } :text

    .data : {
        *(.data)
        . = ALIGN(4096);
    } :data

    .bss : {
        *(.bss)
        . = ALIGN(4096);
    } :bss

    /DISCARD/ : {
        *(.comment)
        *(.eh_frame)
    }
}

PHDRS {
    text PT_LOAD FLAGS (1); /* Read and Execute */
    data PT_LOAD FLAGS (2); /* Read and Write */
    bss  PT_LOAD FLAGS (2); /* Read and Write, BSS might not be in the output */
}
