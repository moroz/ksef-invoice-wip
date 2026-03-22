package invoice_test

import (
	"encoding/xml"
	"ksef-go/lib/ksef/invoice"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerializeInvoice(t *testing.T) {
	one := 1

	instance := invoice.Faktura{
		Naglowek: invoice.Naglowek{
			KodFormularza: invoice.KodFormularza{
				Value:        "FA",
				KodSystemowy: "FA (3)",
				WersjaSchemy: "1-0E",
			},
			WariantFormularza: 3,
			DataWytworzeniaFa: time.Now().Format(time.RFC3339),
			SystemInfo:        "Aplikacja Podatnika KSeF",
		},
		Podmiot1: invoice.Podmiot1{
			DaneIdentyfikacyjne: invoice.Podmiot1Ident{
				NIP:   "6692435667",
				Nazwa: "Wyrób Robotów CG",
			},
			Adres: invoice.Adres{
				KodKraju: "PL",
				AdresL1:  "Test Street 123",
			},
		},
		Podmiot2: invoice.Podmiot2{
			DaneIdentyfikacyjne: invoice.Podmiot2Ident{
				KodUE:   "DE",
				NrVatUE: "123456123",
				Nazwa:   "Test Buyer",
			},
			Adres: &invoice.Adres{
				KodKraju: "DE",
				AdresL1:  "Teststraße 123",
			},
			JST: 2,
			GV:  2,
		},
		Fa: invoice.Fa{
			KodWaluty: "EUR",
			P1:        "2026-03-22",
			P1M:       "Pcim Dolny",
			P2:        "01/01/2026",
			P6:        "2026-03-22",
			P1310:     "2137",
			P15:       "2137",
			Adnotacje: invoice.Adnotacje{
				P16:  2,
				P17:  2,
				P18:  1,
				P18A: 2,
				Zwolnienie: invoice.Zwolnienie{
					P19N: &one,
				},
				NoweSrodkiTransportu: invoice.NoweSrodkiTransportu{
					P22N: &one,
				},
				P23: 2,
				PMarzy: invoice.PMarzy{
					PPMarzyN: &one,
				},
			},
			RodzajFaktury: "VAT",
			FaWiersz: []invoice.FaWiersz{
				{
					NrWierszaFa: 1,
					P7:          "Nazwa towaru 1",
					P8A:         "szt.",
					P8B:         "1",
					P9A:         "2137",
					P11:         "2137",
					P12:         "oo",
					KursWaluty:  "4.24",
				},
			},
			Platnosc: &invoice.Platnosc{
				RachunekBankowy: []invoice.RachunekBankowy{
					{
						NrRB:  "PL123456789",
						SWIFT: "REVOLT21",
					},
				},
			},
		},
	}

	asXml, err := xml.Marshal(instance)
	assert.NoError(t, err)
	assert.NotEmpty(t, asXml)
}
