<h1 align=center>scrabbler</h1>

<p align=center></p>
<p align=center>
    <img alt="GitHub Release (latest SemVer)" src="https://img.shields.io/github/v/release/wI2L/scrabbler">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/wI2L/scrabbler/ci.yml">
</p>

![](https://github.com/wI2L/scrabbler/blob/master/scrabbler.gif)

## Features

- Automatic and retryable draw
- Set custom word length (default to `7`)
- Custom draw configuration (vowels and consonants count per word)
- Tile distribution for various languages:
  - [French](https://en.wikipedia.org/wiki/Scrabble_letter_distributions#French)
  - [English](https://en.wikipedia.org/wiki/Scrabble_letter_distributions#English)
  - *...open an issue/pull request to request/add yours*
- Word insights from dictionaries:
  - 🇫🇷 ODS8
  - 🇬🇧 SOPWODS (UK)
  - 🇺🇸 TWL06 (US)

## Installation

### Manually

#### From source

```shell
go install github.com/wI2L/scrabbler@latest
```

> **Note**
> This will install the `scrabbler` binary in `$GOBIN`, which defaults to `$GOPATH/bin` or `$HOME/go/bin` if the `GOPATH` environment variable is not set.

#### Pre-compiled binaries

Download pre-compiled binaries from the [Releases](https://github.com/wI2L/scrabbler/releases) page.

## Usage

### User interface

The application presents itself as a terminal user-interface with two *views*:

- The first one, the *draw view*, displays the tiles that are randomly drawn from the selected distribution. You can either accept the draw or refuse it to generate a new one. Note that tiles that were not played in the previous round are kept and are not drawn again.
- The second, the *play view*, allows you to enter the tiles that were played (the order you type them in doesn't matter). You cannot enter unavailable tiles.

The game goes on, round after round, until tiles are exhausted.

### Options

#### Change the tile distribution

By default, the application starts with the [French](https://en.wikipedia.org/wiki/Scrabble_letter_distributions#French) tiles distribution. You can change it by using the `--distribution` flag.

The list of know distributions are:
- [`french`](https://github.com/wI2L/scrabbler/blob/master/cmd/distribution.go#L23)
- [`english`](https://github.com/wI2L/scrabbler/blob/master/cmd/distribution.go#L68)

See the [distribution.go](https://github.com/wI2L/scrabbler/blob/master/cmd/distribution.go) file, which define the tiles distribution for each language.

#### Use a custom dictionary

By default, the application loads the French *ODS8* dictionary. You can specify the path to your own dictionary with the `--dictionary` flag.

A valid dictionary must be a text file containing one word per line. *The words don't need to be sorted.*

### Key bindings

| Key binding                                           | Description                  |
|-------------------------------------------------------|------------------------------|
| <kbd>Control</kbd> + <kbd>C</kbd> / <kbd>Escape</kbd> | Exit application             |
| <kbd>?</kbd>                                          | Toggle help                  |
| <kbd>Enter</kbd>                                      | Validate selection/play word |
| <kbd>←</kbd> / <kbd>Y</kbd>                           | Select *Yes* option          |
| <kbd>→</kbd> / <kbd>N</kbd>                           | Select *No* option           |
| <kbd>Tab</kbd>                                        | Toggle selection             |
| <kbd>Control</kbd> + <kbd>C</kbd>                     | Show insights                |

## License

`scrabbler` is licensed under the **MIT** license. See the [LICENSE](LICENSE) file.