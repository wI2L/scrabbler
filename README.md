<h1 align=center>scrabbler</h1>

<p align=center><b>Pick your tiles, but not yourself!</b></p>
<p align=center>Automatic draw TUI for <i>duplicate</i> Scrabble games</p>
<p align=center>
    <img alt="French" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/FR.svg">
    <img alt="English (US)" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/US.svg">
    <img alt="English (UK)" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/GB-NIR.svg">
    <img alt="German" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/DE.svg">
    <img alt="Italian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/IT.svg">
    <img alt="Dutch" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/NL.svg">
    <img alt="Czech" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/CZ.svg">
    <img alt="Icelandic" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/IS.svg">
    <img alt="Afrikaans" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/ZA.svg">
    <img alt="Bulgarian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/BG.svg">
    <img alt="Danish" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/DK.svg">
    <img alt="Estonian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/EE.svg">
    <img alt="Finnish" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/FI.svg">
    <img alt="Greek" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/GR.svg">
    <img alt="Indonesian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/ID.svg">
    <img alt="Latvian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/LV.svg">
    <img alt="Lithuanian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/LT.svg">
</p>
<br/>
<p align=center>
    <img alt="GitHub Release (latest SemVer)" src="https://img.shields.io/github/v/release/wI2L/scrabbler">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/wI2L/scrabbler/ci.yml">
</p>

![](https://github.com/wI2L/scrabbler/blob/master/scrabbler.gif)

## Features

- Automatic and retryable random draw
- Custom word length
- Custom draw configuration (*vowels or consonants requirements*)
- [Letter distributions for many languages](#letter-distributions)
- [Word insights from dictionaries](#custom-dictionary)
- Game timer

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

- **Draw view**: displays the tiles that are randomly drawn from the selected distribution. You can either accept the draw or refuse it to generate a new one. *Note that tiles that were not played in the previous round are kept and are not drawn again.*
- **Play view**: allows you to enter the tiles that were played (the order you type them in doesn't matter). You cannot enter unavailable tiles.

### Options

```text
scrabbler — pick your tiles, but not yourself

Usage:
  scrabbler [flags]

Flags:
  -v, --debug                 enable debug logging to a file
  -d, --dictionary string     custom dictionary file path
  -l, --distribution string   letter distribution language (default "french")
  -h, --help                  help for scrabbler
  -p, --show-points           show tile points
  -t, --timer duration[=5m]   enable play timer
  -w, --word-length uint8     the number of tiles to draw (default 7)
```

#### Word length

The official rule define a word of seven (`7`) letters, but you can change it using the `-w`/`--word-length` flags:

```shell
scrabbler --word-length=8
```

#### Tile points

The tiles of the draw can optionally show the points of each letter using the flags `-p`/`--show-points`. This option is disabled by default.

> **Important**
> The subscript characters [U+2080](https://www.compart.com/en/unicode/U+2080) to [U+2089](https://www.compart.com/en/unicode/U+2089) are used to represent the digits from 0 to 9. The number 10 is represented using the [U+2093](https://www.compart.com/en/unicode/U+2093) (*Latin Subscript Small Letter X*) to preserve equal spacing.
>
> Make sure to use a font that support those characters, such as *SF Mono* on macOS.

#### Game timer

It is possible to show a timer during the *play* phase, once a draw have been accepted. To use the default timer duration of 5 minutes, simply use the `-t`/`--timer` flags without specifying a value:

```shell
scrabbler --timer
```

To change the duration, set the flag with a custom value:

```shell
scrabbler --timer=3m
```

> **Note**
> The duration value must be expressed as a stringified Go `Duration`, as defined by the documentation of the [`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration) function.
>
> Examples:
>
> - *20 seconds*: `20s`
> - *1 minute*: `1m`
> - *3 minutes and 20 seconds*: `3m20s`

#### Letter distributions

> Editions of the word board game Scrabble in different languages have differing letter distributions of the tiles, because the frequency of each letter of the alphabet is different for every language. As a general rule, the rarer the letter, the more points it is worth.
>
> Most languages use sets of 100 tiles, since the original distribution of ninety-eight tiles was later augmented with two blank tiles.

By default, the application starts with the [French](https://en.wikipedia.org/wiki/Scrabble_letter_distributions#French) letter distribution.

You can change it with the `-l`/`--distribution` flags:

```shell
scrabbler --distribution=english
```

Below is the list of official distributions that are supported:

- `afrikaans`
- `bulgarian`
- `czech`
- `danish`
- `dutch`
- `english`
- `estonian`
- `finnish`
- `french`
- `german`
- `greek`
- `icelandic`
- `indonesian`
- `italian`
- `latvian`
- `lithuanian`

Alternate distributions are also available:

- `krafla`: Alternate Icelandic distribution, sanctioned by Iceland's Scrabble clubs for their tournaments and for the national championship

> **Note**
> All information are compiled from the [Scrabble letter distributions](https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Indonesian) Wikipedia page.

See the [distribution.go](https://github.com/wI2L/scrabbler/blob/master/cmd/distribution.go) file, which define the letter distribution for each language.

#### Custom dictionary

By default, the application loads the official French Scrabble dictionary (ODS8), which is embedded into the binary with the Go `embed` package.

Alternatively, you can specify the path of a dictionary of your choice with the `-d`/`--dictionary` flags:

```shell
scrabbler --dictionary=cmd/dictionaries/english/twl06.txt.gz
```

A valid dictionary is a text file which contain one word per line (*the words don't need to be sorted*). The file can optionally be *gzipped*.

See the [dictionaries](https://github.com/wI2L/scrabbler/tree/master/cmd/dictionaries) directory, which already contains some official and non-official dictionaries for several languages:

- :fr: [ODS8](https://en.wikipedia.org/wiki/L%27Officiel_du_jeu_Scrabble)
- :uk: [SOPWODS](https://en.wikipedia.org/wiki/Collins_Scrabble_Words)
- :us: [TWL06](https://en.wikipedia.org/wiki/NASPA_Word_List)
- :it: [ZINGA](https://www.listediparole.it/tutteleparole.txt)
- :de: [hippler/german-wordlist](https://github.com/hippler/german-wordlist)
- :iceland: [vthorsteinsson/Skrafl](https://github.com/vthorsteinsson/Skrafl)

### Key bindings

| Keys                                       | Description                                                                                                                                          |
|:-------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------|
| <kbd>Control + C</kbd> / <kbd>Escape</kbd> | Exit the application (:warning: **no confirmation is requested**)                                                                                    |
| <kbd>?</kbd>                               | Toggle help (switch between short and full views)                                                                                                    |
| <kbd>Enter</kbd>                           | Validate selection or play word                                                                                                                      |
| <kbd>←</kbd> / <kbd>Y</kbd>                | Select the *yes* option                                                                                                                              |
| <kbd>→</kbd> / <kbd>N</kbd>                | Select the *no* option                                                                                                                               |
| <kbd>Tab</kbd>                             | Toggle option selection                                                                                                                              |
| <kbd>Control + G</kbd>                     | Press once to show word insights (whether one or more *scrabble(s)* have been found with the tiles of the draw), press twice to show the words found |

## License

`scrabbler` is licensed under the **MIT** license. See the [LICENSE](LICENSE) file.