<h1 align=center>scrabbler</h1>

<p align=center><b>Pick your tiles, but not yourself!</b></p>
<p align=center>Automatic draw TUI for your Scrabble <i>duplicate</i> games</p>
<p align=center>
    <img alt="French" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/l/FR.svg">
    <img alt="English (US)" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/l/US.svg">
    <img alt="English (UK)" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/l/GB-NIR.svg">
    <img alt="German" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/l/DE.svg">
    <img alt="Italian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/l/IT.svg">
    <img alt="Dutch" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/l/NL.svg">
    <img alt="Czech" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/l/CZ.svg">
    <img alt="Icelandic" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/l/IS.svg">
</p>
<p align=center>
    <img alt="GitHub Release (latest SemVer)" src="https://img.shields.io/github/v/release/wI2L/scrabbler">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/wI2L/scrabbler/ci.yml">
</p>
<br/>

![](https://github.com/wI2L/scrabbler/blob/master/scrabbler.gif)

## Features

- Automatic and retryable random draw
- Custom word length
- Custom draw configuration (vowels and consonants requirements)
- [Tiles distribution for 7 languages](#tiles-distribution)
- [Word insights from dictionaries](#custom-dictionary)

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

- **draw view**: displays the tiles that are randomly drawn from the selected distribution. You can either accept the draw or refuse it to generate a new one. *Note that tiles that were not played in the previous round are kept and are not drawn again.*
- **play view**: allows you to enter the tiles that were played (the order you type them in doesn't matter). You cannot enter unavailable tiles.

### Options

```text
scrabbler — pick your tiles, but not yourself

Usage:
  scrabbler [flags]

Flags:
  -v, --debug                 enable debug logging to a file
  -d, --dictionary string     custom dictionary file path
  -l, --distribution string   tiles distribution language (default "french")
  -h, --help                  help for scrabbler
  -p, --show-points           show tile points
  -w, --word-length uint8     the number of tiles to draw (default 7)
```

#### Word length

The official rule define a word of seven (`7`) letters, but you can change it using the `-w`/`--word-length` flags:

```shell
scrabbler --word-length=8
```

#### Show tile points

The tiles of the draw can optionally show the points of each letter using the flags `-p`/`--show-points`. This option is disabled by default.

> **Important**
> The subscript characters [U+2080](https://www.compart.com/en/unicode/U+2080) to [U+2089](https://www.compart.com/en/unicode/U+2089) are used to represent the digits from 0 to 9. The number 10 is represented using the [U+2093](https://www.compart.com/en/unicode/U+2093) (*Latin Subscript Small Letter X*) to preserve equal spacing.
>
> Make sure to use a font that support those characters, such as *SF Mono* on macOS.

#### Tiles distribution

By default, the application starts with the [French](https://en.wikipedia.org/wiki/Scrabble_letter_distributions#French) tiles distribution.

You can change it with the `-d`/`--distribution` flags:

```shell
scrabbler --distribution=english
```

The list of know distributions are:

- `french`
- `english`
- `german`
- `italian`
- `dutch`
- `czech`
- `icelandic`: The standard Icelandic distribution, from Tinderbox games under license from Mattel
- `krafla`: Alternate Icelandic distribution, sanctioned by Iceland's Scrabble clubs for their tournaments and for the national championship

> **Note**
> All information are compiled from the [Scrabble letter distributions](https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Indonesian) Wikipedia page.

See the [distribution.go](https://github.com/wI2L/scrabbler/blob/master/cmd/distribution.go) file, which define the tiles distribution for each language.

#### Custom dictionary

By default, the application loads the French *ODS8* dictionary, which is embedded into the binary with the Go `embed` package.

Alternatively, you can specify the path of a dictionary of your choice with the `-w`/`--dictionary` flags:

```shell
scrabbler --dictionary=dictionaries/english/twl06.txt.gz
```

A valid dictionary is a text file which contain one word per line (*the words don't need to be sorted)*. The file can optionally be gzipped (*.txt.gz*)

See the [dictionaries](https://github.com/wI2L/scrabbler/tree/master/cmd/dictionaries) directory, which already contains some official and non-official dictionaries for several languages:

- 🇫🇷 [ODS8](https://en.wikipedia.org/wiki/L%27Officiel_du_jeu_Scrabble)
- 🇬🇧 [SOPWODS](https://en.wikipedia.org/wiki/Collins_Scrabble_Words)
- 🇺🇸 [TWL06](https://en.wikipedia.org/wiki/NASPA_Word_List)
- 🇮🇹 [ZINGA](https://www.listediparole.it/tutteleparole.txt)
- 🇩🇪 [hippler/german-wordlist](https://github.com/hippler/german-wordlist)
- 🇮🇸 [vthorsteinsson/Skrafl](https://github.com/vthorsteinsson/Skrafl)

### Key bindings

| Keys                                                  | Description                                                  |
|:------------------------------------------------------|:-------------------------------------------------------------|
| <kbd>Control</kbd> + <kbd>C</kbd> / <kbd>Escape</kbd> | Exit application                                             |
| <kbd>?</kbd>                                          | Toggle help                                                  |
| <kbd>Enter</kbd>                                      | Validate selection/play word                                 |
| <kbd>←</kbd> / <kbd>Y</kbd>                           | Select *Yes* option                                          |
| <kbd>→</kbd> / <kbd>N</kbd>                           | Select *No* option                                           |
| <kbd>Tab</kbd>                                        | Toggle selection                                             |
| <kbd>Control</kbd> + <kbd>G</kbd>                     | *press once*: Show insights<br>*press twice*: Show scrabbles |

## License

`scrabbler` is licensed under the **MIT** license. See the [LICENSE](LICENSE) file.