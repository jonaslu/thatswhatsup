import bytecode_ext as vm


def defineInnerFunction():
    codeBlock = vm.VMCodeBlock()
    codeBlock.setCode([vm.PUSH_CONST, 0,
                       vm.PUSH_CONST, 1,
                       vm.ADD,
                       vm.RETURN])

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
    rootCodeBlock.setCode([vm.PUSH_CONST, 0,
                           vm.CALL, 1,
                           vm.PUSH_CONST, 2,
                           vm.ADD,
                           vm.PRINT,
                           vm.RETURN])
    rootCodeBlock.setConstPool(
        [innerFunctionArgument, innerCodeBlock, localAddArgument])

    return rootCodeBlock


innerFunctionCodeBlock = defineInnerFunction()
rootCodeBlock = makeOuterFunction(innerFunctionCodeBlock)

vm.vm_run(rootCodeBlock)
