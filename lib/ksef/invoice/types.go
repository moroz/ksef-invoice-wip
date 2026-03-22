package invoice

import "encoding/xml"

const Namespace = "http://crd.gov.pl/wzor/2025/06/25/13775/"

// Invoice is the root element of a KSeF FA(3) invoice.
type Invoice struct {
	XMLName          xml.Name          `xml:"http://crd.gov.pl/wzor/2025/06/25/13775/ Faktura"`
	Header           Header            `xml:"Naglowek"`
	Seller           Seller            `xml:"Podmiot1"`
	Buyer            Buyer             `xml:"Podmiot2"`
	ThirdParties     []ThirdParty      `xml:"Podmiot3,omitempty"`
	AuthorizedEntity *AuthorizedEntity `xml:"PodmiotUpowazniony,omitempty"`
	Body             InvoiceBody       `xml:"Fa"`
	Footer           *Footer           `xml:"Stopka,omitempty"`
	Attachment       *Attachment       `xml:"Zalacznik,omitempty"`
}

// --- Header ---

type Header struct {
	FormCode     FormCode `xml:"KodFormularza"`
	FormVariant  int8     `xml:"WariantFormularza"`
	CreationDate string   `xml:"DataWytworzeniaFa"`
	SystemInfo   string   `xml:"SystemInfo,omitempty"`
}

type FormCode struct {
	Value         string `xml:",chardata"`
	SystemCode    string `xml:"kodSystemowy,attr"`
	SchemaVersion string `xml:"wersjaSchemy,attr"`
}

// --- Address ---

type Address struct {
	CountryCode  string `xml:"KodKraju"`
	AddressLine1 string `xml:"AdresL1"`
	AddressLine2 string `xml:"AdresL2,omitempty"`
	GLN          string `xml:"GLN,omitempty"`
}

// --- Contact ---

type ContactData struct {
	Email string `xml:"Email,omitempty"`
	Phone string `xml:"Telefon,omitempty"`
}

// --- Seller ---

type Seller struct {
	// TaxpayerPrefix is the EU VAT prefix code (e.g. "PL"), required when the seller uses
	// an intra-Community VAT number per art. 97(10)(2)(3) or art. 136(1)(3). Fixed to "PL".
	TaxpayerPrefix     string        `xml:"PrefiksPodatnika,omitempty"`
	EORINumber         string        `xml:"NrEORI,omitempty"`
	IdentData          SellerIdent   `xml:"DaneIdentyfikacyjne"`
	Address            Address       `xml:"Adres"`
	CorrespondenceAddr *Address      `xml:"AdresKoresp,omitempty"`
	ContactData        []ContactData `xml:"DaneKontaktowe,omitempty"`
	TaxpayerStatusInfo *int          `xml:"StatusInfoPodatnika,omitempty"`
}

type SellerIdent struct {
	NIP  string `xml:"NIP"`
	Name string `xml:"Nazwa"`
}

// --- Buyer ---

type Buyer struct {
	EORINumber         string        `xml:"NrEORI,omitempty"`
	IdentData          BuyerIdent    `xml:"DaneIdentyfikacyjne"`
	Address            *Address      `xml:"Adres,omitempty"`
	CorrespondenceAddr *Address      `xml:"AdresKoresp,omitempty"`
	ContactData        []ContactData `xml:"DaneKontaktowe,omitempty"`
	CustomerNumber     string        `xml:"NrKlienta,omitempty"`
	// BuyerID is a unique key linking buyer data across correction invoice lines,
	// used when buyer details changed between the original and correcting invoice.
	BuyerID string `xml:"IDNabywcy,omitempty"`
	// JST flags that the invoice concerns a subordinate unit of a local government body (JST).
	// Value 1 = yes. When set, populate ThirdParties with the unit's NIP or InternalID and role 8.
	JST int `xml:"JST"`
	// GV flags that the invoice concerns a VAT group member.
	// Value 1 = yes (populate ThirdParties with role 10); value 2 = no.
	GV int `xml:"GV"`
}

// BuyerIdent represents the buyer identification.
// The XSD uses a choice between NIP, EUCode+EUVatNumber, CountryCode+IDNumber, or NoID.
// All fields are optional; populate exactly one group.
type BuyerIdent struct {
	NIP         string `xml:"NIP,omitempty"`
	EUCode      string `xml:"KodUE,omitempty"`
	EUVatNumber string `xml:"NrVatUE,omitempty"`
	CountryCode string `xml:"KodKraju,omitempty"`
	IDNumber    string `xml:"NrID,omitempty"`
	// NoID indicates the entity has no tax identifier or it does not appear on the invoice. Set to 1.
	NoID *int   `xml:"BrakID,omitempty"`
	Name string `xml:"Nazwa,omitempty"`
}

