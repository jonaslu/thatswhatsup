const parseargv = require('./parseargv');

function assertTrue(expected, value) {
  if (expected !== value) {
    throw new Error(`Assert failed, expected ${expected} actual ${value}`);
  }
}

function assertErrorThrown(fun, expectedMessage) {
  try {
    fun()
  } catch (e) {
    assertTrue(e.message, expectedMessage);
    return;
  }

  throw new Error('Assert failed');
}

const tests = {
  testShortOption() {
    parseargv.option('-p, --pizza', 'your flavour of pizza');
    parseargv.parse(['node', 'programName', '-p']);

    assertTrue(parseargv.matchedOptions['pizza'], true);
  },

  testLongOption() {
    parseargv.option('-p, --pizza', 'your flavour of pizza');
    parseargv.parse(['node', 'programName', '--pizza']);

    assertTrue(parseargv.matchedOptions['pizza'], true);
  },

  testTwoOptionsOneGiven() {
    parseargv.option('-p, --pizza', 'your flavour of pizza');
    parseargv.option('-n, --noodles', 'your flavour of noodles');
    parseargv.parse(['node', 'programName', '--pizza']);

    assertTrue(parseargv.matchedOptions['pizza'], true);
    assertTrue(parseargv.matchedOptions['noodles'], undefined);
  },

  testTwoOptionsTwoConsumed() {
    parseargv.option('-p, --pizza', 'your flavour of pizza');
    parseargv.option('-n, --noodles', 'your flavour of noodles');
    parseargv.parse(['node', 'programName', '--pizza']);

    assertTrue(parseargv.matchedOptions['pizza'], true);
    assertTrue(parseargv.matchedOptions['noodles'], undefined);
  },

  testErrorOnNoLongOption() {
    assertErrorThrown(() =>
      parseargv.option('-p', 'your flavour of pizza'),
      'No long option was given for format: -p');
  },

  tesOnlyLongOptionOk() {
    parseargv.option('--pizza', 'your flavour of pizza');
    parseargv.parse(['node', 'programName', '--pizza']);

    assertTrue(parseargv.matchedOptions['pizza'], true);
  },

  testErrorOnUnkownOption() {
    parseargv.option('-p, --pizza', 'your flavour of pizza');

    assertErrorThrown(() =>
      parseargv.parse(['node', 'programName', '--fluffy']),
      'Unknown option given: --fluffy');
  }
};

Object.keys(tests).forEach(testKey => {
    try {
      tests[testKey]();
    } catch (e) {
      console.error(`Error in test ${testKey}: ${e.message}`)
    }
  });
