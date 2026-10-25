package br

// NormalizeCPF returns the digits of s and drops every other character: the
// dots, the hyphen, the spaces a person or a spreadsheet put around them.
//
// It is for text that already left the form -- a column, a file, a pasted
// list -- and has to be compared with the stored spelling. It does not
// validate, does not check the length and does not restore a zero a number
// lost, because a CPF is text: keep it as text from the form to the column.
// Taking a value in is ParseCPF, which refuses what NormalizeCPF would have
// forgiven.
func NormalizeCPF(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			out = append(out, s[i])
		}
	}
	return string(out)
}

// NormalizeCNPJ returns the digits and the letters of s, the letters in
// capitals, and drops every other character: the dots, the slash, the hyphen
// and the spaces. A CNPJ may carry letters in its first twelve positions, so
// keeping only the digits would turn one into another number.
//
// It does not validate and does not check the length; taking a value in is
// ParseCNPJ. The result is text, as the document is: a CNPJ read as a number
// loses its leading zeros and cannot hold a letter.
func NormalizeCNPJ(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c >= '0' && c <= '9', c >= 'A' && c <= 'Z':
			out = append(out, c)
		case c >= 'a' && c <= 'z':
			out = append(out, c-'a'+'A')
		}
	}
	return string(out)
}