// --- ThirdParty ---

type ThirdParty struct {
	// BuyerID is a unique key linking buyer data across correction invoice lines,
	// used when buyer details changed between the original and correcting invoice.
	BuyerID            string          `xml:"IDNabywcy,omitempty"`
	EORINumber         string          `xml:"NrEORI,omitempty"`
	IdentData          ThirdPartyIdent `xml:"DaneIdentyfikacyjne"`
	Address            *Address        `xml:"Adres,omitempty"`
	CorrespondenceAddr *Address        `xml:"AdresKoresp,omitempty"`
	ContactData        []ContactData   `xml:"DaneKontaktowe,omitempty"`
	Role               *int            `xml:"Rola,omitempty"`
	OtherRole          *int            `xml:"RolaInna,omitempty"`
	RoleDescription    string          `xml:"OpisRoli,omitempty"`
	Share              string          `xml:"Udzial,omitempty"`
	CustomerNumber     string          `xml:"NrKlienta,omitempty"`
}

type ThirdPartyIdent struct {
	NIP         string `xml:"NIP,omitempty"`
	// InternalID is an internal identifier derived from the NIP (used for JST subordinate units).
	InternalID  string `xml:"IDWew,omitempty"`
	EUCode      string `xml:"KodUE,omitempty"`
	EUVatNumber string `xml:"NrVatUE,omitempty"`
	CountryCode string `xml:"KodKraju,omitempty"`
	IDNumber    string `xml:"NrID,omitempty"`
	// NoID indicates the entity has no tax identifier or it does not appear on the invoice. Set to 1.
	NoID *int   `xml:"BrakID,omitempty"`
	Name string `xml:"Nazwa,omitempty"`
}

// --- AuthorizedEntity ---

type AuthorizedEntity struct {
	EORINumber         string                    `xml:"NrEORI,omitempty"`
	IdentData          SellerIdent               `xml:"DaneIdentyfikacyjne"`
	Address            Address                   `xml:"Adres"`
	CorrespondenceAddr *Address                  `xml:"AdresKoresp,omitempty"`
	ContactData        []AuthorizedEntityContact `xml:"DaneKontaktowe,omitempty"`
	Role               int                       `xml:"RolaPU"`
}

type AuthorizedEntityContact struct {
	Email string `xml:"EmailPU,omitempty"`
	Phone string `xml:"TelefonPU,omitempty"`
}

// --- InvoiceBody ---

