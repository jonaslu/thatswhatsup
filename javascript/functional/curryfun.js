function curry(fun, self) {
  const noArguments = fun.length;
  let passedArguments = [];

  return function innerCurry() {
    console.log(this);
    passedArguments.push.apply(passedArguments, [].slice.call(arguments));

    if (passedArguments.length >= noArguments) {
      return fun.apply(self, passedArguments);
    } else {
      return innerCurry;
    }
  }
}

const test = curry(function (a,b) { return "There will be " + (a + b) + " at the party, sir!"; });

console.log(test(1)(2));
console.log(test(1,2));
