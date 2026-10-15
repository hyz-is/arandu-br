package br

import "strings"

// A Redaction is how a page writes a document with part of it hidden, spelled
// the way the document is spelled: '#' shows the character in that position,
// and every other character is written as it stands.
//
//	"***.###.###-**"  the fourth to the ninth digits of a CPF, punctuated
//	"***######**"     the same digits, bare
//	"•••.###.###-••"  the same digits, hidden behind bullets
//
// A Redaction has the length, in characters, of either spelling of the
// document: the punctuated one (14 for a CPF, 18 for a CNPJ) or the bare one
// (11 and 14). Where the punctuated spelling holds a separator, the redaction
// writes its own character there, so it chooses the punctuation as well as
// the hiding. The length alone says which spelling a redaction is, so one
// written for a document of another kind is read as whichever of the two it
// measures, or refused.
//
// The named redactions are the ones Brazilian sources publish and pages
// print. A redaction written as a literal is just as good: any string of one of
// the two lengths is one.
type Redaction string

// The redactions a page usually needs. Each is a constant of the type, so a
// call reads br.RedactCPF(cpf, br.KeepCPFMiddle).
const (
	// KeepCPFMiddle shows the fourth to the ninth digits of a CPF, punctuated:
	// ***.456.789-**. It is how the Portal da Transparência publishes the CPF
	// of a public agent.
	KeepCPFMiddle Redaction = "***.###.###-**"
	// KeepCPFMiddleBare shows the same six digits with no punctuation:
	// ***456789**. It is how the Receita Federal's open CNPJ data publishes the
	// CPF of a partner.
	KeepCPFMiddleBare Redaction = "***######**"
	// KeepCPFCheckDigits shows only the two check digits of a CPF, which
	// identify nobody: ***.***.***-00.
	KeepCPFCheckDigits Redaction = "***.***.***-##"
	// KeepCPFLastFive shows the last five digits of a CPF, as a profile page
	// reminds its owner which CPF it holds: ***.***.789-00.
	KeepCPFLastFive Redaction = "***.***.###-##"
	// KeepCNPJEnds shows the first two and the last two characters of a CNPJ:
	// 12.***.***/****-35.
	KeepCNPJEnds Redaction = "##.***.***/****-##"
)

// shown is the character of a Redaction that shows a position.
const shown = '#'

// RedactCPF writes a CPF with the positions r hides replaced, as r spells it.
//
// It takes the eleven digits bare or punctuated, as ParseCPF reads them, and
// does not check the check digits: a page that hides a stored value hides it
// whether or not it was ever valid. For anything else -- a value of another
// shape, or a redaction of neither length -- it returns "", and never the value
// it was given: a redaction that fell back to its input would print the
// document it exists to hide.
func RedactCPF(s string, r Redaction) string { return redact(s, cpfMask, r) }

// RedactCNPJ writes a CNPJ with the positions r hides replaced, as r spells it.
//
// It takes the fourteen characters bare or punctuated, as ParseCNPJ reads
// them, a lower-case letter shown as its capital, and does not check the check
// digits. For anything else, or a redaction of neither length, it returns "".
func RedactCNPJ(s string, r Redaction) string { return redact(s, cnpjMask, r) }

// redact reads s against the mask and writes it through the redaction.
func redact(s, mask string, r Redaction) string {
	bare, ok := read(s, mask)
	if !ok {
		return ""
	}
	pattern := []rune(string(r))
	punctuated := len(pattern) == len(mask)
	if !punctuated && len(pattern) != len(bare) {
		return ""
	}

	var b strings.Builder
	b.Grow(len(r))
	j := 0
	for i, p := range pattern {
		if punctuated && !slot(mask[i]) {
			if p == shown {
				b.WriteByte(mask[i])
			} else {
				b.WriteRune(p)
			}
			continue
		}
		if p == shown {
			b.WriteByte(bare[j])
		} else {
			b.WriteRune(p)
		}
		j++
	}
	return b.String()
}
