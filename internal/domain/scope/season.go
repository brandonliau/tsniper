package scope

import (
	"strings"
)

type Season struct {
	Term Term
	Year string
}

func ParseSeason(s string) (Season, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return Season{}, ErrSeasonInvalid
	}
	term, err := ParseTerm(parts[0])
	if err != nil {
		return Season{}, ErrSeasonInvalid
	}
	return Season{Term: term, Year: parts[1]}, nil
}

func (s Season) Code() string {
	return string(s.Term) + ":" + s.Year
}

func (s Season) DisplayName() string {
	return s.Term.DisplayName() + " " + s.Year
}
