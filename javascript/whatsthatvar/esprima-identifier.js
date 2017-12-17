const esprima = require("esprima");
const klaw = require("klaw");
const fs = require("fs-extra");
const path = require("path");

const identifiers = {};
const moduleRegExp = /(^|(\s)+)((export)|(import))\s+/;
const crunchBangRegExp = /^#!\s+/;

function getAllIdentifiers(fileContents) {
  let parseMethod = esprima.parseScript;

  if (moduleRegExp.test(fileContents)) {
    parseMethod = esprima.parseModule;
  }

  // Remove any starting crunch-bang
  if (fileContents.startsWith("#!")) {
    firstNewline = fileContents.indexOf("\n");

    if (firstNewline > -1) {
      fileContents = fileContents.substr(firstNewline + 1);
    }
  }

  try {
    parseMethod(fileContents, {}, node => {
      // TODO Filter out requires and module.exports
      // TODO Add function names

      if (node.type === "Identifier") {
        const { name } = node;
        identifiers[name] = identifiers[name] ? identifiers[name] + 1 : 1;
      }
    });
  } catch (e) {
    console.error(`Error ${e}`);
    throw e;
  }
}

function parseFile(fileName) {
  return fs.readFile(fileName, "utf-8").then(fileContents => {
    try {
      getAllIdentifiers(fileContents);
    } catch (e) {
      console.error(`In file ${fileName}`);
    }
  });
}

const jsFileNames = [];

function enqueueFile(fileName) {
  if (path.extname(fileName) === ".js") {
    jsFileNames.push(fileName);
  }
}

new Promise((resolve, reject) => {
  klaw(process.argv[2])
    .on("data", item => enqueueFile(item.path))
    .on("end", resolve);
})
  .then(() => {
    const allFilesParsedPromises = jsFileNames.map(parseFile);

    return Promise.all(allFilesParsedPromises);
  })
  .then(() => {
    console.log(identifiers);
  });
