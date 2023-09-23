package cmd

import (
	"bytes"
	"compress/gzip"
	"slices"

	"golang.org/x/text/language"

	en "github.com/wI2L/scrabbler/dictionaries/english"
	fr "github.com/wI2L/scrabbler/dictionaries/french"
)

const blank = "?"

// letter represents a Scrabble letter.
// It can either be a single rune or a sequence of
// several runes to represent a digraph.
type letter struct {
	L         string
	frequency uint
	points    uint
}

// distribution maps the letters of a Scrabble game for
// a particular language to their frequency and points.
type distribution struct {
	lang      language.Tag
	name      string
	dict      []byte
	letters   []letter
	tileCount int
}

func (d distribution) dictionary() (indexedDict, error) {
	if d.dict == nil {
		return nil, nil
	}
	r, err := gzip.NewReader(bytes.NewReader(d.dict))
	if err != nil {
		return nil, err
	}
	return parseDictionary(r, d.lang)
}

func (d distribution) alphabet() []string {
	a := make([]string, 0, len(d.letters))
	for _, v := range d.letters {
		if v.L == blank {
			continue
		}
		a = append(a, v.L)
	}
	slices.Sort(a)

	return a
}

// french represents the distribution of letters for the
// French edition. It contains 102 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#French
// +----+-----------+---------+-----+----+-------------+----+----+-----+
// |    | ×1        | ×2      | ×3  | ×5 | ×6          | ×8 | ×9 | ×15 |
// +----+-----------+---------+-----+----+-------------+----+----+-----+
// | 0  |           | [blank] |     |    |             |    |    |     |
// | 1  |           |         |     | L  | N O R S T U | I  | A  | E   |
// | 2  |           | G       | D M |    |             |    |    |     |
// | 3  |           | B C P   |     |    |             |    |    |     |
// | 4  |           | F H V   |     |    |             |    |    |     |
// | 8  | J Q       |         |     |    |             |    |    |     |
// | 10 | K W X Y Z |         |     |    |             |    |    |     |
// +----+-----------+---------+-----+----+-------------+----+----+-----+
var french = distribution{
	lang: language.French,
	name: "Français",
	dict: fr.ODS8,
	letters: []letter{
		{blank, 2, 0},
		{"A", 9, 1},
		{"B", 2, 3},
		{"C", 2, 3},
		{"D", 3, 2},
		{"E", 15, 1},
		{"F", 2, 4},
		{"G", 2, 2},
		{"H", 2, 4},
		{"I", 8, 1},
		{"J", 1, 8},
		{"K", 1, 10},
		{"L", 5, 1},
		{"M", 3, 2},
		{"N", 6, 1},
		{"O", 6, 1},
		{"P", 2, 3},
		{"Q", 1, 8},
		{"R", 6, 1},
		{"S", 6, 1},
		{"T", 6, 1},
		{"U", 6, 1},
		{"V", 2, 4},
		{"W", 1, 10},
		{"X", 1, 10},
		{"Y", 1, 10},
		{"Z", 1, 10},
	},
	tileCount: 102,
}

// english represents the distribution of letters for the
// standard English edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#English
// +----+-----+-----------+----+-------+-------+----+-----+-----+
// |    | ×1  | ×2        | ×3 | ×4    | ×6    | ×8 | ×9  | ×12 |
// +----+-----+-----------+----+-------+-------+----+-----+-----+
// | 0  |     | [blank]   |    |       |       |    |     |     |
// | 1  |     |           |    | L S U | N R T | O  | A I | E   |
// | 2  |     |           | G  | D     |       |    |     |     |
// | 3  |     | B C M P   |    |       |       |    |     |     |
// | 4  |     | F H V W Y |    |       |       |    |     |     |
// | 5  | K   |           |    |       |       |    |     |     |
// | 8  | J X |           |    |       |       |    |     |     |
// | 10 | Q Z |           |    |       |       |    |     |     |
// +----+-----+-----------+----+-------+-------+----+-----+-----+
var english = distribution{
	lang: language.English,
	name: "English",
	dict: en.SOWPODS,
	letters: []letter{
		{blank, 2, 0},
		{"A", 9, 1},
		{"B", 2, 3},
		{"C", 2, 3},
		{"D", 4, 2},
		{"E", 12, 1},
		{"F", 2, 4},
		{"G", 3, 2},
		{"H", 2, 4},
		{"I", 9, 1},
		{"J", 1, 8},
		{"K", 1, 5},
		{"L", 4, 1},
		{"M", 2, 3},
		{"N", 6, 1},
		{"O", 8, 1},
		{"P", 2, 3},
		{"Q", 1, 10},
		{"R", 6, 1},
		{"S", 4, 1},
		{"T", 6, 1},
		{"U", 4, 1},
		{"V", 2, 4},
		{"W", 2, 4},
		{"X", 1, 8},
		{"Y", 2, 4},
		{"Z", 1, 10},
	},
	tileCount: 100,
}

// german represents the distribution of letters for the
// standard German edition. It contains 102 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#German
// +----+---------+---------+-------+----+----+---------+----+----+-----+
// |    | ×1      | ×2      | ×3    | ×4 | ×5 | ×6      | ×7 | ×9 | ×15 |
// +----+---------+---------+-------+----+----+---------+----+----+-----+
// | 0  |         | [blank] |       |    |    |         |    |    |     |
// | 1  |         |         |       | D  | A  | I R T U | S  | N  | E   |
// | 2  |         |         | G L O | H  |    |         |    |    |     |
// | 3  | W Z     | B       |       | M  |    |         |    |    |     |
// | 4  | P       | C F K   |       |    |    |         |    |    |     |
// | 6  | Ä J Ü V |         |       |    |    |         |    |    |     |
// | 8  | Ö X     |         |       |    |    |         |    |    |     |
// | 10 | Q Y     |         |       |    |    |         |    |    |     |
// +----+---------+---------+-------+----+----+---------+----+----+-----+
var german = distribution{
	lang: language.German,
	name: "Deutsch",
	letters: []letter{
		{blank, 2, 0},
		{"A", 5, 1},
		{"B", 2, 3},
		{"C", 2, 4},
		{"D", 4, 1},
		{"E", 15, 1},
		{"F", 2, 4},
		{"G", 3, 2},
		{"H", 4, 2},
		{"I", 6, 1},
		{"J", 1, 6},
		{"K", 2, 4},
		{"L", 3, 2},
		{"M", 4, 3},
		{"N", 9, 1},
		{"O", 3, 2},
		{"P", 1, 6},
		{"Q", 1, 10},
		{"R", 6, 1},
		{"S", 7, 1},
		{"T", 6, 1},
		{"U", 6, 1},
		{"V", 1, 6},
		{"W", 1, 3},
		{"X", 1, 8},
		{"Y", 1, 10},
		{"Z", 1, 3},
		{"Ä", 1, 6},
		{"Ö", 1, 8},
		{"Ü", 1, 6},
	},
	tileCount: 102,
}

