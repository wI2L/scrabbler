package cmd

const blank = '?'

// distribution maps the tiles of a Scrabble game for
// a particular language to their frequency and points.
type distribution map[rune]letterProps

// french represents the distribution of tiles for the
// French language edition. It contains 102 tiles.
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
	blank: {2, 0},
	'A':   {9, 1},
	'B':   {2, 3},
	'C':   {2, 3},
	'D':   {3, 2},
	'E':   {15, 1},
	'F':   {2, 4},
	'G':   {2, 2},
	'H':   {2, 4},
	'I':   {8, 1},
	'J':   {1, 8},
	'K':   {1, 10},
	'L':   {5, 1},
	'M':   {3, 2},
	'N':   {6, 1},
	'O':   {6, 1},
	'P':   {2, 3},
	'Q':   {1, 8},
	'R':   {6, 1},
	'S':   {6, 1},
	'T':   {6, 1},
	'U':   {6, 1},
	'V':   {2, 4},
	'W':   {1, 10},
	'X':   {1, 10},
	'Y':   {1, 10},
	'Z':   {1, 10},
}

// english represents the distribution of tiles for the
// English language edition. It contains 100 tiles.
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
	blank: {2, 0},
	'A':   {9, 1},
	'B':   {2, 3},
	'C':   {2, 3},
	'D':   {4, 2},
	'E':   {12, 1},
	'F':   {2, 4},
	'G':   {3, 2},
	'H':   {2, 4},
	'I':   {9, 1},
	'J':   {1, 8},
	'K':   {1, 5},
	'L':   {4, 1},
	'M':   {2, 3},
	'N':   {6, 1},
	'O':   {8, 1},
	'P':   {2, 3},
	'Q':   {1, 10},
	'R':   {6, 1},
	'S':   {4, 1},
	'T':   {6, 1},
	'U':   {4, 1},
	'V':   {2, 4},
	'W':   {2, 4},
	'X':   {1, 8},
	'Y':   {2, 4},
	'Z':   {1, 10},
}

var distributions = map[string]distribution{
	"french":  french,
	"english": english,
}
