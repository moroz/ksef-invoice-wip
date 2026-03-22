package invoice

import "encoding/xml"

const Namespace = "http://crd.gov.pl/wzor/2025/06/25/13775/"

// Faktura is the root element of a KSeF FA(3) invoice.
type Faktura struct {
	XMLName            xml.Name            `xml:"http://crd.gov.pl/wzor/2025/06/25/13775/ Faktura"`
	Naglowek           Naglowek            `xml:"Naglowek"`
	Podmiot1           Podmiot1            `xml:"Podmiot1"`
	Podmiot2           Podmiot2            `xml:"Podmiot2"`
	Podmiot3           []Podmiot3Element   `xml:"Podmiot3,omitempty"`
	PodmiotUpowazniony *PodmiotUpowazniony `xml:"PodmiotUpowazniony,omitempty"`
	Fa                 Fa                  `xml:"Fa"`
	Stopka             *Stopka             `xml:"Stopka,omitempty"`
	Zalacznik          *Zalacznik          `xml:"Zalacznik,omitempty"`
}

// --- Header ---

type Naglowek struct {
	KodFormularza     KodFormularza `xml:"KodFormularza"`
	WariantFormularza int8          `xml:"WariantFormularza"`
	DataWytworzeniaFa string        `xml:"DataWytworzeniaFa"`
	SystemInfo        string        `xml:"SystemInfo,omitempty"`
}

type KodFormularza struct {
	Value        string `xml:",chardata"`
	KodSystemowy string `xml:"kodSystemowy,attr"`
	WersjaSchemy string `xml:"wersjaSchemy,attr"`
}

// --- Address ---

type Adres struct {
	KodKraju string `xml:"KodKraju"`
	AdresL1  string `xml:"AdresL1"`
	AdresL2  string `xml:"AdresL2,omitempty"`
	GLN      string `xml:"GLN,omitempty"`
}

// --- Contact ---

type DaneKontaktowe struct {
	Email   string `xml:"Email,omitempty"`
	Telefon string `xml:"Telefon,omitempty"`
}

// --- Podmiot1 (Seller) ---

type Podmiot1 struct {
	PrefiksPodatnika    string           `xml:"PrefiksPodatnika,omitempty"`
	NrEORI              string           `xml:"NrEORI,omitempty"`
	DaneIdentyfikacyjne Podmiot1Ident    `xml:"DaneIdentyfikacyjne"`
	Adres               Adres            `xml:"Adres"`
	AdresKoresp         *Adres           `xml:"AdresKoresp,omitempty"`
	DaneKontaktowe      []DaneKontaktowe `xml:"DaneKontaktowe,omitempty"`
	StatusInfoPodatnika *int             `xml:"StatusInfoPodatnika,omitempty"`
}

type Podmiot1Ident struct {
	NIP   string `xml:"NIP"`
	Nazwa string `xml:"Nazwa"`
}

// --- Podmiot2 (Buyer) ---

type Podmiot2 struct {
	NrEORI              string           `xml:"NrEORI,omitempty"`
	DaneIdentyfikacyjne Podmiot2Ident    `xml:"DaneIdentyfikacyjne"`
	Adres               *Adres           `xml:"Adres,omitempty"`
	AdresKoresp         *Adres           `xml:"AdresKoresp,omitempty"`
	DaneKontaktowe      []DaneKontaktowe `xml:"DaneKontaktowe,omitempty"`
	NrKlienta           string           `xml:"NrKlienta,omitempty"`
	IDNabywcy           string           `xml:"IDNabywcy,omitempty"`
	JST                 int              `xml:"JST"`
	GV                  int              `xml:"GV"`
}

// Podmiot2Ident represents the buyer identification.
// The XSD uses a choice between NIP, KodUE+NrVatUE, KodKraju+NrID, or BrakID.
// All fields are optional; populate exactly one group.
type Podmiot2Ident struct {
	NIP      string `xml:"NIP,omitempty"`
	KodUE    string `xml:"KodUE,omitempty"`
	NrVatUE  string `xml:"NrVatUE,omitempty"`
	KodKraju string `xml:"KodKraju,omitempty"`
	NrID     string `xml:"NrID,omitempty"`
	BrakID   *int   `xml:"BrakID,omitempty"`
	Nazwa    string `xml:"Nazwa,omitempty"`
}

