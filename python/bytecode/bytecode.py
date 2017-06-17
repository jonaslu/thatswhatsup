what_to_execute = {
  "instructions": [
    # a = 3
    # b = 4
    # print(x+y)
    ("LOAD_VALUE", 0),
    ("STORE_NAME", 0),
    ("LOAD_VALUE", 1),
    ("STORE_NAME", 1),
    ("LOAD_NAME", 0),
    ("LOAD_NAME", 1),
    ("ADD_TWO_VALUES", None),
    ("PRINT", None)
  ],
  "constants": [3, 4],
  "names": ["a", "b"]
 }

stack = []
environment = {}

for step in what_to_execute["instructions"]:
    instruction, value = step

    if instruction == "LOAD_VALUE":
        num = what_to_execute["constants"][value]
        stack.append(num)
    elif instruction == "LOAD_NAME":
        var_name = what_to_execute["names"][value]
        var_value = environment[var_name]
        stack.append(var_value)
    elif instruction == "STORE_NAME":
        var_name = what_to_execute["names"][value]
        environment[var_name] = stack.pop(0)
    elif instruction == "ADD_TWO_VALUES":
        op1, op2 = stack.pop(0), stack.pop(0)
        stack.append(op1 + op2)
    elif instruction == "PRINT":
        print(stack.pop(0))
