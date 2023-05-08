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
	"github.com/mehix/klarna-go/klarna/service/account"
	"github.com/mehix/klarna-go/klarna/service/txs"
	"github.com/spf13/cobra"
)

var (
	insightsConsumerID string
	size               int64
	days               int
	debug              bool
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

			transactions, err := txsSvc.FetchLatest(context.Background(), insightsConsumerID, size)
			if err != nil {
				log.Println(err)
				os.Exit(2)
			}

			b, _ := json.MarshalIndent(transactions, " ", "  ")
			fmt.Println(string(b))
		},
	}

	var cmdAccountsInfo = &cobra.Command{
		Use:   "accounts",
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

	cmdTransactions.Flags().Int64VarP(&size, "limit", "l", 1, "number of transactions to return")
	cmdAccountsBalanceOverTime.Flags().IntVar(&days, "days", 1000, "days of history to retrieve")

	var rootCmd = &cobra.Command{Use: "klarna"}
	rootCmd.PersistentFlags().StringVar(&insightsConsumerID, "id", "", "InsightsConsumerID for the customer's consent")
	rootCmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "Enable debug")

	cmdAccounts.AddCommand(cmdAccountsInfo, cmdAccountsBalances, cmdAccountsBalanceOverTime)
	cmdView.AddCommand(cmdTransactions, cmdAccounts)
	rootCmd.AddCommand(cmdView)

	rootCmd.Execute()
}
