import re
from mal_types import *
import utils

is_string_regex = re.compile("^\"(.*)\"$")


class NoNextToken(Exception):
    pass


class Reader:
    def __init__(self, tokens):
        self.tokens = tokens
        self.max_position = len(tokens) - 1
        self.end_reached = False
        self.position = 0

    def next(self):
        if(self.end_reached):
            raise NoNextToken

        curr_pos = self.position

        if(self.position == self.max_position):
            self.end_reached = True
        else:
            self.position += 1

        return self.tokens[curr_pos]

    def peek(self):
        return self.tokens[self.position]


def tokenizer(program):
    tokens = re.findall(
        r'''[\s,]*(~@|[\[\]{}()\\'`~^@]|"(?:\\.|[^\\"])*"|;.*|[^\s\[\]{}('"`,;)]*)''',
        program
    )

    tokens = filter(lambda s: s != '', tokens)
    tokens = filter(lambda s: s[0] != ';', tokens)

    return list(tokens)


def read_list(reader, end_list_token=")"):
    return_list = []

    while(reader.peek() != end_list_token):
        return_list.append(read_form(reader))

    reader.next()
    return return_list


def read_map(reader):
    return_map = {}

    key_value_list =  []

    while(reader.peek() != '}'):
        key_value_list.append(read_form(reader))

    # Skip }
    reader.next()

    key_values = utils.chunk_list(key_value_list, 2)
    for key, value in key_values:
        return_map[key] = value

    return return_map


def read_atom(reader):
    current_token = reader.next()

    if current_token == "true":
        return True

    if current_token == "false":
        return False

    if current_token == "nil":
        return None

    possible_keyword = utils.get_keyword(current_token)
    if possible_keyword:
        return possible_keyword

    is_string_match = is_string_regex.match(current_token)
    if is_string_match:
        unescaped_string = is_string_match.group(1)
        return (
            unescaped_string
            .replace("\\n", "\n")
            .replace("\\\\", "\\")
            .replace("\\\"", "\"")
        )

    try:
        return int(current_token)
    except ValueError:
        pass

    return MalSymbol(current_token)

# Forms are either (fun args)
# So a form is expected to have
# a function first OR
# 3 <- is a self evaluating
# form (special case of a symbol
# evaluating to itself. Whenever there
# is a list the fun goes first).


def read_form(reader):
    next_token = reader.peek()
    if(next_token == "("):
        reader.next()
        return read_list(reader)
    elif(next_token == "["):
        reader.next()
        return MalVector(read_list(reader, "]"))
    elif(next_token == '@'):
        reader.next()
        return [MalSymbol('deref'), read_form(reader)]
    elif(next_token == "{"):
        reader.next()
        return read_map(reader)
    elif(next_token == "'"):
        reader.next()
        return [MalSymbol("quote"), read_form(reader)]
    elif(next_token == "`"):
        reader.next()
        return [MalSymbol("quasiquote"), read_form(reader)]
    elif(next_token == "~"):
        reader.next()
        return [MalSymbol("unquote"), read_form(reader)]
    elif(next_token == "~@"):
        reader.next()
        return [MalSymbol("splice-unquote"), read_form(reader)]
    else:
        return read_atom(reader)


# Could add check for balanced parenthesis here
def read_str(program):
    if not program:
        print("Empty program")
    else:
        tokens = tokenizer(program)

        # If everything is a comment
        if not tokens:
            return None

        reader = Reader(tokens)
        try:
            ast = read_form(reader)
            return ast

        except NoNextToken:
            print("No next token found after: " + reader.peek())
