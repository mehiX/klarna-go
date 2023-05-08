package domain

import "fmt"

type Amount struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

func (a Amount) String() string {
	return fmt.Sprintf("%s %.2f", a.Currency, float64(a.Amount)/100)
}

type addressData struct {
	StreetAddress  string `json:"street_address"`
	StreetAddress2 string `json:"street_address2"`
	PostalCode     string `json:"postalcode"`
	City           string `json:"city"`
	Region         string `json:"region"`
	// Country code (2)
	Country string `json:"country"`
}

type Pagination struct {
	Count int64          `json:"count"`
	URL   string         `json:"url"`
	Next  PaginationNext `json:"next"`
}

type PaginationNext struct {
	Offset string `json:"offset"`
}

type AllowedAccounts struct {
	AccountTypes []AccountType `json:"account_types"`
	AccountIDs   []string      `json:"account_ids"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Brand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ErrorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Errors  []ErrorDetails `json:"errors,omitempty"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s - %s", e.Code, e.Message)
}

type ErrorDetails struct {
	Code     string `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	Location string `json:"location,omitempty"`
}

type BankReferences struct {
	Unstructured string `json:"unstructured"`
	Structured   string `json:"structured"`
	EndToEnd     string `json:"end_to_end"`
}

type BankTransactionCode struct {
	Code        string `json:"code"`
	SubCode     string `json:"sub_code"`
	Description string `json:"description"`
}

type CounterParty struct {
	IBAN          string        `json:"iban"`
	AccountNumber string        `json:"account_number"`
	BIC           string        `json:"bic"`
	BankCode      string        `json:"bank_code"`
	HolderName    string        `json:"holder_name"`
	HolderAddress HolderAddress `json:"holder_address"`
}

type HolderAddress struct {
	Country string `json:"country"`
}
