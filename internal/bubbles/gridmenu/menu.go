package gridmenu

import (
	"math"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Choice struct {
	Name        string
	Description string
}

type Styles struct {
	Choice    lipgloss.Style
	Selection lipgloss.Style
}

type matrix struct {
	cols  int
	rows  int
	width int
	grid  [][]*Choice
}

type margins struct {
	vertical   int
	horizontal int
}

type Model struct {
	Width     int
	choices   []Choice
	selection *Choice
	matrixes  []matrix
	grid      [][]*Choice
	keys      keyMap
	help      help.Model
	margins   margins
	styles    Styles
	limit     int
	rows      int
	cols      int
	posX      int
	posY      int
}

func New(choices []Choice, maxColumns int) Model {
	m := Model{
		choices:   choices,
		selection: &choices[0],
		limit:     maxColumns,
		keys:      keys,
		help:      help.New(),
		styles: Styles{
			Choice:    lipgloss.NewStyle().Faint(true),
			Selection: lipgloss.NewStyle().Underline(true),
		},
	}
	m.help.FullSeparator = strings.Repeat(" ", 3)
	m.initMatrixes()
	m.setActiveGrid()

	return m
}

func (m Model) Selection() string {
	return m.selection.Name
}

func (m *Model) Margin(i ...int) {
	// 0: apply default values
	// 1: apply to all directions
	// 2: horizontal -> vertical
	switch len(i) {
	case 0:
		m.margins.horizontal = 2
		m.margins.vertical = 1
	case 1:
		m.margins.horizontal = clampColumns(i[0])
		m.margins.vertical = clampColumns(i[0])
	case 2:
		m.margins.horizontal = clampColumns(i[0])
		m.margins.vertical = clampColumns(i[1])
	}
	m.initMatrixes()
	m.setActiveGrid()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.Width = msg.Width
		m.setActiveGrid()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			// If we aren't in the top row, move up.
			if m.posY > 0 {
				m.posY--
			}
		case key.Matches(msg, m.keys.Down):
			// If we aren't in the bottom row, and
			// the item below isn't nil, move down.
			if m.posY < m.rows-1 && m.grid[m.posY+1][m.posX] != nil {
				m.posY++
			}
		case key.Matches(msg, m.keys.Left):
			// If we aren't in the leftmost column, move left.
			if m.posX > 0 {
				m.posX--
			}
		case key.Matches(msg, m.keys.Right):
			// If we aren't in the rightmost column,
			// and the item on the right isn't nil,
			// move right.
			if m.posX < m.cols-1 && m.grid[m.posY][m.posX+1] != nil {
				m.posX++
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	m.selection = m.grid[m.posY][m.posX]

	return m, nil
}

func (m Model) View() string {
	cols := make([]string, 0, m.cols)

	for x := 0; x < m.cols; x++ {
		var items []string

		for y := 0; y < m.rows; y++ {
			c := m.grid[y][x]
			if c == nil {
				continue
			}
			var s string
			if c == m.selection {
				s = m.styles.Selection.Render(c.Description)
			} else {
				s = m.styles.Choice.Render(c.Description)
			}
			items = append(items, s)
		}
		// Join the items vertically with newlines.
		s := strings.Join(
			items,
			strings.Repeat("\n", m.margins.vertical),
		)
		if x != m.cols-1 {
			// Calculate inner column padding equal to
			// the length of the longest item minus the
			// length of the last item.
			pad := longestString(items) - lipgloss.Width(items[len(items)-1])

			// Add inner column padding and horizontal margin.
			s += strings.Repeat(" ", pad+m.margins.horizontal)
		}
		cols = append(cols, s)
	}
	sb := strings.Builder{}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cols...))

	var hv string
	if !m.help.ShowAll {
		hv = m.help.ShortHelpView(m.keys.ShortHelp())
	} else {
		hs := m.help.Styles.FullSeparator.Render(m.help.FullSeparator)
		hw := lipgloss.Width(hs)
		st := lipgloss.NewStyle().MarginLeft(hw)
		hv = st.Render(m.help.FullHelpView(m.keys.FullHelp()))
	}
	sb.WriteString(strings.Repeat("\n", 3))
	sb.WriteString(hv)

	return sb.String()
}

func (m *Model) initMatrixes() {
	cols := m.limit
	if cols < 1 {
		cols = 1
	}
	m.matrixes = make([]matrix, 0, cols)

	// Create N matrixes from 1 column to max
	// columns and record the maximum width occupied
	// as the sum of the longest element of each
	// column plus margins.
	for ; cols > 0; cols-- {
		rows := int(math.Ceil(float64(len(m.choices)) / float64(cols)))
		table := make([][]*Choice, rows)

		for i := range table {
			table[i] = make([]*Choice, cols)
		}
		var x, y int
		for i := range m.choices {
			table[y][x] = &m.choices[i]
			x++
			if x > cols-1 {
				y++
				x = 0
			}
		}
		var w int
		for x := 0; x < cols; x++ {
			var colWidth int

			for y := 0; y < rows; y++ {
				c := table[y][x]
				if c == nil {
					continue
				}
				if w := lipgloss.Width(c.Description); w > colWidth {
					colWidth = w
				}
			}
			w += colWidth
		}
		w += (cols - 1) * m.margins.horizontal

		m.matrixes = append(m.matrixes, matrix{
			cols:  cols,
			rows:  rows,
			width: w,
			grid:  table,
		})
	}
	sort.Slice(m.matrixes, func(i, j int) bool {
		return m.matrixes[i].width < m.matrixes[j].width
	})
}

func (m *Model) setActiveGrid() {
	for i, v := range m.matrixes {
		if v.width > m.Width && i != 0 {
			break
		}
		m.cols = v.cols
		m.rows = v.rows
		m.grid = v.grid
	}
	m.posX, m.posY = m.getCurrentCoords()
}

func (m Model) getCurrentCoords() (int, int) {
	for y, row := range m.grid {
		for x, c := range row {
			if c == m.selection {
				return x, y
			}
		}
	}
	return 0, 0
}

func longestString(strs []string) int {
	var length int
	for _, s := range strs {
		if n := lipgloss.Width(s); n > length {
			length = n
		}
	}
	return length
}

func clampColumns(i int) int {
	if i < 1 {
		return 1
	}
	return i
}
