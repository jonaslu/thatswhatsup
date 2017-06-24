import sys
from env import *
from mal_types import *
import reader
import printer
import utils

# Instantiate core function Env
repl_env = Env()
repl_env.set('+', lambda a, b: a + b)
repl_env.set('-', lambda a, b: a - b)
repl_env.set('*', lambda a, b: a * b)
repl_env.set('/', lambda a, b: int(a / b))

def READ(program):
    return reader.read_str(program)

# TODO Possibly refactor pr_str and this so we have
# Same switches? Or let the types themselves
# print? Unless pretty-printing needs extra capabilities

# Evaluation / lookup phase
def eval_ast(maltype, repl_env):
    if (type(maltype) is list):
        # This causes sub-expressions to be evaluated and applied
        # Eg (+ 1 (+ 2 3)< ) < that guy needs EVAL:ing to 5 first
        return list(map(lambda item: EVAL(item, repl_env), maltype))
    elif (type(maltype) is int):
        return maltype
    elif (type(maltype) is MalSymbol):
        symbol = maltype.value
        return repl_env.get(symbol)

# Apply phase if called with list, eval if not list
def EVAL(ast, repl_env):
    if (type(ast) is list):
        if not ast:
            return ast

        function_symbol = ast[0].value

        # (def! a 3)
        if function_symbol == "def!":
            symbol = ast[1].value
            value = EVAL(ast[2], repl_env)
            repl_env.set(symbol, value)
            return value

        # (let* (a 2 c (+ 1 3) (+ a c))
        elif function_symbol == "let*":
            let_env = Env(repl_env)
            let_definitions, body = ast[1], ast[2]
            for wrapped_symbol, expression in utils.chunk_list(let_definitions, 2):
                symbol = wrapped_symbol.value
                value = EVAL(expression, let_env)
                let_env.set(symbol, value)

            return EVAL(body, let_env)

        # Regular function application (+ 1 2)
        else:
            function_lambda, *arguments = eval_ast(ast, repl_env)
            return function_lambda(*arguments)
    else:
        return eval_ast(ast, repl_env)

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
        print(exception)

if __name__ == "__main__":
    while True:
        utils.prn("user> ")
        program = sys.stdin.readline().strip()
        if program:
            rep(program)
        else:
            print("\nBye!")
            sys.exit(0)
