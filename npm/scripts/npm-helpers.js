const os = require('os');
const path = require('path');
const tar = require('tar');
const spawnSync = require('child_process').spawnSync;
const semver = require('semver');
const packageMetadata = require('../package.json');
const { default: axios } = require('axios');
const { existsSync, mkdirSync, createWriteStream } = require('fs');
const unzipper = require('unzipper');
const { pipeline } = require('stream/promises');
const { Readable } = require('stream');

// Maps the architectures between Nodejs and GitHub release binaries
const architectureMap = {
  'ia32': 'i386',
  'x64': 'x86_64',
  'arm': 'arm64',
  'arm64': 'arm64'
};

// Maps the platforms between Nodejs and GitHub release binaries
const platformMap = {
  'darwin': 'Darwin',
  'linux': 'Linux',
  'win32': 'Windows',
};

class InitiumExecutable {
  constructor() {
    this.name = packageMetadata.binaryName;
    this.organization = 'nearform';
    this.sourceRepo = packageMetadata.main;
    this.version = packageMetadata.version;
    this.releaseFileExtension = os.platform === 'win32' ? 'zip' : 'tar.gz';
    this.releaseFileNameFull = `${this.sourceRepo}_${platformMap[os.platform]}_${architectureMap[os.arch]}.${this.releaseFileExtension}`;
    this.releaseUrl = `${packageMetadata.repository.url}/releases/download/v${this.version}/${this.releaseFileNameFull}`;
    this.executableExtension = os.platform === 'win32' ? '.exe' : '';
    this.relativeInstallDirectory = path.join('node_modules', '.bin');
    this.installDirectory = path.join(process.cwd(), this.relativeInstallDirectory);
    this.executablePath = `${path.join(this.installDirectory, this.name)}${this.executableExtension}`
  }

  async checkForUpdates() {
    let releases = `https://api.github.com/repos/${this.organization}/${this.sourceRepo}/releases/latest`
    console.log(`Fetching releases from ${releases}`)
    const latestReleaseResponse = await axios(releases);
    const latestVersion = latestReleaseResponse.data.tag_name.replace('v', '');
    if (semver.gt(latestVersion, this.version)) {
      console.log(`There is a new ${this.name} version available!\n\nCurrent version: ${this.version}\nLatest version: ${latestVersion}\n\nConsider upgrading using npm update.\n`);
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
      try {
        console.log(`Fetching release from ${this.releaseUrl}`)
        const releaseFileData = await axios(this.releaseUrl, { responseType: 'arraybuffer' });
        const releaseFileDataReadable = Readable.from(Buffer.from(releaseFileData.data));
        if (this.releaseFileExtension === 'tar.gz') {
          await pipeline(
            releaseFileDataReadable,
            tar.x({ cwd: this.installDirectory }, [`${this.name}${this.executableExtension}`])
          );
        } else {
          await pipeline(
            releaseFileDataReadable,
            unzipper.ParseOne(`${this.name}${this.executableExtension}`),
            createWriteStream(path.join(this.installDirectory, `${this.name}${this.executableExtension}`))
          )
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
        throw new Error(result.error);
      }

      process.exit(result.status);
    } catch (error) {
      console.error(`Error when running ${this.name}: ${error}.`);
      process.exit(1);
    }
  }
}

module.exports = InitiumExecutable
