package scope

import (
	"slices"
)

type ActiveScope struct {
	defaultScope AcademicScope
	scopes       []AcademicScope
	campuses     []Campus
	seasons      []Season
}

func NewActiveScope(campuses []Campus, seasons []Season) ActiveScope {
	defaultScope := AcademicScope{
		Campus: campuses[0],
		Term:   seasons[0].Term,
		Year:   seasons[0].Year,
	}

	var scopes []AcademicScope
	for _, campus := range campuses {
		for _, szn := range seasons {
			scopes = append(scopes, AcademicScope{
				Campus: campus,
				Term:   szn.Term,
				Year:   szn.Year,
			})
		}
	}

	return ActiveScope{
		defaultScope: defaultScope,
		scopes:       scopes,
		campuses:     campuses,
		seasons:      seasons,
	}
}

func (s *ActiveScope) Resolve(campus *Campus, season *Season) (AcademicScope, error) {
	c := s.defaultScope.Campus
	if campus != nil {
		c = *campus
	}

	term, year := s.defaultScope.Term, s.defaultScope.Year
	if season != nil {
		term = season.Term
		year = season.Year
	}

	scp := AcademicScope{Campus: c, Term: term, Year: year}
	if !slices.Contains(s.scopes, scp) {
		return AcademicScope{}, ErrScopeInvalid
	}
	return scp, nil
}

func (s *ActiveScope) Validate(scp AcademicScope) error {
	if !slices.Contains(s.scopes, scp) {
		return ErrScopeInvalid
	}
	return nil
}

func (s *ActiveScope) Scopes() []AcademicScope {
	return s.scopes
}

func (s *ActiveScope) Campuses() []Campus {
	return s.campuses
}

func (s *ActiveScope) Seasons() []Season {
	return s.seasons
}
