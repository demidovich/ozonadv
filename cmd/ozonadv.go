package main

import (
	"fmt"
	"log"
	"ozonadv/config"
	"ozonadv/internal/ozon"
	"ozonadv/internal/stat"
	"ozonadv/internal/storage"

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(0)
	fmt.Println("")

	cfg := config.NewOrFail("config.yml")
	storage := storage.New()
	ozonClient := ozon.NewClient(cfg.Ozon)
	statUsecases := stat.New(storage, ozonClient)

	rootCmd := &cobra.Command{
		Use:   "ozonadv",
		Short: "Консольное приложение выгрузки статистики рекламных кабинетов Озон",
	}

	initFetchCommand(rootCmd, statUsecases)

	fmt.Println("")
	rootCmd.Execute()
}

func initFetchCommand(rootCmd *cobra.Command, statUsecases stat.Usecases) {
	cmd := &cobra.Command{
		Use:     "fetch",
		Short:   "Запрос на формирование отчетов статистики по кампаниям",
		Example: "ozonadv fetch --from-date 2025-01-01 --to-date 2025-01-02",
		RunE: func(cmd *cobra.Command, args []string) error {
			options := stat.FetchOptions{}
			options.FromDate, _ = cmd.PersistentFlags().GetString("from-date")
			options.ToDate, _ = cmd.PersistentFlags().GetString("to-date")

			return statUsecases.Fetch(options)
		},
	}

	cmd.PersistentFlags().String("from-date", "", "Начало периода")
	cmd.PersistentFlags().String("to-date", "", "Окончание периода")

	rootCmd.AddCommand(cmd)
}
