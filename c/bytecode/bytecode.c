#include <stdio.h>
#include <stdlib.h>

static const char PUSH_CONST = 1;
static const char CALL = 2;
static const char ADD = 3;
static const char RETURN = 4;
static const char PRINT = 5;

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

#define MAX_STACK_SLOTS 256
#define MAX_STACK_FRAMES 256

typedef struct {
  Object* slots[MAX_STACK_SLOTS];
  int currentStackSlot;

  CodeBlock* codeBlock;
  int ip;
} StackFrame;

typedef struct {
  StackFrame* stackFrames[MAX_STACK_FRAMES];

  int currentFrame;
} VM;

static void pushStackValue(StackFrame* currentStackFrame, Object* value) {
  currentStackFrame->slots[currentStackFrame->currentStackSlot++] = value;
}

static Object* popStackValue(StackFrame* currentStackFrame) {
  return currentStackFrame->slots[--currentStackFrame->currentStackSlot];
}

static StackFrame* getCurrentStackFrame(VM* vm) {
  return vm->stackFrames[vm->currentFrame];
}

static unsigned char getNextInstruction(StackFrame* currentStackFrame) {
  return currentStackFrame->codeBlock->code[currentStackFrame->ip++];
}

static Object* getConstantPoolValue(StackFrame* currentStackFrame) {
  unsigned char index = getNextInstruction(currentStackFrame);
  return currentStackFrame->codeBlock->constPool[index];
}

static StackFrame* createStackFrame(CodeBlock *codeBlock) {
  StackFrame *stackFrame = malloc(sizeof(StackFrame));

  stackFrame->codeBlock = codeBlock;
  stackFrame->currentStackSlot = 0;
  stackFrame->ip = 0;

  return stackFrame;
}

static VM* createVM(StackFrame *baseStackFrame) {
  VM* vm = malloc(sizeof(VM));

  vm->stackFrames[0] = baseStackFrame;
  vm->currentFrame = 0;

  return vm;
}

static void pushStackFrame(VM* vm, StackFrame* nextStackFrame) {
  vm->stackFrames[++vm->currentFrame] = nextStackFrame;
}

static StackFrame* popStackFrame(VM* vm) {
  return vm->stackFrames[--vm->currentFrame];
}

static void interpret(VM* vm) {
  StackFrame* currentStackFrame = getCurrentStackFrame(vm);

  for(;;) {
    unsigned char opcode = getNextInstruction(currentStackFrame);

    if (opcode == PUSH_CONST) {
      Object* constObject = getConstantPoolValue(currentStackFrame);

      pushStackValue(currentStackFrame, constObject);
    }

    if (opcode == ADD) {
      Object* value1 = popStackValue(currentStackFrame);
      Object* value2 = popStackValue(currentStackFrame);

      if (value1->type == TYPE_INT && value2->type == TYPE_INT) {
        Object* addedValue = malloc(sizeof(Object));

        addedValue->type = TYPE_INT;
        addedValue->value = value1->value + value2->value;

        pushStackValue(currentStackFrame, addedValue);
      }
      // TODO - error on non-addable types
    }

    if (opcode == PRINT) {
      Object* argument = popStackValue(currentStackFrame);

      // TODO Switch on type of argument when printing
      printf("%d", argument->value);
    }

    if (opcode == CALL) {
      Object* codeObject = getConstantPoolValue(currentStackFrame);

      // TODO Assert that const pool object is of TYPE_CODE

      CodeBlock* subRoutineCodeBlock = codeObject->codeBlock;
      StackFrame* nextStackFrame = createStackFrame(codeObject->codeBlock);

      // Put arguments on the far end of the const pool
      for (int i = 0; i < subRoutineCodeBlock->numberArguments; i++) {
        Object* subroutineArgument = popStackValue(currentStackFrame);

        nextStackFrame->codeBlock->constPool[nextStackFrame->codeBlock->numberConsts + i] = subroutineArgument;
      }

      pushStackFrame(vm, nextStackFrame);
      currentStackFrame = nextStackFrame;
    }

    if (opcode == RETURN) {
      if (vm->currentFrame <= 0) {
        free(currentStackFrame);
        return;
      }

      StackFrame* previousStackFrame = popStackFrame(vm);

      Object* returnValue = popStackValue(currentStackFrame);
      pushStackValue(previousStackFrame, returnValue);

      free(currentStackFrame);

      currentStackFrame = previousStackFrame;
    }
  }
}

static void closeVM(VM* vm) {
  free(vm);
}

// Main callable from outside
void run(CodeBlock* rootCodeBlock) {
  StackFrame* baseStackFrame = createStackFrame(rootCodeBlock);
  VM* vm = createVM(baseStackFrame);

  interpret(vm);
  closeVM(vm);
}

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

static void print_code(unsigned char* test) {
  printf("Code:");

  int i=0;
  while (test[i] != RETURN) {
    printf(" %d", test[i++]);
  }
  printf("\n");
}

// Debugging, prints code
void print_code_block(CodeBlock* codeBlock) {
  printf("Num consts %d\n", codeBlock->numberConsts);
  printf("Num arguments %d\n", codeBlock->numberArguments);

  print_code(codeBlock->code);

  for(int i=0; i<codeBlock->numberConsts; i++) {
    printf("\nConst pool index %d\n", i);
    printType(codeBlock->constPool[i]);
  }
 }
