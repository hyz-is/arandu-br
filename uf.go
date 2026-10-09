package br

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/arandu-io/hesape/enum"
)

// UF is one of the twenty-seven federative units of Brazil -- the twenty-six
// states and the Federal District -- stored as its two-letter code.
//
// It is a closed set in the shape a generated enum has: constants that are the
// only valid values, Valid to ask, Label for what a form shows, and a
// Value/Scan pair that refuses a code outside the set on the way to and from a
// column. So hesape/enum reads it like any other: enum.Options(br.UFValues()...)
// fills a select, and validation.EnumOf(br.UFValues()...) is the rule.
type UF string

// The units, in the alphabetical order of their names, which is the order a
// list of them is read in. Stored exactly as written: the code is what the
// Correios, the Receita Federal and IBGE all print.
const (
	UFAC UF = "AC"
	UFAL UF = "AL"
	UFAP UF = "AP"
	UFAM UF = "AM"
	UFBA UF = "BA"
	UFCE UF = "CE"
	UFDF UF = "DF"
	UFES UF = "ES"
	UFGO UF = "GO"
	UFMA UF = "MA"
	UFMT UF = "MT"
	UFMS UF = "MS"
	UFMG UF = "MG"
	UFPA UF = "PA"
	UFPB UF = "PB"
	UFPR UF = "PR"
	UFPE UF = "PE"
	UFPI UF = "PI"
	UFRJ UF = "RJ"
	UFRN UF = "RN"
	UFRS UF = "RS"
	UFRO UF = "RO"
	UFRR UF = "RR"
	UFSC UF = "SC"
	UFSP UF = "SP"
	UFSE UF = "SE"
	UFTO UF = "TO"
)

// ErrInvalidUF is what ParseUF, Value and Scan return for a code outside the
// set.
var ErrInvalidUF = errors.New("br: not a UF")

// UFValues lists the twenty-seven units in declaration order.
//
// It is what fills a select and what a rule is built from, so the set is
// written once, here, and every list of units derives from it.
func UFValues() []UF {
	return []UF{
		UFAC, UFAL, UFAP, UFAM, UFBA, UFCE, UFDF, UFES, UFGO,
		UFMA, UFMT, UFMS, UFMG, UFPA, UFPB, UFPR, UFPE, UFPI,
		UFRJ, UFRN, UFRS, UFRO, UFRR, UFSC, UFSP, UFSE, UFTO,
	}
}

// Valid reports whether v is one of the units.
func (v UF) Valid() bool { return v.Label() != "" }

// String is the code, so fmt and a view print it unchanged.
func (v UF) String() string { return string(v) }

// Label is the name of the unit, which is what a form shows: São Paulo for SP.
// It is empty for a code outside the set.
//
// It is separate from String on purpose. The code is a contract with every
// column and every message that holds it; the name is a contract with nobody,
// and changing how one is spelled must not be able to change a row.
func (v UF) Label() string {
	switch v {
	case UFAC:
		return "Acre"
	case UFAL:
		return "Alagoas"
	case UFAP:
		return "Amapá"
	case UFAM:
		return "Amazonas"
	case UFBA:
		return "Bahia"
	case UFCE:
		return "Ceará"
	case UFDF:
		return "Distrito Federal"
	case UFES:
		return "Espírito Santo"
	case UFGO:
		return "Goiás"
	case UFMA:
		return "Maranhão"
	case UFMT:
		return "Mato Grosso"
	case UFMS:
		return "Mato Grosso do Sul"
	case UFMG:
		return "Minas Gerais"
	case UFPA:
		return "Pará"
	case UFPB:
		return "Paraíba"
	case UFPR:
		return "Paraná"
	case UFPE:
		return "Pernambuco"
	case UFPI:
		return "Piauí"
	case UFRJ:
		return "Rio de Janeiro"
	case UFRN:
		return "Rio Grande do Norte"
	case UFRS:
		return "Rio Grande do Sul"
	case UFRO:
		return "Rondônia"
	case UFRR:
		return "Roraima"
	case UFSC:
		return "Santa Catarina"
	case UFSP:
		return "São Paulo"
	case UFSE:
		return "Sergipe"
	case UFTO:
		return "Tocantins"
	}
	return ""
}

// ParseUF turns request input into a UF, or wraps ErrInvalidUF.
//
// It reads the code exactly as the set spells it, in capitals: the value a
// select sends back is the one Options gave it, and a value that differs from
// every option did not come from the form.
func ParseUF(s string) (UF, error) {
	if v := UF(s); v.Valid() {
		return v, nil
	}
	return "", fmt.Errorf("%w: %q", ErrInvalidUF, s)
}

// Compile-time proof that the type is what hesape/enum and a repository read
// it as.
var (
	_ enum.Enum     = UF("")
	_ driver.Valuer = UF("")
	_ sql.Scanner   = (*UF)(nil)
)

// Value implements driver.Valuer. It refuses a code outside the set, which is
// what keeps a zero-valued struct from putting an empty string in the column.
func (v UF) Value() (driver.Value, error) {
	if !v.Valid() {
		return nil, fmt.Errorf("%w: refusing to write %q", ErrInvalidUF, string(v))
	}
	return string(v), nil
}

// Scan implements sql.Scanner. A code the set does not hold is an error at the
// row that has it, rather than a value that behaves like a unit three layers
// up.
func (v *UF) Scan(src any) error {
	var s string
	switch raw := src.(type) {
	case string:
		s = raw
	case []byte:
		s = string(raw)
	default:
		return fmt.Errorf("%w: cannot read %T from the database", ErrInvalidUF, src)
	}
	parsed, err := ParseUF(s)
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
