<h1 align=center>scrabbler</h1>

<h3 align=center>Pick tiles, but not yourself!</h3>
<p align=center>Automatic draw TUI for your <i>duplicate</i> Scrabble games</p>
<br/>
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
    <br/>
    <img alt="Greek" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/GR.svg">
    <img alt="Indonesian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/ID.svg">
    <img alt="Latvian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/LV.svg">
    <img alt="Lithuanian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/LT.svg">
    <img alt="Malay" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/MY.svg">
    <img alt="Norwegian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/NO.svg">
    <img alt="Polish" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/PL.svg">
    <img alt="Portuguese" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/PT.svg">
    <img alt="Romanian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/RO.svg">
    <img alt="Slovak" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/SK.svg">
    <img alt="Slovenian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/SI.svg">
    <img alt="Swedish" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/SE.svg">
    <img alt="Ukrainian" src="https://raw.githubusercontent.com/Yummygum/flagpack-core/main/svg/m/UA.svg">
</p>
<br/>
<p align=center>
    <img alt="GitHub Release (latest SemVer)" src="https://img.shields.io/github/v/release/wI2L/scrabbler">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/wI2L/scrabbler/ci.yml">
</p>

![](https://github.com/wI2L/scrabbler/blob/master/scrabbler.gif)

### Features

- 🎲 **Retryable draw** — redo the draw if you are not satisfied with the picks
- 📏 **Custom word length** — play with standard 7 tiles or for any variant
- ⚙️ **Custom draw configuration** — minimum number of vowels and/or consonants per draw
- 🌐 [**Letter distributions**](#letter-distributions) — play in the language of your choice
- 💡 [**Word insights from dictionaries**](#custom-dictionary) — reveal the number of *scrabbles* and the words
- ⏱️ **Game timer** — thinking time per play

---

## Installation

### Manually

#### From source

```shell
go install github.com/wI2L/scrabbler@latest
```

> [!NOTE]
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
Usage:
  scrabbler [flags]

Flags:
      --consonants uint8      number of required consonant letters
  -v, --debug                 enable debug logging
  -d, --dictionary string     custom dictionary file path
  -l, --distribution string   letter distribution language
  -h, --help                  help for scrabbler
  -p, --show-points           show letter points
  -t, --timer duration[=5m]   enable play timer (default 5m)
      --vowels uint8          number of required vowel letters
  -w, --word-length uint8     the number of tiles to draw (default 7)
```

#### Word length

The official rule define a word of seven (`7`) letters, but you can change it using the `-w`/`--word-length` flags:

```shell
scrabbler --word-length=8
```

#### Draw requirements

The official [duplicate scrabble rules](https://en.wikipedia.org/wiki/Duplicate_Scrabble#Rules) states that a draw must always contain one vowel and one consonant. You can use the `--vowels` and `--consonants` flags to configure this behavior (disabled by default, the draw is completely random).

Unlike the official rules, the draws won't stop automatically once there are no more vowels or consonants to pick. However, you can choose to stop the game yourself.

```shell
scrabbler --vowels=1 --consonants=1
```

> [!IMPORTANT]
> The sum of required vowels and consonants cannot exceed the configured word length.

#### Tile points

The tiles of the draw can optionally show the points of each letter using the flags `-p`/`--show-points`. This option is disabled by default.

> [!IMPORTANT]
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

> [!NOTE]
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
>
> &mdash;&mdash; [Wikipedia](https://en.wikipedia.org/wiki/Scrabble_letter_distributions)

By default, the application loads the [French](https://en.wikipedia.org/wiki/Scrabble_letter_distributions#French) letter distribution.

You can change it with the `-l`/`--distribution` flags:

```shell
scrabbler --distribution=english
```

Below is the list of official distributions that are supported:

- `afrikaans` — *Afrikaans*
- `bulgarian` — *Български*
- `czech` — *Čeština*
- `danish` — *Dansk*
- `dutch` — *Nederlands*
- `english` — *English*
- `estonian` — *Eesti*
- `finnish` — *suomi*
- `french` — *Français*
- `german` — *Deutsch*
- `greek` — *Ελληνικά*
- `icelandic` — *Íslenska*
- `indonesian` — *Bahasa Indonesia*
- `italian` — *Italiano*
- `latvian` — *Latviešu*
- `lithuanian` — *Lietuvių*
- `malay` — *Bahasa Melayu*
- `norwegian` — *Norsk*
- `polish` —  *Polski*
- `portuguese` — *Português*
- `romanian` — *Română*
- `slovak` — *Slovenčina*
- `slovenian` — *Slovenščina*
- `swedish` — *Svenska*
- `ukrainian` — *Українська*

Alternate distributions are also available:

- `krafla`: Alternate Icelandic distribution, sanctioned by Iceland's Scrabble clubs for their tournaments and for the national championship

> [!NOTE]
> All information are compiled from the [Scrabble letter distributions](https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Indonesian) Wikipedia page.

See the [distribution.go](https://github.com/wI2L/scrabbler/blob/master/cmd/distribution.go) file, which define the letter distribution for each language.

> [!IMPORTANT]
> Several editions such as Spanish, Catalan, Hungarian or Welsh (to cite a few) are not included because they use [digraphs](https://en.wikipedia.org/wiki/Digraph_(orthography)), which are challenging to deal with in a text-based UI.
>
> *They might be added in the future, once I have figured out an intuitive way to handle them*.

#### Custom dictionary

By default, the application loads the official French Scrabble dictionary (ODS8), which is embedded into the binary with the Go `embed` package.

Alternatively, you can specify the path of a dictionary of your choice with the `-d`/`--dictionary` flags:

```shell
scrabbler --dictionary=dictionaries/english/twl06.txt.gz
```

A valid dictionary is a text file which contain one word per line (*the words don't need to be sorted*).

The file can optionally be *gzipped* (the file extension doesn't matter, the detection is [header-based](https://pkg.go.dev/net/http#DetectContentType)).

Browse the [dictionaries](https://github.com/wI2L/scrabbler/tree/master/dictionaries) directory, which already contains some official and non-official dictionaries for several languages:

- :fr: [ODS8](https://en.wikipedia.org/wiki/L%27Officiel_du_jeu_Scrabble): The 8th version of the official dictionary for Francophone Scrabble.
- :uk: [SOPWODS](https://en.wikipedia.org/wiki/Collins_Scrabble_Words): Official word list used in English-language tournament Scrabble in most countries except the US, Thailand and Canada.
- :us: [TWL06](https://en.wikipedia.org/wiki/NASPA_Word_List): Official word authority for tournament Scrabble in the USA and Canada.
- :it: [listediparole](https://www.listediparole.it/tutteleparole.txt): Unofficial word list extracted from the [listediparole.it](https://www.listediparole.it) website.
- :de: [hippler/german-wordlist](https://github.com/hippler/german-wordlist): Unofficial word list compiled by [Stefan Hippler](https://github.com/hippler).
- :iceland: [vthorsteinsson/Skrafl](https://github.com/vthorsteinsson/Skrafl): Unofficial word list compiled by [Vilhjalmur Thorsteinsson](https://github.com/vthorsteinsson) from the *Database of Icelandic Morphology* (DIM, BÍN) for the Icelandic crossword game [Netskrafl](https://github.com/mideind/Netskrafl).

### Key bindings

| Keys                                      | Description                                                                                                                                          |
|:------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------|
| <kbd>Control+C</kbd> or <kbd>Escape</kbd> | Exit the application *without confirmation*                                                                                                          |
| <kbd>?</kbd>                              | Toggle help view (switch between short and extended)                                                                                                 |
| <kbd>Enter</kbd>                          | Validate selection or play word                                                                                                                      |
| <kbd>←</kbd>/<kbd>y</kbd>                 | Select the *yes* option                                                                                                                              |
| <kbd>→</kbd>/<kbd>n</kbd>                 | Select the *no* option                                                                                                                               |
| <kbd>Tab</kbd>                            | Toggle option selection                                                                                                                              |
| <kbd>Control+G</kbd>                      | Press once to show word insights (whether one or more *scrabble(s)* have been found with the tiles of the draw), press twice to show the words found |

## License

The code is licensed under the **MIT** license. [Read this](https://www.tldrlegal.com/license/mit-license) or see the [LICENSE](LICENSE) file.