// italian represents the distribution of letters for the
// standard Italian edition. It contains 120 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Italian
// +----+----+---------+-----------+---------+---------+-----+-----+-----+-----+
// |    | ×1 | ×2      | ×3        | ×5      | ×6      | ×11 | ×12 | ×14 | ×15 |
// +----+----+---------+-----------+---------+---------+-----+-----+-----+-----+
// | 0  |    | [blank] |           |         |         |     |     |     |     |
// | 1  |    |         |           |         |         | E   | I   | A   | O   |
// | 2  |    |         |           |         | C R S T |     |     |     |     |
// | 3  |    |         |           | L M N U |         |     |     |     |     |
// | 5  |    |         | B D F P V |         |         |     |     |     |     |
// | 8  |    | G H Z   |           |         |         |     |     |     |     |
// | 10 | Q  |         |           |         |         |     |     |     |     |
// +----+----+---------+-----------+---------+---------+-----+-----+-----+-----+
var italian = distribution{
	lang: language.Italian,
	name: "Italiano",
	letters: []letter{
		{blank, 2, 0},
		{"A", 14, 1},
		{"B", 3, 5},
		{"C", 6, 2},
		{"D", 3, 5},
		{"E", 11, 1},
		{"F", 3, 5},
		{"G", 2, 8},
		{"H", 2, 8},
		{"I", 12, 1},
		{"L", 5, 3},
		{"M", 5, 3},
		{"N", 5, 3},
		{"O", 15, 1},
		{"P", 3, 5},
		{"Q", 1, 10},
		{"R", 6, 2},
		{"S", 6, 2},
		{"T", 6, 2},
		{"U", 5, 3},
		{"V", 3, 5},
		{"Z", 2, 8},
	},
	tileCount: 120,
}

// dutch represents the distribution of letters for the
// standard Dutch edition. It contains 102 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Dutch
// +----+-----+-----------+---------+----+---------+-----+-----+-----+
// |    | ×1  | ×2        | ×3      | ×4 | ×5      | ×6  | ×10 | ×18 |
// +----+-----+-----------+---------+----+---------+-----+-----+-----+
// | 0  |     | [blank]   |         |    |         |     |     |     |
// | 1  |     |           |         | I  |         | A O | N   | E   |
// | 2  |     |           |         |    | D R S T |     |     |     |
// | 3  |     | B P       | G K L M |    |         |     |     |     |
// | 4  |     | F H J V Z | U       |    |         |     |     |     |
// | 5  |     | C W       |         |    |         |     |     |     |
// | 8  | X Y |           |         |    |         |     |     |     |
// | 10 | Q   |           |         |    |         |     |     |     |
// +----+-----+-----------+---------+----+---------+-----+-----+-----+

var dutch = distribution{
	lang: language.Dutch,
	name: "Nederlands",
	letters: []letter{
		{blank, 2, 0},
		{"A", 6, 1},
		{"B", 2, 3},
		{"C", 2, 5},
		{"D", 5, 2},
		{"E", 18, 1},
		{"F", 2, 4},
		{"G", 3, 3},
		{"H", 2, 4},
		{"I", 4, 1},
		{"J", 2, 4},
		{"K", 3, 3},
		{"L", 3, 3},
		{"M", 3, 3},
		{"N", 10, 1},
		{"O", 6, 1},
		{"P", 2, 3},
		{"Q", 1, 10},
		{"R", 5, 2},
		{"S", 5, 2},
		{"T", 5, 2},
		{"U", 3, 4},
		{"V", 2, 4},
		{"W", 2, 5},
		{"X", 1, 8},
		{"Y", 1, 8},
		{"Z", 2, 4},
	},
	tileCount: 102,
}

// czech represents the distribution of letters for the
// standard Czech edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Czech
// +----+-------+---------+-----------+---------+-------+----+
// |    | ×1    | ×2      | ×3        | ×4      | ×5    | ×6 |
// +----+-------+---------+-----------+---------+-------+----+
// | 0  |       | [blank] |           |         |       |    |
// | 1  |       |         | D K L P R | I S T V | A E N | O  |
// | 2  |       | Á J Y Z | C H Í M U |         |       |    |
// | 3  |       | B É Ě   |           |         |       |    |
// | 4  | Č Ů Ž | Ř Š Ý   |           |         |       |    |
// | 5  | F G Ú |         |           |         |       |    |
// | 6  | Ň     |         |           |         |       |    |
// | 7  | Ó Ť   |         |           |         |       |    |
// | 8  | Ď     |         |           |         |       |    |
// | 10 | X     |         |           |         |       |    |
// +----+-------+---------+-----------+---------+-------+----+
var czech = distribution{
	lang: language.Czech,
	name: "Čeština",
	letters: []letter{
		{blank, 2, 0},
		{"A", 5, 1},
		{"B", 2, 3},
		{"C", 3, 2},
		{"D", 3, 1},
		{"E", 5, 1},
		{"F", 1, 5},
		{"G", 1, 5},
		{"H", 3, 2},
		{"I", 4, 1},
		{"J", 2, 2},
		{"K", 3, 1},
		{"L", 3, 1},
		{"M", 3, 2},
		{"N", 5, 1},
		{"O", 6, 1},
		{"P", 3, 1},
		{"R", 3, 1},
		{"S", 4, 1},
		{"T", 4, 1},
		{"U", 3, 1},
		{"V", 4, 1},
		{"X", 1, 10},
		{"Y", 2, 2},
		{"Z", 2, 2},
		{"Á", 2, 2},
		{"É", 2, 3},
		{"Í", 3, 2},
		{"Ó", 1, 7},
		{"Ú", 1, 5},
		{"Ý", 2, 4},
		{"Č", 1, 4},
		{"Ď", 1, 8},
		{"Ě", 2, 3},
		{"Ň", 1, 6},
		{"Ř", 2, 4},
		{"Š", 2, 4},
		{"Ť", 1, 7},
		{"Ů", 1, 4},
		{"Ž", 1, 4},
	},
	tileCount: 100,
}

