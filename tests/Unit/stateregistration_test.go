package unit_test

import (
	"testing"

	br "github.com/hyz-is/arandu-br"
)

// The numbers below are the ones the PHP packages test with: one valid for
// each unit and the variants of the rule that each unit has -- the digit that
// ten turns into zero, the older and the newer number of a unit, the prefix
// that is not the unit's, a digit too many. The PHP packages publish them as
// specimens from number generators; none is read from a record.
type stateCase struct {
	uf, in string
	valid  bool
}

var stateCases = []stateCase{
	{"AC", "0108368143106", true},
	{"AC", "0187634580933", false},
	{"AC", "0187634580924", false},
	{"AC", "0018763458000", false},
	{"AC", "01018763458064", false},
	{"AL", "248659758", true},
	{"AL", "247424170", true},
	{"AL", "248659759", false},
	{"AL", "258659750", false},
	{"AL", "2486597584", false},
	{"AP", "036029572", true},
	{"AP", "030123459", true},
	{"AP", "030000080", true},
	{"AP", "030000160", true},
	{"AP", "030170011", true},
	{"AP", "030170020", true},
	{"AP", "030170071", true},
	{"AP", "036029573", false},
	{"AP", "0306029570", false},
	{"AP", "003060292", false},
	{"AM", "036029572", true},
	{"AM", "000000019", true},
	{"AM", "046893830", true},
	{"AM", "036029573", false},
	{"AM", "0036029572", false},
	{"BA", "12345663", true},
	{"BA", "74219145", true},
	{"BA", "038343081", true},
	{"BA", "100000306", true},
	{"BA", "778514741", true},
	{"BA", "078771760", true},
	{"BA", "039474751", true},
	{"BA", "090529323", true},
	{"BA", "04772253", true},
	{"BA", "12345636", false},
	{"BA", "74219154", false},
	{"BA", "038343001", false},
	{"BA", "778514731", false},
	{"BA", "0012345636", false},
	{"CE", "853511942", true},
	{"CE", "853511943", false},
	{"CE", "0853511942", false},
	{"DF", "0754002000176", true},
	{"DF", "0801949600112", true},
	{"DF", "0754002000508", true},
	{"DF", "07008368143094", false},
	{"DF", "0754002000175", false},
	{"ES", "639191444", true},
	{"ES", "639191445", false},
	{"ES", "0639191444", false},
	{"GO", "109161793", true},
	{"GO", "101031131", true},
	{"GO", "101030940", true},
	{"GO", "109161794", false},
	{"GO", "121031131", false},
	{"GO", "0101030940", false},
	{"MA", "120000008", true},
	{"MA", "120000040", true},
	{"MA", "120000130", true},
	{"MA", "0120000008", false},
	{"MA", "109161793", false},
	{"MA", "120000007", false},
	{"MT", "00000000000", true},
	{"MT", "0000000000", false},
	{"MT", "12345678901", false},
	{"MS", "280000006", true},
	{"MS", "280000090", true},
	{"MS", "280000030", true},
	{"MS", "0280000006", false},
	{"MS", "853511942", false},
	{"MS", "280000031", false},
	{"MG", "4333908330177", true},
	{"MG", "4333908330410", true},
	{"MG", "4333908332560", true},
	{"MG", "04333908330177", false},
	{"MG", "4333908330176", false},
	{"MG", "4333908330167", false},
	{"PA", "150000006", true},
	{"PA", "150000030", true},
	{"PA", "150000260", true},
	{"PA", "0150000006", false},
	{"PA", "120000008", false},
	{"PA", "150000007", false},
	{"PB", "853511942", true},
	{"PB", "853511950", true},
	{"PB", "853512230", true},
	{"PB", "853511943", false},
	{"PB", "0853511942", false},
	{"PR", "4447953604", true},
	{"PR", "4447953640", false},
	{"PR", "04447953604", false},
	{"PE", "288625706", true},
	{"PE", "925870101", false},
	{"PE", "0925870110", false},
	{"PI", "052364534", true},
	{"RJ", "62545372", true},
	{"RJ", "62545380", true},
	{"RJ", "62545470", true},
	{"RJ", "20441620", false},
	{"RJ", "020441623", false},
	{"RN", "2007693232", true},
	{"RN", "2003569880", true},
	{"RN", "203569881", true},
	{"RN", "0203569881", false},
	{"RN", "20356988104", false},
	{"RN", "2007693231", false},
	{"RS", "0305169149", true},
	{"RS", "1202762660", true},
	{"RS", "1202762120", true},
	{"RS", "02007693230", false},
	{"RS", "2007693232", false},
	{"RO", "01078042249629", true},
	{"RO", "01078042249670", true},
	{"RO", "01078042249751", true},
	{"RO", "001078042249627", false},
	{"RO", "01078042249756", false},
	{"RR", "240061536", true},
	{"RR", "024006150", false},
	{"RR", "2400615366", false},
	{"RR", "240061537", false},
	{"SC", "330430572", true},
	{"SP", "110042490114", true},
	{"SP", "1110042494114", false},
	{"SP", "110042490113", false},
	{"SP", "110042498113", false},
	{"SE", "017682606", true},
	{"TO", "01027737427", true},
	{"TO", "294467696", true},
	{"TO", "294150870", true},
	{"TO", "01047737427", false},
	{"TO", "099999916599", false},
	{"TO", "99999916598", false},
	{"TO", "294467690", false},
}

// the one number the PHP library accepts and this package does not: eleven
// zeros, a repetition of one digit, whose check digit computes.
const zeroedMatoGrosso = "00000000000"

