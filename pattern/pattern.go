package pattern

import (
	"bytes"
	"errors"
)

var (
	ErrCaptureNotAllowed        = errors.New("named captures are not allowed")
	ErrUnnamedCaptureNotAllowed = errors.New("unnamed captures are not allowed")
	ErrInvalidExpr              = errors.New("invalid expression")
)

type Matcher struct {
	e     expr
	names []string
}

func New(in string) (*Matcher, error) {
	e, err := parseExpr(in)
	if err != nil {
		return nil, err
	}
	if err := e.validate(); err != nil {
		return nil, err
	}
	return &Matcher{
		e:     e,
		names: e.captures(),
	}, nil
}

func ParseLineFilter(in []byte) (*Matcher, error) {
	if len(in) == 0 {
		return new(Matcher), nil
	}
	e, err := parseExprBytes(in)
	if err != nil {
		return nil, err
	}
	if err = e.validateNoConsecutiveCaptures(); err != nil {
		return nil, err
	}
	return &Matcher{e: e}, nil
}

func ParseLiterals(in string) ([][]byte, error) {
	e, err := parseExpr(in)
	if err != nil {
		return nil, err
	}
	lit := make([][]byte, 0, len(e))
	for _, n := range e {
		if l, ok := n.(literals); ok {
			lit = append(lit, l)
		}
	}
	return lit, nil
}

func ParseNodes(in string) ([][]byte, []string, error) {
	if len(in) == 0 {
		return [][]byte{}, []string{}, nil
	}
	e, err := parseExpr(in)
	if err != nil {
		return nil, nil, err
	}
	if err = e.validateNoUnnamedCaptures(); err != nil {
		return nil, nil, err
	}
	lit := make([][]byte, len(e))
	names := make([]string, len(e))
	for i, n := range e {
		if l, ok := n.(literals); ok {
			lit[i] = l
		} else if c, ok := n.(capture); ok {
			names[i] = string(c)
		}
	}
	return lit, names, nil
}

// Matches matches the given line with the provided pattern.
// Matches invalidates the previous returned captures array.
func (m *Matcher) Matches(in []byte) [][]byte {
	if len(in) == 0 {
		return nil
	}
	if len(m.e) == 0 {
		return nil
	}
	var result [][]byte
	expr := m.e
	if ls, ok := expr[0].(literals); ok {
		i := bytes.Index(in, ls)
		if i != 0 {
			return nil
		}
		in = in[len(ls):]
		expr = expr[1:]
	}
	if len(expr) == 0 {
		return nil
	}
	// from now we have capture - literals - capture ... (literals)?
	for i := 0; i < len(expr); i += 2 {
		capt := expr[i].(capture)
		if i+1 >= len(expr) { // we're ending on a capture.
			if !capt.isUnnamed() {
				result = append(result, in)
			}
			return result
		}
		ls := expr[i+1].(literals)
		i := bytes.Index(in, ls)
		if i == -1 {
			// if a capture is missed we return up to the end as the capture.
			if !capt.isUnnamed() {
				result = append(result, in)
			}
			return result
		}
		if capt.isUnnamed() {
			in = in[len(ls)+i:]
			continue
		}
		result = append(result, in[:i])
		in = in[len(ls)+i:]
	}

	return result
}

func (m *Matcher) Names() []string {
	return m.names
}

func (m *Matcher) Test(in []byte) bool {
	if len(in) == 0 || len(m.e) == 0 {
		// An empty line can only match an empty pattern.
		return len(in) == 0 && len(m.e) == 0
	}
	var off int
	for i := range m.e {
		lit, ok := m.e[i].(literals)
		if !ok {
			continue
		}
		j := bytes.Index(in[off:], lit)
		if j == -1 {
			return false
		}
		if i != 0 && j == 0 {
			// This means we either have repetitive literals, or an empty
			// capture. Either way, the line does not match the pattern.
			return false
		}
		off += j + len(lit)
	}
	// If we end up on a literal, we only consider the test successful if
	// the remaining input is empty. Otherwise, if we end up on a capture,
	// the remainder (the captured text) must not be empty.
	//
	// For example, "foo bar baz" does not match "<_> bar", but it matches
	// "<_> baz" and "foo <_>".
	//
	// Empty captures are not allowed as well: " bar " does not match
	// "<_> bar <_>", but matches "<_>bar<_>".
	_, reqRem := m.e[len(m.e)-1].(capture)
	hasRem := off != len(in)
	return reqRem == hasRem
}
