import test_struct

intObj = test_struct.VMIntObject()
intObj.value = 5

innerCodeBlock = test_struct.VMCodeBlock()
innerCodeBlock.code = [6,5,4,3,2,1,0]
innerCodeBlock.constPool = [intObj]

codeBlock = test_struct.VMRootCodeBlock()
codeBlock.setCode([1,2,3,4,5,0])
codeBlock.setConstPool([intObj, intObj, innerCodeBlock])

test_struct.vm_print_code_block(codeBlock)