package main

import (
	"fmt"
	"ozonadv/config"
	"ozonadv/internal/ozon"
	"ozonadv/internal/stat"

	"github.com/spf13/cobra"
)

func main() {
	fmt.Println("")

	cfg := config.NewOrFail("config.yml")
	ozonClient := ozon.NewClient(cfg.Ozon)

	rootCmd := &cobra.Command{
		Use:   "ozonadv",
		Short: "Консольное приложение выгрузки статистики рекламных кабинетов Озон",
	}

	initFetchCommand(rootCmd, ozonClient)

	fmt.Println("")
	rootCmd.Execute()
}

func initFetchCommand(rootCmd *cobra.Command, ozonClient *ozon.Client) {
	cmd := &cobra.Command{
		Use:     "stat",
		Short:   "Статистика по кампаниям",
		Example: "ozonadv stat --from-date 2025-01-01 --to-date 2025-01-02",
		RunE: func(cmd *cobra.Command, args []string) error {
			options := stat.HandleOptions{}
			options.FromDate, _ = cmd.PersistentFlags().GetString("from-date")
			options.ToDate, _ = cmd.PersistentFlags().GetString("to-date")

			return stat.Handle(ozonClient, options)
		},
	}

	cmd.PersistentFlags().String("from-date", "", "Начало периода")
	cmd.PersistentFlags().String("to-date", "", "Окончание периода")

	rootCmd.AddCommand(cmd)
}
