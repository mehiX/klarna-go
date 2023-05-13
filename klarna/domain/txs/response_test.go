package txs

import "testing"

func TestTimeOfDay(t *testing.T) {

	type scenario struct {
		Transaction            CategorizedTransaction
		ReturnTime             string
		ReturnPartOfDay        string
		ReturnPartOfDaySuccess bool
	}

	scenarios := []scenario{
		{
			Transaction:            CategorizedTransaction{Reference: "BEA, Betaalpas OCCIDENTAL PUERTO BANU,PAS257 NR:12010159, 09.05.23/18:32 MARBELLA, Land: ESP"},
			ReturnTime:             "18:32",
			ReturnPartOfDay:        "evening",
			ReturnPartOfDaySuccess: true,
		},
		{
			Transaction:            CategorizedTransaction{Reference: "eCom, Betaalpas BOLT.EU/O/2305051750 05.05.23/19.50 Tallinn"},
			ReturnTime:             "19:50",
			ReturnPartOfDay:        "evening",
			ReturnPartOfDaySuccess: true,
		},
		{
			Transaction:            CategorizedTransaction{Reference: "SEPA iDEAL IBAN: NL04ADYB2017400157 BIC: ADYBNL2A Naam: GoFundMeGroupInc Omschrijving: BD5MP8RC2PC5Q4C2JE 98 0180312529043507 Ironman 70.3 Marbella - For Stre Kenmerk: 05-05-2023 14:58 018031 2529043507"},
			ReturnTime:             "14:58",
			ReturnPartOfDay:        "afternoon",
			ReturnPartOfDaySuccess: true,
		},
		{
			Transaction:            CategorizedTransaction{Reference: "SEPA Incasso algemeen doorlopend Incassant: DE5603800002197951 Naam: GC re ZWIFT Machtiging: HNK6P3S IBAN: DE81503104000437760000 Kenmerk: PA03TY2QZYSENF"},
			ReturnTime:             "",
			ReturnPartOfDay:        "",
			ReturnPartOfDaySuccess: false,
		},
		{
			Transaction:            CategorizedTransaction{Reference: "Betaalautomaat 00:26 pasnr. 021"},
			ReturnTime:             "00:26",
			ReturnPartOfDay:        "night",
			ReturnPartOfDaySuccess: true,
		},
		{
			Transaction:            CategorizedTransaction{Reference: "Loon dag mei 2023 322794 VER-VH0TvnN52e"},
			ReturnTime:             "",
			ReturnPartOfDay:        "",
			ReturnPartOfDaySuccess: false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.Transaction.Reference, func(t *testing.T) {
			s := s
			t.Parallel()

			hh := s.Transaction.Time()
			partOfDay, ok := s.Transaction.PartOfDay()

			if hh != s.ReturnTime {
				t.Fatalf("getting time of transaction. expected: %s, got: %s", s.ReturnTime, hh)
			}
			if ok != s.ReturnPartOfDaySuccess {
				t.Fatalf("assessing part of day. expected: %v, got: %v", s.ReturnPartOfDaySuccess, ok)
			}
			if s.ReturnPartOfDay != partOfDay {
				t.Fatalf("getting part of day. expected: %s, got: %s", s.ReturnPartOfDay, partOfDay)
			}
		})
	}
}
