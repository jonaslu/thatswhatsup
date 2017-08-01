// Stack

/*
PUSH_CONST 1
CALL 0
PUSH_CONST 2
PUSH_CONST 3
ADD
RETURN
ADD
PRINT
*/

/*
typedef enum {
  PUSH_CONST,
  CALL,
  ADD,
  RETURN,
  PRINT
} OP_CODE;
*/

#include <stdio.h>
#include <stdlib.h>

static const char PUSH_CONST = 1;
static const char CALL = 2;
static const char ADD = 3;
static const char RETURN = 4;
static const char PRINT = 5;

typedef struct {
  char* code;
  char* constPool;
} CodeBlock;

#define MAX_STACK_SLOTS 256
#define MAX_STACK_FRAMES 256

typedef struct {
  char slots[MAX_STACK_SLOTS];
  int currentStackSlot;

  CodeBlock* codeBlock;
  int ip;
} StackFrame;

// VM
typedef struct {
  StackFrame* stackFrames[MAX_STACK_FRAMES];

  int currentFrame;
} VM;

void pushStackValue(StackFrame* currentStackFrame, char value) {
  currentStackFrame->slots[currentStackFrame->currentStackSlot++] = value;
}

char popStackValue(StackFrame* currentStackFrame) {
  return currentStackFrame->slots[--currentStackFrame->currentStackSlot];
}

StackFrame* getCurrentStackFrame(VM* vm) {
  return vm->stackFrames[vm->currentFrame];
}

char getNextInstruction(StackFrame* currentStackFrame) {
  return currentStackFrame->codeBlock->code[currentStackFrame->ip++];
}

StackFrame* createStackFrame(CodeBlock *codeBlock) {
  StackFrame *stackFrame = malloc(sizeof(StackFrame));

  stackFrame->codeBlock = codeBlock;
  stackFrame->currentStackSlot = 0;
  stackFrame->ip = 0;
}

VM* createVM(StackFrame *baseStackFrame) {
  VM* vm = malloc(sizeof(VM));

  vm->stackFrames[0] = baseStackFrame;
  vm->currentFrame = 0;
}

void interpret(VM* vm) {
  StackFrame* currentStackFrame = getCurrentStackFrame(vm);

  int exitLoop = 0;
  do {
    char opcode = getNextInstruction(currentStackFrame);

    if (opcode == PUSH_CONST) {
      char constant = getNextInstruction(currentStackFrame);
      pushStackValue(currentStackFrame, constant);
    }

    if (opcode == ADD) {
      char value1 = popStackValue(currentStackFrame);
      char value2 = popStackValue(currentStackFrame);

      pushStackValue(currentStackFrame, value1 + value2);
    }

    if (opcode == PRINT) {
      char argument = popStackValue(currentStackFrame);
      printf("%d", argument);
    }

    if (opcode == RETURN) {
      return;
    }
  } while(!exitLoop);
}

void closeVM(VM* vm) {
  while(vm->currentFrame-- > 0) {
    StackFrame *stackFrame = getCurrentStackFrame(vm);

    free(stackFrame->codeBlock->code);
    free(stackFrame->codeBlock->constPool);
    free(stackFrame);
  }

  free(vm);
}

int main(void) {
  CodeBlock* codeBlock = malloc(sizeof(CodeBlock));

  char code[] = { 1, 1, 1, 2, 3, 5, 4 };
  codeBlock->code = &code;

  char constPool[] = { 1, 2 };
  codeBlock->constPool = &constPool;

  StackFrame *baseStackFrame = createStackFrame(codeBlock);
  VM* vm = createVM(baseStackFrame);

  interpret(vm);
  closeVM(vm);

  return 0;
}