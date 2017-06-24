import sys
from mal_types import *
import reader
import printer
import utils

# Possibly refactor out to a message exception base class
class SymbolNotFound(Exception):
    def __init__(self, message):
        self.message = message

    def __repr__(self):
        return message

repl_env = {'+': lambda a,b: a+b,
            '-': lambda a,b: a-b,
            '*': lambda a,b: a*b,
            '/': lambda a,b: int(a/b)}

def READ(program):
    return reader.read_str(program)

# TODO Possibly refactor pr_str and this so we have
# Same switches? Or let the types themselves
# print? Unless pretty-printing needs extra capabilities

# Evaluation / lookup phase
def eval_ast(maltype, repl_env):
    if (type(maltype) is list):
        return list(map(lambda item: EVAL(item, repl_env), maltype))
    elif (type(maltype) is int):
        return maltype
    elif (type(maltype) is MalSymbol):
        try:
            return repl_env[maltype.value]
        except KeyError:
            raise SymbolNotFound(maltype.value)

# Apply phase
def EVAL(maltype, repl_env):
    if (type(maltype) is list):
        if (not maltype):
            return maltype

        func, *args = eval_ast(maltype, repl_env)
        # Here is the magic! The function goes first
        # and the arguments are applied to the func
        return func(*args)
    else:
        return eval_ast(maltype, repl_env)

def PRINT(result):
    print(printer.pr_str(result))

def rep(program):
    # TODO Move exception handling here, else we'll print result
    # even if read or eval phase fails
    try:
        ast = READ(program)
        result = EVAL(ast, repl_env)
        PRINT(result)
    except SymbolNotFound as exception:
        print("Symbol not found", exception)

if __name__ == "__main__":
    while True:
        utils.prn("user> ")
        program = sys.stdin.readline().strip()
        if program:
            rep(program)
        else:
            print("\nBye!")
            sys.exit(0)
