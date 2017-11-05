import bytecode_ext as vm

def defineInnerFunction():
  codeBlock = vm.VMCodeBlock()
  codeBlock.setCode([1,0,1,1,3,4])

  localAddIntArgument = vm.VMIntObject()
  localAddIntArgument.value = 2

  codeBlock.setConstPool([localAddIntArgument])
  codeBlock.numberArguments = 1

  return codeBlock

def makeOuterFunction(innerCodeBlock):
  rootCodeBlock = vm.VMRootCodeBlock()

  innerFunctionArgument = vm.VMIntObject()
  innerFunctionArgument.value = 2

  localAddArgument = vm.VMIntObject()
  localAddArgument.value = 1

  # PUSH_CONST 0
  rootCodeBlock.setCode([1,0,2,1,1,2,3,5,4])
  rootCodeBlock.setConstPool([innerFunctionArgument, innerCodeBlock, localAddArgument])

  return rootCodeBlock

innerFunctionCodeBlock = defineInnerFunction()
rootCodeBlock = makeOuterFunction(innerFunctionCodeBlock)

vm.vm_run(rootCodeBlock)
