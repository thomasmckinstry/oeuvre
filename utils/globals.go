package utils

import "strings"

type NavMsg bool

type NewWorkMsg []string

// TODO: Define these with reading from the db
const (
	Pending int = iota
	Started
	Hiatus
	Completed
	Dropped
)

const ( // Order that the forms are in
	TitleForm int = iota
	YearForm
	TagsForm
	MediumForm
	StatusForm
	EnterForm
)

const ( // Sort Option
	Title int = iota
	Tags
	Medium
	Status
	Released
)

var statusName = map[int]string{
	Pending:   "Pending",
	Started:   "Started",
	Hiatus:    "Hiatus",
	Completed: "Completed",
	Dropped:   "Dropped",
}

var statusInt = map[string]int{
	"Pending":   Pending,
	"Started":   Started,
	"Hiatus":    Hiatus,
	"Completed": Completed,
	"Dropped":   Dropped,
}

const (
	Anime int = iota
	Manga
	Movie
	Book
	Comic
	Show
	Animated
	Live_Action
)

var mediumName = map[int]string{
	Anime:       "Anime",
	Manga:       "Manga",
	Movie:       "Movie",
	Book:        "Book",
	Comic:       "Comic",
	Show:        "Show",
	Animated:    "Animated",
	Live_Action: "Live Action",
}

var mediumInt = map[string]int{
	"Anime":       Anime,
	"Manga":       Manga,
	"Movie":       Movie,
	"Book":        Book,
	"Comic":       Comic,
	"Show":        Show,
	"Animated":    Animated,
	"Live Action": Live_Action,
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
