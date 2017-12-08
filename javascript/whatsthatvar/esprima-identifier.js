const esprima = require("esprima");
const klaw = require("klaw");
const fs = require("fs-extra");
const path = require("path");

const identifiers = {};

function getAllIdentifiers(fileContents) {
  // TODO Handle import modules (parseModule)
  esprima.parseScript(fileContents, {}, node => {
    // TODO Filter out requires and module.exports
    // TODO Add function names

    if (node.type === "Identifier") {
      const { name } = node;
      identifiers[name] = identifiers[name] ? identifiers[name] + 1 : 1;
    }
  });
}

function parseFile(fileName) {
  return fs.readFile(fileName, "utf-8").then(getAllIdentifiers);
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
  .then(() => console.log(identifiers));