type InvoiceBody struct {
	CurrencyCode string   `xml:"KodWaluty"`
	// P1 is the invoice issue date (art. 106na(1)).
	P1  string   `xml:"P_1"`
	// P1M is the place of invoice issue.
	P1M string   `xml:"P_1M,omitempty"`
	// P2 is the sequential invoice number uniquely identifying the invoice within one or more series.
	P2  string   `xml:"P_2"`
	// WZ are warehouse release document numbers (wydanie na zewnątrz) associated with this invoice.
	WZ  []string `xml:"WZ,omitempty"`

	// P6 is the date of delivery/completion of goods or services, or the date of prepayment receipt
	// (art. 106b(1)(4)), when it differs from the issue date and is the same for all invoice lines.
	// Use Period instead when the invoice covers a date range, or P6A on individual lines when
	// dates differ per line.
	P6     string         `xml:"P_6,omitempty"`
	Period *InvoicePeriod `xml:"OkresFa,omitempty"`

	// VAT rate summary fields (P_13_x = net sales total, P_14_x = VAT amount).
	// For advance invoices: net/VAT amounts of the advance payment.
	// For correction invoices: the difference per art. 106j(2)(5).

	// P131/P141: Standard rate (currently 23% or 22%).
	P131  string `xml:"P_13_1,omitempty"`
	P141  string `xml:"P_14_1,omitempty"`
	// P141W: P141 converted to PLN for foreign-currency invoices (per Chapter VI / art. 106e(11)).
	P141W string `xml:"P_14_1W,omitempty"`

	// P132/P142: First reduced rate (currently 8% or 7%).
	P132  string `xml:"P_13_2,omitempty"`
	P142  string `xml:"P_14_2,omitempty"`
	// P142W: P142 converted to PLN for foreign-currency invoices.
	P142W string `xml:"P_14_2W,omitempty"`

	// P133/P143: Second reduced rate (currently 5%).
	P133  string `xml:"P_13_3,omitempty"`
	P143  string `xml:"P_14_3,omitempty"`
	// P143W: P143 converted to PLN for foreign-currency invoices.
	P143W string `xml:"P_14_3W,omitempty"`

	// P134/P144: Taxi flat-rate scheme (ryczałt dla taksówek osobowych).
	P134  string `xml:"P_13_4,omitempty"`
	P144  string `xml:"P_14_4,omitempty"`
	// P144W: P144 converted to PLN for foreign-currency invoices.
	P144W string `xml:"P_14_4W,omitempty"`

	// P135/P145: OSS special procedure (Chapter XII, Section 6a).
	P135 string `xml:"P_13_5,omitempty"`
	P145 string `xml:"P_14_5,omitempty"`

	// 0% rate buckets (no corresponding P_14 as there is no VAT due):
	// P1361: 0% domestic sales (excluding WDT and exports).
	P1361 string `xml:"P_13_6_1,omitempty"`
	// P1362: 0% intra-Community supply of goods (WDT).
	P1362 string `xml:"P_13_6_2,omitempty"`
	// P1363: 0% exports.
	P1363 string `xml:"P_13_6_3,omitempty"`

	// P137: VAT-exempt sales.
	P137  string `xml:"P_13_7,omitempty"`
	// P138: Sales outside Poland's territory, excluding OSS (P135) and intra-Community services (P139).
	P138  string `xml:"P_13_8,omitempty"`
	// P139: Intra-Community services under art. 100(1)(4) (reported in EC sales list).
	P139  string `xml:"P_13_9,omitempty"`
	// P1310: Reverse-charge sales where the buyer is the taxpayer (art. 17(1)(7)(8)) and other
	// domestic reverse-charge cases.
	P1310 string `xml:"P_13_10,omitempty"`
	// P1311: Margin scheme sales (art. 119 – travel agencies; art. 120 – used goods/art/antiques).
	P1311 string `xml:"P_13_11,omitempty"`

	// P15 is the total amount due. For advance invoices: the payment amount documented by this invoice.
	// For art. 106f(3) invoices: the remaining amount to pay. For corrections: the correction amount.
	P15 string `xml:"P_15"`
	// ExchangeRateZ is the exchange rate used to calculate VAT on advance payment invoices
	// (art. 106b(1)(4)) in cases covered by Chapter VI.
	ExchangeRateZ string `xml:"KursWalutyZ,omitempty"`

	Annotations Annotations `xml:"Adnotacje"`
	// InvoiceType is the invoice kind. Use TRodzajFaktury values:
	// VAT (standard), KOR (correction), ZAL (advance), ROZ (art. 106f(3) settlement),
	// UPR (art. 106e(5)(3) simplified), KOR_ZAL (advance correction), KOR_ROZ (settlement correction).
	InvoiceType string `xml:"RodzajFaktury"`

	// Correction-related fields
	CorrectionReason  string                 `xml:"PrzyczynaKorekty,omitempty"`
	CorrectionType    *int                   `xml:"TypKorekty,omitempty"`
	CorrectedInvoices []CorrectedInvoiceData `xml:"DaneFaKorygowanej,omitempty"`
	CorrectedPeriod   string                 `xml:"OkresFaKorygowanej,omitempty"`
	CorrectedNumber   string                 `xml:"NrFaKorygowany,omitempty"`
	SellerCorrection  *SellerCorrection      `xml:"Podmiot1K,omitempty"`
	BuyerCorrections  []BuyerCorrection      `xml:"Podmiot2K,omitempty"`
	// P15ZK is the total amount due before correction (used on advance/art.106f(3) corrections).
	P15ZK          string                 `xml:"P_15ZK,omitempty"`
	// ExchangeRateZK is the exchange rate used to calculate VAT before the correction.
	ExchangeRateZK string                 `xml:"KursWalutyZK,omitempty"`

	// PartialAdvances holds individual advance payment data when the invoice documents more than
	// one prepayment (art. 106b(1)(4)). Each entry corresponds to one payment installment.
	PartialAdvances []PartialAdvance `xml:"ZaliczkaCzesciowa,omitempty"`
	// FP marks an invoice linked to a fiscal receipt per art. 109(3d). Set to 1.
	FP *int `xml:"FP,omitempty"`
	// TP flags a related-party transaction per JPK_VAT regulations (§10(4)(3)). Set to 1.
	TP             *int             `xml:"TP,omitempty"`
	AdditionalDesc []KeyValue       `xml:"DodatkowyOpis,omitempty"`
	// AdvanceInvoices lists the advance invoices (by number or KSeF number) that this
	// settlement invoice (ROZ) or correction is finalising.
	AdvanceInvoices []AdvanceInvoice `xml:"FakturaZaliczkowa,omitempty"`
	// ExciseTaxReturn is additional information required for farmers claiming a refund of
	// excise tax included in the price of diesel fuel. Set to 1.
	ExciseTaxReturn *int             `xml:"ZwrotAkcyzy,omitempty"`

	Lines            []InvoiceLine          `xml:"FaWiersz,omitempty"`
	Settlement       *Settlement            `xml:"Rozliczenie,omitempty"`
	Payment          *Payment               `xml:"Platnosc,omitempty"`
	TransactionConds *TransactionConditions `xml:"WarunkiTransakcji,omitempty"`
	Order            *Order                 `xml:"Zamowienie,omitempty"`
}

