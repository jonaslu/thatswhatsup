from cpython.mem cimport PyMem_Malloc, PyMem_Free
from cpython cimport array
import array
from sys import exit

cdef extern from "bytecode.c":
  cdef unsigned int _RETURN "RETURN"
  cdef unsigned int _PUSH_CONST "PUSH_CONST"
  cdef unsigned int _CALL "CALL"
  cdef unsigned int _ADD "ADD"
  cdef unsigned int _PRINT "PRINT"

  ctypedef struct CodeBlock:
    unsigned char* code
    Object** constPool

    int numberConsts
    int numberArguments

  ctypedef enum Type:
    TYPE_INT,
    TYPE_CODE

  ctypedef struct Object:
    Type type

    int value
    CodeBlock* codeBlock

  cdef void print_code_block(CodeBlock* codeBlock)

  cdef void run(CodeBlock* codeBlock)


cdef class VMRootCodeBlock:
  """The main function should use this. This does not hold an Object* reference
     to itsef and thusly cannot be added into the constPool of another CodeBlock"""
  cdef CodeBlock codeBlock
  cdef array.array codePtr

  def __setattr__(self, name, value):
    if name == 'code':
      self.setCode(value)
    elif name == 'constPool':
      self.setConstPool(value)
    elif name == 'numberArguments':
      self.codeBlock.numberArguments = value
    else:
      super().__setattr__(name, value)

  def setCode(self, value):
    self.codePtr = array.array('B', value)
    self.codeBlock.code = self.codePtr.data.as_uchars

  def __dealloc__(self):
    if self.codeBlock.constPool:
      PyMem_Free(self.codeBlock.constPool)

  def setConstPool(self, val):
    numberOfConsts = len(val)

    self.codeBlock.numberConsts = numberOfConsts
    self.codeBlock.constPool = <Object**> PyMem_Malloc(numberOfConsts * sizeof(Object*))

    for i in range(numberOfConsts):
      object = val[i]

      if isinstance(object, VMIntObject):
        self.__setIntValue(object, i)
      elif isinstance(object, VMCodeBlock):
        self.__setConstPoolCodeBlock(object, i)
      else:
        print("Unknown type", object)
        exit(1)

  def __setConstPoolCodeBlock(self, VMCodeBlock vmCodeBlock, index):
    self.codeBlock.constPool[index] = &vmCodeBlock.object

  def __setIntValue(self, VMIntObject intValue, index):
    self.codeBlock.constPool[index] = &intValue.object


cdef class VMCodeBlock(VMRootCodeBlock):
  """This holds a reference to itself via an Object* and
     this is what's added when this codeblock is inserted
     into another CodeBlocks const pool"""
  cdef Object object

  def __cinit__(self):
    self.object.type = TYPE_CODE
    self.object.codeBlock = &self.codeBlock


cdef class VMIntObject:
  cdef Object object

  def __cinit__(self):
    self.object.type = TYPE_INT

  def __setattr__(self, name, val):
    if name == 'value':
      self.object.value = val

def vm_print_code_block(VMRootCodeBlock cb):
  print_code_block(&cb.codeBlock)

def vm_run(VMRootCodeBlock cb):
  run(&cb.codeBlock)

RETURN = _RETURN
PUSH_CONST = _PUSH_CONST
CALL = _CALL
ADD = _ADD
PRINT = _PRINT
