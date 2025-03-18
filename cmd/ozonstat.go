package main

import (
	"fmt"
	"os"
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
		Short: "CLI версия приложения",
	}

	initFetchCommand(rootCmd, ozonClient)

	fmt.Println("")
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initFetchCommand(rootCmd *cobra.Command, ozonClient *ozon.Client) {
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Извлечь данные по кампаниям",
		RunE: func(cmd *cobra.Command, args []string) error {
			options := stat.FetchOptions{}
			options.Days, _ = cmd.Flags().GetUint("days")

			return stat.Fetch(ozonClient, options)
		},
	}

	cmd.Flags().Uint("days", 0, "Интервал дней")

	rootCmd.AddCommand(cmd)
}
