import printer
import reader
from mal_types import *
from utils import decode_keyword, chunk_list


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
    if isinstance(item1, list) and isinstance(item2, list):
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


def __is_value(compare_value):
    return lambda value: value == compare_value


def __throw(message):
    raise MalException(message)


def __is_symbol(value):
    return type(value) is MalSymbol


def __symbol(value):
    return MalSymbol(value)


def __is_keyword(value):
    return type(value) is str and decode_keyword(value) is not None


def __keyword(value):
    if value[:1] == '\u029E':
        return value

    return "\u029E" + value


def __is_sequential(value):
    return type(value) is list or type(value) is MalVector


def __map(func, value):
    map_func = func

    if type(func) is ResultingLambda:
        map_func = func.lambda_fn

    return list(map(map_func, value))


def __apply(func, *rest):
    last_elem = rest[-1]
    but_last = list(rest[:-1])

    apply_func = func

    if type(func) is ResultingLambda:
        apply_func = func.lambda_fn

    all_args = but_last + last_elem

    return apply_func(*all_args)


def __get_hashmap(hashmap, key):
    if hashmap and type(hashmap) is dict:
        return hashmap[key]

    return None


def __assoc(hashmap, *keyvalues):
    ret_val = dict(hashmap)
    for key, value in chunk_list(keyvalues, 2):
        ret_val[key] = value

    return ret_val

def __dissoc(hashmap, *remove_keys):
    ret_val = dict(hashmap)
    for key in remove_keys:
        if key in ret_val:
            del ret_val[key]

    return ret_val

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
    'nil?': __is_value(None),
    'true?': __is_value(True),
    'false?': __is_value(False),
    'symbol?': __is_symbol,
    'symbol': __symbol,
    'keyword?': __is_keyword,
    'sequential?': __is_sequential,
    'keyword': __keyword,
    'map': __map,
    'apply': __apply,
    'vector?': lambda x: type(x) is MalVector,
    'map?': lambda x: type(x) is dict,
    'vector': lambda *x: MalVector(list(x)),
    'hash-map': lambda *values: dict(chunk_list(values, 2)),
    'get': __get_hashmap,
    'assoc': __assoc,
    'dissoc': __dissoc,
    'contains?': lambda hashmap, key: key in hashmap,
    'keys': lambda hashmap: list(hashmap.keys()),
    'vals': lambda hashmap: list(hashmap.values()),
    'split-str': lambda string, separator=None: string.split(separator)
}


def import_core_functions(env):
    for func, lambdar in core_functions.items():
        env.set(func, lambdar)
