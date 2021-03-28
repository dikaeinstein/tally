package tournament

import (
	"strings"
	"testing"
)

func TestParseInput(t *testing.T) {
	input := buildInput(t)

	matches, err := ParseInput(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	expectedLen := 6
	if len(matches) != expectedLen {
		t.Errorf("ParseInput failed: want len(matches) => %d; got %d",
			expectedLen, len(matches))
	}

	expectedFirstMatch := Match{
		Home:    "Allegoric Alaskans",
		Away:    "Blithering Badgers",
		Outcome: Win,
	}
	if matches[0] != expectedFirstMatch {
		t.Errorf("ParseInput failed: want %v; got %v", expectedFirstMatch, matches[0])
	}
}

func TestTally(t *testing.T) {
	input := buildInput(t)
	matches, err := ParseInput(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	expectedTable := buildTable(t, Desc)

	table := Tally(matches)
	if len(table.Rows) != len(expectedTable.Rows) {
		t.Errorf("Tally failed: len of TableRows don't match: want %d; got %d",
			len(expectedTable.Rows), len(table.Rows))
	}

	teamNames := mapToTeamNames(t, table.Rows)
	expectedTeamNames := mapToTeamNames(t, expectedTable.Rows)

	if teamNames != expectedTeamNames {
		t.Errorf("Tally failed: table not sorted correctly: want %#v, got %#v",
			expectedTeamNames, teamNames)
	}
}

func TestSortTable(t *testing.T) {
	testCases := []struct {
		desc   string
		input  Table
		order  string
		output Table
	}{
		{
			desc:   "Ascending",
			input:  buildTable(t, Desc),
			order:  Asc,
			output: buildTable(t, Asc),
		},
		{
			desc:   "Descending",
			input:  buildTable(t, Asc),
			order:  Desc,
			output: buildTable(t, Desc),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			table := tC.input.Sort(tC.order)

			teamNames := mapToTeamNames(t, table.Rows)
			expectedTeamNames := mapToTeamNames(t, tC.output.Rows)

			if teamNames != expectedTeamNames {
				t.Errorf("Tally failed: table not sorted correctly: want %#v, got %#v",
					expectedTeamNames, teamNames)
			}
		})
	}
}

func buildInput(t *testing.T) string {
	t.Helper()

	return `
Allegoric Alaskans;Blithering Badgers;win
Devastating Donkeys;Courageous Californians;draw
Devastating Donkeys;Allegoric Alaskans;win
Courageous Californians;Blithering Badgers;loss
Blithering Badgers;Devastating Donkeys;loss
Allegoric Alaskans;Courageous Californians;win
`
}

func buildTable(t *testing.T, order string) Table {
	expectedTableRows := []TableRow{
		{
			Team:         "Devastating Donkeys",
			MatchPlayed:  3,
			MatchesWon:   2,
			MatchesDrawn: 1,
			MatchesLost:  0,
			Points:       7,
		},
		{
			Team:         "Allegoric Alaskans",
			MatchPlayed:  3,
			MatchesWon:   2,
			MatchesDrawn: 0,
			MatchesLost:  1,
			Points:       6,
		},
		{
			Team:         "Blithering Badgers",
			MatchPlayed:  3,
			MatchesWon:   1,
			MatchesDrawn: 0,
			MatchesLost:  0,
			Points:       3,
		},
		{
			Team:         "Courageous Californians",
			MatchPlayed:  3,
			MatchesWon:   0,
			MatchesDrawn: 1,
			MatchesLost:  2,
			Points:       1,
		},
	}

	if order == Asc {
		// reverse expectedTableRows because it is sorted in
		// desc order by default
		for i := len(expectedTableRows)/2 - 1; i >= 0; i-- {
			opp := len(expectedTableRows) - 1 - i
			expectedTableRows[i], expectedTableRows[opp] =
				expectedTableRows[opp], expectedTableRows[i]
		}
	}

	return Table{Rows: expectedTableRows}
}

func mapToTeamNames(t *testing.T, rows []TableRow) [4]string {
	t.Helper()

	var teams [4]string
	for i, r := range rows {
		teams[i] = r.Team
	}

	return teams
}
