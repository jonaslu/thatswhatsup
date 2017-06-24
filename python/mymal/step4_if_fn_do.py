import sys
import argparse
import readline

from env import *
from mal_types import *
import core
import reader
import printer
import utils

class ItemNotCallable(Exception):
    pass

# Instantiate core function Env
initial_repl_env = Env()
core.import_core_functions(initial_repl_env)

def READ(program):
    return reader.read_str(program)

# TODO Possibly refactor pr_str and this so we have
# Same switches? Or let the types themselves
# print? Unless pretty-printing needs extra capabilities

# Evaluation / lookup phase
def eval_ast(maltype, repl_env):
    # This is only here because of the call from regular application
    # at the far end of EVAL. Could possibly be moved and this
    # could be called eval-atom instead
    if (type(maltype) is list):
        # This causes sub-expressions to be evaluated and applied
        # Eg (+ 1 (+ 2 3)< ) < that guy needs EVAL:ing to 5 first
        return list(map(lambda item: EVAL(item, repl_env), maltype))

    if (type(maltype) is MalVector):
        return MalVector(list(map(lambda item: EVAL(item, repl_env), maltype)))

    if (type(maltype) is MalSymbol):
        symbol = maltype.value
        return repl_env.get(symbol)

    # Catch-all for when the MalType is a native primitive already
    # int, True, False
    return maltype

def EVAL(ast, repl_env):
    # Primitive case - ast is either a native type or a MalSymbol
    # in which it's native "value" is returned - which is symbols
    # since true, false, nil, numbers and "" are returned
    # as python primitives
    if (type(ast) is not list):
        return eval_ast(ast, repl_env)

    if (type(ast) is list):
        # Empty list
        if not ast:
            return ast

        if type(ast[0]) is MalSymbol:
            function_symbol = ast[0].value

            ############################################################
            # Special forms
            ############################################################

            # (do (+ 1 2) (+ 3 4))
            if function_symbol == "do":
                # Here we could also add error checking that do has one
                # argument - and that's a list
                rest = ast[1:]
                return list(map(lambda item: EVAL(item, repl_env), rest))[-1]

            elif function_symbol == "if":
                condition = ast[1]
                condition_result = EVAL(condition, repl_env)

                if (condition_result is None or
                        condition_result is False):
                    try:
                        else_branch = ast[3]
                        return EVAL(else_branch, repl_env)
                    except IndexError:
                        # return nill which will internally evaluate to None
                        return None
                else:
                    true_branch = ast[2]
                    return EVAL(true_branch, repl_env)

            # (def! a 3)
            elif function_symbol == "def!":
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

            elif function_symbol == "fn*":
                lambda_params = list(map(lambda item: item.value, ast[1]))
                lambda_body = ast[2]

                def lambdar(*args):
                    lambda_env = Env(repl_env)
                    lambda_env.bind(lambda_params, list(args))

                    # If args > lambda_params there is an
                    # opportunity to return a partial application

                    return EVAL(lambda_body, lambda_env)

                return lambdar

        # AST is a list here, and in it's first position
        # is either a symbol (which will be looked up in the env)
        # or a lambda definition ((fn* (a b) (+ a b) ))

        # The call to eval_ast here causes inner to outer
        # evaluation since the function is applied when all
        # other expressions have evalated -> if
        # ast is only primitives.

        # Regular function application (+ 1 2)
        # or a lambda ((fn* (a b) (+ a b) 1 2)
        function_lambda, *arguments = eval_ast(ast, repl_env)

        if not callable(function_lambda):
            raise ItemNotCallable(ast[0])

        return function_lambda(*arguments)

def PRINT(result):
    print(printer.pr_str(result))

def rep(program):
    # TODO Move exception handling here, else we'll print result
    # even if read or eval phase fails
    try:
        ast = READ(program)
        result = EVAL(ast, initial_repl_env)
        PRINT(result)
    except SymbolNotFound as exception:
        print("Symbol not found", exception)
    except ItemNotCallable as exception:
        print("Item not callable", exception)

# Define not in mal itself
rep("(def! not (fn* (a) (if a false true)))")

if __name__ == "__main__":
    readline.parse_and_bind("")
    parser = argparse.ArgumentParser()
    parser.add_argument('-d', '--debug', dest='debug', action='store_true')
    args = parser.parse_args()

    if args.debug:
        utils.setDebug()

    while True:
        program = input("user> ").strip()
        if program:
            if program == "env":
                print(initial_repl_env)
            else:
                rep(program)
        else:
            print("\nBye!")
            sys.exit(0)
