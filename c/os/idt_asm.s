extern c_interrupt_handler

%macro no_error_code_handler 1
global interrupt_handler_%1
interrupt_handler_%1:
  cli
  push dword 0        ; Have error-code 0 when there is no error
  push dword %1       ; Push the interrupt number on the stack
  jmp common_interrupt_handler
%endmacro

%macro error_code_handler 1
global interrupt_handler_%1
interrupt_handler_%1:
  cli
  push dword %1       ; Only the interrupt number, error is already on the stack
  jmp common_interrupt_handler
%endmacro

common_interrupt_handler:
  pusha               ; push edi, esi, ebp, esp, ebx, edx, ecx, eax
  mov ax, ds          ; DS = current segment selector before interrupt, save it off to eax
  push eax

  mov ax, 0x10        ; kernel data-segment selector
  mov ds, ax
  mov es, ax
  mov fs, ax
  mov gs, ax

  call c_interrupt_handler

  pop eax             ; Move back the previous data-segment selector
  mov ds, ax
  mov es, ax
  mov fs, ax
  mov gs, ax

  popa                ; pop edi, esi, ebp, esp, ebx, edx, ecx, eax
  add esp, 8          ; Removed any pushed error-code (ours or theirs)
  sti                 ; set interrupt flag
  iret

no_error_code_handler 0
no_error_code_handler 1
no_error_code_handler 2
no_error_code_handler 3
no_error_code_handler 4
no_error_code_handler 5
no_error_code_handler 6
no_error_code_handler 7
error_code_handler 8
no_error_code_handler 9
error_code_handler 10
error_code_handler 11
error_code_handler 12
error_code_handler 13
error_code_handler 14
no_error_code_handler 15
no_error_code_handler 16
no_error_code_handler 17
no_error_code_handler 18
no_error_code_handler 19
no_error_code_handler 20
no_error_code_handler 21
no_error_code_handler 22
no_error_code_handler 23
no_error_code_handler 24
no_error_code_handler 25
no_error_code_handler 26
no_error_code_handler 27
no_error_code_handler 28
no_error_code_handler 29
no_error_code_handler 30
no_error_code_handler 31
