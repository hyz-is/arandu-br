package unit_test

import (
	"testing"

	"github.com/arandu-io/hesape/validation"
	br "github.com/hyz-is/arandu-br"
)

// installed is the factory an application builds at boot, with the rules
// installed on it. The tests that name a rule other than the seven registered
// when the package is imported compile their sets from it, after Install, as
// an application does.
var installed = func() *validation.Factory {
	f := validation.NewFactory(nil)
	br.Install(f)
	return f
}()

// example is one value and what the rule says of it. note, when set, is where
// the answer departs from the PHP package the example comes from.
type example struct {
	in    string
	valid bool
	note  string
}

func ok(in ...string) []example {
	out := make([]example, len(in))
	for i, v := range in {
		out[i] = example{in: v, valid: true}
	}
	return out
}

func no(in ...string) []example {
	out := make([]example, len(in))
	for i, v := range in {
		out[i] = example{in: v}
	}
	return out
}

func join(parts ...[]example) []example {
	var out []example
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// refused is a value the PHP package accepts and this one does not, with why.
func refused(note string, in ...string) []example {
	out := no(in...)
	for i := range out {
		out[i].note = note
	}
	return out
}

// accepted is a value the PHP package refuses and this one accepts, with why.
func accepted(note string, in ...string) []example {
	out := ok(in...)
	for i := range out {
		out[i].note = note
	}
	return out
}

const (
	eightDigits = "the PHP package takes the eight-digit number a mobile line had before 2016; since then every mobile number has nine digits, and the eight-digit one is a landline"
	otherCode   = "the PHP package takes any country code; the area code and the number are checked against Brazil's plan, so only +55"
	anyBlank    = "the PHP package forgives any character that is not a digit; this one takes the value bare or in its one punctuation"
)

// examples are the valid and invalid values the two PHP packages test their
// rules with, rewritten as data, and the same rules held to values of their
// own. The key is the rule as a rule string writes it.
var examples = map[string][]example{
	"cpf": join(
		ok("094.050.986-59", "98136622809", "981.366.228-09", "00000000191"),
		no("99800-1926", "08136622809", "000.366.228-09", "11111111111", "111.111.111-11", "", "123", "0940509865", "094.050.986-5", "094.050.986-599"),
		refused(anyBlank, "094 050 986 59", "094-050-986.59", " 094.050.986-59"),
	),
	"cnpj": join(
		ok("T6.JSP.XPS/0001-11", "T6.JSP.XPS/J84K-69", "D2.M97.AA0/0001-63", "W9.7VY.JMY/0001-81", "E5SGVHX9000190", "E5SGVHX960LR53",
			"12ABC34501DE35", "12.ABC.345/01DE-35", "12abc34501de35", "12.abc.345/01de-35", "32.332.643/0001-29", "32.332.643/3060-40", "38725833000192",
			"38725833223060", "33.109.809/0001-05", "53.084.587/0001-20", "49851807000127", "56.611.605/0001-73", "00.000.000/0001-91"),
		no("12.345.678/0001-00", "00.000.000/0000-00", "32.332.643/0001-9!", "12.ABC.345/01DE-36", "33.109.809/0001-06", "51.084.587/0001-20",
			"11111111111111", "22222222222222", "11.111.111/1111-11", "00000000000000", "", "12ABC34501DEA5", "12ABC34501DE3", "12ABC34501DE355",
			"12.ABC.345/01DE35", "12ABC34501DE-35"),
		refused(anyBlank, " 12ABC34501DE35", "12 ABC 345 01DE 35"),
	),
	"cpf_ou_cnpj": join(
		ok("981.366.228-09", "56.611.605/0001-73", "49851807000127", "53.084.587/0001-20", "63.657.343/0001-43", "12ABC34501DE35"),
		no("000.366.228-09", "11.611.605/0001-73", "22851807000127", "99800-1926", ""),
	),
	"cpf_cnpj": join(
		ok("981.366.228-09", "56.611.605/0001-73", "53.084.587/0001-20"),
		no("000.366.228-09", "11.611.605/0001-73", "99800-1926"),
	),
	"cpf_or_cnpj": join(ok("981.366.228-09", "49851807000127"), no("000.366.228-09", "22851807000127")),
	"formato_cpf": join(
		ok("981.366.228-09", "094.050.986-59", "000.000.000-00", "111.111.111-11"),
		no("08136622809", "094.050.986-591", "981.366.228-0", "981.366.228 09", "981-366-228.09", "98a.366.228-09", ""),
	),
	"formato_cnpj": join(
		ok("12.ABC.345/01DE-35", "33.109.809/0001-05", "53.084.587/0001-20", "00.000.000/0000-00"),
		no("12.ABC.345/01DE35", "33.109.809/000105", "51.084.587/000120", "00000000000000", "12ABC34501DE35", "12.abc.345/01de-35", "12.ABC.345/01DE-3A", "12.ABC.345/01DE-355", ""),
	),
	"formato_cpf_ou_cnpj": join(
		ok("981.366.228-09", "000.000.000-00", "56.611.605/0001-73", "111.111.111-11", "63.657.343/0001-43"),
		no("0000.366.228-09", "11.6211.605/0001-73", "22851807000127", "11111111111", "63.657.343/0001", ""),
	),
	"formato_cpf_cnpj": join(ok("111.111.111-11", "63.657.343/0001-43"), no("11111111111", "63.657.343/0001")),
	"cep": join(
		ok("01000-000", "70040-010", "70040010", "32400-000", "07550000"),
		no("00000-000", "00999-999", "7004001", "700400100", "70040-0100", "70.040-010", "70040 010", " 70040-010", ""),
	),
	"formato_cep": join(
		ok("32400-000", "32.400-000", "07.550-000", "30.150-150", "00000-000"),
		no("32400.000", "32.400-0000", "0.400-000", "300.40-000", "300.400-000", "32400000", "32400-00", ""),
	),
	"celular": join(
		ok("98899-4444", "988994444", "99999-5555"),
		no("900-1926", "8899-4444", "90000-0000", "(31) 98899-4444", ""),
		refused(eightDigits, "9800-1936"),
	),
	"celular_com_ddd": join(
		ok("(99)98899-4444", "(99) 98899-4444", "(31)99999-5555", "31999995555", "(11) 91234-5678"),
		no("(99)800-1926", "(10)98899-4444", "(11) 81234-5678", "(11) 90000-0000", "98899-4444", "+55(31)99999-5555", ""),
		refused(eightDigits, "(31)9999-5555", "(99)9800-1936"),
	),
	"celular_com_codigo": join(
		ok("+55(99)98899-4444", "+55 (99) 98899-4444", "+55(31)99999-5555", "+5531999995555", "+55 31 999995555"),
		no("+5(99)800-1926", "+55(99)800-1926", "+55(10)98899-4444", "(99)98899-4444", "5599988994444", ""),
		refused(otherCode, "+1(31)99999-5555", "+351(31)99999-5555"),
	),
	"celular_com_codigo_sem_mascara": join(
		ok("+5514999999999", "+55 14 999999999", "+55 14999999999", "+5599988994444"),
		no("+55(14)99999-9999", "+5514 99999-9999", "5514999999999", "+551499999999", "+5510999999999", ""),
		refused(otherCode, "+1 31 999995555"),
	),
	"telefone": join(
		ok("3598-4550", "35984550", "9999-5555"),
		no("99800-1926", "598-4550", "1598-4550", "(11) 3598-4550", ""),
	),
	"telefone_com_ddd": join(
		ok("(99)3500-4444", "(31)9999-5555", "(77)9999-3333", "(11) 3000-0000", "1130000000"),
		no("(99)9-1926", "(10)3500-4444", "(99)99800-1926", "3500-4444", "+55(99)3500-4444", ""),
	),
	"telefone_com_codigo": join(
		ok("+55(99)3500-4444", "+55 (99) 3500-4444", "+559935004444"),
		no("+5(99)9-1926", "+55(99)9-1926", "(99)3500-4444", ""),
		refused(otherCode, "+66(31)9999-5555"),
	),
	"ddd": join(ok("11", "28", "99", "79"), no("10", "20", "00", "1", "011", "(11)", "")),
	"placa": join(
		ok("P15186", "ABC-1234", "abc-1234", "ABC1234", "aBc1234", "abc1234", "BEE4R22", "FUM-0B05", "FUM-5L58", "1234567", "AA1234", "ABC12D3"),
		no("X1234", "a2c-1234", "abc-12ed", "abc 1234", "ãBC1234", "ABCD1234", "ABC12345", "ab1234a", "abc+1234", ""),
	),
	"formato_placa_de_veiculo": join(
		ok("ABC-1234", "abc-1234", "ABC1234", "aBc1234", "abc1234", "BEE4R22", "FUM-0B05", "FUM-5L58"),
		no("a2c-1234", "abc-12ed", "abc 1234", "ãBC1234", "ABCD1234", "ABC12345", "ab1234", "ab123a4", "abc+1234", "XXXBEE4R22", ""),
	),
	"renavam": join(
		ok("32094074362", "23478829239", "34145742746", "41833820181", "63988496222", "197073212", "00197073212"),
		no("11111111111", "32094074212", "62128843267", "1234555582", "00000000000", ""),
		refused(anyBlank, "1970.73212"),
	),
	"cnh": join(
		ok("96784547943", "04463004100"),
		no("96784547999", "11111111111", "00000000000", "9678454794", ""),
	),
	"titulo_eleitor": join(
		ok("3021260"),
		no("1000101230190", "0000000000", ""),
	),
	"nis": join(
		ok("201.73374.34-9", "20173374349", "690.30244.88-6", "042.33768.05-2", "971.78508.77-5"),
		no("201.73374.34-0", "290.30244.88-5", "11111111111", "2017337434", ""),
	),
	"pis": join(
		ok("690.30244.88-6", "042.33768.05-2", "971.78508.77-5", "69030244886"),
		no("290.30244.88-5", "11111111111", "2017337434", ""),
	),
	"formato_nis": join(ok("201.73374.34-9", "276.96730.83-0"), no("276.96730.830", "20173374349", "")),
	"formato_pis": join(ok("276.96730.83-0", "201.73374.34-9"), no("276.96730.830", "27696730830", "")),
	"cns": join(
		ok("116 3876 9194 0009", "892 1623 5477 0008", "272004493990007", "140776230420006", "254446765170004", "912176122180009", "174295560290007",
			"103409299850000", "773398431180002"),
		no("116 5698 9194 0009", "892 2641 5477 0008", "072004493990007", "000000000000000", "111111111111111", "222222222222222",
			"999999999999999", "123456789123456", ""),
	),
	"certidao": join(
		ok("659447 02 55 9015 1 99679 468 2559590-16", "659447.02.55.9015.1.99679.468.2559590-16", "65944702559015199679468255959016"),
		no("659447 02 55 2015 1 27861 468 2559590-32", "00000000000000000000000000000000", ""),
	),
	"formato_certidao": join(
		ok("434546.02.55.2019.1.71037.134.6484858-10", "434546 02 55 2019 1 71037 134 6484858 10", "659447 02 55 9015 1 99679 468 2559590-16"),
		no("434546.02.55.2019.1.71037.134.6484858", "43454602552019171037134648485810", "434546/02/55/2019/1/71037/134/6484858-10", ""),
	),
	"uf": join(
		ok("AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA", "MS", "MT", "MG", "PA", "PB", "PR", "PE", "PI", "RJ", "RN", "RS", "RO", "RR", "SC", "SP", "SE", "TO"),
		no("FALSE", "mg", "sp", "MO", "CB", "", "S", "SPP"),
	),
	"passaporte": join(ok("GB139485", "gb139485"), no("A3948", "GB13948", "")),
	"data": join(
		ok("31/12/2026", "29/02/2028", "01/01/2000"),
		no("31/04/2026", "29/02/2027", "2026-12-31", "31122026", "1/1/2026", "", "00/00/0000"),
	),
	"horario": join(ok("00:00", "23:59", "09:30:15"), no("24:00", "12:60", "9:30", "0930", "12:30 PM", "")),
	"brl":     join(ok("R$ 1.234,56", "1234.56", "47"), no("1,234.56", "R$", "-5", "")),
	"moeda":   join(ok("R$ 1.234,56", "1234,56"), no("1,234.56", "")),
}

// stateExamples are the inscrição estadual numbers geekcom/validator-docs
// tests with, under the rule string that names the unit.
var stateExamples = map[string][]example{
	"inscricao_estadual:AC": join(ok("0184765932153"), no("0184765932154")),
	"inscricao_estadual:AL": join(ok("248308335"), no("248308336")),
	"inscricao_estadual:AP": join(ok("039895661"), no("039895662")),
	"inscricao_estadual:AM": join(ok("292278012"), no("292278013")),
	"inscricao_estadual:BA": join(ok("96338555", "82679341", "826793-41"), no("82679342", "96338556")),
	"inscricao_estadual:CE": join(ok("980874165"), no("980874166")),
	"inscricao_estadual:DF": join(ok("0740769300107"), no("0740769300108")),
	"inscricao_estadual:ES": join(ok("342048090"), no("342048091")),
	"inscricao_estadual:GO": join(ok("107900459"), no("107900450")),
	"inscricao_estadual:MA": join(ok("122917510"), no("122917511")),
	"inscricao_estadual:MT": join(ok("55160385510"), no("55160385511")),
	"inscricao_estadual:MS": join(ok("283814659"), no("283814650")),
	"inscricao_estadual:MG": join(ok("0526366324132"), no("0526366324133")),
	"inscricao_estadual:PA": join(ok("158330005"), no("158330006")),
	"inscricao_estadual:PB": join(ok("439622301"), no("439622302")),
	"inscricao_estadual:PR": join(ok("5714953410"), no("5714953411")),
	"inscricao_estadual:PE": join(ok("920145698"), no("920145699")),
	"inscricao_estadual:PI": join(ok("758505183"), no("758505184")),
	"inscricao_estadual:RJ": join(ok("61804315"), no("61804316")),
	"inscricao_estadual:RN": join(ok("208627715"), no("208627716")),
	"inscricao_estadual:RS": join(ok("3821957672"), no("3821957673")),
	"inscricao_estadual:RO": join(ok("52059985926850"), no("52059985926851")),
	"inscricao_estadual:RR": join(ok("249977060"), no("249977061")),
	"inscricao_estadual:SC": join(ok("696192667"), no("696192668")),
	"inscricao_estadual:SP": join(ok("653172024009"), no("653172024000", "82679341")),
	"inscricao_estadual:SE": join(ok("646597361"), no("646597362")),
	"inscricao_estadual:TO": join(ok("75036274184"), no("75036274185")),
	// The unit is the argument, in either case, and it has to be one.
	"inscricao_estadual:ba":    ok("82679341"),
	"inscricao_estadual:ZA":    no("82679341"),
	"inscricao_estadual:y":     no("82679341"),
	"inscricao_estadual:12":    no("82679341"),
	"inscricao_estadual":       no("82679341"),
	"inscricao_estadual:BA,SP": no("82679341"),
}

func passes(t *testing.T, rule, value string) bool {
	t.Helper()
	set, err := validation.Compile(validation.Rules{"field": rule})
	if err != nil {
		t.Fatalf("compiling %q: %v", rule, err)
	}
	return installed.Make(validation.Data{"field": value}, set).Passes()
}

func TestEveryRuleAgreesWithTheExamplesOfThePHPPackages(t *testing.T) {
	t.Parallel()
	for _, table := range []map[string][]example{examples, stateExamples} {
		for rule, cases := range table {
			for _, c := range cases {
				// A blank value is not checked by any rule: required says it is missing.
				if c.in == "" {
					if !passes(t, rule, c.in) {
						t.Errorf("%s refused a blank value; that is required's to say", rule)
					}
					continue
				}
				if got := passes(t, rule, c.in); got != c.valid {
					t.Errorf("%s on %q = %v, want %v %s", rule, c.in, got, c.valid, c.note)
				}
			}
		}
	}
}

func TestEveryRuleHasAtLeastOneExampleOfEachKind(t *testing.T) {
	t.Parallel()
	covered := map[string]bool{}
	for rule, cases := range examples {
		valid, invalid := 0, 0
		for _, c := range cases {
			if c.valid {
				valid++
			} else if c.in != "" {
				invalid++
			}
		}
		covered[rule] = true
		if valid == 0 || invalid == 0 {
			t.Errorf("%s has %d valid and %d invalid examples", rule, valid, invalid)
		}
	}
	for _, rule := range br.Rules() {
		for _, name := range rule.Names() {
			if !covered[name] && name != "inscricao_estadual" {
				t.Errorf("rule %s has no examples", name)
			}
		}
	}
	for uf := range stateExamples {
		covered[uf] = true
	}
	for _, uf := range br.UFValues() {
		if !covered["inscricao_estadual:"+string(uf)] {
			t.Errorf("no state registration example for %s", uf)
		}
	}
}

func TestARuleWithoutParametersRefusesAnArgumentAndANonString(t *testing.T) {
	t.Parallel()
	for _, rule := range br.Rules() {
		if len(rule.Parameters) != 0 {
			continue
		}
		for _, name := range rule.Names() {
			set, err := validation.Compile(validation.Rules{"field": name + ":x"})
			if err != nil {
				t.Fatalf("compiling %s:x: %v", name, err)
			}
			// Find a value the rule accepts: the first valid example.
			var good string
			for _, c := range examples[name] {
				if c.valid {
					good = c.in
					break
				}
			}
			if installed.Make(validation.Data{"field": good}, set).Passes() {
				t.Errorf("%s:x accepted %q; an argument the rule ignores is a typo nothing reports", name, good)
			}
			plain := validation.MustCompile(validation.Rules{"field": name, "other": "nullable|" + name})
			if installed.Make(validation.Data{"field": 12345678901}, plain).Passes() {
				t.Errorf("%s accepted a number", name)
			}
			if errs := failuresOf(installed, plain, validation.Data{"field": "", "other": nil}); errs.Any() {
				t.Errorf("%s refused a blank value and a null beside nullable: %v", name, errs)
			}
		}
	}
}

// failuresOf runs the set through a factory and returns the messages by field.
func failuresOf(f *validation.Factory, set *validation.Set, data validation.Data) validation.Errors {
	v := f.Make(data, set)
	v.Passes()
	return v.Errors()
}
