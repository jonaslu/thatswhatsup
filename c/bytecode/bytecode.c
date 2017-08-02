#include <stdio.h>
#include <stdlib.h>

static const char PUSH_CONST = 1;
static const char CALL = 2;
static const char ADD = 3;
static const char RETURN = 4;
static const char PRINT = 5;

struct Object;

typedef struct {
  char* code;
  struct Object** constPool;

  int numberArguments;
} CodeBlock;

typedef enum {
  TYPE_INT,
  TYPE_CODE
} Type;

typedef struct {
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

void pushStackValue(StackFrame* currentStackFrame, Object* value) {
  currentStackFrame->slots[currentStackFrame->currentStackSlot++] = value;
}

Object* popStackValue(StackFrame* currentStackFrame) {
  return currentStackFrame->slots[--currentStackFrame->currentStackSlot];
}

StackFrame* getCurrentStackFrame(VM* vm) {
  return vm->stackFrames[vm->currentFrame];
}

char getNextInstruction(StackFrame* currentStackFrame) {
  return currentStackFrame->codeBlock->code[currentStackFrame->ip++];
}

Object* getConstantPoolValue(StackFrame* currentStackFrame) {
  char index = getNextInstruction(currentStackFrame);
  return currentStackFrame->codeBlock->constPool[index];
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

void pushStackFrame(VM* vm, StackFrame* nextStackFrame) {
  vm->stackFrames[++vm->currentFrame] = nextStackFrame;
}

StackFrame* popStackFrame(VM* vm) {
  return vm->stackFrames[--vm->currentFrame];
}

void interpret(VM* vm) {
  StackFrame* currentStackFrame = getCurrentStackFrame(vm);

  for(;;) {
    char opcode = getNextInstruction(currentStackFrame);

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
      printf("%d", argument->value);
    }

    if (opcode == CALL) {
      Object* constObject = getConstantPoolValue(currentStackFrame);

      // TODO Assert that const pool object is of TYPE_CODE

      CodeBlock* subRoutineCodeBlock = constObject->codeBlock;
      StackFrame* nextStackFrame = createStackFrame(constObject->codeBlock);

      for (int i = 0; i < subRoutineCodeBlock->numberArguments; i++) {
        Object* subroutineArgument = popStackValue(currentStackFrame);
        pushStackValue(nextStackFrame, subroutineArgument);
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

void closeVM(VM* vm) {
  // TODO Invoke garbage collection
  free(vm);
}

int main(void) {
  /*
    // Stack frame now contains first argument (2)
    int raj(a) {
      return a + 2
    }

    PUSH_CONST 0
    ADD
    RETURN
  */
  CodeBlock* subCodeBlock = malloc(sizeof(CodeBlock));

  char subCode[] = { 1, 0, 3, 4 };
  Object subConstPool[] = { { TYPE_INT, 2 } };

  subCodeBlock->numberArguments = 1;
  subCodeBlock->code = &subCode;
  subCodeBlock->constPool = malloc(sizeof(Object*));
  subCodeBlock->constPool[0] = &subConstPool[0];

  /*
    print(1 + raj(2))

    PUSH_CONST 0
    CALL 1

    PUSH_CONST 2
    ADD
    PRINT
    RETURN
  */

  CodeBlock* outerCodeBlock = malloc(sizeof(CodeBlock));
  char code[] = { 1, 0, 2, 1, 1, 2, 3, 5, 4 };
  outerCodeBlock->code = &code;

  Object outerConstPool[] = { { TYPE_INT, 2 }, { TYPE_CODE, 0 }, { TYPE_INT, 1} };
  outerConstPool[1].codeBlock = subCodeBlock;

  outerCodeBlock->constPool = malloc(sizeof(Object*) * 3);

  outerCodeBlock->constPool[0] = &outerConstPool[0];
  outerCodeBlock->constPool[1] = &outerConstPool[1];
  outerCodeBlock->constPool[2] = &outerConstPool[2];

  Object* test = outerCodeBlock->constPool[0];

  StackFrame* baseStackFrame = createStackFrame(outerCodeBlock);
  VM* vm = createVM(baseStackFrame);

  interpret(vm);
  closeVM(vm);

  return 0;
}