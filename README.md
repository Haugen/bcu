# bcu, the Branch Cleaner Upper

Tired of manually cleaning up and deleting your local git branches? I know, me too. `bcu` aims to be a simple programs that makes selecting and deleting local git branches a breeze.

## Installation

Binaries for different OS are available in the Releases.

### Homebrew

You can install using my private Homebrew Cask:

```bash
brew install --cask haugen/tap/bcu
```

Since I don't have an Apple Developer Account you might have issues runing the installed binary. To allow `bcu` to run on MacOS you can remove the quarantine like this

```bash
xattr -d com.apple.quarantine /opt/homebrew/bin/bcu
```