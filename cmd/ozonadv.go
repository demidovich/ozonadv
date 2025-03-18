package main

import (
	"fmt"
	"log"
	"ozonadv/config"
	"ozonadv/internal/ozon"
	"ozonadv/internal/stat"
	"ozonadv/internal/storage"
	"ozonadv/pkg/console"

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(0)
	fmt.Println("")
	defer fmt.Println("")

	storage := storage.New()
	defer storage.SaveState()

	cfg := config.NewOrFail("config.yml")
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
			if statUsecases.HasIncompletedStatistics() {
				fmt.Println("Найдены незагруженные отчеты.")
				if console.Ask("Продолжить?") == false {
					return nil
				}
			}

			options := stat.FetchOptions{}
			options.FromDate, _ = cmd.PersistentFlags().GetString("from-date")
			options.ToDate, _ = cmd.PersistentFlags().GetString("to-date")
			options.CampaignsPerRequest, _ = cmd.PersistentFlags().GetInt("campaigns-per-request")

			return statUsecases.Fetch(options)
		},
	}

	cmd.PersistentFlags().String("from-date", "", "Начало периода")
	cmd.PersistentFlags().String("to-date", "", "Окончание периода")
	cmd.PersistentFlags().Int("campaigns-per-request", 10, "Количество кампаний на один запрос")

	rootCmd.AddCommand(cmd)
}
