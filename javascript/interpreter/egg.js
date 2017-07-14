function skipWitespace(program) {
  const nonWhitespaceCharacter = program.search(/\S/);
  return program.substring(nonWhitespaceCharacter);
}

function parseExpression(program) {
  program = skipWitespace(program);

  const stringValue = /^"([^"]*)"/.exec(program);
  if (stringValue) {
    return {
      expression: {
        type: 'string',
        value: stringValue[1]
      },
      // +2 for the quote chars ("(something)")
      program: program.substr(stringValue[1].length + 2)
    }
  }

  const numberValue = /^(\d+)/.exec(program);
  if (numberValue) {
    return {
      expression: {
        type: 'integer',
        value: parseInt(numberValue[1], 10)
      },
      program: program.substr(numberValue[1].length)
    }
  }

  const word = /^([^\(\),\s]+)/.exec(program);

  if (word) {
    program = skipWitespace(program.substr(word[1].length));

    if (program[0] === '(') {
      return parseApply(word[1], program.substr(1))
    }

    return {
      expression: {
        type: 'word',
        value: word[1]
      },
      program
    }
  }

  throw new Error(`Token not recognized: ${program}`);
}

function parseApply(applyToken, program) {
  if (program[0] === '(') {
    // When we call ourselves for returning functions that are immediatly called
    program = program.substr(1);
  }

  const args = [];

  for (; ;) {
    const result = parseExpression(program);

    let expression = result.expression;
    program = result.program;

    args.push(expression);

    program = skipWitespace(program);

    if (program[0] === ')') {
      break;
    }

    if (program[0] === ',') {
      program = skipWitespace(program.substr(1));
    } else {
      throw new Error(`Failed, expected comma as argument separator: ${program} ${JSON.stringify(expression)}`);
    }
  }

  const expression = {
    type: 'apply',
    name: applyToken,
    args
  };

  // Advance beyond ')'
  program = skipWitespace(program.substr(1));
  if (program[0] === '(') {
    // When interpreting, check if name is of type apply
    return parseApply(expression, program)
  }

  return {
    expression,
    program
  }
}

// ***********************
// TESTS
// ***********************
function prettyPrint(expression) {
  console.log(JSON.stringify(expression, null, 2));
}

function testSingleExpression() {
  prettyPrint(parseExpression("   \"yaya\""));
  prettyPrint(parseExpression("   1234"));
  prettyPrint(parseExpression("   a,b"));
  prettyPrint(parseExpression("   a,   b"));
}

function testParseApply() {
  prettyPrint(parseExpression("do(if(true, +(1,2), false))"));
  prettyPrint(parseExpression("multiplier(2)(1)"));
}

testSingleExpression();
testParseApply();