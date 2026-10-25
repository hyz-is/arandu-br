package unit_test

import (
	"strings"
	"testing"

	"github.com/arandu-io/hesape/validation"
	br "github.com/hyz-is/arandu-br"
)

// linesTranslator is a validation.Translator over a map of lines, which is how
// an application's language file reaches the validator.
type linesTranslator map[string]string

func (l linesTranslator) Get(key string, _ map[string]any, _ string) any {
	if line, ok := l[key]; ok {
		return line
	}
	return key
}

func (l linesTranslator) Choice(key string, _ int, _ map[string]any, _ string) string { return key }

// refusedValue is a value no rule here accepts, and a rule string for each
// name that makes it a string the rule is asked about.
func refusedValue(rule string) string {
	switch {
	case strings.HasPrefix(rule, "inscricao_estadual"):
		return "123"
	}
	return "!"
}

func TestEveryDeclaredNameResolvesAndSaysItsSentenceWithTheAttribute(t *testing.T) {
	t.Parallel()

	seen := map[string]bool{}
	for _, rule := range br.Rules() {
		if rule.Name == "" || rule.Check == nil {
			t.Fatalf("a rule has no name or no check: %+v", rule)
		}
		message := rule.Message("pt-BR")
		if !strings.HasPrefix(message, "O campo :attribute ") || !strings.HasSuffix(message, ".") {
			t.Errorf("%s: sentence %q does not read 'O campo :attribute ... .'", rule.Name, message)
		}
		if rule.Message("xx-XX") != message {
			t.Errorf("%s: an unknown locale does not fall back to pt-BR", rule.Name)
		}
		for _, parameter := range rule.Parameters {
			if !strings.Contains(message, ":"+parameter) {
				t.Errorf("%s: sentence %q does not hold :%s", rule.Name, message, parameter)
			}
		}
		for _, name := range rule.Names() {
			if seen[name] {
				t.Errorf("%s is declared twice", name)
			}
			seen[name] = true

			ruleString := name
			if len(rule.Parameters) > 0 {
				ruleString += ":BA"
			}
			set, err := validation.Compile(validation.Rules{"documento": ruleString})
			if err != nil {
				t.Errorf("%s does not resolve in the catalogue: %v", ruleString, err)
				continue
			}
			errs := failuresOf(installed, set, validation.Data{"documento": refusedValue(name)})
			got := errs.Get("documento")
			want := strings.ReplaceAll(strings.ReplaceAll(message, ":attribute", "documento"), ":uf", "BA")
			if len(got) != 1 || got[0] != want {
				t.Errorf("%s: messages = %q, want %q", ruleString, got, want)
			}
		}
	}
}

func TestEveryNameTheTaskNamesIsDeclared(t *testing.T) {
	t.Parallel()
	// The names of geekcom/validator-docs and laravellegends/pt-br-validator,
	// and the ones the Arandu validation mission lists.
	want := []string{
		"cpf", "cnpj", "cpf_ou_cnpj", "formato_cpf", "formato_cnpj", "formato_cpf_ou_cnpj", "cep", "formato_cep",
		"celular", "celular_com_ddd", "celular_com_codigo", "celular_com_codigo_sem_mascara",
		"telefone", "telefone_com_ddd", "telefone_com_codigo", "placa", "formato_placa_de_veiculo", "renavam", "cnh",
		"titulo_eleitor", "nis", "pis", "formato_nis", "formato_pis", "cns", "certidao", "formato_certidao", "uf", "ddd",
		"passaporte", "inscricao_estadual", "data", "horario", "moeda", "brl",
		// the spellings that were or are another name for one of them
		"cpf_or_cnpj", "cpf_cnpj", "formato_cpf_cnpj",
	}
	declared := map[string]bool{}
	for _, rule := range br.Rules() {
		for _, name := range rule.Names() {
			declared[name] = true
		}
	}
	for _, name := range want {
		if !declared[name] {
			t.Errorf("rule %s is not declared", name)
		}
	}
}

func TestTheSentencesGoThroughAnApplicationsLanguageFile(t *testing.T) {
	t.Parallel()

	f := validation.NewFactory(linesTranslator{
		"validation.cpf":         "Informe um CPF em :attribute.",
		"validation.formato_cep": "O CEP de :attribute vai com hífen.",
	})
	br.Install(f)
	set := validation.MustCompile(validation.Rules{"documento": "cpf", "endereco": "formato_cep", "unidade": "uf"})
	errs := failuresOf(f, set, validation.Data{"documento": "1", "endereco": "70040010", "unidade": "XX"})
	for field, want := range map[string]string{
		"documento": "Informe um CPF em documento.",
		"endereco":  "O CEP de endereco vai com hífen.",
		"unidade":   "O campo unidade não é uma UF válida.",
	} {
		if got := errs.Get(field); len(got) != 1 || got[0] != want {
			t.Errorf("%s: messages = %q, want %q", field, got, want)
		}
	}
}

