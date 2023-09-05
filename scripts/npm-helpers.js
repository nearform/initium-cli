const os = require('os');
const path = require('path');
const tar = require('tar');
const zlib = require('zlib');
const spawnSync = require('child_process').spawnSync;
const semver = require('semver');
const packageMetadata = require('../package.json');
const { default: axios } = require('axios');
const { existsSync, mkdirSync, createWriteStream } = require('fs');
const { unzip } = require('unzipper');
const { pipeline } = require('stream/promises');

// Maps the architectures between Nodejs and GitHub release binaries
const architectureMap = {
  'ia32': 'i386',
  'x64': 'x86_64',
  'arm': 'arm64'
};

// Maps the platforms between Nodejs and GitHub release binaries
const platformMap = {
  'darwin': 'Darwin',
  'linux': 'Linux',
  'win32': 'Windows',
};

class InitiumExecutable {
  constructor() {
    this.name = packageMetadata.name;
    this.version = packageMetadata.version;
    this.releaseFileExtension = os.platform === 'win32' ? 'zip' : 'tar.gz';
    this.releaseUrl = `${packageMetadata.repository.url}/releases/download/v${this.version}/${this.name}_${platformMap[os.platform]}_${architectureMap[os.arch]}.${this.releaseFileExtension}`;
    this.executableExtension = os.platform === 'win32' ? '.exe' : ''
    this.installDirectory = path.join(__dirname, 'node_modules', '.bin');
    this.executablePath = `${path.join(this.installDirectory, this.name)}${this.executableExtension}`
  }

  async checkForUpdates() {
    const latestReleaseResponse = await axios(`https://api.github.com/repos/nearform/${this.name}/releases/latest`);
    console.log(latestReleaseResponse.data);
    const latestVersion = latestReleaseResponse.data['tag_name'].replace('v', '');
    if (semver.gt(latestVersion), this.version) {
      console.log(`There's a new version available!\n\nCurrent version: ${this.version}\nLatest version: ${latestVersion}\n\nConsider upgrading using npm update.`);
    }
  }

  async install(runOnly = false) {
    try {
      await this.checkForUpdates();
    } catch (error) {
      console.error(`Error when checking for updates: ${error}.`);
    }
    if (!existsSync(this.executablePath)) {
      mkdirSync(this.installDirectory, { recursive: true });
      const executableData = await axios(this.releaseUrl, { responseType: "stream" });
      try {
        if (this.releaseFileExtension === 'tar.gz') {
          const gzUnzip = zlib.createUnzip();
          const tarUnzip = tar.Extract();
          const fileOutput = createWriteStream(this.executablePath);
          await pipeline(executableData, gzUnzip, tarUnzip, fileOutput);
        } else {
          const fileOutput = createWriteStream(this.executablePath);
          await pipeline(executableData, unzip, fileOutput);
        }
      } catch (error) {
        console.error(`Error when extracting file: ${error}`);
        process.exit(1);
      }
    } else {
      if (!runOnly) {
        console.log(`${this.name} has already been installed. Run it with npx ${this.name}.`);
      }
    }
  }

  async run() {
    try {
      await this.install(true);
      const [, , ...args] = process.argv;

      const options = { cwd: process.cwd(), stdio: "inherit" };

      const result = spawnSync(this.executablePath, args, options);

      if (result.error) {
        error(result.error);
      }

      process.exit(result.status);
    } catch (error) {
      console.error(`Error when running ${this.name}: ${error}.`);
      process.exit(1);
    }
  }
}

module.exports = InitiumExecutable
