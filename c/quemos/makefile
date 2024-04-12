# Compiler and linker options
CC = gcc
AS = nasm
LD = ld
CFLAGS = -m32 -ffreestanding -nostdlib -nostartfiles -nodefaultlibs
LDFLAGS = -m elf_i386 -T linker.ld

# Source files
ASM_SRC = boot.asm
C_SRC = kmain.c

# Object files
ASM_OBJ = $(ASM_SRC:.asm=.o)
C_OBJ = $(C_SRC:.c=.o)

# Output file
OUTPUT = kernel.elf

all: $(OUTPUT)

$(OUTPUT): $(ASM_OBJ) $(C_OBJ)
	$(LD) $(LDFLAGS) $(ASM_OBJ) $(C_OBJ) -o $(OUTPUT)

%.o: %.asm
	$(AS) -f elf32 $< -o $@

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(ASM_OBJ) $(C_OBJ) $(OUTPUT)
