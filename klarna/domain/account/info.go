package account

type RequestInfo struct {
	InsightsConsumerID         string   `json:"insights_consumer_id"`                    // String
	InsightsAccountIDs         []string `json:"insights_account_ids,omitempty"`          // ?Array<String>
	ExcludedInsightsAccountIDs []string `json:"excluded_insights_account_ids,omitempty"` // ?Array<String>
}

type ResponseInfo struct {
	Data DataInfo `json:"data"`
}

type DataInfo struct {
	Reports []Info `json:"reports"`
}

type Info struct {
	Type                      string            `json:"type"`                         // "ACCOUNT_INFOS"
	InsightsConsumerID        string            `json:"insights_consumer_id"`         // String
	InsightsAccountIDs        []string          `json:"insights_account_ids"`         // Array<String>
	AccountAlias              string            `json:"account_alias"`                // ?String
	HolderName                string            `json:"holder_name"`                  // ?String
	Currency                  string            `json:"currency"`                     // ?String
	Iban                      string            `json:"iban"`                         // ?String
	AccountNumber             string            `json:"account_number"`               // ?String
	Bic                       string            `json:"bic"`                          // ?String
	BankCode                  string            `json:"bank_code"`                    // ?String
	BankName                  string            `json:"bank_name"`                    // ?String
	Country                   string            `json:"country"`                      // ?String
	AccountType               string            `json:"account_type"`                 // ?Enum<'DEFAULT', 'SAVING', 'CREDITCARD', 'DEPOT'>
	TransferType              string            `json:"transfer_type"`                // ?Enum<'NONE', 'DOMESTIC', 'FULL', 'REFERENCE', 'RESTRICTED'>
	RefreshAt                 string            `json:"refreshed_at"`                 // DateTime
	AccountsGroupID           string            `json:"accounts_group_id"`            // ?String
	OldestTransactionDate     string            `json:"oldest_transaction_date"`      // ?Date
	ConsentExpectedExpiryDate string            `json:"consent_expected_expiry_date"` // ?Date
	Labels                    map[string]string `json:"labels"`                       // ?{ <label_key>: String }
}
