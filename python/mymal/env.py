import printer

# Possibly refactor out to a message exception base class


class SymbolNotFound(Exception):
    def __init__(self, symbol):
        self.symbol = symbol

    def __repr__(self):
        return "'" + self.symbol + "' not found"


class Env:
    def __init__(self, outer=None):
        self.outer = outer
        self.data = {}

    def __repr__(self):
        items = []
        for key, value in self.data.items():
            items.append(key + ": " + repr(value))

        result = "\n".join(items)

        if (self.outer):
            result += "\nOuter:\n"
            result += (repr(self.outer))

        return result

    def set(self, symbol, native_value):
        self.data[symbol] = native_value

    def find(self, symbol):
        try:
            return self.data[symbol]
        except KeyError:
            if self.outer:
                return self.outer.find(symbol)
            else:
                raise SymbolNotFound(symbol)

    def get(self, symbol):
        return self.find(symbol)

    def bind(self, symbols, values):
        for i in range(0, len(symbols)):
            symbol = symbols[i]

            if symbol == "&":
                # Set all the following symbol to the rest
                # of the argument list
                self.set(symbols[i + 1], values[i:])
                return
            else:
                self.set(symbols[i], values[i])
