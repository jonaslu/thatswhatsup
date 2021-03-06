from mal_types import *
import utils


def print_lists(native_value, print_readably, start_symbol="(", end_symbol=")"):
    return_expr = start_symbol
    return_expr += " ".join(list(map(lambda str: pr_str(str,
                                                        print_readably), native_value)))
    return_expr += end_symbol
    return return_expr


def pr_str(native_value, print_readably=True):
    if (type(native_value) is list):
        return print_lists(native_value, print_readably)

    if (type(native_value) is MalVector):
        return print_lists(native_value, print_readably, "[", "]")

    if native_value is None:
        return "nil"
    if native_value is True:
        return "true"
    if native_value is False:
        return "false"

    if type(native_value) is str:
        decoded_keyword = utils.decode_keyword(native_value)
        if decoded_keyword:
            return decoded_keyword

        # TODO make this reversible
        if print_readably:
            return (
                "\"" +
                native_value
                .replace("\\", "\\\\")
                .replace("\n", "\\n")
                .replace("\"", "\\\"") +
                "\""
            )
        else:
            return native_value

    if callable(native_value):
        return "#function"

    if type(native_value) is ResultingLambda:
        return "#function"

    if type(native_value) is MalSymbol:
        return native_value.value

    if type(native_value) is Atom:
        return "(atom " + str(native_value.value) + ")"

    if type(native_value) is dict:
        map_string_value = "{"
        key_values = []

        for key, value in native_value.items():
            key_values.append(pr_str(key, print_readably) + \
                " " + pr_str(value, print_readably))

        map_string_value += " ".join(key_values)
        map_string_value += "}"
        return map_string_value

    else:
        # Numbers are all that is left
        return str(native_value)