// --- Podmiot3 (Third parties) ---

type Podmiot3Element struct {
	IDNabywcy           string           `xml:"IDNabywcy,omitempty"`
	NrEORI              string           `xml:"NrEORI,omitempty"`
	DaneIdentyfikacyjne Podmiot3Ident    `xml:"DaneIdentyfikacyjne"`
	Adres               *Adres           `xml:"Adres,omitempty"`
	AdresKoresp         *Adres           `xml:"AdresKoresp,omitempty"`
	DaneKontaktowe      []DaneKontaktowe `xml:"DaneKontaktowe,omitempty"`
	Rola                *int             `xml:"Rola,omitempty"`
	RolaInna            *int             `xml:"RolaInna,omitempty"`
	OpisRoli            string           `xml:"OpisRoli,omitempty"`
	Udzial              string           `xml:"Udzial,omitempty"`
	NrKlienta           string           `xml:"NrKlienta,omitempty"`
}

type Podmiot3Ident struct {
	NIP      string `xml:"NIP,omitempty"`
	IDWew    string `xml:"IDWew,omitempty"`
	KodUE    string `xml:"KodUE,omitempty"`
	NrVatUE  string `xml:"NrVatUE,omitempty"`
	KodKraju string `xml:"KodKraju,omitempty"`
	NrID     string `xml:"NrID,omitempty"`
	BrakID   *int   `xml:"BrakID,omitempty"`
	Nazwa    string `xml:"Nazwa,omitempty"`
}

// --- PodmiotUpowazniony ---

type PodmiotUpowazniony struct {
	NrEORI              string             `xml:"NrEORI,omitempty"`
	DaneIdentyfikacyjne Podmiot1Ident      `xml:"DaneIdentyfikacyjne"`
	Adres               Adres              `xml:"Adres"`
	AdresKoresp         *Adres             `xml:"AdresKoresp,omitempty"`
	DaneKontaktowe      []DaneKontaktowePU `xml:"DaneKontaktowe,omitempty"`
	RolaPU              int                `xml:"RolaPU"`
}

type DaneKontaktowePU struct {
	EmailPU   string `xml:"EmailPU,omitempty"`
	TelefonPU string `xml:"TelefonPU,omitempty"`
}

// --- Fa (Invoice body) ---

type Fa struct {
	KodWaluty string   `xml:"KodWaluty"`
	P1        string   `xml:"P_1"`            // Data wystawienia
	P1M       string   `xml:"P_1M,omitempty"` // Miejsce wystawienia
	P2        string   `xml:"P_2"`            // Numer faktury
	WZ        []string `xml:"WZ,omitempty"`

	// Date of delivery/service (choice: single date or period)
	P6      string   `xml:"P_6,omitempty"`
	OkresFa *OkresFa `xml:"OkresFa,omitempty"`

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

	P15         string `xml:"P_15"` // Kwota należności ogółem
	KursWalutyZ string `xml:"KursWalutyZ,omitempty"`

	Adnotacje     Adnotacje `xml:"Adnotacje"`
	RodzajFaktury string    `xml:"RodzajFaktury"`

	// Correction-related fields
	PrzyczynaKorekty   string              `xml:"PrzyczynaKorekty,omitempty"`
	TypKorekty         *int                `xml:"TypKorekty,omitempty"`
	DaneFaKorygowanej  []DaneFaKorygowanej `xml:"DaneFaKorygowanej,omitempty"`
	OkresFaKorygowanej string              `xml:"OkresFaKorygowanej,omitempty"`
	NrFaKorygowany     string              `xml:"NrFaKorygowany,omitempty"`
	Podmiot1K          *Podmiot1K          `xml:"Podmiot1K,omitempty"`
	Podmiot2K          []Podmiot2K         `xml:"Podmiot2K,omitempty"`
	P15ZK              string              `xml:"P_15ZK,omitempty"`
	KursWalutyZK       string              `xml:"KursWalutyZK,omitempty"`

	ZaliczkaCzesciowa []ZaliczkaCzesciowa `xml:"ZaliczkaCzesciowa,omitempty"`
	FP                *int                `xml:"FP,omitempty"`
	TP                *int                `xml:"TP,omitempty"`
	DodatkowyOpis     []KluczWartosc      `xml:"DodatkowyOpis,omitempty"`
	FakturaZaliczkowa []FakturaZaliczkowa `xml:"FakturaZaliczkowa,omitempty"`
	ZwrotAkcyzy       *int                `xml:"ZwrotAkcyzy,omitempty"`

	FaWiersz          []FaWiersz         `xml:"FaWiersz,omitempty"`
	Rozliczenie       *Rozliczenie       `xml:"Rozliczenie,omitempty"`
	Platnosc          *Platnosc          `xml:"Platnosc,omitempty"`
	WarunkiTransakcji *WarunkiTransakcji `xml:"WarunkiTransakcji,omitempty"`
	Zamowienie        *Zamowienie        `xml:"Zamowienie,omitempty"`
}

