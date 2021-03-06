import os
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

# Add special cases so we don't have to export EVAL to the core
initial_repl_env.set('eval', lambda ast: EVAL(ast, initial_repl_env))


def is_macro_call(ast, env):
    # Here
    if type(ast) is not list or len(ast) == 0:
        return False

    if type(ast[0]) is not MalSymbol:
        return False

    try:
        env_value = env.get(ast[0].value)
        if type(env_value) is not ResultingLambda:
            return False

        return env_value.is_macro
    except SymbolNotFound:
        return False


def macroexpand(ast, env):
    macro_expanded = False

    if is_macro_call(ast, env):
        macro_expanded = True

    while(is_macro_call(ast, env)):
        macro_function = env.get(ast[0].value)
        ast = macro_function.lambda_fn(*ast[1:])

    if utils.debug() and macro_expanded:
        print("Macro expanded >")
        utils.print_ast(ast)

    return ast


def is_pair(arg):
    # Here
    return (type(arg) is list or type(arg) is MalVector) and len(arg) > 0


def quasi_quote(ast):
    if not is_pair(ast):
        return [MalSymbol("quote"), ast]

    # List with length > 0
    head = ast[0]

    if type(head) is MalSymbol:
        function_symbol = head.value
        if function_symbol == "unquote":
            return ast[1]

    elif is_pair(ast[0]) and type(ast[0][0]) is MalSymbol and ast[0][0].value == "splice-unquote":
        return [MalSymbol("concat"), ast[0][1], quasi_quote(ast[1:])]

    return [MalSymbol("cons"), quasi_quote(ast[0]), quasi_quote(ast[1:])]


def READ(program):
    return reader.read_str(program)

# TODO Possibly refactor pr_str and this so we have
# Same switches? Or let the types themselves
# print? Unless pretty-printing needs extra capabilities

# Evaluation / lookup phase


