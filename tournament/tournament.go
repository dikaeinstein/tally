package tournament

import (
	"encoding/csv"
	"io"
	"sort"
)

// Outcome refers to the result of a match(win|draw|loss)
// The value of an outcome is the point(s) earned.
type Outcome int

const (
	Loss Outcome = 0
	Draw Outcome = 1
	Win  Outcome = 3
)

const (
	Asc  = "asc"
	Desc = "desc"
)

type Match struct {
	Home string
	Away string
	// refers to the result of the home team
	Outcome
}

type TableRow struct {
	Team         string
	MatchPlayed  int
	MatchesWon   int
	MatchesDrawn int
	MatchesLost  int
	Points       int
}

// Table represents the tournament table ranking teams
// by points earned
type Table struct {
	Rows []TableRow
}

func (t *Table) Sort(order string) *Table {
	sort.SliceStable(t.Rows, func(i, j int) bool {
		if t.Rows[i].Points == t.Rows[j].Points {
			if order == Asc {
				return t.Rows[i].Team < t.Rows[j].Team
			} else {
				return t.Rows[i].Team > t.Rows[j].Team
			}
		}

		if order == Asc {
			return t.Rows[i].Points < t.Rows[j].Points
		}
		return t.Rows[i].Points > t.Rows[j].Points
	})

	return t
}

func ParseInput(input io.Reader) ([]Match, error) {
	r := csv.NewReader(input)
	r.Comma = ';'
	r.FieldsPerRecord = 3

	var matches []Match

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return matches, err
		}

		match := parseRecord(record)
		matches = append(matches, match)
	}

	return matches, nil
}

func parseRecord(record []string) Match {
	return Match{
		Home:    record[0],
		Away:    record[1],
		Outcome: parseOutcome(record[2]),
	}
}

func parseOutcome(outcome string) Outcome {
	outcomes := map[string]Outcome{
		"draw": Draw,
		"loss": Loss,
		"win":  Win,
	}

	return outcomes[outcome]
}

// Tally builds the tournament table from match results and sorts
// the table based on points earned by teams in descending order.
// It is a semantic helper. If you want more control, you can
// build the table with `tournament.BuildTable` and sort the resulting table
// with Table.Sort in any order(Desc|Asc) that you choose.
func Tally(matches []Match) *Table {
	t := BuildTable(matches)
	return t.Sort(Desc)
}

func BuildTable(matches []Match) *Table {
	// build TableRows from matches
	seen := map[string]TableRow{}
	for _, m := range matches {
		if homeRow, ok := seen[m.Home]; !ok {
			homeRow.Team = m.Home
			homeRow.MatchPlayed = 1
			homeRow = setRowOutcome(m.Outcome, homeRow)

			seen[m.Home] = homeRow
		} else {
			homeRow.MatchPlayed += 1
			homeRow = setRowOutcome(m.Outcome, homeRow)

			seen[m.Home] = homeRow
		}

		if awayRow, ok := seen[m.Away]; !ok {
			awayRow.Team = m.Away
			awayRow.MatchPlayed = 1
			awayRow = setRowOutcome(awayOutcome(m.Outcome), awayRow)

			seen[m.Away] = awayRow
		} else {
			awayRow.MatchPlayed += 1
			awayRow = setRowOutcome(awayOutcome(m.Outcome), awayRow)

			seen[m.Away] = awayRow
		}
	}

	// build table
	var t Table
	for _, row := range seen {
		t.Rows = append(t.Rows, row)
	}

	return &t
}

func awayOutcome(outcome Outcome) Outcome {
	mapToAwayOutcome := map[Outcome]Outcome{
		Win:  Loss,
		Draw: Draw,
		Loss: Win,
	}

	return mapToAwayOutcome[outcome]
}

func isExpectedOutcome(outcome, expected Outcome) int {
	if outcome == expected {
		return 1
	}

	return 0
}

func setRowOutcome(outcome Outcome, row TableRow) TableRow {
	isWin := isExpectedOutcome(outcome, Win)
	isDraw := isExpectedOutcome(outcome, Draw)
	isLoss := isExpectedOutcome(outcome, Loss)

	row.MatchesWon += isWin
	row.MatchesDrawn += isDraw
	row.MatchesLost += isLoss

	if isWin == 1 {
		row.Points += int(Win)
	} else if isDraw == 1 {
		row.Points += int(Draw)
	} else {
		row.Points += int(Loss)
	}

	return row
}