type InvoicePeriod struct {
	// From is the start date of the billing period.
	From string `xml:"P_6_Od"`
	// To is the end date of the billing period (date of delivery/service completion).
	To   string `xml:"P_6_Do"`
}

// --- Annotations ---

// Annotations (Adnotacje) are mandatory boolean/choice flags required on every invoice.
// Most fields use a 1/2 choice: 1 = condition applies, 2 = condition does not apply.
type Annotations struct {
	// P16: Cash-method indicator. Set 1 if the tax point arises under art. 19a(5)(1) or art. 21(1)
	// (the invoice must then carry the words "metoda kasowa"); otherwise set 2.
	P16 int `xml:"P_16"`
	// P17: Self-billing indicator. Set 1 if the invoice is issued by the buyer per art. 106d(1)
	// (the invoice must carry "samofakturowanie"); otherwise set 2.
	P17 int `xml:"P_17"`
	// P18: Reverse-charge indicator. Set 1 if the buyer is liable for VAT on this supply
	// (the invoice must carry "odwrotne obciążenie"); otherwise set 2.
	P18 int `xml:"P_18"`
	// P18A: Split-payment (MPP) indicator. Set 1 if the invoice total exceeds 15,000 PLN and
	// covers goods/services listed in Annex 15 to the VAT Act (the invoice must carry
	// "mechanizm podzielonej płatności"); otherwise set 2.
	P18A         int               `xml:"P_18A"`
	Exemption    Exemption         `xml:"Zwolnienie"`
	NewTransport NewTransportMeans `xml:"NoweSrodkiTransportu"`
	// P23: Simplified triangular trade indicator. Set 1 if this invoice is issued by the second
	// party in a triangular chain per art. 135(1)(4)(b)(c) and carries the annotations required
	// by art. 136(1); otherwise set 2.
	P23    int    `xml:"P_23"`
	Margin Margin `xml:"PMarzy"`
}

// Exemption covers the VAT-exemption annotation block. Exactly one of P19/P19N must be set.
type Exemption struct {
	// P19: Set to 1 if any line is exempt from VAT under art. 43(1), art. 113(1)(9),
	// art. 82(3), or other provisions. When set, populate at least one of P19A/P19B/P19C.
	P19  *int   `xml:"P_19,omitempty"`
	// P19A: When P19=1, the national statutory provision granting the exemption.
	P19A string `xml:"P_19A,omitempty"`
	// P19B: When P19=1, the EU Directive 2006/112/EC provision granting the exemption.
	P19B string `xml:"P_19B,omitempty"`
	// P19C: When P19=1, any other legal basis for the exemption.
	P19C string `xml:"P_19C,omitempty"`
	// P19N: Set to 1 when no VAT-exempt supply is present (mutually exclusive with P19).
	P19N *int   `xml:"P_19N,omitempty"`
}

// NewTransportMeans covers the intra-Community supply of new means of transport block.
// Exactly one of P22/P22N must be set.
type NewTransportMeans struct {
	// P22: Set to 1 if the invoice covers an intra-Community supply of a new means of transport.
	P22   *int                    `xml:"P_22,omitempty"`
	// P42_5: Set 1 if the obligation under art. 42(5) applies (0% WDT documentation), 2 otherwise.
	P42_5 *int                    `xml:"P_42_5,omitempty"`
	Items []NewTransportMeansItem `xml:"NowySrodekTransportu,omitempty"`
	// P22N: Set to 1 when no intra-Community supply of a new means of transport is present.
	P22N  *int                    `xml:"P_22N,omitempty"`
}

