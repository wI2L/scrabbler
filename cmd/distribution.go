package cmd

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"slices"

	"golang.org/x/text/language"
)

const (
	blank          = "?"
	defaultDistrib = "french" // 🇫🇷
)

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
	lang    language.Tag
	dict    []byte
	letters []letter
}

func (d distribution) dictionary() (indexedDict, error) {
	if d.dict == nil {
		return nil, nil
	}
	r, err := gzip.NewReader(bytes.NewReader(d.dict))
	if err != nil {
		return nil, err
	}
	return parseDictionary(r)
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
	dict: ods8,
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
	dict: sowpods,
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
}

// czech represents the distribution of letters for the
// standard Czech edition. It contains 102 tiles.
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
}

// afrikaans represents the distribution of letters for the
// standard Afrikaans edition. It contains 104 tiles.
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
}

// bulgarian represents the distribution of letters for the
// standard Bulgarian edition. It contains 102 tiles.
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
}

// sorted by addition time
var distributions = map[string]distribution{
	"french":    french,
	"english":   english,
	"german":    german,
	"italian":   italian,
	"dutch":     dutch,
	"czech":     czech,
	"icelandic": icelandic,
	"krafla":    krafla,
	"afrikaans": afrikaans,
	"bulgarian": bulgarian,
}
