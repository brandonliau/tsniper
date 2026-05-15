package view

import (
	"tsniper/internal/domain/scope"
)

type SeasonOption struct {
	Name  string
	Value string
}

type CampusOption struct {
	Name string
	Code string
}

func FromSeason(s scope.Season) SeasonOption {
	return SeasonOption{
		Name:  s.DisplayName(),
		Value: s.Code(),
	}
}

func FromSeasons(seasons []scope.Season) []SeasonOption {
	out := make([]SeasonOption, 0, len(seasons))
	for _, s := range seasons {
		out = append(out, FromSeason(s))
	}
	return out
}

func FromCampus(c scope.Campus) CampusOption {
	return CampusOption{
		Name: c.DisplayName(),
		Code: c.Code(),
	}
}

func FromCampuses(campuses []scope.Campus) []CampusOption {
	out := make([]CampusOption, 0, len(campuses))
	for _, c := range campuses {
		out = append(out, FromCampus(c))
	}
	return out
}