type NewTransportMeansItem struct {
	// P22A is the type of new means of transport (land vehicle, watercraft, or aircraft).
	P22A string `xml:"P_22A"`
	// LineNumber is the invoice line number (FaWiersz) that documents this transport item.
	LineNumber int    `xml:"P_NrWierszaNST"`
	P22BMK     string `xml:"P_22BMK,omitempty"`
	P22BMD     string `xml:"P_22BMD,omitempty"`
	P22BK      string `xml:"P_22BK,omitempty"`
	P22BNR     string `xml:"P_22BNR,omitempty"`
	P22BRP     string `xml:"P_22BRP,omitempty"`
	// Land vehicle
	P22B  string `xml:"P_22B,omitempty"`
	P22B1 string `xml:"P_22B1,omitempty"`
	P22B2 string `xml:"P_22B2,omitempty"`
	P22B3 string `xml:"P_22B3,omitempty"`
	P22B4 string `xml:"P_22B4,omitempty"`
	P22BT string `xml:"P_22BT,omitempty"`
	// Watercraft
	P22C  string `xml:"P_22C,omitempty"`
	P22C1 string `xml:"P_22C1,omitempty"`
	// Aircraft
	P22D  string `xml:"P_22D,omitempty"`
	P22D1 string `xml:"P_22D1,omitempty"`
}

// Margin covers the VAT margin scheme annotation block.
// Exactly one of the presence flags (PPMarzy or PPMarzyN) must be set.
type Margin struct {
	// PPMarzy: Set to 1 if any margin procedure (art. 119 or 120) applies.
	// When set, also set exactly one of PPMarzy2/PPMarzy31/PPMarzy32/PPMarzy33.
	PPMarzy   *int `xml:"P_PMarzy,omitempty"`
	// PPMarzy2: Travel agency margin scheme (art. 119) – invoice must carry "procedura marży dla biur podróży".
	PPMarzy2  *int `xml:"P_PMarzy_2,omitempty"`
	// PPMarzy31: Used goods margin scheme (art. 120) – invoice must carry "procedura marży - towary używane".
	PPMarzy31 *int `xml:"P_PMarzy_3_1,omitempty"`
	// PPMarzy32: Works of art margin scheme (art. 120) – invoice must carry "procedura marży - dzieła sztuki".
	PPMarzy32 *int `xml:"P_PMarzy_3_2,omitempty"`
	// PPMarzy33: Collectibles/antiques margin scheme (art. 120) – invoice must carry "procedura marży - przedmioty kolekcjonerskie i antyki".
	PPMarzy33 *int `xml:"P_PMarzy_3_3,omitempty"`
	// PPMarzyN: Set to 1 when no margin procedure applies (mutually exclusive with PPMarzy).
	PPMarzyN  *int `xml:"P_PMarzyN,omitempty"`
}

// --- Invoice line items ---

