package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"ozonadv/internal/application"
	"ozonadv/internal/stat"
	"ozonadv/pkg/console"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(0)
	defer fmt.Println("")

	app := application.New()
	defer app.Shutdown()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		app.Shutdown()
		os.Exit(1)
	}()

	rootCmd := &cobra.Command{
		Use:   "ozonadv",
		Short: "Консольное приложение выгрузки статистики рекламных кабинетов Озон",
	}

	initFindCampaignsCommand(rootCmd, app)
	initStatCommand(rootCmd, app)
	initStatInfoCommand(rootCmd, app)
	initStatResetCommand(rootCmd, app)

	fmt.Println("")
	rootCmd.Execute()
}

func initFindCampaignsCommand(rootCmd *cobra.Command, app *application.Application) {
	cmd := &cobra.Command{
		Use:     "find:campaigns",
		Short:   "Поиск кампаний в Озон",
		Example: "ozonadv find:campaigns",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cmd.Short)
			fmt.Println("")

			findUsecases := app.FindUsecases()
			fmt.Println("")

			return findUsecases.Campaigns()
		},
	}

	rootCmd.AddCommand(cmd)
}

func initStatCommand(rootCmd *cobra.Command, app *application.Application) {
	cmd := &cobra.Command{
		Use:     "stat",
		Short:   "Формирование и загрузка статистики по кампаниям",
		Example: "ozonadv stat --date-from 2025-01-01 --date-to 2025-01-02",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cmd.Short)
			fmt.Println("")

			statUsecases := app.StatUsecases()
			fmt.Println("")

			if statUsecases.HasIncompleteProcessing() {
				fmt.Println("Найдена незавершенная обработка кампаний")
				if console.Ask("Продолжить ее?") == true {
					fmt.Println("")
					return statUsecases.StatContinue()
				}
			}

			fmt.Println("")
			options := stat.StatOptions{}
			options.DateFrom, _ = cmd.PersistentFlags().GetString("date-from")
			options.DateTo, _ = cmd.PersistentFlags().GetString("date-to")
			options.CampaignId, _ = cmd.Flags().GetString("campaign-id")
			options.ExportFile, _ = cmd.PersistentFlags().GetString("export-file")
			options.GroupBy = "DATE"

			return statUsecases.StatNew(options)
		},
	}

	cmd.Flags().StringP("config", "c", "", "Конфигурационный файл")
	cmd.PersistentFlags().StringP("date-from", "f", "", "Начало периода, обязательный")
	cmd.PersistentFlags().StringP("date-to", "t", "", "Окончание периода, обязательный")
	cmd.Flags().StringP("campaign-id", "i", "", "ID кампании")
	cmd.PersistentFlags().StringP("export-file", "e", "", "Файл для экспорта данных, обязательный")

	rootCmd.AddCommand(cmd)
}

func initStatInfoCommand(rootCmd *cobra.Command, app *application.Application) {
	cmd := &cobra.Command{
		Use:     "stat:info",
		Short:   "Статус формирования статистики",
		Example: "ozonadv stat:info",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cmd.Short)
			fmt.Println("")

			statUsecases := app.StatUsecases()
			fmt.Println("")

			return statUsecases.StatInfo()
		},
	}

	rootCmd.AddCommand(cmd)
}

func initStatResetCommand(rootCmd *cobra.Command, app *application.Application) {
	cmd := &cobra.Command{
		Use:     "stat:reset",
		Short:   "Удалить незавершенное формирование статистики",
		Example: "ozonadv stat:reset",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cmd.Short)
			fmt.Println("")

			statUsecases := app.StatUsecases()
			fmt.Println("")

			return statUsecases.StatReset()
		},
	}

	rootCmd.AddCommand(cmd)
}
