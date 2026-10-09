package br

// A mask is how a document is written when it is punctuated: '9' marks a
// position that holds a digit, 'A' one that holds a digit or a letter, and
// every other character stands for itself.
//
// One mask answers both spellings a form receives. The punctuated one is the
// mask filled in; the bare one is the marked positions alone. The two lengths
// differ for every mask here, so the length says which one arrived.
const (
	cpfMask  = "999.999.999-99"
	cnpjMask = "AA.AAA.AAA/AAAA-99"
	cepMask  = "99999-999"
)

// slot reports whether a mask character marks a position rather than standing
// for itself.
func slot(m byte) bool { return m == '9' || m == 'A' }

// slots is the number of positions a mask marks, which is the length of the
// bare spelling.
func slots(mask string) int {
	n := 0
	for i := 0; i < len(mask); i++ {
		if slot(mask[i]) {
			n++
		}
	}
	return n
}

// fill reads one character into a position of the given kind. A letter is
// upper-cased, because the letters of a document are capitals and a lower-case
// one typed into a form is the same character.
func fill(c, kind byte) (byte, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c, true
	case kind == 'A' && c >= 'A' && c <= 'Z':
		return c, true
	case kind == 'A' && c >= 'a' && c <= 'z':
		return c - 'a' + 'A', true
	}
	return 0, false
}

// read returns the characters s holds in the positions the mask marks, when s
// is either the mask filled in or those characters alone.
//
// Nothing else is accepted: not another punctuation, not a separator in the
// wrong place, and not whitespace around the value. A form that trimmed or
// re-punctuated what it was given would accept a value the person did not
// type, and the error they would see is about the value they did.
func read(s, mask string) (string, bool) {
	out := make([]byte, 0, slots(mask))
	switch len(s) {
	case len(mask):
		for i := 0; i < len(mask); i++ {
			if !slot(mask[i]) {
				if s[i] != mask[i] {
					return "", false
				}
				continue
			}
			c, ok := fill(s[i], mask[i])
			if !ok {
				return "", false
			}
			out = append(out, c)
		}
	case cap(out):
		j := 0
		for i := 0; i < len(mask); i++ {
			if !slot(mask[i]) {
				continue
			}
			c, ok := fill(s[j], mask[i])
			if !ok {
				return "", false
			}
			out = append(out, c)
			j++
		}
	default:
		return "", false
	}
	return string(out), true
}

// punctuate writes the characters of a bare value into the mask.
func punctuate(bare, mask string) string {
	out := make([]byte, 0, len(mask))
	j := 0
	for i := 0; i < len(mask); i++ {
		if slot(mask[i]) {
			out = append(out, bare[j])
			j++
			continue
		}
		out = append(out, mask[i])
	}
	return string(out)
}

// format punctuates s when it has the shape of the mask, punctuated or bare,
// and returns it unchanged when it has not.
func format(s, mask string) string {
	bare, ok := read(s, mask)
	if !ok {
		return s
	}
	return punctuate(bare, mask)
}

// repeated reports whether every character of s is the same one.
func repeated(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}
