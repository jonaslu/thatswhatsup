function compile(string) {
  const replaced = string.replace(/{{\s*([^\s]+)\s*}}/, '\' + data\.$1 + \'');

  const evalStr = "(function (data) { return '" + replaced + "'})";
  const evaled = eval(evalStr);

  return function(data) {
    return evaled(data);
  }
}

const result = compile('<h1>yahoo {{name}}</h1>')({name: 'jesus'});
console.log(result);