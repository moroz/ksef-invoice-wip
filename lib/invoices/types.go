package invoices

import (
	"encoding/xml"
	"time"
)

type Invoice struct {
	XMLName  xml.Name       `xml:"http://crd.gov.pl/wzor/2025/06/25/13775/ Faktura"`
	XMLNSXsi string         `xml:"xmlns:xsi,attr"`
	XMLNSXsd string         `xml:"xmlns:xsd,attr"`
	Header   Header         `xml:"Naglowek"`
	Entity1  Entity         `xml:"Podmiot1"`
	Entity2  Entity         `xml:"Podmiot2"`
	Invoice  InvoiceDetails `xml:"Fa"`
}

type Header struct {
	FormCode     FormCode  `xml:"KodFormularza"`
	FormVariant  string    `xml:"WariantFormularza"`
	CreationDate time.Time `xml:"DataWytworzeniaFa"`
	SystemInfo   string    `xml:"SystemInfo"`
}

type FormCode struct {
	Code          string `xml:",chardata"`
	SystemCode    string `xml:"kodSystemowy,attr"`
	SchemaVersion string `xml:"wersjaSchemy,attr"`
}

type Entity struct {
	IdentificationData IdentificationData `xml:"DaneIdentyfikacyjne"`
	Address            Address            `xml:"Adres"`
	JST                string             `xml:"JST,omitempty"`
	GV                 string             `xml:"GV,omitempty"`
}

type IdentificationData struct {
	TaxID   string `xml:"NIP,omitempty"`
	EUCode  string `xml:"KodUE,omitempty"`
	EUVatID string `xml:"NrVatUE,omitempty"`
	Name    string `xml:"Nazwa"`
}

type Address struct {
	CountryCode  string `xml:"KodKraju"`
	AddressLine1 string `xml:"AdresL1"`
	AddressLine2 string `xml:"AdresL2,omitempty"`
}

type InvoiceDetails struct {
	CurrencyCode string        `xml:"KodWaluty"`
	P1           string        `xml:"P_1"`
	P2           string        `xml:"P_2"`
	P6           string        `xml:"P_6"`
	P13_10       string        `xml:"P_13_10"`
	P15          string        `xml:"P_15"`
	Annotations  Annotations   `xml:"Adnotacje"`
	InvoiceType  string        `xml:"RodzajFaktury"`
	InvoiceLines []InvoiceLine `xml:"FaWiersz"`
}

type Annotations struct {
	P16          string       `xml:"P_16"`
	P17          string       `xml:"P_17"`
	P18          string       `xml:"P_18"`
	P18A         string       `xml:"P_18A"`
	Exemption    Exemption    `xml:"Zwolnienie"`
	NewTransport NewTransport `xml:"NoweSrodkiTransportu"`
	P23          string       `xml:"P_23"`
	Margin       Margin       `xml:"PMarzy"`
}

type Exemption struct {
	P19N string `xml:"P_19N"`
}

type NewTransport struct {
	P22N string `xml:"P_22N"`
}

type Margin struct {
	PMarzyN string `xml:"P_PMarzyN"`
}

type InvoiceLine struct {
	LineNumber   string `xml:"NrWierszaFa"`
	P7           string `xml:"P_7"`
	P8A          string `xml:"P_8A"`
	P8B          string `xml:"P_8B"`
	P9A          string `xml:"P_9A"`
	P11          string `xml:"P_11"`
	P12          string `xml:"P_12"`
	ExchangeRate string `xml:"KursWaluty"`
}
