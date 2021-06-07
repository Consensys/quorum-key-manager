const fs = require('fs');
const path = require('path');

const core = require('@actions/core');
const github = require('@actions/github');

const latestVersionLabel = 'latest';

const specsSource = core.getInput('specs-source');
const ghPagesBranch = core.getInput('gh-pages-branch');
const devBranch = core.getInput('dev-branch');
const specTargetPrefix = core.getInput('spec-target-prefix');
const versionsFile = core.getInput('versions-file');

const log = (...args) => console.log(...args); // eslint-disable-line no-console

async function main() {
  try {
    const cfg = getConfig();
    log("config", cfg);

    await writeVersionsFile(cfg);
    writeSpecFiles(specsSource, cfg);

  } catch (error) {
    core.setFailed(error);
  }
}

async function writeVersionsFile(cfg) {
  const versions = JSON.parse(fs.readFileSync(path.join(ghPagesBranch, versionsFile), 'utf8'));
  log("Original versions file", versions);
  if (cfg.spec.isReleaseVersion) {
    const updatedVersions = updateVersions(versions, cfg.spec.version);
    saveVersionsJson(updatedVersions, path.join(ghPagesBranch, versionsFile));
  }
}

function writeSpecFiles(specsSource, cfg){
  //write json files
  const specJsonString = fs.readFileSync(specsSource, 'utf8');
  const spec = JSON.parse(specJsonString);
  if(cfg.spec.isReleaseVersion) {
    spec.info.version = cfg.spec.version;
    saveSpecJson(spec, cfg.spec.releaseDist);
  }else{
    spec.info.version = `${cfg.spec.version}-${github.context.sha.substring(0,8)}`;
  }
  saveSpecJson(spec, cfg.spec.latestDist);
}

function updateVersions(versions, specVersion) {
  versions[specVersion] = {
    spec: specVersion,
    source: specVersion,
  };
  versions["stable"] = {
    spec: specVersion,
    source: specVersion,
  };
  return versions;
}

function saveVersionsJson(versions, versionsDist) {
  fs.writeFileSync(versionsDist, JSON.stringify(versions, null, 1));
}

function saveSpecJson(spec, specDist) {
  fs.writeFileSync(specDist, JSON.stringify(spec));
  log("writing spec file to", specDist);
}

function getConfig() {

  return {
    spec: calculateSpecDetails(specsSource),
    distDir: ghPagesBranch,
  };
}

function calculateSpecDetails(specFile) {

  const specVersion = calculateSpecVersion(github.context.payload.ref);
  const release = isReleaseVersion(specVersion);
  const latestDist = destinationPath(specFile, latestVersionLabel);
  const releaseDist = destinationPath(specFile, specVersion);

  return {
    path: specFile,
    version: specVersion,
    isReleaseVersion: release,
    latestDist: latestDist,
    releaseDist: releaseDist,
  };
}

function calculateSpecVersion(ref) {
  if(ref === `refs/heads/${devBranch}`){
    return latestVersionLabel;
  }else{
    const releasePattern = /^refs\/tags\/v(.+)?/;
    let match = ref.match( releasePattern );
    return match != null && match.length >= 2 ? match[1] : latestVersionLabel;
  }
}

function isReleaseVersion(version) {
  return version !== latestVersionLabel;
}

function destinationPath(specFile, version) {
  const extension = path.extname(specFile);
  return path.join(ghPagesBranch, `${specTargetPrefix}.${version}${extension}`);
}

main();
