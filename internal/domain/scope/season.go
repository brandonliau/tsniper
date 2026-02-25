package scope

type Season struct {
	Term Term
	Year string
}

func (s Season) DisplayName() string {
	return s.Term.DisplayName() + " " + s.Year
}