// icelandic represents the distribution of letters for the
// standard Icelandic edition. It contains 104 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Icelandic
// +----+-----------+-------------+-----+-------+----+-------+----+-----+-----+
// |    | ×1        | ×2          | ×3  | ×4    | ×5 | ×6    | ×7 | ×8  | ×10 |
// +----+-----------+-------------+-----+-------+----+-------+----+-----+-----+
// | 0  |           | [blank]     |     |       |    |       |    |     |     |
// | 1  |           |             |     |       | T  | E S U | I  | N R | A   |
// | 2  |           |             |     | Ð G L |    |       |    |     |     |
// | 3  |           |             | F K | M     |    |       |    |     |     |
// | 4  |           | Á D H Í O V |     |       |    |       |    |     |     |
// | 5  | Þ         |             |     |       |    |       |    |     |     |
// | 6  | B J Ó Y Æ |             |     |       |    |       |    |     |     |
// | 8  | É P Ú Ö   |             |     |       |    |       |    |     |     |
// | 9  | Ý         |             |     |       |    |       |    |     |     |
// | 10 | X         |             |     |       |    |       |    |     |     |
// +----+-----------+-------------+-----+-------+----+-------+----+-----+-----+
var icelandic = distribution{
	lang: language.Icelandic,
	name: "Íslenska",
	letters: []letter{
		{blank, 2, 0},
		{"A", 10, 1},
		{"B", 1, 6},
		{"D", 2, 4},
		{"E", 6, 1},
		{"F", 3, 3},
		{"G", 4, 2},
		{"H", 2, 4},
		{"I", 7, 1},
		{"J", 1, 6},
		{"K", 3, 3},
		{"L", 4, 2},
		{"M", 4, 3},
		{"N", 8, 1},
		{"O", 2, 4},
		{"P", 1, 8},
		{"R", 8, 1},
		{"S", 6, 1},
		{"T", 5, 1},
		{"U", 6, 1},
		{"V", 2, 4},
		{"X", 1, 10},
		{"Y", 1, 6},
		{"Á", 2, 4},
		{"Æ", 1, 6},
		{"É", 1, 8},
		{"Í", 2, 4},
		{"Ð", 4, 2},
		{"Ó", 1, 6},
		{"Ö", 1, 8},
		{"Ú", 1, 8},
		{"Ý", 1, 9},
		{"Þ", 1, 5},
	},
	tileCount: 104,
}

// krafla represents the distribution of letters for the
// non-standard Krafla Icelandic edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Icelandic
// +----+-------------+---------+-------+-----+----+-----+-------+----+-----+
// |    | ×1          | ×2      | ×3    | ×4  | ×5 | ×6  | ×7    | ×8 | ×11 |
// +----+-------------+---------+-------+-----+----+-----+-------+----+-----+
// | 0  |             | [blank] |       |     |    |     |       |    |     |
// | 1  |             |         |       |     |    |     | I N S | R  | A   |
// | 2  |             |         | M     | Ð K | L  | T U |       |    |     |
// | 3  |             | Á Ó     | E F G |     |    |     |       |    |     |
// | 4  | H Í Ú       | Æ       |       |     |    |     |       |    |     |
// | 5  | B D O P V Ý |         |       |     |    |     |       |    |     |
// | 6  | J Y Ö       |         |       |     |    |     |       |    |     |
// | 7  | É Þ         |         |       |     |    |     |       |    |     |
// | 10 | X           |         |       |     |    |     |       |    |     |
// +----+-------------+---------+-------+-----+----+-----+-------+----+-----+
var krafla = distribution{
	lang: language.Icelandic,
	name: "Íslenska (Krafla)",
	letters: []letter{
		{blank, 2, 0},
		{"A", 11, 1},
		{"B", 1, 5},
		{"D", 1, 5},
		{"E", 3, 3},
		{"F", 3, 3},
		{"G", 3, 3},
		{"H", 1, 4},
		{"I", 7, 1},
		{"J", 1, 6},
		{"K", 4, 2},
		{"L", 5, 2},
		{"M", 3, 3},
		{"N", 7, 1},
		{"O", 1, 5},
		{"P", 1, 5},
		{"R", 8, 1},
		{"S", 7, 1},
		{"T", 6, 2},
		{"U", 6, 2},
		{"V", 1, 5},
		{"X", 1, 10},
		{"Y", 1, 6},
		{"Á", 2, 3},
		{"Æ", 2, 4},
		{"É", 1, 7},
		{"Í", 1, 4},
		{"Ð", 4, 2},
		{"Ó", 2, 3},
		{"Ö", 1, 6},
		{"Ú", 1, 4},
		{"Ý", 1, 5},
		{"Þ", 1, 7},
	},
	tileCount: 100,
}

// afrikaans represents the distribution of letters for the
// standard Afrikaans edition. It contains 104 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Afrikaans
// +----+-----+---------+-----+----+-----------+-----+----+-----+
// |    | ×1  | ×2      | ×3  | ×4 | ×6        | ×8  | ×9 | ×16 |
// +----+-----+---------+-----+----+-----------+-----+----+-----+
// | 0  |     | [blank] |     |    |           |     |    |     |
// | 1  |     |         |     |    | D O R S T | I N | A  | E   |
// | 2  |     |         | H L | G  |           |     |    |     |
// | 3  |     |         | K W |    |           |     |    |     |
// | 4  |     | M U Y   |     |    |           |     |    |     |
// | 5  |     | P V     |     |    |           |     |    |     |
// | 8  | B F |         |     |    |           |     |    |     |
// | 10 | J   |         |     |    |           |     |    |     |
// +----+-----+---------+-----+----+-----------+-----+----+-----+
var afrikaans = distribution{
	lang: language.Afrikaans,
	name: "Afrikaans",
	letters: []letter{
		{blank, 2, 0},
		{"A", 9, 1},
		{"B", 1, 8},
		{"D", 6, 1},
		{"E", 16, 1},
		{"F", 1, 8},
		{"G", 4, 2},
		{"H", 3, 2},
		{"I", 8, 1},
		{"J", 1, 10},
		{"K", 3, 3},
		{"L", 3, 2},
		{"M", 2, 4},
		{"N", 8, 1},
		{"O", 6, 1},
		{"P", 2, 5},
		{"R", 6, 1},
		{"S", 6, 1},
		{"T", 6, 1},
		{"U", 2, 4},
		{"V", 2, 5},
		{"W", 3, 3},
		{"Y", 2, 4},
	},
	tileCount: 102, // wikipedia says 104, but the total is 102
}

// bulgarian represents the distribution of letters for the
// standard Bulgarian edition. It contains 102 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Bulgarian
// +----+-------+---------+-------+---------+----+-----+-----+
// |    | ×1    | ×2      | ×3    | ×4      | ×5 | ×8  | ×9  |
// +----+-------+---------+-------+---------+----+-----+-----+
// | 0  |       | [blank] |       |         |    |     |     |
// | 1  |       |         |       | Н П Р С | Т  | Е И | А О |
// | 2  |       |         | Б К Л | В Д М   |    |     |     |
// | 3  |       | Ъ       | Г     |         |    |     |     |
// | 4  |       | Ж З     |       |         |    |     |     |
// | 5  | Й Х   | Ч Я     | У     |         |    |     |     |
// | 8  | Ц Ш Ю |         |       |         |    |     |     |
// | 10 | Ф Щ Ь |         |       |         |    |     |     |
// +----+-------+---------+-------+---------+----+-----+-----+
var bulgarian = distribution{
	lang: language.Bulgarian,
	name: "Български",
	letters: []letter{
		{blank, 2, 0},
		{"А", 9, 1},
		{"Б", 3, 2},
		{"В", 4, 2},
		{"Г", 3, 3},
		{"Д", 4, 2},
		{"Е", 8, 1},
		{"Ж", 2, 4},
		{"З", 2, 4},
		{"И", 8, 1},
		{"Й", 1, 5},
		{"К", 3, 2},
		{"Л", 3, 2},
		{"М", 4, 2},
		{"Н", 4, 1},
		{"О", 9, 1},
		{"П", 4, 1},
		{"Р", 4, 1},
		{"С", 4, 1},
		{"Т", 5, 1},
		{"У", 3, 5},
		{"Ф", 1, 10},
		{"Х", 1, 5},
		{"Ц", 1, 8},
		{"Ч", 2, 5},
		{"Ш", 1, 8},
		{"Щ", 1, 10},
		{"Ъ", 2, 3},
		{"Ь", 1, 10},
		{"Ю", 1, 8},
		{"Я", 2, 5},
	},
	tileCount: 102,
}