type OkresFa struct {
	P6Od string `xml:"P_6_Od"`
	P6Do string `xml:"P_6_Do"`
}

// --- Adnotacje ---

type Adnotacje struct {
	P16                  int                  `xml:"P_16"`
	P17                  int                  `xml:"P_17"`
	P18                  int                  `xml:"P_18"`
	P18A                 int                  `xml:"P_18A"`
	Zwolnienie           Zwolnienie           `xml:"Zwolnienie"`
	NoweSrodkiTransportu NoweSrodkiTransportu `xml:"NoweSrodkiTransportu"`
	P23                  int                  `xml:"P_23"`
	PMarzy               PMarzy               `xml:"PMarzy"`
}

type Zwolnienie struct {
	P19  *int   `xml:"P_19,omitempty"`
	P19A string `xml:"P_19A,omitempty"`
	P19B string `xml:"P_19B,omitempty"`
	P19C string `xml:"P_19C,omitempty"`
	P19N *int   `xml:"P_19N,omitempty"`
}

type NoweSrodkiTransportu struct {
	P22                  *int                   `xml:"P_22,omitempty"`
	P42_5                *int                   `xml:"P_42_5,omitempty"`
	NowySrodekTransportu []NowySrodekTransportu `xml:"NowySrodekTransportu,omitempty"`
	P22N                 *int                   `xml:"P_22N,omitempty"`
}

