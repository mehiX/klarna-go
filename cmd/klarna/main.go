package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mehix/klarna-go/klarna"
	"github.com/mehix/klarna-go/klarna/domain/report"
	domainTxs "github.com/mehix/klarna-go/klarna/domain/txs"
	"github.com/mehix/klarna-go/klarna/service/account"
	"github.com/mehix/klarna-go/klarna/service/txs"
	"github.com/spf13/cobra"
)

var (
	insightsConsumerID    string
	size                  int64
	days                  int
	debug                 bool
	onlyCredit, onlyDebit bool
	fromDate, toDate      string
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	var cmdView = &cobra.Command{
		Use:   "view",
		Short: "Fetch information from Klarna",
	}

	var cmdAccounts = &cobra.Command{
		Use:   "accounts",
		Short: "Show information about an account",
	}

	var cmdTransactions = &cobra.Command{
		Use:   "txs [# limit]",
		Short: "List transactions",
		Run: func(cmd *cobra.Command, args []string) {
			kc := klarna.New(os.Getenv("KLARNA_BASE_URL"), os.Getenv("KLARNA_TOKEN"), klarna.WithDebug(debug))
			txsSvc := txs.NewService(kc)

			if (fromDate != "" && toDate == "") || (fromDate == "" && toDate != "") {
				log.Println("fromDate and toDate should be either both empty or both provided")
				os.Exit(3)
			}

			var transactions any
			var err error

			if onlyCredit {
				transactions, err = txsSvc.FetchLatestCredit(context.TODO(), insightsConsumerID, size)
			} else if onlyDebit {
				if fromDate != "" {
					transactions, err = txsSvc.FetchDebitForPeriod(context.TODO(), insightsConsumerID, fromDate, toDate)
				} else {
					transactions, err = txsSvc.FetchLatestDebit(context.TODO(), insightsConsumerID, size)
				}
			} else {
				transactions, err = txsSvc.FetchLatest(context.TODO(), insightsConsumerID, size)
			}

			if err != nil {
				log.Println(err)
				os.Exit(2)
			}

			b, _ := json.MarshalIndent(transactions, " ", "  ")
			fmt.Println(string(b))
		},
	}

	var cmdReportMonthly = &cobra.Command{
		Use:   "report-monthly",
		Short: "show a monthly report for credit and debit",
		Run: func(cmd *cobra.Command, args []string) {
			kc := klarna.New(os.Getenv("KLARNA_BASE_URL"), os.Getenv("KLARNA_TOKEN"), klarna.WithDebug(debug))
			txsSvc := txs.NewService(kc)

			transactions, err := txsSvc.ReportMonthlyCreditBalance(context.Background(), domainTxs.Filter{InsightsConsumerID: insightsConsumerID})
			if err != nil {
				log.Println(err)
				os.Exit(2)
			}

			b, _ := json.MarshalIndent(transactions, " ", "  ")
			fmt.Println(string(b))
		},
	}

	var cmdReportDaily = &cobra.Command{
		Use:   "report-daily",
		Short: "show a daily spending report, split by parts of day where possible",
		Run: func(cmd *cobra.Command, args []string) {
			kc := klarna.New(os.Getenv("KLARNA_BASE_URL"), os.Getenv("KLARNA_TOKEN"), klarna.WithDebug(debug))
			txsSvc := txs.NewService(kc)

			if (fromDate != "" && toDate == "") || (fromDate == "" && toDate != "") {
				log.Println("fromDate and toDate should be either both empty or both provided")
				os.Exit(3)
			}

			var transactions []report.DailySpending
			var err error

			if fromDate != "" {
				transactions, err = txsSvc.ReportDailySpendingForPeriod(context.TODO(), insightsConsumerID, fromDate, toDate)
			} else {
				transactions, err = txsSvc.ReportDailySpending(context.Background(), insightsConsumerID)
			}

			if err != nil {
				log.Println(err)
				os.Exit(2)
			}

			b, _ := json.MarshalIndent(transactions, " ", "  ")
			fmt.Println(string(b))
		},
	}

	var cmdAccountsInfo = &cobra.Command{
		Use:   "info",
		Short: "Show information for accounts associated to the consumer ID",
		Run: func(cmd *cobra.Command, args []string) {

			kc := klarna.New(os.Getenv("KLARNA_BASE_URL"), os.Getenv("KLARNA_TOKEN"), klarna.WithDebug(debug))
			accSvc := account.NewService(kc)

			info, err := accSvc.Info(context.Background(), insightsConsumerID)
			if err != nil {
				log.Println(err)
				os.Exit(2)
			}

			b, _ := json.MarshalIndent(info, " ", "  ")
			fmt.Println(string(b))
		},
	}

	var cmdAccountsBalances = &cobra.Command{
		Use:   "balances",
		Short: "Show balances for accounts associated to the consumer ID",
		Run: func(cmd *cobra.Command, args []string) {

			kc := klarna.New(os.Getenv("KLARNA_BASE_URL"), os.Getenv("KLARNA_TOKEN"), klarna.WithDebug(debug))
			accSvc := account.NewService(kc)

			balances, err := accSvc.Balances(context.Background(), insightsConsumerID)
			if err != nil {
				log.Println(err)
				os.Exit(2)
			}

			b, _ := json.MarshalIndent(balances, " ", "  ")
			fmt.Println(string(b))
		},
	}

	var cmdAccountsBalanceOverTime = &cobra.Command{
		Use:   "balanceOverTime",
		Short: "Show balance over time for accounts associated to the consumer ID",
		Run: func(cmd *cobra.Command, args []string) {

			kc := klarna.New(os.Getenv("KLARNA_BASE_URL"), os.Getenv("KLARNA_TOKEN"), klarna.WithDebug(debug))
			accSvc := account.NewService(kc)

			balances, err := accSvc.BalanceOverTime(context.Background(), insightsConsumerID, time.Now().Add(-time.Duration(days)*24*time.Hour), time.Now())
			if err != nil {
				log.Println(err)
				os.Exit(2)
			}

			b, _ := json.MarshalIndent(balances, " ", "  ")
			fmt.Println(string(b))
		},
	}

	cmdTransactions.PersistentFlags().Int64VarP(&size, "limit", "l", 1, "number of transactions to return")
	cmdTransactions.Flags().BoolVar(&onlyCredit, "only-credit", false, "show only credit transactions")
	cmdTransactions.Flags().BoolVar(&onlyDebit, "only-debit", false, "show only debit transactions")
	cmdTransactions.PersistentFlags().StringVar(&fromDate, "from", "", "begin of the reporting period")
	cmdTransactions.PersistentFlags().StringVar(&toDate, "to", "", "end of the reporting period")

	cmdAccountsBalanceOverTime.Flags().IntVar(&days, "days", 1000, "days of history to retrieve")

	var rootCmd = &cobra.Command{Use: "klarna"}
	rootCmd.PersistentFlags().StringVar(&insightsConsumerID, "id", "", "InsightsConsumerID for the customer's consent")
	rootCmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "Enable debug")

	cmdAccounts.AddCommand(cmdAccountsInfo, cmdAccountsBalances, cmdAccountsBalanceOverTime)
	cmdTransactions.AddCommand(cmdReportMonthly, cmdReportDaily)
	cmdView.AddCommand(cmdTransactions, cmdAccounts)
	rootCmd.AddCommand(cmdView)

	rootCmd.Execute()
}
