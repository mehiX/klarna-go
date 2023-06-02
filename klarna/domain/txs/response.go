package txs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mehix/klarna-go/klarna/domain"
)

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Reports []CategorizedTransactionReport `json:"reports"`
}

type CategorizedTransactionReport struct {
	Type               string                   `json:"type"`
	InsightsConsumerID string                   `json:"insights_consumer_id"`
	InsightsAccountIDs []string                 `json:"insights_account_ids"`
	Transactions       []CategorizedTransaction `json:"transactions"`
	FromDate           string                   `json:"from_date"`
	ToDate             string                   `json:"to_date"`
	RefreshedAt        time.Time                `json:"refreshed_at"`
	DataAvailability   dataAvailability         `json:"data_availability"`
}

type CategorizedTransaction struct {
	InsightsTransactionID string                     `json:"insights_transaction_id"`
	InsightsAccountID     string                     `json:"insights_account_id"`
	OBtransactionID       string                     `json:"ob_transaction_id"`
	ValueDate             string                     `json:"value_date"`
	BookingDate           string                     `json:"booking_date"`
	Amount                domain.Amount              `json:"amount"`
	Reference             string                     `json:"reference"`
	BankReferences        domain.BankReferences      `json:"bank_references"`
	State                 string                     `json:"state"`  // Enum<'PROCESSED', 'PENDING', 'CANCELED', 'FAILED'>
	Type                  string                     `json:"type"`   // Enum<'DEBIT', 'CREDIT'>
	Method                string                     `json:"method"` // Enum<'TRANSFER', 'DIRECT_DEBIT', 'INSTANT', 'UNKNOWN'>
	BankTransactionCode   domain.BankTransactionCode `json:"bank_transaction_code"`
	CounterParty          domain.CounterParty        `json:"counter_party"`
	Categories            []domain.Category          `json:"categories"`
	Brand                 domain.Brand               `json:"brand"`
	Labels                map[string]string          `json:"labels"`
	ExtraCategories       []domain.Category          `json:"extra_categories"` // categorized added by client libraries
}

func (t CategorizedTransaction) IsBetaalautomaat() bool {
	// IF(OR(LEFT([@[bank_references_unstructured]];14)="Betaalautomaat";LEFT([@[bank_references_unstructured]];3)="BEA");1;0)
	return strings.HasPrefix(t.BankReferences.Unstructured, "Betaalautomat") || strings.HasPrefix(t.BankReferences.Unstructured, "BEA")
}

func (t CategorizedTransaction) Time() string {

	var reTime = regexp.MustCompile("[0-9]{2}[:.]{1}[0-9]{2} ")
	// IF(LEFT([@reference];3)="BEA";REPLACE(MID([@reference];SEARCH("/";[@reference];1)+1;5);3;1;":");IF([@[T_Betaalautomaat]]=1;MID([@reference];16;5);""))
	return strings.Replace(strings.TrimSpace(reTime.FindString(t.Reference)), ".", ":", 1)
}

func (t CategorizedTransaction) IsCredit() bool {
	return t.Type == "CREDIT"
}

func (t CategorizedTransaction) IsDebit() bool {
	return t.Type == "DEBIT"
}

// Hour returns the hour from the transaction.
// If the time cannot be found, it returns -1
func (t CategorizedTransaction) hour() int {
	hhmm := t.Time()
	fmt.Printf("finding hour for time: %s\n", hhmm)
	if hhmm == "" {
		return -1
	}

	hh := strings.Split(hhmm, ":")[0]

	n, err := strconv.Atoi(hh)
	if err != nil {
		return -1
	}

	return n
}

func (t CategorizedTransaction) PartOfDay() (string, bool) {
	h := t.hour()
	if h < 0 {
		return "", false
	}

	switch {
	case h >= 6 && h < 12:
		return "morning", true
	case h >= 12 && h < 18:
		return "afternoon", true
	case h >= 18 && h < 22:
		return "evening", true
	default:
		return "night", true
	}
}

func (t CategorizedTransaction) AmountStr() string {

	// IF([@type]="CREDIT";[@[amount_amount]]/100;-1*([@[amount_amount]]/100))
	sign := 1.0
	if t.Type != "CREDIT" {
		sign = -sign
	}
	return fmt.Sprintf("%.2f", sign*float64(t.Amount.Amount)/100.0)
}

func (t CategorizedTransaction) TimePerDaypart() string {

	// IFERROR(IF(HOUR([@[T_Time]])<6;"4 Night";IF(HOUR([@[T_Time]])<12;"1 Morning";IF(HOUR([@[T_Time]])<18;"2 Afternoon";"3 Evening")));"")
	if hhmm, err := time.Parse("15:04", t.Time()); err == nil {
		switch h := hhmm.Hour(); {
		case h < 12:
			return "1 Morning"
		case h < 18:
			return "2 Afternoon"
		default:
			return "3 Evening"
		}
	}

	return ""
}

func (t CategorizedTransaction) Weekday() string {

	// WEEKDAY([@[value_date]];2)
	vd, err := time.Parse("2006-01-02", t.ValueDate)
	if err == nil {
		return vd.Format("Monday")
	}

	return err.Error()
}

// YearMonth returns the ValueDate as `yyyy-mm`
func (t CategorizedTransaction) YearMonth() string {
	if len(t.ValueDate) < 7 {
		return t.ValueDate
	}

	return string(t.ValueDate[:7])
}