// danish represents the distribution of letters for the
// standard Danish edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Danish
// +---+-----+---------------+-----------+-------+-----------+-----+----+----+
// |   | ×1  | ×2            | ×3        | ×4    | ×5        | ×6  | ×7 | ×9 |
// +---+-----+---------------+-----------+-------+-----------+-----+----+----+
// | 0 |     | [blank]       |           |       |           |     |    |    |
// | 1 |     |               |           |       |           | N R | A  | E  |
// | 2 |     |               |           |       | D L O S T |     |    |    |
// | 3 |     |               | F G M U V | B I K |           |     |    |    |
// | 4 |     | H J P Y Æ Ø Å |           |       |           |     |    |    |
// | 8 | X Z | C             |           |       |           |     |    |    |
// +---+-----+---------------+-----------+-------+-----------+-----+----+----+
var danish = distribution{
	lang: language.Danish,
	name: "Dansk",
	letters: []letter{
		{blank, 2, 0},
		{"A", 7, 1},
		{"B", 4, 3},
		{"C", 2, 8},
		{"D", 5, 2},
		{"E", 9, 1},
		{"F", 3, 3},
		{"G", 3, 3},
		{"H", 2, 4},
		{"I", 4, 3},
		{"J", 2, 4},
		{"K", 4, 3},
		{"L", 5, 2},
		{"M", 3, 3},
		{"N", 6, 1},
		{"O", 5, 2},
		{"P", 2, 4},
		{"R", 6, 1},
		{"S", 5, 2},
		{"T", 5, 2},
		{"U", 3, 3},
		{"V", 3, 3},
		{"X", 1, 8},
		{"Y", 2, 4},
		{"Z", 1, 8},
		{"Å", 2, 4},
		{"Æ", 2, 4},
		{"Ø", 2, 4},
	},
	tileCount: 100,
}

// estonian represents the distribution of letters for the
// standard Estonian edition. It contains 102 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Estonian
// +----+-------+---------+-------+---------+----+----+-----+-----+
// |    | ×1    | ×2      | ×4    | ×5      | ×7 | ×8 | ×9  | ×10 |
// +----+-------+---------+-------+---------+----+----+-----+-----+
// | 0  |       | [blank] |       |         |    |    |     |     |
// | 1  |       |         |       | K L O U | T  | S  | E I | A   |
// | 2  |       | R       | D M N |         |    |    |     |     |
// | 3  |       | G V     |       |         |    |    |     |     |
// | 4  | B     | H J P Õ |       |         |    |    |     |     |
// | 5  |       | Ä Ü     |       |         |    |    |     |     |
// | 6  |       | Ö       |       |         |    |    |     |     |
// | 8  | F     |         |       |         |    |    |     |     |
// | 10 | Š Z Ž |         |       |         |    |    |     |     |
// +----+-------+---------+-------+---------+----+----+-----+-----+
var estonian = distribution{
	lang: language.Estonian,
	name: "Eesti",
	letters: []letter{
		{blank, 2, 0},
		{"A", 10, 1},
		{"B", 1, 4},
		{"D", 4, 2},
		{"E", 9, 1},
		{"F", 1, 8},
		{"G", 2, 3},
		{"H", 2, 4},
		{"I", 9, 1},
		{"J", 2, 4},
		{"K", 5, 1},
		{"L", 5, 1},
		{"M", 4, 2},
		{"N", 4, 2},
		{"O", 5, 1},
		{"P", 2, 4},
		{"R", 2, 2},
		{"S", 8, 1},
		{"T", 7, 1},
		{"U", 5, 1},
		{"V", 2, 3},
		{"Z", 1, 10},
		{"Ä", 2, 5},
		{"Õ", 2, 4},
		{"Ö", 2, 6},
		{"Ü", 2, 5},
		{"Š", 1, 10},
		{"Ž", 1, 10},
	},
	tileCount: 102,
}

// finnish represents the distribution of letters for the
// standard Finnish edition. It contains 101 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Finnish
// +----+---------+-------------+----+----+---------+----+----+-----+-----+
// |    | ×1      | ×2          | ×3 | ×4 | ×5      | ×7 | ×8 | ×9  | ×10 |
// +----+---------+-------------+----+----+---------+----+----+-----+-----+
// | 0  |         | [blank]     |    |    |         |    |    |     |     |
// | 1  |         |             |    |    |         | S  | E  | N T | A I |
// | 2  |         |             |    |    | K L O Ä |    |    |     |     |
// | 3  |         |             | M  | U  |         |    |    |     |     |
// | 4  |         | H J P R V Y |    |    |         |    |    |     |     |
// | 7  | D Ö     |             |    |    |         |    |    |     |     |
// | 8  | B F G W |             |    |    |         |    |    |     |     |
// | 10 | C       |             |    |    |         |    |    |     |     |
// +----+---------+-------------+----+----+---------+----+----+-----+-----+
var finnish = distribution{
	lang: language.Finnish,
	name: "suomi",
	letters: []letter{
		{blank, 2, 0},
		{"A", 10, 1},
		{"B", 1, 8},
		{"C", 1, 19},
		{"D", 1, 7},
		{"E", 8, 1},
		{"F", 1, 8},
		{"G", 1, 8},
		{"H", 2, 4},
		{"I", 10, 1},
		{"J", 2, 4},
		{"K", 5, 2},
		{"L", 5, 2},
		{"M", 3, 3},
		{"N", 9, 1},
		{"O", 5, 2},
		{"P", 2, 4},
		{"R", 2, 4},
		{"S", 7, 1},
		{"T", 9, 1},
		{"U", 4, 3},
		{"V", 2, 4},
		{"W", 1, 8},
		{"Y", 2, 4},
		{"Ä", 5, 2},
		{"Ö", 1, 7},
	},
	tileCount: 101,
}

