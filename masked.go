package br

// This file holds the checks of the formato_ rules: whether a value is
// written in the shape of a document, with its punctuation in place. They
// look at the spelling and at nothing else -- not the check digits, not the
// range -- because that is the question those rules ask, and a form asks it
// when it wants the punctuated spelling and no other.

// punctuated reports whether s is exactly the mask filled in: the
// punctuation of the mask in place, a digit where it marks '9', and a digit or
// a capital letter where it marks 'A'. The bare spelling is not punctuated.
func punctuated(s, mask string) bool {
	if len(s) != len(mask) {
		return false
	}
	for i := 0; i < len(mask); i++ {
		c := s[i]
		switch mask[i] {
		case '9':
			if c < '0' || c > '9' {
				return false
			}
		case 'A':
			if (c < '0' || c > '9') && (c < 'A' || c > 'Z') {
				return false
			}
		default:
			if c != mask[i] {
				return false
			}
		}
	}
	return true
}

// maskedCPF reports whether s is written 000.000.000-00.
func maskedCPF(s string) bool { return punctuated(s, cpfMask) }

// maskedCNPJ reports whether s is written 00.000.000/0000-00, with digits or
// capital letters in the first twelve positions.
func maskedCNPJ(s string) bool { return punctuated(s, cnpjMask) }

// maskedCPFOrCNPJ reports whether s is written as either of the two.
func maskedCPFOrCNPJ(s string) bool { return maskedCPF(s) || maskedCNPJ(s) }

// maskedCEP reports whether s is written 00000-000 or 00.000-000, the second
// being how a CEP is often typed.
func maskedCEP(s string) bool { return punctuated(s, cepMask) || punctuated(s, "99.999-999") }

// maskedNIS reports whether s is written 000.00000.00-0.
func maskedNIS(s string) bool { return punctuated(s, nisMask) }

// maskedCivilCertificate reports whether s is the number of a certificate in
// its nine groups, in any of the separators ParseCivilCertificate reads.
func maskedCivilCertificate(s string) bool {
	_, ok := readCertificate(s)
	return ok && len(s) != 32
}

// platePatternShape reports whether s is three letters, an optional hyphen, a
// digit, a digit or a letter, and two digits: the plate of 1990 to 2018 and the
// Mercosul plate alike.
func platePatternShape(s string) bool {
	if len(s) == 8 && s[3] == '-' {
		s = s[:3] + s[4:]
	}
	if len(s) != 7 {
		return false
	}
	for i := 0; i < 7; i++ {
		c := s[i]
		letter := c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z'
		digit := c >= '0' && c <= '9'
		switch {
		case i < 3 && !letter, i == 3 && !digit, i == 4 && !letter && !digit, i > 4 && !digit:
			return false
		}
	}
	return true
}