type InvoiceLine struct {
	LineNumber int    `xml:"NrWierszaFa"`
	UUID       string `xml:"UU_ID,omitempty"`
	// P6A is the per-line delivery/service date or prepayment date (art. 106b(1)(4)).
	// Use this when individual lines have different dates; use InvoiceBody.P6 when all lines share one date.
	P6A           string `xml:"P_6A,omitempty"`
	// P7 is the name/description of the goods or service.
	P7            string `xml:"P_7,omitempty"`
	Index         string `xml:"Indeks,omitempty"`
	GTIN          string `xml:"GTIN,omitempty"`
	PKWiU         string `xml:"PKWiU,omitempty"`
	CN            string `xml:"CN,omitempty"`
	PKOB          string `xml:"PKOB,omitempty"`
	// P8A is the unit of measure (e.g. "szt.", "kg", "godz.").
	P8A string `xml:"P_8A,omitempty"`
	// P8B is the quantity of goods delivered or scope of services performed.
	P8B string `xml:"P_8B,omitempty"`
	// P9A is the unit net price (excluding VAT). Used in the standard case.
	P9A string `xml:"P_9A,omitempty"`
	// P9B is the unit gross price (including VAT), used when art. 106e(7)(8) applies.
	P9B string `xml:"P_9B,omitempty"`
	// P10 is the total discount or price reduction not already reflected in the unit price.
	P10 string `xml:"P_10,omitempty"`
	// P11 is the net value of goods/services on this line (excluding VAT).
	P11    string `xml:"P_11,omitempty"`
	// P11A is the gross value of goods/services on this line, used when art. 106e(7)(8) applies.
	P11A   string `xml:"P_11A,omitempty"`
	// P11Vat is the VAT amount for the case described in art. 106e(10).
	P11Vat string `xml:"P_11Vat,omitempty"`
	// P12 is the VAT rate code. Valid values: 23, 22, 8, 7, 5, 4, 3, "0 KR", "0 WDT", "0 EX", "zw", "oo", "np I", "np II".
	P12    string `xml:"P_12,omitempty"`
	// P12XII is the VAT rate under the OSS special procedure (Chapter XII, Section 6a).
	P12XII   string `xml:"P_12_XII,omitempty"`
	// P12Zal15: Set to 1 if this line item is listed in Annex 15 to the VAT Act
	// (mandatory split-payment goods/services).
	P12Zal15      *int   `xml:"P_12_Zal_15,omitempty"`
	ExciseTax     string `xml:"KwotaAkcyzy,omitempty"`
	// GTU is the goods and services classification code (GTU_01–GTU_13) used in JPK_VAT reporting.
	GTU           string `xml:"GTU,omitempty"`
	Procedure     string `xml:"Procedura,omitempty"`
	// ExchangeRate is the exchange rate applied to calculate VAT for this line (Chapter VI cases).
	ExchangeRate  string `xml:"KursWaluty,omitempty"`
	// PreviousState: Set to 1 on correction invoices when this line shows the "before" state,
	// used in the before/after correction method where both states appear as separate lines.
	PreviousState *int   `xml:"StanPrzed,omitempty"`
}

// --- Correction data ---

type CorrectedInvoiceData struct {
	IssueDate      string `xml:"DataWystFaKorygowanej"`
	Number         string `xml:"NrFaKorygowanej"`
	KSeFNumber     *int   `xml:"NrKSeF,omitempty"`
	KSeFCorrNumber string `xml:"NrKSeFFaKorygowanej,omitempty"`
	NoKSeFNumber   *int   `xml:"NrKSeFN,omitempty"`
}

type SellerCorrection struct {
	// TaxpayerPrefix is the EU VAT prefix code for the seller before correction.
	TaxpayerPrefix string      `xml:"PrefiksPodatnika,omitempty"`
	IdentData      SellerIdent `xml:"DaneIdentyfikacyjne"`
	Address        Address     `xml:"Adres"`
}

type BuyerCorrection struct {
	IdentData BuyerIdent `xml:"DaneIdentyfikacyjne"`
	Address   *Address   `xml:"Adres,omitempty"`
	// BuyerID is a unique key linking buyer data across correction invoice lines.
	BuyerID   string     `xml:"IDNabywcy,omitempty"`
}

// --- Advance payments ---

// PartialAdvance holds data for a single advance payment installment when an invoice documents
// multiple prepayments (art. 106b(1)(4)). Use multiple entries in InvoiceBody.PartialAdvances.
type PartialAdvance struct {
	// P6Z is the date the advance payment was received.
	P6Z          string `xml:"P_6Z"`
	// P15Z is the payment amount that makes up part of the total in InvoiceBody.P15.
	P15Z         string `xml:"P_15Z"`
	// ExchangeRate is the rate used to calculate VAT on this advance (Chapter VI cases).
	ExchangeRate string `xml:"KursWalutyZW,omitempty"`
}

type AdvanceInvoice struct {
	NoKSeFNumber *int   `xml:"NrKSeFZN,omitempty"`
	Number       string `xml:"NrFaZaliczkowej,omitempty"`
	KSeFNumber   string `xml:"NrKSeFFaZaliczkowej,omitempty"`
}

// --- Key-Value ---

type KeyValue struct {
	LineNumber *int   `xml:"NrWiersza,omitempty"`
	Key        string `xml:"Klucz"`
	Value      string `xml:"Wartosc"`
}

// --- Settlement ---

type Settlement struct {
	Charges         []ChargeDeduction `xml:"Obciazenia,omitempty"`
	TotalCharges    string            `xml:"SumaObciazen,omitempty"`
	Deductions      []ChargeDeduction `xml:"Odliczenia,omitempty"`
	TotalDeductions string            `xml:"SumaOdliczen,omitempty"`
	AmountDue       string            `xml:"DoZaplaty,omitempty"`
	AmountToSettle  string            `xml:"DoRozliczenia,omitempty"`
}

