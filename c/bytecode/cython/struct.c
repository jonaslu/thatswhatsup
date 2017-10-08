#include <stdio.h>

struct Object;

typedef struct {
  unsigned char* code;
  struct Object** constPool;

  int numberConsts;
  int numberArguments;
} CodeBlock;

typedef enum {
   TYPE_INT,
   TYPE_CODE
 } Type;

 typedef struct Object {
   Type type;

   union {
     int value;
     CodeBlock* codeBlock;
   };
 } Object;


static void printIntType(Object* intType) {
  printf("TYPE_INT\n");
  printf("Int value %d\n", intType->value);
}

void print_code_block(CodeBlock* codeBlock);

static void printCodeType(Object* codeType) {
  printf("TYPE_CODE\n");
  print_code_block(codeType->codeBlock);
}

static void printType(Object* object) {
  switch(object->type) {
    case TYPE_INT: printIntType(object); return;
    case TYPE_CODE: printCodeType(object); return;
  }
}

void print_code(unsigned char* test) {
  printf("Code:");

  int i=0;
  while (test[i] != 0) {
    printf(" %d", test[i++]);
  }
  printf("\n");
}

void print_code_block(CodeBlock* codeBlock) {
  printf("Num consts %d\n", codeBlock->numberConsts);
  printf("Num arguments %d\n", codeBlock->numberArguments);

  print_code(codeBlock->code);

  for(int i=0; i<codeBlock->numberConsts; i++) {
    printf("\nConst pool index %d\n", i);
    printType(codeBlock->constPool[i]);
  }
 }
