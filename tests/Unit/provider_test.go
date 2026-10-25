package unit_test

import (
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/translation"
	"github.com/arandu-io/hesape/validation"
	br "github.com/hyz-is/arandu-br"
)

// These tests are the acceptance of the module: a real Factory of hesape, the
// real pt-BR translator, and the rule sets an application writes -- one in a
// package-level variable, compiled before the Factory exists, as the
// application's request types are, with the option that lets it name the
// rules a module declares later.

// storeCustomer is the criterion: a controller validates these four, in the
// list form and the string form, mixing the standard rules and the module's.
var storeCustomer = validation.MustCompile(validation.Rules{
	"cpf":     []any{"required", "cpf"},
	"cnpj":    "required|cnpj",
	"cep":     []any{"required", "formato_cep"},
	"celular": "required|celular_com_ddd",
	"nome":    "nullable|string|max:20",
}, validation.WithModuleRules())

// everyRule names every rule and alias of the module, once.
var everyRule = func() *validation.Set {
	rules := validation.Rules{}
	for _, def := range br.Defs() {
		for _, name := range append([]string{def.Name}, def.Aliases...) {
			rule := name
			if def.Params.Max != 0 {
				rule += ":SP"
			}
			rules["f_"+name] = rule
		}
	}
	return validation.MustCompile(rules, validation.WithModuleRules())
}()

func ptBR() *translation.Translator { return translation.New(nil, "pt-BR", "en") }

func realFactory(tr *translation.Translator) *validation.Factory {
	return validation.NewFactory(validation.FromTranslation(tr)).Use(br.Rules())
}

func TestTheCriterionRunsOnARealFactoryInPortuguese(t *testing.T) {
	t.Parallel()

	factory := realFactory(ptBR())
	if err := factory.Check(storeCustomer, everyRule); err != nil {
		t.Fatalf("Check: %v", err)
	}

	cpf, cnpj := generatedCPFs(1)[0], generatedCNPJs(2)[1]
	valid := validation.Data{
		"cpf": punctuatedCPF(cpf), "cnpj": strings.ToLower(cnpj), "cep": "70.040-010", "celular": "(11) 91234-5678",
		"is_admin": "true", // nothing declares it
	}
	input, err := factory.Validate(valid, storeCustomer, validation.WithLocale("pt-BR"))
	if err != nil {
		t.Fatalf("a valid request failed: %v", err)
	}
	var kept []string
	for key := range input.Data() {
		kept = append(kept, key)
	}
	sort.Strings(kept)
	if strings.Join(kept, ",") != "celular,cep,cnpj,cpf" {
		t.Errorf("validated holds %v, want only the four declared fields that were sent", kept)
	}
	if input.String("cpf") != punctuatedCPF(cpf) {
		t.Errorf("validated cpf = %q", input.String("cpf"))
	}

	_, err = factory.Validate(validation.Data{"cpf": "123", "cnpj": "12.345.678/0001-00", "cep": "70040010", "celular": "(11) 3000-0000", "nome": strings.Repeat("a", 30)},
		storeCustomer, validation.WithLocale("pt-BR"))
	var invalid *validation.ValidationException
	if !errors.As(err, &invalid) {
		t.Fatalf("err = %v, want a validation exception", err)
	}
	errs := invalid.Errors()
	for field, want := range map[string]string{
		"cpf":     "O campo cpf não é um CPF válido.",
		"cnpj":    "O campo cnpj não é um CNPJ válido.",
		"cep":     "O campo cep não possui um formato válido de CEP.",
		"celular": "O campo celular não é um celular com DDD válido.",
	} {
		if got := errs[field]; len(got) != 1 || got[0] != want {
			t.Errorf("%s: messages = %q, want %q", field, got, want)
		}
	}
	if got := errs["nome"]; len(got) != 1 || !strings.Contains(got[0], "nome") {
		t.Errorf("a standard rule beside the module's: %q", got)
	}

	// A required field that is absent says so, in Portuguese, from hesape's catalogue.
	_, err = factory.Validate(validation.Data{}, storeCustomer, validation.WithLocale("pt-BR"))
	if !errors.As(err, &invalid) || invalid.Errors()["cpf"][0] != "O campo cpf é obrigatório." {
		t.Errorf("required says %v", err)
	}
}