type ChargeDeduction struct {
	Amount string `xml:"Kwota"`
	Reason string `xml:"Powod"`
}

// --- Payment ---

type Payment struct {
	// Full payment
	Paid     *int   `xml:"Zaplacono,omitempty"`
	PaidDate string `xml:"DataZaplaty,omitempty"`

	// Partial payment
	PartialPaymentFlag *int             `xml:"ZnacznikZaplatyCzesciowej,omitempty"`
	PartialPayments    []PartialPayment `xml:"ZaplataCzesciowa,omitempty"`

	PaymentTerms       []PaymentTerm `xml:"TerminPlatnosci,omitempty"`
	PaymentForm        *int          `xml:"FormaPlatnosci,omitempty"`
	OtherPayment       *int          `xml:"PlatnoscInna,omitempty"`
	PaymentDescription string        `xml:"OpisPlatnosci,omitempty"`
	BankAccounts       []BankAccount `xml:"RachunekBankowy,omitempty"`
	FactorBankAccounts []BankAccount `xml:"RachunekBankowyFaktora,omitempty"`
	Discount           *Discount     `xml:"Skonto,omitempty"`
	PaymentLink        string        `xml:"LinkDoPlatnosci,omitempty"`
	KSeFPaymentID      string        `xml:"IPKSeF,omitempty"`
}

type PartialPayment struct {
	Amount             string `xml:"KwotaZaplatyCzesciowej"`
	Date               string `xml:"DataZaplatyCzesciowej"`
	PaymentForm        *int   `xml:"FormaPlatnosci,omitempty"`
	OtherPayment       *int   `xml:"PlatnoscInna,omitempty"`
	PaymentDescription string `xml:"OpisPlatnosci,omitempty"`
}

type PaymentTerm struct {
	Deadline    string           `xml:"Termin,omitempty"`
	Description *PaymentTermDesc `xml:"TerminOpis,omitempty"`
}

type PaymentTermDesc struct {
	Quantity   int    `xml:"Ilosc"`
	Unit       string `xml:"Jednostka"`
	StartEvent string `xml:"ZdarzeniePoczatkowe"`
}

type Discount struct {
	Conditions string `xml:"WarunkiSkonta"`
	Amount     string `xml:"WysokoscSkonta"`
}

type BankAccount struct {
	AccountNumber  string `xml:"NrRB"`
	SWIFT          string `xml:"SWIFT,omitempty"`
	OwnBankAccount *int   `xml:"RachunekWlasnyBanku,omitempty"`
	BankName       string `xml:"NazwaBanku,omitempty"`
	AccountDesc    string `xml:"OpisRachunku,omitempty"`
}

// --- Transaction conditions ---

type TransactionConditions struct {
	Contracts          []Contract  `xml:"Umowy,omitempty"`
	Orders             []OrderRef  `xml:"Zamowienia,omitempty"`
	GoodsLotNumbers    []string    `xml:"NrPartiiTowaru,omitempty"`
	DeliveryConditions string      `xml:"WarunkiDostawy,omitempty"`
	ContractRate       string      `xml:"KursUmowny,omitempty"`
	ContractCurrency   string      `xml:"WalutaUmowna,omitempty"`
	Transport          []Transport `xml:"Transport,omitempty"`
	IntermediaryEntity *int        `xml:"PodmiotPosredniczacy,omitempty"`
}

type Contract struct {
	Date   string `xml:"DataUmowy,omitempty"`
	Number string `xml:"NrUmowy,omitempty"`
}

type OrderRef struct {
	Date   string `xml:"DataZamowienia,omitempty"`
	Number string `xml:"NrZamowienia,omitempty"`
}

