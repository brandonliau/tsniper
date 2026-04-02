package scope

type Campus string

const (
	CampusNewBrunswick Campus = "NB"
	CampusNewark       Campus = "NK"
	CampusCamden       Campus = "CM"
)

func ParseCampus(s string) (Campus, error) {
	switch Campus(s) {
	case CampusNewBrunswick, CampusNewark, CampusCamden:
		return Campus(s), nil
	default:
		return "", ErrCampusInvalid
	}
}

func (c Campus) Code() string {
	return string(c)
}

func (c Campus) DisplayName() string {
	switch c {
	case CampusNewBrunswick:
		return "new brunswick"
	case CampusNewark:
		return "newark"
	case CampusCamden:
		return "camden"
	default:
		return string(c)
	}
}
