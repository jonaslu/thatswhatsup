class MalSymbol:
    def __init__(self, symbol):
        self.value = symbol

    def __repr__(self):
        return "Symbol:" + self.value

    def __eq__(self, other):
        if type(other) is MalSymbol:
            return self.value == other.value

        return False

class MalVector(list):
    def __init__(self, *args):
        super(MalVector, self).__init__(*args)
        self.meta = None

class MalList(list):
    def __init__(self, *args):
        super(MalList, self).__init__(*args)
        self.meta = None

class Atom:
    def __init__(self, value):
        self.value = value


class ResultingLambda:
    def __init__(self, lambda_fn, lambda_body, lambda_params, lambda_env):
        self.lambda_fn = lambda_fn
        self.lambda_body = lambda_body
        self.lambda_params = lambda_params
        self.env = lambda_env
        self.is_macro = False
        self.meta = None

    def __repr__(self):
        return "Lambda body: " + str(self.lambda_body)


class MalException(Exception):
    def __init__(self, value):
        self.value = value