// greek represents the distribution of letters for the
// standard Greek edition. It contains 104 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Greek
// +----+---------+---------+-------+-------+----+----+-----+-------+----+-----+
// |    | ×1      | ×2      | ×3    | ×4    | ×5 | ×6 | ×7  | ×8    | ×9 | ×12 |
// +----+---------+---------+-------+-------+----+----+-----+-------+----+-----+
// | 0  |         | [blank] |       |       |    |    |     |       |    |     |
// | 1  |         |         |       |       |    | Ν  | Η Σ | Ε Ι Τ | Ο  | Α   |
// | 2  |         |         |       | Κ Π Υ | Ρ  |    |     |       |    |     |
// | 3  |         |         | Λ Μ Ω |       |    |    |     |       |    |     |
// | 4  |         | Γ Δ     |       |       |    |    |     |       |    |     |
// | 8  | Β Φ Χ   |         |       |       |    |    |     |       |    |     |
// | 10 | Ζ Θ Ξ Ψ |         |       |       |    |    |     |       |    |     |
// +----+---------+---------+-------+-------+----+----+-----+-------+----+-----+
var greek = distribution{
	lang: language.Greek,
	name: "Ελληνικά",
	letters: []letter{
		{blank, 2, 0},
		{"Α", 12, 1},
		{"Β", 1, 8},
		{"Γ", 2, 4},
		{"Δ", 2, 4},
		{"Ε", 8, 1},
		{"Ζ", 1, 10},
		{"Η", 7, 1},
		{"Θ", 1, 10},
		{"Ι", 8, 1},
		{"Κ", 4, 2},
		{"Λ", 3, 3},
		{"Μ", 3, 3},
		{"Ν", 6, 1},
		{"Ξ", 1, 10},
		{"Ο", 9, 1},
		{"Π", 4, 2},
		{"Ρ", 5, 2},
		{"Σ", 7, 1},
		{"Τ", 8, 1},
		{"Υ", 4, 2},
		{"Φ", 1, 8},
		{"Χ", 1, 8},
		{"Ψ", 1, 10},
		{"Ω", 3, 3},
	},
	tileCount: 104,
}

// indonesian represents the distribution of letters for the
// standard Indonesian edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Indonesian
// +----+-----+---------+-----+----+-----+-----+----+-----+
// |    | ×1  | ×2      | ×3  | ×4 | ×5  | ×8  | ×9 | ×19 |
// +----+-----+---------+-----+----+-----+-----+----+-----+
// | 0  |     | [blank] |     |    |     |     |    |     |
// | 1  |     |         | O S | R  | T U | E I | N  | A   |
// | 2  |     |         | K M |    |     |     |    |     |
// | 3  |     |         | G   | D  |     |     |    |     |
// | 4  |     | H P     | L   |    |     |     |    |     |
// | 5  | F V | Y       |     | B  |     |     |    |     |
// | 8  | W   |         | C   |    |     |     |    |     |
// | 10 | J Z |         |     |    |     |     |    |     |
// +----+-----+---------+-----+----+-----+-----+----+-----+
var indonesian = distribution{
	lang: language.Indonesian,
	name: "Bahasa Indonesia",
	letters: []letter{
		{blank, 2, 0},
		{"A", 19, 1},
		{"B", 4, 5},
		{"C", 3, 8},
		{"D", 4, 3},
		{"E", 8, 1},
		{"F", 1, 5},
		{"G", 3, 3},
		{"H", 2, 4},
		{"I", 8, 1},
		{"J", 1, 10},
		{"K", 3, 2},
		{"L", 3, 4},
		{"M", 3, 2},
		{"N", 9, 1},
		{"O", 3, 1},
		{"P", 2, 4},
		{"R", 4, 1},
		{"S", 3, 1},
		{"T", 5, 1},
		{"U", 5, 1},
		{"V", 1, 5},
		{"W", 1, 8},
		{"Y", 2, 5},
		{"Z", 1, 10},
	},
	tileCount: 100,
}

// latvian represents the distribution of letters for the
// standard Latvian edition. It contains 104 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Latvian
// +----+-----------+---------+-------+---------+-----+-----+----+----+-----+
// |    | ×1        | ×2      | ×3    | ×4      | ×5  | ×6  | ×8 | ×9 | ×11 |
// +----+-----------+---------+-------+---------+-----+-----+----+----+-----+
// | 0  |           | [blank] |       |         |     |     |    |    |     |
// | 1  |           |         |       |         | R U | E T | S  | I  | A   |
// | 2  |           |         | L P   | Ā K M N |     |     |    |    |     |
// | 3  |           | Z       | D O V |         |     |     |    |    |     |
// | 4  |           | Ē Ī J   |       |         |     |     |    |    |     |
// | 5  | B C G     |         |       |         |     |     |    |    |     |
// | 6  | Ņ Š Ū     |         |       |         |     |     |    |    |     |
// | 8  | Ļ Ž       |         |       |         |     |     |    |    |     |
// | 10 | Č F Ģ H Ķ |         |       |         |     |     |    |    |     |
// +----+-----------+---------+-------+---------+-----+-----+----+----+-----+
var latvian = distribution{
	lang: language.Latvian,
	name: "Latviešu",
	letters: []letter{
		{blank, 2, 0},
		{"A", 11, 1},
		{"B", 1, 5},
		{"C", 1, 5},
		{"D", 3, 3},
		{"E", 6, 1},
		{"F", 1, 10},
		{"G", 1, 5},
		{"H", 1, 10},
		{"I", 9, 1},
		{"J", 2, 4},
		{"K", 4, 2},
		{"L", 3, 2},
		{"M", 4, 2},
		{"N", 4, 2},
		{"O", 3, 3},
		{"P", 3, 2},
		{"R", 5, 1},
		{"S", 8, 1},
		{"T", 6, 1},
		{"U", 5, 1},
		{"V", 3, 3},
		{"Z", 2, 3},
		{"Ā", 4, 2},
		{"Č", 1, 10},
		{"Ē", 2, 4},
		{"Ģ", 1, 10},
		{"Ī", 2, 4},
		{"Ķ", 1, 10},
		{"Ļ", 1, 8},
		{"Ņ", 1, 6},
		{"Š", 1, 6},
		{"Ū", 1, 6},
		{"Ž", 1, 8},
	},
	tileCount: 104,
}

// lithuanian represents the distribution of letters for the
// standard Lithuanian edition. It contains 104 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Lithuanian
// +----+-----------+---------+-------+-----+-------+-----+----+-----+-----+
// |    | ×1        | ×2      | ×3    | ×4  | ×5    | ×6  | ×8 | ×12 | ×13 |
// +----+-----------+---------+-------+-----+-------+-----+----+-----+-----+
// | 0  |           | [blank] |       |     |       |     |    |     |     |
// | 1  |           |         |       | K U | E N R | O T | S  | A   | I   |
// | 2  | B         |         | D L M |     |       |     |    |     |     |
// | 3  |           |         | P     |     |       |     |    |     |     |
// | 4  |           | Ė G J V |       |     |       |     |    |     |     |
// | 5  | Š Y       |         |       |     |       |     |    |     |     |
// | 6  | Ų Ž       |         |       |     |       |     |    |     |     |
// | 8  | Ą Č Į Ū   |         |       |     |       |     |    |     |     |
// | 10 | C Ę F H Z |         |       |     |       |     |    |     |     |
// +----+-----------+---------+-------+-----+-------+-----+----+-----+-----+
var lithuanian = distribution{
	lang: language.Lithuanian,
	name: "Lietuvių",
	letters: []letter{
		{blank, 2, 0},
		{"A", 12, 1},
		{"B", 1, 2},
		{"C", 1, 10},
		{"D", 3, 2},
		{"E", 5, 1},
		{"F", 1, 10},
		{"G", 2, 4},
		{"H", 1, 10},
		{"I", 13, 1},
		{"J", 2, 4},
		{"K", 4, 1},
		{"L", 3, 2},
		{"M", 3, 2},
		{"N", 5, 1},
		{"O", 6, 1},
		{"P", 3, 3},
		{"R", 5, 1},
		{"S", 8, 1},
		{"T", 6, 1},
		{"U", 4, 1},
		{"V", 2, 4},
		{"Y", 1, 5},
		{"Z", 1, 10},
		{"Ą", 1, 8},
		{"Č", 1, 8},
		{"Ė", 2, 4},
		{"Ę", 1, 10},
		{"Į", 1, 8},
		{"Š", 1, 5},
		{"Ū", 1, 8},
		{"Ų", 1, 6},
		{"Ž", 1, 6},
	},
	tileCount: 104,
}

