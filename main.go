package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"ozonadv/internal/app"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(0)
	defer fmt.Println("")

	app := app.New()
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

	// initCampaignsCommand(rootCmd, app)

	// initStatCommand(rootCmd, app)
	// initStatContinueCommand(rootCmd, app)
	// initStatInfoCommand(rootCmd, app)
	// initStatExportCommand(rootCmd, app)
	// initStatResetCommand(rootCmd, app)

	// initObjectStatCommand(rootCmd, app)
	// initObjectStatContinueCommand(rootCmd, app)
	// initObjectStatInfoCommand(rootCmd, app)
	// initObjectStatExportCommand(rootCmd, app)
	// initObjectStatResetCommand(rootCmd, app)

	fmt.Println("")
	rootCmd.Execute()
}

// func initCampaignsCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "campaigns",
// 		Short:   "Поиск кампаний",
// 		Example: "ozonadv campaigns",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println(cmd.Short)
// 			fmt.Println("")

// 			campaignsUsecases := app.CampaignsUsecases()
// 			fmt.Println("")

// 			return campaignsUsecases.Select()
// 		},
// 	}

// 	rootCmd.AddCommand(cmd)
// }

// func initStatCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "stat",
// 		Short:   "Сформировать общую статистику по кампаниям",
// 		Example: "ozonadv stat --date-from 2025-01-01 --date-to 2025-01-02",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println("Формирование статистики по кампаниям")
// 			fmt.Println("")

// 			statUsecases := app.StatUsecases()

// 			fmt.Println("")
// 			options := stat.StatOptions{}
// 			options.DateFrom, _ = cmd.PersistentFlags().GetString("date-from")
// 			options.DateTo, _ = cmd.PersistentFlags().GetString("date-to")
// 			options.CampaignId, _ = cmd.Flags().GetString("campaign-id")
// 			options.GroupBy = "DATE"

// 			return statUsecases.StatNew(options)
// 		},
// 	}

// 	cmd.Flags().StringP("config", "c", "", "Конфигурационный файл")
// 	cmd.PersistentFlags().StringP("date-from", "f", "", "Начало периода, обязательный")
// 	cmd.PersistentFlags().StringP("date-to", "t", "", "Окончание периода, обязательный")
// 	cmd.Flags().StringP("campaign-id", "i", "", "ID кампании")

// 	rootCmd.AddCommand(cmd)
// }

// func initStatContinueCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "stat:continue",
// 		Short:   "Продолжить прерваное формирования статистики",
// 		Example: "ozonadv stat:continue",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println("Продолжение прерваного формирования статистики")
// 			fmt.Println("")

// 			statUsecases := app.StatUsecases()
// 			fmt.Println("")

// 			return statUsecases.StatContinue()
// 		},
// 	}

// 	rootCmd.AddCommand(cmd)
// }

// func initStatInfoCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "stat:info",
// 		Short:   "Состояние формирования статистики",
// 		Example: "ozonadv stat:info",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println(cmd.Short)
// 			fmt.Println("")

// 			statUsecases := app.StatUsecases()
// 			fmt.Println("")

// 			return statUsecases.StatInfo()
// 		},
// 	}

// 	rootCmd.AddCommand(cmd)
// }

// func initStatExportCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "stat:export",
// 		Short:   "Экспорт сформированной статистики",
// 		Example: "ozonadv stat:export --file ./export_file.csv",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println(cmd.Short)
// 			fmt.Println("")

// 			options := stat.StatExportOptions{}
// 			options.File, _ = cmd.PersistentFlags().GetString("file")

// 			statUsecases := app.StatUsecases()
// 			fmt.Println("")

// 			return statUsecases.StatExport(options)
// 		},
// 	}

// 	cmd.PersistentFlags().StringP("file", "f", "", "Файл экспорта")
// 	rootCmd.AddCommand(cmd)
// }

// func initStatResetCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "stat:reset",
// 		Short:   "Удалить незавершенное формирование статистики",
// 		Example: "ozonadv stat:reset",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println("Удаление незавершенного формирования статистики")
// 			fmt.Println("")

// 			statUsecases := app.StatUsecases()
// 			fmt.Println("")

// 			return statUsecases.StatReset()
// 		},
// 	}

// 	rootCmd.AddCommand(cmd)
// }

// func initObjectStatCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "object-stat",
// 		Short:   "Сформировать статистику по рекламным объектам",
// 		Example: "ozonadv object-stat --date-from 2025-01-01 --date-to 2025-01-02",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println("Формирование статистики по рекламным объектам")
// 			fmt.Println("")

// 			objectStatUsecases := app.ObjectStatUsecases()

// 			fmt.Println("")
// 			options := object_stat.StatOptions{}
// 			options.DateFrom, _ = cmd.PersistentFlags().GetString("date-from")
// 			options.DateTo, _ = cmd.PersistentFlags().GetString("date-to")
// 			options.CampaignId, _ = cmd.Flags().GetString("campaign-id")
// 			options.GroupBy = "DATE"

// 			return objectStatUsecases.StatNew(options)
// 		},
// 	}

// 	cmd.Flags().StringP("config", "c", "", "Конфигурационный файл")
// 	cmd.PersistentFlags().StringP("date-from", "f", "", "Начало периода, обязательный")
// 	cmd.PersistentFlags().StringP("date-to", "t", "", "Окончание периода, обязательный")
// 	cmd.Flags().StringP("campaign-id", "i", "", "ID кампании")

// 	rootCmd.AddCommand(cmd)
// }

// func initObjectStatContinueCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "object-stat:continue",
// 		Short:   "Продолжить прерваное формирования статистики",
// 		Example: "ozonadv object-stat:continue",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println("Продолжение прерваного формирования статистики")
// 			fmt.Println("")

// 			statUsecases := app.ObjectStatUsecases()
// 			fmt.Println("")

// 			return statUsecases.StatContinue()
// 		},
// 	}

// 	rootCmd.AddCommand(cmd)
// }

// func initObjectStatInfoCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "object-stat:info",
// 		Short:   "Состояние формирования статистики по рекламным объектам",
// 		Example: "ozonadv stat:info",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println(cmd.Short)
// 			fmt.Println("")

// 			statUsecases := app.ObjectStatUsecases()
// 			fmt.Println("")

// 			return statUsecases.StatInfo()
// 		},
// 	}

// 	rootCmd.AddCommand(cmd)
// }

// func initObjectStatExportCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "object-stat:export",
// 		Short:   "Экспорт сформированной статистики",
// 		Example: "ozonadv object-stat:export --file ./export_file.csv",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println(cmd.Short)
// 			fmt.Println("")

// 			options := object_stat.StatExportOptions{}
// 			options.File, _ = cmd.PersistentFlags().GetString("file")

// 			statUsecases := app.ObjectStatUsecases()
// 			fmt.Println("")

// 			return statUsecases.StatExport(options)
// 		},
// 	}

// 	cmd.PersistentFlags().StringP("file", "f", "", "Файл экспорта")
// 	rootCmd.AddCommand(cmd)
// }

// func initObjectStatResetCommand(rootCmd *cobra.Command, app *app.Application) {
// 	cmd := &cobra.Command{
// 		Use:     "object-stat:reset",
// 		Short:   "Удалить незавершенное формирование статистики",
// 		Example: "ozonadv object-stat:reset",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println("Удаление незавершенного формирования статистики")
// 			fmt.Println("")

// 			statUsecases := app.ObjectStatUsecases()
// 			fmt.Println("")

// 			return statUsecases.StatReset()
// 		},
// 	}

// 	rootCmd.AddCommand(cmd)
// }