def eval_ast(maltype, repl_env):
    if utils.debug():
        utils.prn("eval_ast >\n")
        utils.print_ast(maltype)

    # This is only here because of the call from regular application
    # at the far end of EVAL. Could possibly be moved and this
    # could be called eval-atom instead
    if type(maltype) is list:
        # This causes sub-expressions to be evaluated and applied
        # Eg (+ 1 (+ 2 3)< ) < that guy needs EVAL:ing to 5 first
        return list(map(lambda item: EVAL(item, repl_env), maltype))

    if type(maltype) is MalVector:
        return MalVector(list(map(lambda item: EVAL(item, repl_env), maltype)))

    if type(maltype) is dict:
        ret_val = {}
        for key, value in maltype.items():
            ret_val[key] = EVAL(value, repl_env)

        return ret_val

    if type(maltype) is MalSymbol:
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
    eval_invocation = 0

    while(True):
        if utils.debug():
            eval_invocation += 1
            print("EVAL invocation", eval_invocation)
            # TODO make this take a flag to print builtins (core functions)
            # or not
            # print(repr(repl_env))

            utils.print_ast(ast)

        if (type(ast) is not list):
            return eval_ast(ast, repl_env)

        if (type(ast) is list):
            # Empty list
            if not ast:
                return ast

            if type(ast[0]) is MalSymbol:

                ast = macroexpand(ast, repl_env)
                # print the expanded macro here

                # If macro-expansion has caused the list to not contain macros anymore
                if type(ast) is not list:
                    return eval_ast(ast, repl_env)

                ############################################################
                # Special forms
                ############################################################
                function_symbol = ast[0].value

                if function_symbol == "macroexpand":
                    return macroexpand(*ast[1:], repl_env)

                if function_symbol == "try*":
                    # If it's a native exception such as IndexError
                    # we need to translate it to some mal value such
                    # as a string "Index out of bounds"

                    # If it's a wrapped MalType exception, get the value
                    # (native type) of the exception

                    # E g (throw "fart") -> MalException.value = "fart"
                    # E g (throw nil) -> MalException.value = None
                    try:
                        return EVAL(ast[1], repl_env)
                    except (MalException, Exception) as exception:
                        if not ast[2]:
                            raise exception

                        catch_clause = ast[2]

                        exception_symbol = catch_clause[1].value
                        exception_block = catch_clause[2]

                        if type(exception) is MalException:
                            mal_exception_value = exception.value
                        else:
                            mal_exception_value = repr(exception)

                        repl_env.set(exception_symbol, mal_exception_value)
                        ast = exception_block
                        continue

                if function_symbol == "do":
                    # Here we could also add error checking that do has one
                    # argument - and that's a list

                    last = ast[-1]
                    if len(ast) > 2:
                        rest = ast[1:-1]
                        eval_ast(rest, repl_env)

                    ast = last
                    continue

                elif function_symbol == "if":
                    condition = ast[1]
                    condition_result = EVAL(condition, repl_env)

                    if (condition_result is None or
                            condition_result is False):
                        try:
                            else_branch = ast[3]

                            ast = else_branch
                            continue

                        # There is no else-branch
                        except IndexError:
                            # return None which is mapped to nil in mal
                            return None
                    else:
                        true_branch = ast[2]

                        ast = true_branch
                        continue

                # (def! a 3)
                elif function_symbol == "def!" or function_symbol == "defmacro!":
                    symbol = ast[1].value
                    value = EVAL(ast[2], repl_env)

                    if function_symbol == "defmacro!":
                        if type(value) is not ResultingLambda:
                            raise Exception('Macros must be functions')
                        value.is_macro = True

                    repl_env.set(symbol, value)
                    return value

                # (let* (a 2 c (+ 1 3)) (+ a c))
                elif function_symbol == "let*":
                    let_env = Env(repl_env)
                    let_definitions, body = ast[1], ast[2]

                    for wrapped_symbol, expression in utils.chunk_list(let_definitions, 2):
                        symbol = wrapped_symbol.value
                        value = EVAL(expression, let_env)
                        let_env.set(symbol, value)

                    ast = body
                    repl_env = let_env
                    continue

                elif function_symbol == "quote":
                    return ast[1]

                elif function_symbol == "quasiquote":
                    ast = quasi_quote(ast[1])

                    if (utils.debug()):
                        print("Quasiquiote:", ast)

                    continue

                elif function_symbol == "fn*":
                    if utils.debug():
                        print("fn* -> returning lambda")

                    lambda_params = list(map(lambda item: item.value, ast[1]))
                    lambda_body = ast[2]

                    def lambdar(*args):
                        lambda_env = Env(repl_env)
                        lambda_env.bind(lambda_params, list(args))

                        # If args > lambda_params there is an
                        # opportunity to return a partial application

                        # Add debugging here? Use it some

                        return EVAL(lambda_body, lambda_env)

                    resulting_lambda = ResultingLambda(lambdar,
                                                       lambda_body,
                                                       lambda_params,
                                                       repl_env)

                    return resulting_lambda

            # AST is a list here, and in it's first position
            # is either a symbol (which will be looked up in the env)
            # or a lambda definition ((fn* (a b) (+ a b) ))

            # The call to eval_ast here causes inner to outer
            # evaluation since the function is applied when all
            # other expressions have evalated -> if
            # ast is only primitives.

            # Regular function application (+ 1 2)
            # or a lambda ((fn* (a b) (+ a b) 1 2)
            if utils.debug():
                print("Function expansion")

            function_lambda, *arguments = eval_ast(ast, repl_env)

            if type(function_lambda) is ResultingLambda:
                if utils.debug():
                    print("Calling fn* lambda with arguments", list(arguments))

                lambda_env = Env(function_lambda.env)
                lambda_env.bind(function_lambda.lambda_params, list(arguments))

                ast = function_lambda.lambda_body
                repl_env = lambda_env
                continue

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

        if utils.debug():
            print("Result >")

        PRINT(result)
    except SymbolNotFound as exception:
        print("Symbol not found", exception)
    except ItemNotCallable as exception:
        print("Item not callable", exception)
    except IndexError as exception:
        print("Index out of bounds", exception)
    except Exception as exception:
        print("General exception caught", exception)


EVAL(READ("(def! load-file (fn* (f) (eval (read-string (str \"(do \" (slurp f) \")\")))))"), initial_repl_env)

this_files_dir = os.path.dirname(os.path.realpath(__file__))
stdlib_file_path = os.path.join(this_files_dir, "stdlib.mal")

EVAL(READ("(load-file \"" + stdlib_file_path + "\")"), initial_repl_env)
EVAL(READ("(def! *host-language* \"mymal\")"), initial_repl_env)

if __name__ == "__main__":
    readline.parse_and_bind("")

    parser = argparse.ArgumentParser()

    parser.add_argument('-d', '--debug', dest='debug', action='store_true')
    args, argv = parser.parse_known_args()

    if args.debug:
        utils.setDebug()

    initial_repl_env.set('*ARGV*', [])

    # Run file
    if (len(argv) >= 1):
        filename, *rest = argv
        initial_repl_env.set('*ARGV*', rest)
        EVAL(READ("(load-file \"" + filename + "\")"), initial_repl_env)
        exit(0)

    rep("(println (str \"Mal [\" *host-language* \"]\"))")

    # REPL
    while True:
        try:
            program = input("user> ").strip()

            if program == "env":
                print(repr(initial_repl_env))
            else:
                rep(program)
        except EOFError:
            print("\nBye!")
            sys.exit(0)
