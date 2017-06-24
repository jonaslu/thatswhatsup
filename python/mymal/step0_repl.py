import sys

def prn(output):
    sys.stdout.write(output)
    sys.stdout.flush()

def READ(program):
    return program

def EVAL(ast):
    return ast

def PRINT(result):
    prn(result)

def rep(program):
    ast = READ(program)
    result = EVAL(ast)
    PRINT(result)

if __name__ == "__main__":
    while True:
        prn("user> ")
        program = sys.stdin.readline()
        if program:
            rep(program)
        else:
            print("\nBye!")
            sys.exit(0)
