package domain

type AccountType string

const (
	ACCOUNT_DEFAULT    AccountType = AccountType("DEFAULT")
	ACCOUNT_SAVING     AccountType = AccountType("SAVING")
	ACCOUNT_CREDITCARD AccountType = AccountType("CREDITCARD")
	ACCOUNT_DEPOT      AccountType = AccountType("DEPOT")
)