func TestEveryRuleOfTheModulePassesTheBootCheckAndAnUnknownOneFailsIt(t *testing.T) {
	t.Parallel()

	factory := realFactory(ptBR())
	if err := factory.Check(everyRule); err != nil {
		t.Fatalf("Check: %v", err)
	}

	// Without the module the same set is a boot failure naming every rule.
	bare := validation.NewFactory(validation.FromTranslation(ptBR()))
	err := bare.Check(everyRule)
	if err == nil {
		t.Fatal("Check passed on a factory that never used the module")
	}
	for _, want := range []string{"formato_cep", "placa", "inscricao_estadual"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the boot failure does not name %s: %v", want, err)
		}
	}
}

func TestTheApplicationsSentenceOverridesTheModulesAndTheRestKeepTheirs(t *testing.T) {
	t.Parallel()

	tr := ptBR()
	tr.AddLines(map[string]string{"validation.cpf": "Confira o CPF de :attribute."}, "pt-BR", "")
	factory := realFactory(tr)
	set := factory.MustCompile(validation.Rules{"cpf": "cpf", "cep": "formato_cep"})
	v := factory.Make(validation.Data{"cpf": "1", "cep": "1"}, set, validation.WithLocale("pt-BR"))
	v.Passes()
	if got := v.Errors().Get("cpf"); len(got) != 1 || got[0] != "Confira o CPF de cpf." {
		t.Errorf("cpf: %q, want the application's sentence", got)
	}
	if got := v.Errors().Get("cep"); len(got) != 1 || got[0] != "O campo cep não possui um formato válido de CEP." {
		t.Errorf("cep: %q, want the module's sentence", got)
	}

	// And a custom message for one field beats both.
	v = factory.Make(validation.Data{"cpf": "1", "cep": "1"}, set, validation.WithLocale("pt-BR"),
		validation.WithCustomMessages(map[string]any{"cep.formato_cep": "CEP inválido."}))
	v.Passes()
	if got := v.Errors().Get("cep"); len(got) != 1 || got[0] != "CEP inválido." {
		t.Errorf("cep: %q", got)
	}
}

func TestTheStateRegistrationSentenceNamesTheUnitThroughTheRegistry(t *testing.T) {
	t.Parallel()

	factory := realFactory(ptBR())
	set := factory.MustCompile(validation.Rules{"ie": "inscricao_estadual:sp"})
	v := factory.Make(validation.Data{"ie": "123"}, set, validation.WithLocale("pt-BR"))
	v.Passes()
	if got := v.Errors().Get("ie"); len(got) != 1 || got[0] != "O campo ie não é uma inscrição estadual válida para SP." {
		t.Errorf("messages = %q", got)
	}
}

func TestTheModuleIsDeclaredOnTheFactoryAndNowhereElse(t *testing.T) {
	t.Parallel()

	// A factory that never used the module does not know the rules the module
	// declares. The eight names the import still registers are the exception.
	bare := validation.NewFactory(nil)
	for _, name := range []string{"formato_cep", "placa", "inscricao_estadual:SP", "telefone"} {
		if _, err := bare.Compile(validation.Rules{"f": name}); err == nil {
			t.Errorf("%s compiled on a factory that never used the module", name)
		}
	}
	registry := realFactory(ptBR()).Registry()
	for _, name := range []string{"formato_cep", "placa", "cpf_cnpj", "moeda", "inscricao_estadual"} {
		if !registry.Has(name) {
			t.Errorf("the registry does not declare %s", name)
		}
	}
	// What the import registered is not declared twice.
	if registry.Has("cpf") {
		t.Error("cpf is declared in the registry although the process-wide catalogue holds it")
	}
}