func TestAnAttributeNameIsTheOneTheApplicationGives(t *testing.T) {
	t.Parallel()

	f := validation.NewFactory(nil)
	f.SetAttributeNames(map[string]string{"documento": "CPF do cliente"})
	br.Install(f)
	set := validation.MustCompile(validation.Rules{"documento": "cpf"})
	if got := failuresOf(f, set, validation.Data{"documento": "1"}).Get("documento"); len(got) != 1 || got[0] != "O campo CPF do cliente não é um CPF válido." {
		t.Errorf("messages = %q", got)
	}
}

func TestTheStateRegistrationSentenceNamesTheUnit(t *testing.T) {
	t.Parallel()

	for rule, want := range map[string]string{
		"inscricao_estadual:sp": "O campo ie não é uma inscrição estadual válida para SP.",
		"inscricao_estadual:BA": "O campo ie não é uma inscrição estadual válida para BA.",
		"inscricao_estadual":    "O campo ie não é uma inscrição estadual válida para a UF informada.",
	} {
		set := validation.MustCompile(validation.Rules{"ie": rule})
		if got := failuresOf(installed, set, validation.Data{"ie": "123"}).Get("ie"); len(got) != 1 || got[0] != want {
			t.Errorf("%s: messages = %q, want %q", rule, got, want)
		}
	}
}

func TestMessagesListsEveryNameInTheLocale(t *testing.T) {
	t.Parallel()

	messages := br.Messages("pt-BR")
	for _, rule := range br.Rules() {
		for _, name := range rule.Names() {
			if messages[name] != rule.Message("pt-BR") {
				t.Errorf("Messages(pt-BR)[%q] = %q, want %q", name, messages[name], rule.Message("pt-BR"))
			}
		}
	}
	if messages["cpf_or_cnpj"] != messages["cpf_ou_cnpj"] {
		t.Error("an alias does not say what its rule says")
	}
}

func TestInstallIsRepeatableAndReachesTheCatalogueWithoutAFactory(t *testing.T) {
	t.Parallel()

	// A second factory, and a nil one, find every name already registered and
	// register nothing twice: a panic here is the catalogue refusing a name.
	br.Install(validation.NewFactory(nil))
	br.Install(nil)
}

func TestTheConstantsOfTheTwoRenamedRulesHoldTheNewNames(t *testing.T) {
	t.Parallel()

	if br.RuleCPFOrCNPJ != "cpf_ou_cnpj" || br.RuleMobile != "celular_com_ddd" || br.RuleLocalMobile != "celular" {
		t.Errorf("constants hold %q, %q and %q", br.RuleCPFOrCNPJ, br.RuleMobile, br.RuleLocalMobile)
	}
}

func TestCelularIsTheNumberWithoutTheAreaCodeAndCelularComDDDIsTheOldOne(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		rule, in string
		want     bool
	}{
		// What celular accepted until this release, now under celular_com_ddd.
		{"celular_com_ddd", "(11) 91234-5678", true},
		{"celular_com_ddd", "11912345678", true},
		{"celular_com_ddd", "(11)91234-5678", true},
		{"celular", "(11) 91234-5678", false},
		{"celular", "11912345678", false},
		// What celular means in the PHP packages.
		{"celular", "91234-5678", true},
		{"celular", "912345678", true},
		{"celular_com_ddd", "91234-5678", false},
	} {
		if got := passes(t, c.rule, c.in); got != c.want {
			t.Errorf("%s on %q = %v, want %v", c.rule, c.in, got, c.want)
		}
	}
}

func TestDocumentsAreHandledAsTextWithTheirLeadingZeros(t *testing.T) {
	t.Parallel()

	// A CPF and a CNPJ that start with zeros keep them, and as numbers they
	// could not have been written.
	cpf := "00000000191"
	cnpj := "00000000000191"
	for rule, in := range map[string]string{"cpf": cpf, "cnpj": cnpj, "cpf_ou_cnpj": cnpj, "formato_cpf": "000.000.001-91"} {
		if !passes(t, rule, in) {
			t.Errorf("%s refused %q", rule, in)
		}
	}
	for _, in := range []string{cpf, cnpj} {
		got, err := map[bool]func(string) (string, error){true: br.ParseCNPJ, false: br.ParseCPF}[len(in) == 14](in)
		if err != nil || got != in {
			t.Errorf("parsing %q = %q, %v; the leading zeros are part of the document", in, got, err)
		}
	}
	// A number loses them, and no rule takes a number.
	set := validation.MustCompile(validation.Rules{"cpf": "cpf", "cnpj": "cnpj"})
	if errs := failuresOf(installed, set, validation.Data{"cpf": 191, "cnpj": int64(191)}); len(errs) != 2 {
		t.Errorf("a document sent as a number was not refused: %v", errs)
	}
}
