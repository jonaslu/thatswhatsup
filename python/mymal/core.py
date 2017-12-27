import printer
import reader
from mal_types import *


def print_with_spaces(str):
    return " ".join(map(lambda x: printer.pr_str(x), str))


def __prn(*str):
    print(print_with_spaces(str))

    return None


def __prn_str(*str):
    return(print_with_spaces(str))


def __str(*str):
    return "".join(map(lambda x: printer.pr_str(x, False), str))


def __println(*strs):
    print(" ".join(map(lambda x: printer.pr_str(x, False), strs)))

    return None


def __list(*items):
    return [*items]


def __isList(lst):
    return True if (type(lst) is list) else False


def __isEmpty(lst):
    return True if not lst else False


def __count(lst):
    if isinstance(lst, list):
        return len(lst)

    return 0


def __equals(item1, item2):
    if item1 is list and item2 is list:
        if len(item1) == len(item2):
            for i in range(0, len(item1)):
                if not __equals(item1[i], item2[i]):
                    return False

            return True
        else:
            return False
    else:
        return item1 == item2


def __lt(a, b):
    return a < b


def __lte(a, b):
    return a <= b


def __gt(a, b):
    return a > b


def __gte(a, b):
    return a >= b


def __read_string(str):
    return reader.read_str(str)


def __slurp(filename):
    with open(filename, 'r') as myfile:
        return myfile.read()


def __atom(value):
    return Atom(value)


def __is_atom(value):
    return type(value) is Atom


def __deref(value):
    return value.value


def __reset(atom, value):
    atom.value = value
    return value


def __swap(atom, fn, *args):
    if type(fn) is ResultingLambda:
        new_value = fn.lambda_fn(atom.value, *args)
    else:
        new_value = fn(atom.value, *args)

    atom.value = new_value
    return new_value


def __cons(val1, val2):
    return [val1] + val2


def __concat(*lists):
    return sum(lists, [])


def __nth(lst, index):
    if index >= len(lst):
        raise IndexError

    return lst[index]


def __first(lst):
    if not lst:
        return None

    return lst[0]


def __rest(lst):
    if not lst:
        return []

    return lst[1:]


def __throw(message):
    raise Exception(message)


core_functions = {
    '+': lambda a, b: a + b,
    '-': lambda a, b: a - b,
    '*': lambda a, b: a * b,
    '/': lambda a, b: int(a / b),
    'prn': __prn,
    'pr-str': __prn_str,
    'str': __str,
    'println': __println,
    'list': __list,
    'list?': __isList,
    'empty?': __isEmpty,
    'count': __count,
    '=': __equals,
    '<': __lt,
    '<=': __lte,
    '>': __gt,
    '>=': __gte,
    'read-string': __read_string,
    'slurp': __slurp,
    'atom': __atom,
    'atom?': __is_atom,
    'deref': __deref,
    'reset!': __reset,
    'swap!': __swap,
    'cons': __cons,
    'concat': __concat,
    'nth': __nth,
    'first': __first,
    'rest': __rest,
    'throw': __throw,
}


def import_core_functions(env):
    for func, lambdar in core_functions.items():
        env.set(func, lambdar)
