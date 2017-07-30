#include <stdlib.h> // malloc, free
#include <stdio.h> // printf

// Object
typedef enum {
  TYPE_INT,
  TYPE_PAIR
} ObjectType;

typedef struct sObject {
  ObjectType type;
  unsigned char inUse;
  struct sObject* next;

  union {
    int value;

    struct {
      struct sObject* head;
      struct sObject* tail;
    };
  };
} Object;

#define MAX_STACKSIZE 256

typedef struct {
  Object* stack[MAX_STACKSIZE];

  Object* first;

  int stackSize;
  int objectsAllocated;
} VM;

VM* newVM() {
  VM* vm = malloc(sizeof(VM));
  vm->stackSize = 0;
  vm->objectsAllocated = 0;
  vm->first =  NULL;
  return vm;
}

Object* pushVM(VM* vm, Object* object) {
  vm->stack[vm->stackSize++] = object;
  return object;
}

Object* popVM(VM* vm) {
  return vm->stack[--vm->stackSize];
}

Object* newObject(VM* vm, ObjectType type) {
  Object* newObject = malloc(sizeof(Object));
  newObject->type = type;
  newObject->inUse = 0;

  newObject->next = vm->first;
  vm->first = newObject;

  vm->objectsAllocated++;

  return newObject;
}

Object* pushInt(VM* vm, int value) {
  Object* intObject = newObject(vm, TYPE_INT);
  intObject->value = value;

  return pushVM(vm, intObject);
}

Object* pushPair(VM* vm) {
  Object* pairObject = newObject(vm, TYPE_PAIR);

  pairObject->head = popVM(vm);
  pairObject->tail = popVM(vm);

  return pushVM(vm, pairObject);
}

void mark(Object* object) {
  if (object->inUse) return;

  object->inUse = 1;

  if (object->type == TYPE_PAIR) {
    mark(object->head);
    mark(object->tail);
  }
}

void markAll(VM* vm) {
  for(int i=0; i < vm->stackSize; i++) {
    mark(vm->stack[i]);
  }
}

void sweep(VM* vm) {
  Object** object = &vm->first;

  while (*object) {
    if (!(*object)->inUse) {
      Object* unreached = *object;

      // This threw me off some, it's a pointer to a pointer - so setting
      // what the first pointer points to equals to re-binding vm->first
      *object = unreached->next;
      free(unreached);
      vm->objectsAllocated--;
    } else {
      (*object)->inUse = 0;
      object = &(*object)->next;
    }
  }
}

void printPrimitive();

void printPair(Object* pair) {
   printf("head: ");
   printPrimitive(pair->head);

   printf("tail: ");
   printPrimitive(pair->tail);
}

void printPrimitive(Object* object) {
  if (object->type == TYPE_INT) {
    printf("%d\n", object->value);
  } else {
    printPair(object);
  }
}

void printGCStats(char* message, VM* vm) {
  printf("VM info - %s\n", message);
  printf("  stacksize: %d \n", vm->stackSize);
  printf("  objects allocated: %d\n", vm->objectsAllocated);
}

void gc(VM* vm) {
  markAll(vm);
  sweep(vm);
}

void freeVM(VM* vm) {
  vm->stackSize = 0;
  gc(vm);

  free(vm);
}

void main(void) {
  VM* vm = newVM();

  pushInt(vm, 0);
  pushInt(vm, 1);

  // First on stack [1, 0]
  Object* firstPair = pushPair(vm);

  // [1, 0], 2
  Object* int2 = pushInt(vm, 2);

  // [2, [1, 0]]
  Object* secondPair = pushPair(vm);

  printPrimitive(secondPair);

  // Should do nothing
  printGCStats("Before no change on stack", vm);
  gc(vm);
  printGCStats("After no change on stack", vm);

  // Clears up everything
  popVM(vm);

  gc(vm);
  printGCStats("After popping everything off the stack", vm);

  freeVM(vm);
}