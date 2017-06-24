class MalSymbol:
    def __init__(self, symbol):
        self.value = symbol

    def __repr__(self):
        return "Symbol:" + self.value

class MalVector(list):
    pass
