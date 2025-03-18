package main

import (
	"fmt"
	"log"
	"ozonadv/internal/application"
	"ozonadv/internal/stat"
	"ozonadv/pkg/console"

	"github.com/spf13/cobra"
)

var statUsecases *stat.Usecases

func main() {
	log.SetFlags(0)
	defer fmt.Println("")

	app := application.New()
	defer app.Shutdown()

	rootCmd := &cobra.Command{
		Use:   "ozonadv",
		Short: "Консольное приложение выгрузки статистики рекламных кабинетов Озон",
	}

	initFetchCommand(rootCmd, app)

	fmt.Println("")
	rootCmd.Execute()
}

func initFetchCommand(rootCmd *cobra.Command, app *application.Application) {
	cmd := &cobra.Command{
		Use:     "fetch",
		Short:   "Запрос на формирование отчетов статистики по кампаниям",
		Example: "ozonadv fetch --from-date 2025-01-01 --to-date 2025-01-02",
		RunE: func(cmd *cobra.Command, args []string) error {
			statUsecases := app.StatUsecases()

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