// malay represents the distribution of letters for the
// standard Malay edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Malay
// +----+-----+---------+-----+-----+-------+-----+-----+----+-----+
// |    | ×1  | ×2      | ×3  | ×4  | ×5    | ×6  | ×7  | ×8 | ×19 |
// +----+-----+---------+-----+-----+-------+-----+-----+----+-----+
// | 0  |     | [blank] |     |     |       |     |     |    |     |
// | 1  |     |         |     |     | M R T | K U | E I | N  | A   |
// | 2  |     |         |     | L S |       |     |     |    |     |
// | 3  |     |         | B D | G   |       |     |     |    |     |
// | 4  |     | H O P   |     |     |       |     |     |    |     |
// | 5  | J Y |         |     |     |       |     |     |    |     |
// | 8  | C W |         |     |     |       |     |     |    |     |
// | 10 | F Z |         |     |     |       |     |     |    |     |
// +----+-----+---------+-----+-----+-------+-----+-----+----+-----+
var malay = distribution{
	lang: language.Malay,
	name: "Bahasa Melayu",
	letters: []letter{
		{blank, 2, 0},
		{"A", 19, 1},
		{"B", 3, 3},
		{"C", 1, 8},
		{"D", 3, 3},
		{"E", 7, 1},
		{"F", 1, 10},
		{"G", 4, 3},
		{"H", 2, 4},
		{"I", 7, 1},
		{"J", 1, 5},
		{"K", 6, 1},
		{"L", 4, 2},
		{"M", 5, 1},
		{"N", 8, 1},
		{"O", 2, 4},
		{"P", 2, 4},
		{"R", 5, 1},
		{"S", 4, 2},
		{"T", 5, 1},
		{"U", 6, 1},
		{"W", 1, 8},
		{"Y", 1, 5},
		{"Z", 1, 10},
	},
	tileCount: 100,
}

// norwegian represents the distribution of letters for the
// standard Norwegian edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Norwegian
// +----+-----+---------+-------+---------+-------+---------+----+----+
// |    | ×1  | ×2      | ×3    | ×4      | ×5    | ×6      | ×7 | ×9 |
// +----+-----+---------+-------+---------+-------+---------+----+----+
// | 0  |     | [blank] |       |         |       |         |    |    |
// | 1  |     |         |       |         | D I L | N R S T | A  | E  |
// | 2  |     |         | M     | F G K O |       |         |    |    |
// | 3  |     |         | H     |         |       |         |    |    |
// | 4  |     | J P Å   | B U V |         |       |         |    |    |
// | 5  |     | Ø       |       |         |       |         |    |    |
// | 6  | Y Æ |         |       |         |       |         |    |    |
// | 8  | W   |         |       |         |       |         |    |    |
// | 10 | C   |         |       |         |       |         |    |    |
// +----+-----+---------+-------+---------+-------+---------+----+----+
var norwegian = distribution{
	lang: language.Norwegian,
	name: "Norsk",
	letters: []letter{
		{blank, 2, 0},
		{"A", 7, 1},
		{"B", 3, 4},
		{"C", 1, 10},
		{"D", 5, 1},
		{"E", 9, 1},
		{"F", 4, 2},
		{"G", 4, 2},
		{"H", 3, 3},
		{"I", 5, 1},
		{"J", 2, 4},
		{"K", 4, 2},
		{"L", 5, 1},
		{"M", 3, 2},
		{"N", 6, 1},
		{"O", 4, 2},
		{"P", 2, 4},
		{"R", 6, 1},
		{"S", 6, 1},
		{"T", 6, 1},
		{"U", 3, 4},
		{"V", 3, 4},
		{"W", 1, 8},
		{"Y", 1, 6},
		{"Å", 2, 4},
		{"Æ", 1, 6},
		{"Ø", 2, 5},
	},
	tileCount: 100,
}

// polish represents the distribution of letters for the
// standard Polish edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Polish
// +---+-------------+-------------+---------------+-------+-----+----+----+----+----+
// |   | ×1          | ×2          | ×3            | ×4    | ×5  | ×6 | ×7 | ×8 | ×9 |
// +---+-------------+-------------+---------------+-------+-----+----+----+----+----+
// | 0 |             | [blank]     |               |       |     |    |    |    |    |
// | 1 |             |             |               | R S W | N Z | O  | E  | I  | A  |
// | 2 |             |             | C D K L M P T | Y     |     |    |    |    |    |
// | 3 |             | B G H J Ł U |               |       |     |    |    |    |    |
// | 5 | Ą Ę F Ó Ś Ż |             |               |       |     |    |    |    |    |
// | 6 | Ć           |             |               |       |     |    |    |    |    |
// | 7 | Ń           |             |               |       |     |    |    |    |    |
// | 9 | Ź           |             |               |       |     |    |    |    |    |
// +---+-------------+-------------+---------------+-------+-----+----+----+----+----+
var polish = distribution{
	lang: language.Polish,
	name: "Polski",
	letters: []letter{
		{blank, 2, 0},
		{"A", 9, 1},
		{"B", 2, 3},
		{"C", 3, 2},
		{"D", 3, 2},
		{"E", 7, 1},
		{"F", 1, 5},
		{"G", 2, 3},
		{"H", 2, 3},
		{"I", 8, 1},
		{"J", 2, 3},
		{"K", 3, 2},
		{"L", 3, 2},
		{"M", 3, 2},
		{"N", 5, 1},
		{"O", 6, 1},
		{"P", 3, 2},
		{"R", 4, 1},
		{"S", 4, 1},
		{"T", 3, 2},
		{"U", 2, 3},
		{"W", 4, 1},
		{"Y", 4, 2},
		{"Z", 5, 1},
		{"Ó", 1, 5},
		{"Ą", 1, 5},
		{"Ć", 1, 6},
		{"Ę", 1, 5},
		{"Ł", 2, 3},
		{"Ń", 1, 7},
		{"Ś", 1, 5},
		{"Ź", 1, 9},
		{"Ż", 1, 5},
	},
	tileCount: 100,
}

