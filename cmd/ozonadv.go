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

	initStatCommand(rootCmd, app)
	initStatInfoCommand(rootCmd, app)
	initStatPullCommand(rootCmd, app)

	fmt.Println("")
	rootCmd.Execute()
}

func initStatCommand(rootCmd *cobra.Command, app *application.Application) {
	cmd := &cobra.Command{
		Use:     "stat",
		Short:   "Формирование и загрузка статистики по кампаниям",
		Example: "ozonadv stat --from-date 2025-01-01 --to-date 2025-01-02",
		RunE: func(cmd *cobra.Command, args []string) error {
			statUsecases := app.StatUsecases()

			if statUsecases.HasStatistics() {
				fmt.Println("Найдены незагруженные отчеты.")
				fmt.Println("Предыдущая загрузка была завершена не полностью.")
				fmt.Println("Для завершения загрузки следует выполнить команду stat:pull")
				fmt.Println("")
				fmt.Println("Незагруженные отчеты будут удалены.")
				if console.Ask("Продолжить?") == false {
					return nil
				}
				statUsecases.RemoveAllStatistics()
			}

			options := stat.CreateOptions{}
			options.FromDate, _ = cmd.PersistentFlags().GetString("from-date")
			options.ToDate, _ = cmd.PersistentFlags().GetString("to-date")
			options.ExportFile, _ = cmd.PersistentFlags().GetString("export-file")

			return statUsecases.Create(options)
		},
	}

	cmd.Flags().StringP("config", "c", "", "Конфигурационный файл")
	cmd.PersistentFlags().StringP("from-date", "f", "", "Начало периода")
	cmd.PersistentFlags().StringP("to-date", "t", "", "Окончание периода")
	cmd.PersistentFlags().StringP("export-file", "e", "", "Файл для экспорта данных")

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
		Short:   "Получить незагруженные отчеты",
		Example: "ozonadv stat:pull",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.StatUsecases().Pull()
		},
	}

	rootCmd.AddCommand(cmd)
}
