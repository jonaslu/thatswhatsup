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
        value : stringValue[1]
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
  // ( has already been discarded, collect arguments until )
  const args = [];

  console.log(applyToken, program);

  for(;;) {
    const result = parseExpression(program);

    let expression = result.expression;
    program = result.program;

    args.push(expression);

    console.log(expression, program);

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

  // Advance beyond ')'
  return {
    expression: {
      type: 'apply',
      name: applyToken,
      args
    },
     program: program.substr(1)
  }
}

// ***********************
// TESTS
// ***********************
function testSingleExpression() {
  console.log(parseExpression("   \"yaya\""));
  console.log(parseExpression("   1234"));
  console.log(parseExpression("   do("));
  console.log(parseExpression("   a,b"));
  console.log(parseExpression("   a,   b"));
}

function testParseApply() {
  console.log(parseExpression("do(if(true, +(1,2), false))"))
}

testSingleExpression();
testParseApply();