type NowySrodekTransportu struct {
	P22A          string `xml:"P_22A"`
	PNrWierszaNST int    `xml:"P_NrWierszaNST"`
	P22BMK        string `xml:"P_22BMK,omitempty"`
	P22BMD        string `xml:"P_22BMD,omitempty"`
	P22BK         string `xml:"P_22BK,omitempty"`
	P22BNR        string `xml:"P_22BNR,omitempty"`
	P22BRP        string `xml:"P_22BRP,omitempty"`
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

type PMarzy struct {
	PPMarzy   *int `xml:"P_PMarzy,omitempty"`
	PPMarzy2  *int `xml:"P_PMarzy_2,omitempty"`
	PPMarzy31 *int `xml:"P_PMarzy_3_1,omitempty"`
	PPMarzy32 *int `xml:"P_PMarzy_3_2,omitempty"`
	PPMarzy33 *int `xml:"P_PMarzy_3_3,omitempty"`
	PPMarzyN  *int `xml:"P_PMarzyN,omitempty"`
}

// --- Invoice line items ---

type FaWiersz struct {
	NrWierszaFa int    `xml:"NrWierszaFa"`
	UUID        string `xml:"UU_ID,omitempty"`
	P6A         string `xml:"P_6A,omitempty"`
	P7          string `xml:"P_7,omitempty"`
	Indeks      string `xml:"Indeks,omitempty"`
	GTIN        string `xml:"GTIN,omitempty"`
	PKWiU       string `xml:"PKWiU,omitempty"`
	CN          string `xml:"CN,omitempty"`
	PKOB        string `xml:"PKOB,omitempty"`
	P8A         string `xml:"P_8A,omitempty"`  // unit of measure
	P8B         string `xml:"P_8B,omitempty"`  // quantity
	P9A         string `xml:"P_9A,omitempty"`  // unit price net
	P9B         string `xml:"P_9B,omitempty"`  // unit price gross
	P10         string `xml:"P_10,omitempty"`  // discount
	P11         string `xml:"P_11,omitempty"`  // net value
	P11A        string `xml:"P_11A,omitempty"` // gross value
	P11Vat      string `xml:"P_11Vat,omitempty"`
	P12         string `xml:"P_12,omitempty"` // tax rate
	P12XII      string `xml:"P_12_XII,omitempty"`
	P12Zal15    *int   `xml:"P_12_Zal_15,omitempty"`
	KwotaAkcyzy string `xml:"KwotaAkcyzy,omitempty"`
	GTU         string `xml:"GTU,omitempty"`
	Procedura   string `xml:"Procedura,omitempty"`
	KursWaluty  string `xml:"KursWaluty,omitempty"`
	StanPrzed   *int   `xml:"StanPrzed,omitempty"`
}

// --- Correction data ---

type DaneFaKorygowanej struct {
	DataWystFaKorygowanej string `xml:"DataWystFaKorygowanej"`
	NrFaKorygowanej       string `xml:"NrFaKorygowanej"`
	NrKSeF                *int   `xml:"NrKSeF,omitempty"`
	NrKSeFFaKorygowanej   string `xml:"NrKSeFFaKorygowanej,omitempty"`
	NrKSeFN               *int   `xml:"NrKSeFN,omitempty"`
}

type Podmiot1K struct {
	PrefiksPodatnika    string        `xml:"PrefiksPodatnika,omitempty"`
	DaneIdentyfikacyjne Podmiot1Ident `xml:"DaneIdentyfikacyjne"`
	Adres               Adres         `xml:"Adres"`
}

type Podmiot2K struct {
	DaneIdentyfikacyjne Podmiot2Ident `xml:"DaneIdentyfikacyjne"`
	Adres               *Adres        `xml:"Adres,omitempty"`
	IDNabywcy           string        `xml:"IDNabywcy,omitempty"`
}

// --- Advance payments ---

type ZaliczkaCzesciowa struct {
	P6Z          string `xml:"P_6Z"`
	P15Z         string `xml:"P_15Z"`
	KursWalutyZW string `xml:"KursWalutyZW,omitempty"`
}

type FakturaZaliczkowa struct {
	NrKSeFZN            *int   `xml:"NrKSeFZN,omitempty"`
	NrFaZaliczkowej     string `xml:"NrFaZaliczkowej,omitempty"`
	NrKSeFFaZaliczkowej string `xml:"NrKSeFFaZaliczkowej,omitempty"`
}

// --- Key-Value ---

type KluczWartosc struct {
	NrWiersza *int   `xml:"NrWiersza,omitempty"`
	Klucz     string `xml:"Klucz"`
	Wartosc   string `xml:"Wartosc"`
}

// --- Settlement ---

type Rozliczenie struct {
	Obciazenia    []ObciazenieOdliczenie `xml:"Obciazenia,omitempty"`
	SumaObciazen  string                 `xml:"SumaObciazen,omitempty"`
	Odliczenia    []ObciazenieOdliczenie `xml:"Odliczenia,omitempty"`
	SumaOdliczen  string                 `xml:"SumaOdliczen,omitempty"`
	DoZaplaty     string                 `xml:"DoZaplaty,omitempty"`
	DoRozliczenia string                 `xml:"DoRozliczenia,omitempty"`
}

type ObciazenieOdliczenie struct {
	Kwota string `xml:"Kwota"`
	Powod string `xml:"Powod"`
}

// --- Payment ---

type Platnosc struct {
	// Full payment
	Zaplacono   *int   `xml:"Zaplacono,omitempty"`
	DataZaplaty string `xml:"DataZaplaty,omitempty"`

	// Partial payment
	ZnacznikZaplatyCzesciowej *int               `xml:"ZnacznikZaplatyCzesciowej,omitempty"`
	ZaplataCzesciowa          []ZaplataCzesciowa `xml:"ZaplataCzesciowa,omitempty"`

	TerminPlatnosci        []TerminPlatnosci `xml:"TerminPlatnosci,omitempty"`
	FormaPlatnosci         *int              `xml:"FormaPlatnosci,omitempty"`
	PlatnoscInna           *int              `xml:"PlatnoscInna,omitempty"`
	OpisPlatnosci          string            `xml:"OpisPlatnosci,omitempty"`
	RachunekBankowy        []RachunekBankowy `xml:"RachunekBankowy,omitempty"`
	RachunekBankowyFaktora []RachunekBankowy `xml:"RachunekBankowyFaktora,omitempty"`
	Skonto                 *Skonto           `xml:"Skonto,omitempty"`
	LinkDoPlatnosci        string            `xml:"LinkDoPlatnosci,omitempty"`
	IPKSeF                 string            `xml:"IPKSeF,omitempty"`
}

type ZaplataCzesciowa struct {
	KwotaZaplatyCzesciowej string `xml:"KwotaZaplatyCzesciowej"`
	DataZaplatyCzesciowej  string `xml:"DataZaplatyCzesciowej"`
	FormaPlatnosci         *int   `xml:"FormaPlatnosci,omitempty"`
	PlatnoscInna           *int   `xml:"PlatnoscInna,omitempty"`
	OpisPlatnosci          string `xml:"OpisPlatnosci,omitempty"`
}

type TerminPlatnosci struct {
	Termin     string      `xml:"Termin,omitempty"`
	TerminOpis *TerminOpis `xml:"TerminOpis,omitempty"`
}

type TerminOpis struct {
	Ilosc               int    `xml:"Ilosc"`
	Jednostka           string `xml:"Jednostka"`
	ZdarzeniePoczatkowe string `xml:"ZdarzeniePoczatkowe"`
}

type Skonto struct {
	WarunkiSkonta  string `xml:"WarunkiSkonta"`
	WysokoscSkonta string `xml:"WysokoscSkonta"`
}

type RachunekBankowy struct {
	NrRB                string `xml:"NrRB"`
	SWIFT               string `xml:"SWIFT,omitempty"`
	RachunekWlasnyBanku *int   `xml:"RachunekWlasnyBanku,omitempty"`
	NazwaBanku          string `xml:"NazwaBanku,omitempty"`
	OpisRachunku        string `xml:"OpisRachunku,omitempty"`
}

// --- Transaction conditions ---

type WarunkiTransakcji struct {
	Umowy                []Umowa         `xml:"Umowy,omitempty"`
	Zamowienia           []ZamowienieRef `xml:"Zamowienia,omitempty"`
	NrPartiiTowaru       []string        `xml:"NrPartiiTowaru,omitempty"`
	WarunkiDostawy       string          `xml:"WarunkiDostawy,omitempty"`
	KursUmowny           string          `xml:"KursUmowny,omitempty"`
	WalutaUmowna         string          `xml:"WalutaUmowna,omitempty"`
	Transport            []Transport     `xml:"Transport,omitempty"`
	PodmiotPosredniczacy *int            `xml:"PodmiotPosredniczacy,omitempty"`
}

type Umowa struct {
	DataUmowy string `xml:"DataUmowy,omitempty"`
	NrUmowy   string `xml:"NrUmowy,omitempty"`
}

type ZamowienieRef struct {
	DataZamowienia string `xml:"DataZamowienia,omitempty"`
	NrZamowienia   string `xml:"NrZamowienia,omitempty"`
}

type Transport struct {
	RodzajTransportu       *int        `xml:"RodzajTransportu,omitempty"`
	TransportInny          *int        `xml:"TransportInny,omitempty"`
	OpisInnegoTransportu   string      `xml:"OpisInnegoTransportu,omitempty"`
	Przewoznik             *Przewoznik `xml:"Przewoznik,omitempty"`
	NrZleceniaTransportu   string      `xml:"NrZleceniaTransportu,omitempty"`
	OpisLadunku            *int        `xml:"OpisLadunku,omitempty"`
	LadunekInny            *int        `xml:"LadunekInny,omitempty"`
	OpisInnegoLadunku      string      `xml:"OpisInnegoLadunku,omitempty"`
	JednostkaOpakowania    string      `xml:"JednostkaOpakowania,omitempty"`
	DataGodzRozpTransportu string      `xml:"DataGodzRozpTransportu,omitempty"`
	DataGodzZakTransportu  string      `xml:"DataGodzZakTransportu,omitempty"`
	WysylkaZ               *Adres      `xml:"WysylkaZ,omitempty"`
	WysylkaPrzez           []Adres     `xml:"WysylkaPrzez,omitempty"`
	WysylkaDo              *Adres      `xml:"WysylkaDo,omitempty"`
}

type Przewoznik struct {
	DaneIdentyfikacyjne Podmiot2Ident `xml:"DaneIdentyfikacyjne"`
	AdresPrzewoznika    Adres         `xml:"AdresPrzewoznika"`
}

// --- Order (for advance invoices) ---

type Zamowienie struct {
	WartoscZamowienia string             `xml:"WartoscZamowienia"`
	ZamowienieWiersz  []ZamowienieWiersz `xml:"ZamowienieWiersz"`
}

type ZamowienieWiersz struct {
	NrWierszaZam int    `xml:"NrWierszaZam"`
	UUIDZ        string `xml:"UU_IDZ,omitempty"`
	P7Z          string `xml:"P_7Z,omitempty"`
	IndeksZ      string `xml:"IndeksZ,omitempty"`
	GTINZ        string `xml:"GTINZ,omitempty"`
	PKWiUZ       string `xml:"PKWiUZ,omitempty"`
	CNZ          string `xml:"CNZ,omitempty"`
	PKOBZ        string `xml:"PKOBZ,omitempty"`
	P8AZ         string `xml:"P_8AZ,omitempty"`
	P8BZ         string `xml:"P_8BZ,omitempty"`
	P9AZ         string `xml:"P_9AZ,omitempty"`
	P11NettoZ    string `xml:"P_11NettoZ,omitempty"`
	P11VatZ      string `xml:"P_11VatZ,omitempty"`
	P12Z         string `xml:"P_12Z,omitempty"`
	P12ZXII      string `xml:"P_12Z_XII,omitempty"`
	P12ZZal15    *int   `xml:"P_12Z_Zal_15,omitempty"`
	GTUZ         string `xml:"GTUZ,omitempty"`
	ProceduraZ   string `xml:"ProceduraZ,omitempty"`
	KwotaAkcyzyZ string `xml:"KwotaAkcyzyZ,omitempty"`
	StanPrzedZ   *int   `xml:"StanPrzedZ,omitempty"`
}

// --- Footer ---

type Stopka struct {
	Informacje []Informacje `xml:"Informacje,omitempty"`
	Rejestry   []Rejestr    `xml:"Rejestry,omitempty"`
}

type Informacje struct {
	StopkaFaktury string `xml:"StopkaFaktury,omitempty"`
}

type Rejestr struct {
	PelnaNazwa string `xml:"PelnaNazwa,omitempty"`
	KRS        string `xml:"KRS,omitempty"`
	REGON      string `xml:"REGON,omitempty"`
	BDO        string `xml:"BDO,omitempty"`
}

// --- Attachment ---

type Zalacznik struct {
	BlokDanych []BlokDanych `xml:"BlokDanych"`
}

type BlokDanych struct {
	ZNaglowek string          `xml:"ZNaglowek,omitempty"`
	MetaDane  []MetaDane      `xml:"MetaDane"`
	Tekst     *ZalacznikTekst `xml:"Tekst,omitempty"`
	Tabela    []Tabela        `xml:"Tabela,omitempty"`
}

type MetaDane struct {
	ZKlucz   string `xml:"ZKlucz"`
	ZWartosc string `xml:"ZWartosc"`
}

type ZalacznikTekst struct {
	Akapit []string `xml:"Akapit"`
}

type Tabela struct {
	TMetaDane []TMetaDane    `xml:"TMetaDane,omitempty"`
	Opis      string         `xml:"Opis,omitempty"`
	TNaglowek TabelaNaglowek `xml:"TNaglowek"`
	Wiersz    []TabelaWiersz `xml:"Wiersz"`
	Suma      *TabelaSuma    `xml:"Suma,omitempty"`
}

type TMetaDane struct {
	TKlucz   string `xml:"TKlucz"`
	TWartosc string `xml:"TWartosc"`
}

type TabelaNaglowek struct {
	Kol []TabelaKol `xml:"Kol"`
}

type TabelaKol struct {
	Typ  string `xml:"Typ,attr"`
	NKom string `xml:"NKom"`
}

type TabelaWiersz struct {
	WKom []string `xml:"WKom"`
}

type TabelaSuma struct {
	SKom []string `xml:"SKom"`
}
