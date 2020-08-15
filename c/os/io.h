#ifndef __INCLUDE_IO_H
#define __INCLUDE_IO_H_

#include "gdt.h"

void outb(unsigned short port, unsigned short data);
unsigned char inb(unsigned short port);
void lgdt(unsigned int gdt);
void lidt(unsigned int idt);

#endif