// portuguese represents the distribution of letters for the
// standard Portuguese edition. It contains 120 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Portuguese
// +---+-----+---------+---------+-----+-----+-----+----+----+-----+-----+-----+
// |   | ×1  | ×2      | ×3      | ×4  | ×5  | ×6  | ×7 | ×8 | ×10 | ×11 | ×14 |
// +---+-----+---------+---------+-----+-----+-----+----+----+-----+-----+-----+
// | 0 |     |         | [blank] |     |     |     |    |    |     |     |     |
// | 1 |     |         |         |     | T   | M R | U  | S  | I O | E   | A   |
// | 2 |     |         |         | C P | D L |     |    |    |     |     |     |
// | 3 |     | Ç       | B       | N   |     |     |    |    |     |     |     |
// | 4 |     | F G H V |         |     |     |     |    |    |     |     |     |
// | 5 |     | J       |         |     |     |     |    |    |     |     |     |
// | 6 | Q   |         |         |     |     |     |    |    |     |     |     |
// | 8 | X Z |         |         |     |     |     |    |    |     |     |     |
// +---+-----+---------+---------+-----+-----+-----+----+----+-----+-----+-----+
var portuguese = distribution{
	lang: language.Portuguese,
	name: "Português",
	letters: []letter{
		{blank, 3, 0},
		{"A", 14, 1},
		{"B", 3, 3},
		{"C", 4, 2},
		{"D", 5, 2},
		{"E", 11, 1},
		{"F", 2, 4},
		{"G", 2, 4},
		{"H", 2, 4},
		{"I", 10, 1},
		{"J", 2, 5},
		{"L", 5, 2},
		{"M", 6, 1},
		{"N", 4, 3},
		{"O", 10, 1},
		{"P", 4, 2},
		{"Q", 1, 6},
		{"R", 6, 1},
		{"S", 8, 1},
		{"T", 5, 1},
		{"U", 7, 1},
		{"V", 2, 4},
		{"X", 1, 8},
		{"Z", 1, 8},
		{"Ç", 2, 3},
	},
	tileCount: 120,
}

// romanian represents the distribution of letters for the
// standard Romanian edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Romanian
// +----+-----+---------+----+----+-------+-------+----+----+-----+-----+
// |    | ×1  | ×2      | ×3 | ×4 | ×5    | ×6    | ×7 | ×9 | ×10 | ×11 |
// +----+-----+---------+----+----+-------+-------+----+----+-----+-----+
// | 0  |     | [blank] |    |    |       |       |    |    |     |     |
// | 1  |     |         |    |    | C L U | N R S | T  | E  | A   | I   |
// | 2  |     |         |    | P  | O     |       |    |    |     |     |
// | 3  |     |         |    | D  |       |       |    |    |     |     |
// | 4  |     | F V     | M  |    |       |       |    |    |     |     |
// | 5  |     | B       |    |    |       |       |    |    |     |     |
// | 6  |     | G       |    |    |       |       |    |    |     |     |
// | 8  | H Z |         |    |    |       |       |    |    |     |     |
// | 10 | J X |         |    |    |       |       |    |    |     |     |
// +----+-----+---------+----+----+-------+-------+----+----+-----+-----+
var romanian = distribution{
	lang: language.Romanian,
	name: "Română",
	letters: []letter{
		{blank, 2, 0},
		{"A", 10, 1},
		{"B", 2, 5},
		{"C", 5, 1},
		{"D", 4, 3},
		{"E", 9, 1},
		{"F", 2, 4},
		{"G", 2, 6},
		{"H", 1, 8},
		{"I", 11, 1},
		{"J", 1, 10},
		{"L", 5, 1},
		{"M", 3, 4},
		{"N", 6, 1},
		{"O", 5, 2},
		{"P", 4, 2},
		{"R", 6, 1},
		{"S", 6, 1},
		{"T", 7, 1},
		{"U", 5, 1},
		{"V", 2, 4},
		{"X", 1, 10},
		{"Z", 1, 8},
	},
	tileCount: 100,
}

// slovak represents the distribution of letters for the
// standard Slovak edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Slovak
// +----+-----------+---------+---------+---------+-----+----+-----+
// |    | ×1        | ×2      | ×3      | ×4      | ×5  | ×8 | ×9  |
// +----+-----------+---------+---------+---------+-----+----+-----+
// | 0  |           | [blank] |         |         |     |    |     |
// | 1  |           |         |         | R S T V | I N | E  | A O |
// | 2  |           |         | D K L P | M       |     |    |     |
// | 3  |           | J U     |         |         |     |    |     |
// | 4  | Á C H Y Z | B       |         |         |     |    |     |
// | 5  | Č Í Š Ý Ž |         |         |         |     |    |     |
// | 7  | É Ľ Ť Ú   |         |         |         |     |    |     |
// | 8  | Ď F G Ň Ô |         |         |         |     |    |     |
// | 10 | Ä Ĺ Ó Ŕ X |         |         |         |     |    |     |
// +----+-----------+---------+---------+---------+-----+----+-----+
var slovak = distribution{
	lang: language.Slovak,
	name: "Slovenčina",
	letters: []letter{
		{blank, 2, 0},
		{"A", 9, 1},
		{"B", 2, 4},
		{"C", 1, 4},
		{"D", 3, 2},
		{"E", 8, 1},
		{"F", 1, 8},
		{"G", 1, 8},
		{"H", 1, 4},
		{"I", 5, 1},
		{"J", 2, 3},
		{"K", 3, 2},
		{"L", 3, 2},
		{"M", 4, 2},
		{"N", 5, 1},
		{"O", 9, 1},
		{"P", 3, 2},
		{"R", 4, 1},
		{"S", 4, 1},
		{"T", 4, 1},
		{"U", 2, 3},
		{"V", 4, 1},
		{"X", 1, 10},
		{"Y", 1, 4},
		{"Z", 1, 4},
		{"Á", 1, 4},
		{"Ä", 1, 10},
		{"É", 1, 7},
		{"Í", 1, 5},
		{"Ó", 1, 10},
		{"Ô", 1, 8},
		{"Ú", 1, 7},
		{"Ý", 1, 5},
		{"Č", 1, 5},
		{"Ď", 1, 8},
		{"Ĺ", 1, 10},
		{"Ľ", 1, 7},
		{"Ň", 1, 8},
		{"Ŕ", 1, 10},
		{"Š", 1, 5},
		{"Ť", 1, 7},
		{"Ž", 1, 5},
	},
	tileCount: 100,
}