func TestParseStateRegistrationAgreesWithThePublishedSpecimens(t *testing.T) {
	t.Parallel()
	for _, c := range stateCases {
		want := c.valid
		if c.uf == "MT" && c.in == zeroedMatoGrosso {
			want = false
		}
		got, err := br.ParseStateRegistration(br.UF(c.uf), c.in)
		if (err == nil) != want {
			t.Errorf("ParseStateRegistration(%s, %q) = %q, %v; want valid = %v", c.uf, c.in, got, err, want)
		}
		if err == nil && got != c.in {
			t.Errorf("ParseStateRegistration(%s, %q) = %q, want the digits back", c.uf, c.in, got)
		}
	}
}

func TestParseStateRegistrationForEveryUnit(t *testing.T) {
	t.Parallel()
	// A valid number of each of the twenty-seven units.
	valid := map[br.UF]string{
		br.UFAC: "0184765932153", br.UFAL: "248308335", br.UFAP: "039895661", br.UFAM: "292278012", br.UFBA: "96338555",
		br.UFCE: "980874165", br.UFDF: "0740769300107", br.UFES: "342048090", br.UFGO: "107900459", br.UFMA: "122917510",
		br.UFMT: "55160385510", br.UFMS: "283814659", br.UFMG: "0526366324132", br.UFPA: "158330005", br.UFPB: "439622301",
		br.UFPR: "5714953410", br.UFPE: "920145698", br.UFPI: "758505183", br.UFRJ: "61804315", br.UFRN: "208627715",
		br.UFRS: "3821957672", br.UFRO: "52059985926850", br.UFRR: "249977060", br.UFSC: "696192667", br.UFSP: "653172024009",
		br.UFSE: "646597361", br.UFTO: "75036274184",
	}
	if len(valid) != len(br.UFValues()) {
		t.Fatalf("%d numbers for %d units", len(valid), len(br.UFValues()))
	}
	for uf, number := range valid {
		if got, err := br.ParseStateRegistration(uf, number); err != nil || got != number {
			t.Errorf("ParseStateRegistration(%s, %q) = %q, %v", uf, number, got, err)
		}
		if !br.ValidStateRegistration(uf, number) {
			t.Errorf("ValidStateRegistration(%s, %q) = false", uf, number)
		}
		// The last digit is always an output of the others, so one of the ten
		// digits fits and the other nine are refused.
		fits := 0
		for d := byte('0'); d <= '9'; d++ {
			if br.ValidStateRegistration(uf, number[:len(number)-1]+string(d)) {
				fits++
			}
		}
		if fits != 1 {
			t.Errorf("%s: %d last digits fit %q, want exactly one", uf, fits, number)
		}
	}
}

func TestParseStateRegistrationReadsThePunctuationOfTheUnits(t *testing.T) {
	t.Parallel()
	for _, c := range []struct{ uf, in, want string }{
		{"BA", "826793-41", "82679341"},
		{"BA", "82679341", "82679341"},
		{"SP", "653.172.024.009", "653172024009"},
		{"SP", "653.172.024/009", "653172024009"},
		{"MG", "052.636.632/4132", "0526366324132"},
		{"RO", "5205998-5926850", "52059985926850"},
		{"PR", "571.49534-10", "5714953410"},
	} {
		got, err := br.ParseStateRegistration(br.UF(c.uf), c.in)
		if err != nil || got != c.want {
			t.Errorf("ParseStateRegistration(%s, %q) = %q, %v; want %q", c.uf, c.in, got, err, c.want)
		}
	}
	for _, c := range []struct{ uf, in string }{
		{"BA", "8 2__6-7*9.3/41"}, // the PHP library forgives anything; this one does not
		{"BA", " 82679341"},
		{"BA", "82679341 "},
		{"BA", "-82679341"},
		{"BA", "82679341-"},
		{"BA", "826793--41"},
		{"BA", "826793 41"},
		{"BA", "826793,41"},
		{"BA", ""},
		{"BA", "-"},
		{"BA", "٨٢٦٧٩٣٤١"},
		{"BA", "82679342"}, // the digit is wrong
		{"SP", "82679341"}, // the unit is not Bahia
		{"BA", "00000000"}, // a repetition
		{"ZA", "82679341"}, // no such unit
		{"ba", "82679341"}, // the unit is the code, in capitals
		{"", "82679341"},
		{"12", "82679341"},
		{"BA", "ISENTO"},
	} {
		if br.ValidStateRegistration(br.UF(c.uf), c.in) {
			t.Errorf("ParseStateRegistration(%s, %q) accepted", c.uf, c.in)
		}
		if _, err := br.ParseStateRegistration(br.UF(c.uf), c.in); err != br.ErrInvalidStateRegistration {
			t.Errorf("ParseStateRegistration(%s, %q) error = %v, want ErrInvalidStateRegistration", c.uf, c.in, err)
		}
	}
}

func FuzzParseStateRegistration(f *testing.F) {
	for _, seed := range []string{"", "826793-41", "653.172.024.009", "0526366324132", "-", "8 2"} {
		f.Add("BA", seed)
		f.Add("SP", seed)
	}
	f.Fuzz(func(t *testing.T, uf, s string) {
		if digits, err := br.ParseStateRegistration(br.UF(uf), s); err == nil {
			for i := 0; i < len(digits); i++ {
				if digits[i] < '0' || digits[i] > '9' {
					t.Fatalf("ParseStateRegistration(%q, %q) = %q, not digits", uf, s, digits)
				}
			}
		}
	})
}
