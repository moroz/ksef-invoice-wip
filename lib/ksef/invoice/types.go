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
	BuyerID            string        `xml:"IDNabywcy,omitempty"`
	JST                int           `xml:"JST"`
	GV                 int           `xml:"GV"`
}

// BuyerIdent represents the buyer identification.
// The XSD uses a choice between NIP, KodUE+NrVatUE, KodKraju+NrID, or BrakID.
// All fields are optional; populate exactly one group.
type BuyerIdent struct {
	NIP         string `xml:"NIP,omitempty"`
	EUCode      string `xml:"KodUE,omitempty"`
	EUVatNumber string `xml:"NrVatUE,omitempty"`
	CountryCode string `xml:"KodKraju,omitempty"`
	IDNumber    string `xml:"NrID,omitempty"`
	NoID        *int   `xml:"BrakID,omitempty"`
	Name        string `xml:"Nazwa,omitempty"`
}

// --- ThirdParty ---

type ThirdParty struct {
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
	InternalID  string `xml:"IDWew,omitempty"`
	EUCode      string `xml:"KodUE,omitempty"`
	EUVatNumber string `xml:"NrVatUE,omitempty"`
	CountryCode string `xml:"KodKraju,omitempty"`
	IDNumber    string `xml:"NrID,omitempty"`
	NoID        *int   `xml:"BrakID,omitempty"`
	Name        string `xml:"Nazwa,omitempty"`
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
	P1           string   `xml:"P_1"`            // Data wystawienia
	P1M          string   `xml:"P_1M,omitempty"` // Miejsce wystawienia
	P2           string   `xml:"P_2"`            // Numer faktury
	WZ           []string `xml:"WZ,omitempty"`

	// Date of delivery/service (choice: single date or period)
	P6     string         `xml:"P_6,omitempty"`
	Period *InvoicePeriod `xml:"OkresFa,omitempty"`

	// Tax rate summaries: 23%/22%
	P131  string `xml:"P_13_1,omitempty"`
	P141  string `xml:"P_14_1,omitempty"`
	P141W string `xml:"P_14_1W,omitempty"`

	// 8%/7%
	P132  string `xml:"P_13_2,omitempty"`
	P142  string `xml:"P_14_2,omitempty"`
	P142W string `xml:"P_14_2W,omitempty"`

	// 5%
	P133  string `xml:"P_13_3,omitempty"`
	P143  string `xml:"P_14_3,omitempty"`
	P143W string `xml:"P_14_3W,omitempty"`

	// Taxi ryczalt
	P134  string `xml:"P_13_4,omitempty"`
	P144  string `xml:"P_14_4,omitempty"`
	P144W string `xml:"P_14_4W,omitempty"`

	// OSS (dział XII rozdział 6a)
	P135 string `xml:"P_13_5,omitempty"`
	P145 string `xml:"P_14_5,omitempty"`

	// 0% rates
	P1361 string `xml:"P_13_6_1,omitempty"` // 0% krajowe
	P1362 string `xml:"P_13_6_2,omitempty"` // 0% WDT
	P1363 string `xml:"P_13_6_3,omitempty"` // 0% eksport

	P137  string `xml:"P_13_7,omitempty"`  // zwolnione
	P138  string `xml:"P_13_8,omitempty"`  // poza terytorium
	P139  string `xml:"P_13_9,omitempty"`  // art.100 ust.1 pkt 4
	P1310 string `xml:"P_13_10,omitempty"` // odwrotne obciążenie
	P1311 string `xml:"P_13_11,omitempty"` // marża

	P15           string `xml:"P_15"` // Kwota należności ogółem
	ExchangeRateZ string `xml:"KursWalutyZ,omitempty"`

	Annotations Annotations `xml:"Adnotacje"`
	InvoiceType string      `xml:"RodzajFaktury"`

	// Correction-related fields
	CorrectionReason  string                 `xml:"PrzyczynaKorekty,omitempty"`
	CorrectionType    *int                   `xml:"TypKorekty,omitempty"`
	CorrectedInvoices []CorrectedInvoiceData `xml:"DaneFaKorygowanej,omitempty"`
	CorrectedPeriod   string                 `xml:"OkresFaKorygowanej,omitempty"`
	CorrectedNumber   string                 `xml:"NrFaKorygowany,omitempty"`
	SellerCorrection  *SellerCorrection      `xml:"Podmiot1K,omitempty"`
	BuyerCorrections  []BuyerCorrection      `xml:"Podmiot2K,omitempty"`
	P15ZK             string                 `xml:"P_15ZK,omitempty"`
	ExchangeRateZK    string                 `xml:"KursWalutyZK,omitempty"`

	PartialAdvances []PartialAdvance `xml:"ZaliczkaCzesciowa,omitempty"`
	FP              *int             `xml:"FP,omitempty"`
	TP              *int             `xml:"TP,omitempty"`
	AdditionalDesc  []KeyValue       `xml:"DodatkowyOpis,omitempty"`
	AdvanceInvoices []AdvanceInvoice `xml:"FakturaZaliczkowa,omitempty"`
	ExciseTaxReturn *int             `xml:"ZwrotAkcyzy,omitempty"`

	Lines          []InvoiceLine          `xml:"FaWiersz,omitempty"`
	Settlement     *Settlement            `xml:"Rozliczenie,omitempty"`
	Payment        *Payment               `xml:"Platnosc,omitempty"`
	TransactionConds *TransactionConditions `xml:"WarunkiTransakcji,omitempty"`
	Order          *Order                 `xml:"Zamowienie,omitempty"`
}