type Transport struct {
	Type                   *int      `xml:"RodzajTransportu,omitempty"`
	OtherTransport         *int      `xml:"TransportInny,omitempty"`
	OtherTransportDesc     string    `xml:"OpisInnegoTransportu,omitempty"`
	Carrier                *Carrier  `xml:"Przewoznik,omitempty"`
	TransportOrderNumber   string    `xml:"NrZleceniaTransportu,omitempty"`
	CargoType              *int      `xml:"OpisLadunku,omitempty"`
	OtherCargo             *int      `xml:"LadunekInny,omitempty"`
	OtherCargoDesc         string    `xml:"OpisInnegoLadunku,omitempty"`
	PackagingUnit          string    `xml:"JednostkaOpakowania,omitempty"`
	TransportStartDateTime string    `xml:"DataGodzRozpTransportu,omitempty"`
	TransportEndDateTime   string    `xml:"DataGodzZakTransportu,omitempty"`
	ShipFrom               *Address  `xml:"WysylkaZ,omitempty"`
	ShipThrough            []Address `xml:"WysylkaPrzez,omitempty"`
	ShipTo                 *Address  `xml:"WysylkaDo,omitempty"`
}

type Carrier struct {
	IdentData      BuyerIdent `xml:"DaneIdentyfikacyjne"`
	CarrierAddress Address    `xml:"AdresPrzewoznika"`
}

// --- Order (for advance invoices) ---

type Order struct {
	TotalValue string      `xml:"WartoscZamowienia"`
	Lines      []OrderLine `xml:"ZamowienieWiersz"`
}

type OrderLine struct {
	LineNumber     int    `xml:"NrWierszaZam"`
	UUIDZ          string `xml:"UU_IDZ,omitempty"`
	P7Z            string `xml:"P_7Z,omitempty"`
	IndexZ         string `xml:"IndeksZ,omitempty"`
	GTINZ          string `xml:"GTINZ,omitempty"`
	PKWiUZ         string `xml:"PKWiUZ,omitempty"`
	CNZ            string `xml:"CNZ,omitempty"`
	PKOBZ          string `xml:"PKOBZ,omitempty"`
	P8AZ           string `xml:"P_8AZ,omitempty"`
	P8BZ           string `xml:"P_8BZ,omitempty"`
	P9AZ           string `xml:"P_9AZ,omitempty"`
	P11NettoZ      string `xml:"P_11NettoZ,omitempty"`
	P11VatZ        string `xml:"P_11VatZ,omitempty"`
	P12Z           string `xml:"P_12Z,omitempty"`
	P12ZXII        string `xml:"P_12Z_XII,omitempty"`
	P12ZZal15      *int   `xml:"P_12Z_Zal_15,omitempty"`
	GTUZ           string `xml:"GTUZ,omitempty"`
	ProcedureZ     string `xml:"ProceduraZ,omitempty"`
	ExciseTaxZ     string `xml:"KwotaAkcyzyZ,omitempty"`
	PreviousStateZ *int   `xml:"StanPrzedZ,omitempty"`
}

// --- Footer ---

type Footer struct {
	Info       []FooterInfo `xml:"Informacje,omitempty"`
	Registries []Registry   `xml:"Rejestry,omitempty"`
}

type FooterInfo struct {
	Text string `xml:"StopkaFaktury,omitempty"`
}

type Registry struct {
	FullName string `xml:"PelnaNazwa,omitempty"`
	KRS      string `xml:"KRS,omitempty"`
	REGON    string `xml:"REGON,omitempty"`
	BDO      string `xml:"BDO,omitempty"`
}

// --- Attachment ---

type Attachment struct {
	DataBlocks []DataBlock `xml:"BlokDanych"`
}

type DataBlock struct {
	Header   string          `xml:"ZNaglowek,omitempty"`
	Metadata []Metadata      `xml:"MetaDane"`
	Text     *AttachmentText `xml:"Tekst,omitempty"`
	Tables   []Table         `xml:"Tabela,omitempty"`
}

type Metadata struct {
	Key   string `xml:"ZKlucz"`
	Value string `xml:"ZWartosc"`
}

type AttachmentText struct {
	Paragraphs []string `xml:"Akapit"`
}

type Table struct {
	Metadata []TableMetadata `xml:"TMetaDane,omitempty"`
	Desc     string          `xml:"Opis,omitempty"`
	Header   TableHeader     `xml:"TNaglowek"`
	Rows     []TableRow      `xml:"Wiersz"`
	Sum      *TableSum       `xml:"Suma,omitempty"`
}

type TableMetadata struct {
	Key   string `xml:"TKlucz"`
	Value string `xml:"TWartosc"`
}

type TableHeader struct {
	Columns []TableColumn `xml:"Kol"`
}

type TableColumn struct {
	Type string `xml:"Typ,attr"`
	Name string `xml:"NKom"`
}

type TableRow struct {
	Cells []string `xml:"WKom"`
}

type TableSum struct {
	Cells []string `xml:"SKom"`
}
