package utils

import (
	"charm.land/bubbletea/v2"
	"database/sql"
	"strings"
)

type NavMsg bool

type DeleteWorkMsg int

type NewWorkMsg []string

type ViewMsg int

type ConfirmationMsg struct {
	Function func() tea.Cmd
	Msg      string
}

type WorkDetails []string

const (
	HomePage int = iota
	AddPage
	WorkPage
	ConfirmPage
	PageCount
)

const ( // Order that the forms are in for workForm
	TitleForm int = iota
	YearForm
	TagsForm
	MediumForm
	StatusForm
	EnterForm
)

const ( // Table Column
	Title int = iota
	Tags
	Medium
	Status
	Released
	Id
)

var DirectoryPath string

var statusName = map[int]string{}

var statusInt = map[string]int{}

var mediumName = map[int]string{}

var mediumInt = map[string]int{}

func SetupStatuses(db *sql.DB) {
	row, err := db.Query(`SELECT id, status_name FROM status_table;`)
	CheckError("Failed to query status_table: ", err)
	for row.Next() {
		var (
			name string
			id   int
		)

		err = row.Scan(&id, &name)
		CheckError("Failed to scan row from status_table: ", err)
		statusName[id-1] = name
		statusInt[name] = id
	}
	DebugLog("Statuses: ", statusName)
}

func SetupMediums(db *sql.DB) {
	row, err := db.Query(`SELECT id, type_name FROM media_type_table;`)
	CheckError("Failed to query media_type_table: ", err)
	for row.Next() {
		var (
			name string
			id   int
		)

		err = row.Scan(&id, &name)
		CheckError("Failed to scan row from media_type_table: ", err)
		mediumName[id] = name
		mediumInt[name] = id
	}
	DebugLog("Mediums: ", mediumName)
}

func Status_stoi(status string) int {
	return statusInt[status]
}

func Status_itos(status int) string {
	return statusName[status]
}

func Medium_stoi(medium string) int {
	return mediumInt[medium]
}

func Medium_itos(medium int) string {
	return mediumName[medium]
}

func ConvertMedium(mediums []int) string {
	builder := strings.Builder{}

	for i, val := range mediums {
		medium := Medium_itos(val)
		builder.WriteString(medium)
		if i < len(mediums)-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

func GetTagsString(tags []string) string {
	builder := strings.Builder{}

	for i, tag := range tags {
		builder.WriteString(tag)
		if i < len(tags)-1 {
			builder.WriteString(", ")
		}
	}
	return builder.String()
}

// medium.com/@rohitsangamnerkar1999/understanding-string-truncation-in-go-handling-unicode-safely-7740d24fa6a6
func TruncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
