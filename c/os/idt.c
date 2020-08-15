#include "io.h"
#include "framebuffer.h"
#include "serial.h"

extern void interrupt_handler_0();
extern void interrupt_handler_1();
extern void interrupt_handler_2();
extern void interrupt_handler_3();
extern void interrupt_handler_4();
extern void interrupt_handler_5();
extern void interrupt_handler_6();
extern void interrupt_handler_7();
extern void interrupt_handler_8();
extern void interrupt_handler_9();

extern void interrupt_handler_10();
extern void interrupt_handler_11();
extern void interrupt_handler_12();
extern void interrupt_handler_13();
extern void interrupt_handler_14();
extern void interrupt_handler_15();
extern void interrupt_handler_16();
extern void interrupt_handler_17();
extern void interrupt_handler_18();
extern void interrupt_handler_19();

extern void interrupt_handler_20();
extern void interrupt_handler_21();
extern void interrupt_handler_22();
extern void interrupt_handler_23();
extern void interrupt_handler_24();
extern void interrupt_handler_25();
extern void interrupt_handler_26();
extern void interrupt_handler_27();
extern void interrupt_handler_28();
extern void interrupt_handler_29();

extern void interrupt_handler_30();
extern void interrupt_handler_31();

struct idt_entry
{
  unsigned short base_lo;
  unsigned short segment;
  unsigned char always0;
  unsigned char flags;
  unsigned short base_hi;
} __attribute__((packed));

struct idt_ptr
{
  unsigned short size;
  unsigned int address;
} __attribute__((packed));

struct idt_entry idt[256];
struct idt_ptr idt_ptr;

static void init_idt_entry(unsigned char index, unsigned int base, unsigned short segment, unsigned char flags)
{
  idt[index].base_lo = base & 0xFFFF;
  idt[index].base_hi = (base >> 16) & 0xFFFF;

  idt[index].segment = segment;
  idt[index].always0 = 0;
  idt[index].flags = flags;
}

void init_idt()
{
  idt_ptr.size = (sizeof(struct idt_entry) * 256) - 1;
  idt_ptr.address = (unsigned int)&idt;

  // Clear out all of them
  for (int i = 0; i < 256; i++)
  {
    init_idt_entry(i, 0, 0, 0);
  }

  init_idt_entry(0, (unsigned int)interrupt_handler_0, 0x08, 0x8E);
  init_idt_entry(1, (unsigned int)interrupt_handler_1, 0x08, 0x8E);
  init_idt_entry(2, (unsigned int)interrupt_handler_2, 0x08, 0x8E);
  init_idt_entry(3, (unsigned int)interrupt_handler_3, 0x08, 0x8E);
  init_idt_entry(4, (unsigned int)interrupt_handler_4, 0x08, 0x8E);
  init_idt_entry(5, (unsigned int)interrupt_handler_5, 0x08, 0x8E);
  init_idt_entry(6, (unsigned int)interrupt_handler_6, 0x08, 0x8E);
  init_idt_entry(7, (unsigned int)interrupt_handler_7, 0x08, 0x8E);
  init_idt_entry(8, (unsigned int)interrupt_handler_8, 0x08, 0x8E);
  init_idt_entry(9, (unsigned int)interrupt_handler_9, 0x08, 0x8E);

  init_idt_entry(10, (unsigned int)interrupt_handler_10, 0x08, 0x8E);
  init_idt_entry(11, (unsigned int)interrupt_handler_11, 0x08, 0x8E);
  init_idt_entry(12, (unsigned int)interrupt_handler_12, 0x08, 0x8E);
  init_idt_entry(13, (unsigned int)interrupt_handler_13, 0x08, 0x8E);
  init_idt_entry(14, (unsigned int)interrupt_handler_14, 0x08, 0x8E);
  init_idt_entry(15, (unsigned int)interrupt_handler_15, 0x08, 0x8E);
  init_idt_entry(16, (unsigned int)interrupt_handler_16, 0x08, 0x8E);
  init_idt_entry(17, (unsigned int)interrupt_handler_17, 0x08, 0x8E);
  init_idt_entry(18, (unsigned int)interrupt_handler_18, 0x08, 0x8E);
  init_idt_entry(19, (unsigned int)interrupt_handler_19, 0x08, 0x8E);

  init_idt_entry(20, (unsigned int)interrupt_handler_20, 0x08, 0x8E);
  init_idt_entry(21, (unsigned int)interrupt_handler_21, 0x08, 0x8E);
  init_idt_entry(22, (unsigned int)interrupt_handler_22, 0x08, 0x8E);
  init_idt_entry(23, (unsigned int)interrupt_handler_23, 0x08, 0x8E);
  init_idt_entry(24, (unsigned int)interrupt_handler_24, 0x08, 0x8E);
  init_idt_entry(25, (unsigned int)interrupt_handler_25, 0x08, 0x8E);
  init_idt_entry(26, (unsigned int)interrupt_handler_26, 0x08, 0x8E);
  init_idt_entry(27, (unsigned int)interrupt_handler_27, 0x08, 0x8E);
  init_idt_entry(28, (unsigned int)interrupt_handler_28, 0x08, 0x8E);
  init_idt_entry(29, (unsigned int)interrupt_handler_29, 0x08, 0x8E);

  init_idt_entry(30, (unsigned int)interrupt_handler_30, 0x08, 0x8E);
  init_idt_entry(31, (unsigned int)interrupt_handler_31, 0x08, 0x8E);

  lidt((unsigned int)&idt_ptr);
}

struct registers
{
  unsigned int ds;
  unsigned int edi, esi, ebp, esp, ebx, edx, ecx, eax;
  unsigned int int_no, err_code;
  unsigned int eip, cs, eflags, useresp, ss;
};

void c_interrupt_handler(struct registers reg)
{
  fb_write_text("Received interrupt: ");
  fb_write_dec(reg.int_no);
  fb_write_text("\n");
}