// slovenian represents the distribution of letters for the
// standard Slovenian edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Slovenian
// +----+-----+---------+----+-------+-----+----+----+----+-----+-----+
// |    | ×1  | ×2      | ×3 | ×4    | ×6  | ×7 | ×8 | ×9 | ×10 | ×11 |
// +----+-----+---------+----+-------+-----+----+----+----+-----+-----+
// | 0  |     | [blank] |    |       |     |    |    |    |     |     |
// | 1  |     |         |    | J L T | R S | N  | O  | I  | A   | E   |
// | 2  |     |         |    | D V   |     |    |    |    |     |     |
// | 3  |     | M P U   | K  |       |     |    |    |    |     |     |
// | 4  |     | B G Z   |    |       |     |    |    |    |     |     |
// | 5  | Č H |         |    |       |     |    |    |    |     |     |
// | 6  | Š   |         |    |       |     |    |    |    |     |     |
// | 8  | C   |         |    |       |     |    |    |    |     |     |
// | 10 | F Ž |         |    |       |     |    |    |    |     |     |
// +----+-----+---------+----+-------+-----+----+----+----+-----+-----+
var slovenian = distribution{
	lang: language.Slovenian,
	name: "Slovenščina",
	letters: []letter{
		{blank, 2, 0},
		{"A", 10, 1},
		{"B", 2, 4},
		{"C", 1, 8},
		{"D", 4, 2},
		{"E", 11, 1},
		{"F", 1, 10},
		{"G", 2, 4},
		{"H", 1, 5},
		{"I", 9, 1},
		{"J", 4, 1},
		{"K", 3, 3},
		{"L", 4, 1},
		{"M", 2, 3},
		{"N", 7, 1},
		{"O", 8, 1},
		{"P", 2, 3},
		{"R", 6, 1},
		{"S", 6, 1},
		{"T", 4, 1},
		{"U", 2, 3},
		{"V", 4, 2},
		{"Z", 2, 4},
		{"Č", 1, 5},
		{"Š", 1, 6},
		{"Ž", 1, 10},
	},
	tileCount: 100,
}

// swedish represents the distribution of letters for the
// standard Swedish edition. It contains 100 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Swedish
// +----+-----+---------+-------+-------+----+----+---------+
// |    | ×1  | ×2      | ×3    | ×5    | ×6 | ×7 | ×8      |
// +----+-----+---------+-------+-------+----+----+---------+
// | 0  |     | [blank] |       |       |    |    |         |
// | 1  |     |         |       | D I L | N  | E  | A R S T |
// | 2  |     | H       | G K M | O     |    |    |         |
// | 3  |     | F V Ä   |       |       |    |    |         |
// | 4  |     | B P Ö Å | U     |       |    |    |         |
// | 7  | J Y |         |       |       |    |    |         |
// | 8  | C X |         |       |       |    |    |         |
// | 10 | Z   |         |       |       |    |    |         |
// +----+-----+---------+-------+-------+----+----+---------+
var swedish = distribution{
	lang: language.Swedish,
	name: "Svenska",
	letters: []letter{
		{blank, 2, 0},
		{"A", 8, 1},
		{"B", 2, 4},
		{"C", 1, 8},
		{"D", 5, 1},
		{"E", 7, 1},
		{"F", 2, 3},
		{"G", 3, 2},
		{"H", 2, 2},
		{"I", 5, 1},
		{"J", 1, 7},
		{"K", 3, 2},
		{"L", 5, 1},
		{"M", 3, 2},
		{"N", 6, 1},
		{"O", 5, 2},
		{"P", 2, 4},
		{"R", 8, 1},
		{"S", 8, 1},
		{"T", 8, 1},
		{"U", 3, 4},
		{"V", 2, 3},
		{"X", 1, 8},
		{"Y", 1, 7},
		{"Z", 1, 10},
		{"Ä", 2, 3},
		{"Å", 2, 4},
		{"Ö", 2, 4},
	},
	tileCount: 100,
}

// ukrainian represents the distribution of letters for the
// standard Ukrainian edition. It contains 104 tiles.
// https://en.wikipedia.org/wiki/Scrabble_letter_distributions#Ukrainian
// +----+---------+---------+-------+-------+---------+-----+----+-----+
// |    | ×1      | ×2      | ×3    | ×4    | ×5      | ×7  | ×8 | ×10 |
// +----+---------+---------+-------+-------+---------+-----+----+-----+
// | 0  |         | [blank] |       |       |         |     |    |     |
// | 1  |         |         |       | В     | Е І Т Р | И Н | А  | О   |
// | 2  |         |         | Д П Л | К С М |         |     |    |     |
// | 3  |         |         | У     |       |         |     |    |     |
// | 4  |         | З Я Б Г |       |       |         |     |    |     |
// | 5  | Х Й Ч Ь |         |       |       |         |     |    |     |
// | 6  | Ж Ї Ц Ш |         |       |       |         |     |    |     |
// | 7  | Ю       |         |       |       |         |     |    |     |
// | 8  | Є Ф Щ   |         |       |       |         |     |    |     |
// | 10 | Ґ '     |         |       |       |         |     |    |     |
// +----+---------+---------+-------+-------+---------+-----+----+-----+
var ukrainian = distribution{
	lang: language.Ukrainian,
	name: "Українська",
	letters: []letter{
		{blank, 2, 0},
		{"'", 1, 10},
		{"А", 8, 1},
		{"Б", 2, 4},
		{"В", 4, 1},
		{"Г", 2, 4},
		{"Д", 3, 2},
		{"Е", 5, 1},
		{"Ж", 1, 6},
		{"З", 2, 4},
		{"И", 7, 1},
		{"Й", 1, 5},
		{"К", 4, 2},
		{"Л", 3, 2},
		{"М", 4, 2},
		{"Н", 7, 1},
		{"О", 10, 1},
		{"П", 3, 2},
		{"Р", 5, 1},
		{"С", 4, 2},
		{"Т", 5, 1},
		{"У", 3, 3},
		{"Ф", 1, 8},
		{"Х", 1, 5},
		{"Ц", 1, 6},
		{"Ч", 1, 5},
		{"Ш", 1, 6},
		{"Щ", 1, 8},
		{"Ь", 1, 5},
		{"Ю", 1, 7},
		{"Я", 2, 4},
		{"Є", 1, 8},
		{"І", 5, 1},
		{"Ї", 1, 6},
		{"Ґ", 1, 10},
	},
	tileCount: 104,
}

// sorted by addition time
var distributions = map[string]distribution{
	"french":     french,
	"english":    english,
	"german":     german,
	"italian":    italian,
	"dutch":      dutch,
	"czech":      czech,
	"icelandic":  icelandic,
	"krafla":     krafla,
	"afrikaans":  afrikaans,
	"bulgarian":  bulgarian,
	"danish":     danish,
	"estonian":   estonian,
	"finnish":    finnish,
	"greek":      greek,
	"indonesian": indonesian,
	"latvian":    latvian,
	"lithuanian": lithuanian,
	"malay":      malay,
	"norwegian":  norwegian,
	"polish":     polish,
	"portuguese": portuguese,
	"romanian":   romanian,
	"slovak":     slovak,
	"slovenian":  slovenian,
	"swedish":    swedish,
	"ukrainian":  ukrainian,
}
