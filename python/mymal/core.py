import printer

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
    '>=': __gte
}

def import_core_functions(env):
    for func, lambdar in core_functions.items():
        env.set(func, lambdar)
