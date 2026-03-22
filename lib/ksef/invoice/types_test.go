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

	instance := invoice.Invoice{
		Header: invoice.Header{
			FormCode: invoice.FormCode{
				Value:         "FA",
				SystemCode:    "FA (3)",
				SchemaVersion: "1-0E",
			},
			FormVariant:  3,
			CreationDate: time.Now().Format(time.RFC3339),
			SystemInfo:   "Aplikacja Podatnika KSeF",
		},
		Seller: invoice.Seller{
			IdentData: invoice.SellerIdent{
				NIP:  "6692435667",
				Name: "Wyrób Robotów CG",
			},
			Address: invoice.Address{
				CountryCode:  "PL",
				AddressLine1: "Test Street 123",
			},
		},
		Buyer: invoice.Buyer{
			IdentData: invoice.BuyerIdent{
				EUCode:      "DE",
				EUVatNumber: "123456123",
				Name:        "Test Buyer",
			},
			Address: &invoice.Address{
				CountryCode:  "DE",
				AddressLine1: "Teststraße 123",
			},
			JST: 2,
			GV:  2,
		},
		Body: invoice.InvoiceBody{
			CurrencyCode: "EUR",
			P1:           "2026-03-22",
			P1M:          "Pcim Dolny",
			P2:           "01/01/2026",
			P6:           "2026-03-22",
			P1310:        "2137",
			P15:          "2137",
			Annotations: invoice.Annotations{
				P16:  2,
				P17:  2,
				P18:  1,
				P18A: 2,
				Exemption: invoice.Exemption{
					P19N: &one,
				},
				NewTransport: invoice.NewTransportMeans{
					P22N: &one,
				},
				P23: 2,
				Margin: invoice.Margin{
					PPMarzyN: &one,
				},
			},
			InvoiceType: "VAT",
			Lines: []invoice.InvoiceLine{
				{
					LineNumber:   1,
					P7:           "Nazwa towaru 1",
					P8A:          "szt.",
					P8B:          "1",
					P9A:          "2137",
					P11:          "2137",
					P12:          "oo",
					ExchangeRate: "4.24",
				},
			},
			Payment: &invoice.Payment{
				BankAccounts: []invoice.BankAccount{
					{
						AccountNumber: "PL123456789",
						SWIFT:         "REVOLT21",
					},
				},
			},
		},
	}

	asXml, err := xml.Marshal(instance)
	assert.NoError(t, err)
	assert.NotEmpty(t, asXml)
}