type InvoicePeriod struct {
	From string `xml:"P_6_Od"`
	To   string `xml:"P_6_Do"`
}

// --- Annotations ---

type Annotations struct {
	P16          int               `xml:"P_16"`
	P17          int               `xml:"P_17"`
	P18          int               `xml:"P_18"`
	P18A         int               `xml:"P_18A"`
	Exemption    Exemption         `xml:"Zwolnienie"`
	NewTransport NewTransportMeans `xml:"NoweSrodkiTransportu"`
	P23          int               `xml:"P_23"`
	Margin       Margin            `xml:"PMarzy"`
}

type Exemption struct {
	P19  *int   `xml:"P_19,omitempty"`
	P19A string `xml:"P_19A,omitempty"`
	P19B string `xml:"P_19B,omitempty"`
	P19C string `xml:"P_19C,omitempty"`
	P19N *int   `xml:"P_19N,omitempty"`
}

type NewTransportMeans struct {
	P22   *int                    `xml:"P_22,omitempty"`
	P42_5 *int                    `xml:"P_42_5,omitempty"`
	Items []NewTransportMeansItem `xml:"NowySrodekTransportu,omitempty"`
	P22N  *int                    `xml:"P_22N,omitempty"`
}

type NewTransportMeansItem struct {
	P22A       string `xml:"P_22A"`
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

type Margin struct {
	PPMarzy   *int `xml:"P_PMarzy,omitempty"`
	PPMarzy2  *int `xml:"P_PMarzy_2,omitempty"`
	PPMarzy31 *int `xml:"P_PMarzy_3_1,omitempty"`
	PPMarzy32 *int `xml:"P_PMarzy_3_2,omitempty"`
	PPMarzy33 *int `xml:"P_PMarzy_3_3,omitempty"`
	PPMarzyN  *int `xml:"P_PMarzyN,omitempty"`
}

// --- Invoice line items ---

type InvoiceLine struct {
	LineNumber    int    `xml:"NrWierszaFa"`
	UUID          string `xml:"UU_ID,omitempty"`
	P6A           string `xml:"P_6A,omitempty"`
	P7            string `xml:"P_7,omitempty"`
	Index         string `xml:"Indeks,omitempty"`
	GTIN          string `xml:"GTIN,omitempty"`
	PKWiU         string `xml:"PKWiU,omitempty"`
	CN            string `xml:"CN,omitempty"`
	PKOB          string `xml:"PKOB,omitempty"`
	P8A           string `xml:"P_8A,omitempty"`  // unit of measure
	P8B           string `xml:"P_8B,omitempty"`  // quantity
	P9A           string `xml:"P_9A,omitempty"`  // unit price net
	P9B           string `xml:"P_9B,omitempty"`  // unit price gross
	P10           string `xml:"P_10,omitempty"`  // discount
	P11           string `xml:"P_11,omitempty"`  // net value
	P11A          string `xml:"P_11A,omitempty"` // gross value
	P11Vat        string `xml:"P_11Vat,omitempty"`
	P12           string `xml:"P_12,omitempty"` // tax rate
	P12XII        string `xml:"P_12_XII,omitempty"`
	P12Zal15      *int   `xml:"P_12_Zal_15,omitempty"`
	ExciseTax     string `xml:"KwotaAkcyzy,omitempty"`
	GTU           string `xml:"GTU,omitempty"`
	Procedure     string `xml:"Procedura,omitempty"`
	ExchangeRate  string `xml:"KursWaluty,omitempty"`
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
	TaxpayerPrefix string      `xml:"PrefiksPodatnika,omitempty"`
	IdentData      SellerIdent `xml:"DaneIdentyfikacyjne"`
	Address        Address     `xml:"Adres"`
}

type BuyerCorrection struct {
	IdentData BuyerIdent `xml:"DaneIdentyfikacyjne"`
	Address   *Address   `xml:"Adres,omitempty"`
	BuyerID   string     `xml:"IDNabywcy,omitempty"`
}

// --- Advance payments ---

type PartialAdvance struct {
	P6Z          string `xml:"P_6Z"`
	P15Z         string `xml:"P_15Z"`
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
	Type                   *int     `xml:"RodzajTransportu,omitempty"`
	OtherTransport         *int     `xml:"TransportInny,omitempty"`
	OtherTransportDesc     string   `xml:"OpisInnegoTransportu,omitempty"`
	Carrier                *Carrier `xml:"Przewoznik,omitempty"`
	TransportOrderNumber   string   `xml:"NrZleceniaTransportu,omitempty"`
	CargoType              *int     `xml:"OpisLadunku,omitempty"`
	OtherCargo             *int     `xml:"LadunekInny,omitempty"`
	OtherCargoDesc         string   `xml:"OpisInnegoLadunku,omitempty"`
	PackagingUnit          string   `xml:"JednostkaOpakowania,omitempty"`
	TransportStartDateTime string   `xml:"DataGodzRozpTransportu,omitempty"`
	TransportEndDateTime   string   `xml:"DataGodzZakTransportu,omitempty"`
	ShipFrom               *Address `xml:"WysylkaZ,omitempty"`
	ShipThrough            []Address `xml:"WysylkaPrzez,omitempty"`
	ShipTo                 *Address `xml:"WysylkaDo,omitempty"`
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
