/*
  Data layout

  commands: []
  options: [{
    short: '-p',
    long: '--pizza',
    variableName: 'pizza',
    description: 'text'
  },
  {

  }]
*/

const { log } = console;

let program;
const parsedOptions = [];
const matchedOptions = {};

function option(option, description) {
  let shortOption, longOption, variableName;

  option.split(',').forEach(option => {
    option = option.replace(/\s/g, '');

    if (option.startsWith('--')) {
      longOption = option;
      variableName = option.substring(2);
      return;
    }

    if (option.startsWith('-')) {
      shortOption = option;
      return;
    }

    throw new Error('Unknown option format: ' + option);
  });

  if (!longOption) {
    throw new Error('No long option was given for format: ' + option);
  }

  parsedOptions.push({
    shortOption,
    longOption,
    variableName,
    description
  })
}

function programName(programName) {
  program = programName
}

function parse(argv) {
  let [, calledFile, ...args] = argv;
  program = program || calledFile;

  // Parse arguments until a command is hit, these arguments will be parsed by the command
  // until an unknown is hit

  while (args.length) {
    // Is it a command? Let the command consume args
    const [nextArg, ...rest] = args;

    const optionMatched = parsedOptions.some(option => {
      if (option.shortOption === nextArg || option.longOption === nextArg) {
        matchedOptions[option.variableName] = true;
        return true;
      }

      return false;
    })

    if (!optionMatched) {
      throw new Error('Unknown option given: ' + nextArg);
    }

    args = rest;
  }
  // Execute callbacks on the commands
}

module.exports = {
  option,
  programName,
  parse,
  matchedOptions
}
