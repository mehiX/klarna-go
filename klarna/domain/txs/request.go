package txs

var DefaultRequest = Request{
	RequiredDataAvailability: AVAILABILITY_PARTIAL,
	Order: []ReportOptionsOrderBy{
		{Order: "DESC", Field: "DATE"},
	},
	CombineAccounts: "ALL",
	ReportDays:      365 * 5,
}

// Request as described in `https://docs.openbanking.klarna.com/acin/reports/categorized-transactions.html`
type Request struct {
	InsightsConsumerID       string                 `json:"insights_consumer_id"`
	RequiredDataAvailability dataAvailability       `json:"required_data_availability,omitempty"`
	Order                    []ReportOptionsOrderBy `json:"order,omitempty"`
	CombineAccounts          string                 `json:"combine_accounts,omitempty"` // ?Enum<'NONE', 'ALL'> default 'NONE'
	Size                     int64                  `json:"size,omitempty"`             // Only return a maximum of size transactions in the report. The default will be unlimited.
	ReportDays               int64                  `json:"report_days,omitempty"`
	TransactionType          string                 `json:"transaction_type,omitempty"` // ?Enum<'DEBIT', 'CREDIT', 'DEBIT_AND_ZERO'>
	FromDate                 string                 `json:"from_date,omitempty"`        // Date (String: "YYYY-MM-DD")
	ToDate                   string                 `json:"to_date,omitempty"`          // Date (String: "YYYY-MM-DD")
}

type dataAvailability string

const (
	AVAILABILITY_COMPLETE = dataAvailability("COMPLETE")
	AVAILABILITY_PARTIAL  = dataAvailability("PARTIAL")
	AVAILABILITY_NONE     = dataAvailability("NONE")
)

type ReportOptionsOrderBy struct {
	Order string `json:"order"` // ?Enum<'ASC', 'DESC'> default ASC
	Field string `json:"field"` // Enum<'AMOUNT','DATE','CATEGORY','BRAND','AMOUNT_WITH_TYPE'>
}
