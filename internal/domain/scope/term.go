package scope

type Term string

const (
	TermWinter Term = "0"
	TermSpring Term = "1"
	TermSummer Term = "7"
	TermFall   Term = "9"
)

func ParseTerm(s string) (Term, error) {
	switch Term(s) {
	case TermSpring, TermSummer, TermFall, TermWinter:
		return Term(s), nil
	default:
		return "", ErrTermInvalid
	}
}

func (t Term) DisplayName() string {
	switch t {
	case TermWinter:
		return "winter"
	case TermSpring:
		return "spring"
	case TermSummer:
		return "summer"
	case TermFall:
		return "fall"
	default:
		return string(t)
	}
}
