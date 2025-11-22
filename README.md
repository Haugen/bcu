# bcu, the Branch Cleaner Upper

Tired of manually cleaning up and deleting your local git branches? I know, me too. bcu aims to be a simple programs that makes selecting and deleting local git branches a breeze.

## Usage

Just run the program in any git repository and get those local branches cleaned up!

```bash
bcu
```

## Installation

Binaries for different OS are available in the Releases.

### Homebrew

You can install using my private Homebrew Cask. Since I don't have an Apple Developer Account you might have issues runing the installed binary. The "install and trust" script below helps with removing the imposed Appel quarantine of the program.

Install and trust:

```bash
brew install --cask haugen/tap/bcu && xattr -d com.apple.quarantine /opt/homebrew/bin/bcu
```

Only install:

```bash
brew install --cask haugen/tap/bcu
```