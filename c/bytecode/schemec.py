import vm.bytecode_ext as vm


def compile_program(program):
    int_value = int(program)

    rootCodeBlock = vm.VMRootCodeBlock()

    localAddIntArgument = vm.VMIntObject()
    localAddIntArgument.value = int_value

    rootCodeBlock.setCode([vm.PUSH_CONST, 0,
                           vm.PRINT,
                           vm.RETURN])

    rootCodeBlock.setConstPool([localAddIntArgument])

    return rootCodeBlock


def run_program(rootCodeBlock):
    vm.vm_run(rootCodeBlock)


program = "7"

compiled_program = compile_program(program)
run_program(compiled_program)
