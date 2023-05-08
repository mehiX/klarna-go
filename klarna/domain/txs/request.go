package txs

var DefaultRequest = Request{
	RequiredDataAvailability: AVAILABILITY_PARTIAL,
	Order: []ReportOptionsOrderBy{
		{Order: "DESC", Field: "DATE"},
	},
	CombineAccounts: "ALL",
	ReportDays:      365 * 5,
}

type Request struct {
	InsightsConsumerID       string                 `json:"insights_consumer_id"`
	RequiredDataAvailability dataAvailability       `json:"required_data_availability,omitempty"`
	Order                    []ReportOptionsOrderBy `json:"order,omitempty"`
	CombineAccounts          string                 `json:"combine_accounts,omitempty"` // ?Enum<'NONE', 'ALL'> default 'NONE'
	Size                     int64                  `json:"size"`                       // Only return a maximum of size transactions in the report. The default will be unlimited.
	ReportDays               int64                  `json:"report_days"`
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
