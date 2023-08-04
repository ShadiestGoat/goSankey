package parser

type MultiError struct {
	Errors []error
}

func (m MultiError) Error() string {
	if len(m.Errors) == 0 {
		return ""
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}

	str := "\n" + m.PureErrors()

	return str[:len(str)-1]
}

func (m MultiError) PureErrors() string {
	str := ""

	for _, err := range m.Errors {
		if err == nil {
			continue
		}

		if errP, ok := err.(MultiError); ok {
			str += errP.PureErrors() + "\n"
		} else {
			str += "- " + err.Error() + "\n"
		}
	}

	if str == "" {
		return ""
	}

	return str[:len(str)-1]
}

func (m MultiError) Optional() error {
	for _, v := range m.Errors {
		if v != nil {
			return m
		}
	}

	return nil
}