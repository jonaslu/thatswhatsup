class MalSymbol:
    def __init__(self, symbol):
        self.value = symbol

    def __repr__(self):
        return "Symbol:" + self.value


class MalVector(list):
    pass


class Atom:
    def __init__(self, value):
        self.value = value


class ResultingLambda:
    def __init__(self, lambda_fn, lambda_body, lambda_params, lambda_env):
        self.lambda_fn = lambda_fn
        self.lambda_body = lambda_body
        self.lambda_params = lambda_params
        self.env = lambda_env

    def __repr__(self):
        return "Lambda body: " + str(self.lambda_body)
