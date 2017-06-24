import sys
import reader
import printer
import utils

def READ(program):
    return reader.read_str(program)

def EVAL(ast):
    return ast

def PRINT(result):
    print(printer.pr_str(result))

def rep(program):
    ast = READ(program)
    result = EVAL(ast)
    PRINT(result)

if __name__ == "__main__":
    while True:
        utils.prn("user> ")
        program = sys.stdin.readline().strip()
        if program:
            rep(program)
        else:
            print("\nBye!")
            sys.exit(0)
