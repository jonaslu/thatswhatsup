method_add = {
    "code": [
      # func add(x,y):
      # return x + y
      # STORE_NAME 0
      # STORE_NAME 1
      # LOAD_NAME 0
      # LOAD_NAME 1
      # ADD_TWO_VALUES
      # RET
      ("STORE_NAME", 0),
      ("STORE_NAME", 1),
      ("LOAD_NAME", 0),
      ("LOAD_NAME", 1),
      ("ADD_TWO_VALUES", None),
      ("RET", None)
    ],
    "constants": [],
    "names": ["x", "y"],
    "args": 2
 }

method_main = {
    "code": [
        # a = 3
        # b = 4
        # print(add(a, b))
        ("LOAD_VALUE", 0),
        ("STORE_NAME", 0),
        ("LOAD_VALUE", 1),
        ("STORE_NAME", 1),
        ("LOAD_NAME", 0),
        ("LOAD_NAME", 1),
        ("CALL", 2),
        ("PRINT", None)
    ],
    "constants": [3, 4, method_add],
    "names": ["a", "b"],
    "args": 0
}


class Frame:
    def __init__(self, code_block):
        self.code_block = code_block
        self.stack = []
        self.environment = {}

    def run(self):
        for step in self.code_block["code"]:
            instruction, value = step

            if instruction == "LOAD_VALUE":
                num = self.code_block["constants"][value]
                self.stack.append(num)
            elif instruction == "LOAD_NAME":
                var_name = self.code_block["names"][value]
                var_value = self.environment[var_name]
                self.stack.append(var_value)
            elif instruction == "STORE_NAME":
                var_name = self.code_block["names"][value]
                self.environment[var_name] = self.stack.pop(0)
            elif instruction == "ADD_TWO_VALUES":
                op1, op2 = self.stack.pop(0), self.stack.pop(0)
                self.stack.append(op1 + op2)
            elif instruction == "PRINT":
                print(self.stack.pop(0))
            elif instruction == "CALL":
                code_block = self.code_block["constants"][value]
                next_frame = Frame(code_block)

                next_frame.stack = self.stack[-2:]
                self.stack = self.stack[:-2]

                next_frame.run()
                if len(next_frame.stack) > 0:
                    self.stack.append(next_frame.stack[0])
            elif instruction == "RET":
                break


main_frame = Frame(method_main)
main_frame.run()
