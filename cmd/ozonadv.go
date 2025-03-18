package main

import (
	"fmt"
	"log"
	"ozonadv/internal/application"
	"ozonadv/internal/stat"
	"ozonadv/pkg/console"

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(0)
	defer fmt.Println("")

	app := application.New()
	defer app.Shutdown()

	rootCmd := &cobra.Command{
		Use:   "ozonadv",
		Short: "Консольное приложение выгрузки статистики рекламных кабинетов Озон",
	}

	initStatCreateCommand(rootCmd, app)
	initStatInfoCommand(rootCmd, app)
	initStatPullCommand(rootCmd, app)

	fmt.Println("")
	rootCmd.Execute()
}

func initStatCreateCommand(rootCmd *cobra.Command, app *application.Application) {
	cmd := &cobra.Command{
		Use:     "stat:create",
		Short:   "Запрос формирования отчетов статистики по кампаниям",
		Example: "ozonadv stat:create --from-date 2025-01-01 --to-date 2025-01-02",
		RunE: func(cmd *cobra.Command, args []string) error {
			statUsecases := app.StatUsecases()

			if statUsecases.HasIncompletedStatistics() {
				fmt.Println("Найдены незагруженные отчеты.")
				if console.Ask("Продолжить?") == false {
					return nil
				}
			}

			options := stat.CreateOptions{}
			options.FromDate, _ = cmd.PersistentFlags().GetString("from-date")
			options.ToDate, _ = cmd.PersistentFlags().GetString("to-date")
			options.CampaignsPerRequest, _ = cmd.PersistentFlags().GetInt("campaigns-per-request")

			return statUsecases.Create(options)
		},
	}

	cmd.PersistentFlags().String("from-date", "", "Начало периода")
	cmd.PersistentFlags().String("to-date", "", "Окончание периода")
	cmd.PersistentFlags().Int("campaigns-per-request", 10, "Количество кампаний на один запрос")

	rootCmd.AddCommand(cmd)
}

func initStatInfoCommand(rootCmd *cobra.Command, app *application.Application) {
	cmd := &cobra.Command{
		Use:     "stat:info",
		Short:   "Статус формирования отчетов",
		Example: "ozonadv stat:info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.StatUsecases().Info()
		},
	}

	rootCmd.AddCommand(cmd)
}

func initStatPullCommand(rootCmd *cobra.Command, app *application.Application) {
	cmd := &cobra.Command{
		Use:     "stat:pull",
		Short:   "Получить отчеты",
		Example: "ozonadv stat:pull",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.StatUsecases().Pull()
		},
	}

	rootCmd.AddCommand(cmd)
}
