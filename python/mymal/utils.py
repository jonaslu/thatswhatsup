import sys
import re

is_keyword_regex = re.compile("^:(.+)$")

def get_keyword(token):
    is_keyword_match = is_keyword_regex.match(token)
    if is_keyword_match:
        return "\u029E" + is_keyword_match.group(1)

    return None

def decode_keyword(token):
    if token[:1] == '\u029E':
        return ":" + token[1:]

    return None

def prn(output):
    sys.stdout.write(output)
    sys.stdout.flush()

def chunk_list(l, size):
    return [l[i:i + size] for i in range(0, len(l), size)]

class Debug:
    debug = False

def setDebug():
    Debug.debug = True

def debug():
    return Debug.debug
