package stat_export

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
	"ozonadv/pkg/utils"
)

type statExport struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) statExport {
	return statExport{
		storage: storage,
	}
}

func (s statExport) ToFile(file string) error {
	summaryStat := []StatRow{}
	for _, campaign := range s.storage.StatCampaigns().All() {
		err := s.addCampaignStat(&summaryStat, campaign)
		if err != nil {
			logCampaign(campaign, err)
		}
	}

	if len(summaryStat) == 0 {
		fmt.Println("В отчетах нет данных")
		return nil
	}

	fmt.Println("")
	fmt.Printf("Строк в отчете: %d\n", len(summaryStat))

	err := s.writeCsvFile(file, &summaryStat)
	if err != nil {
		return err
	}

	fmt.Printf("Создан csv файл: %s\n", file)

	return nil
}

func (s statExport) writeCsvFile(filepath string, summaryStat *[]StatRow) error {
	if len(*summaryStat) == 0 {
		return errors.New("нет данных")
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	err = csvWriter.Write(statRowCsvHeaders())
	if err != nil {
		return err
	}

	for _, row := range *summaryStat {
		err = csvWriter.Write(statRowCsvValues(row))
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}

func (s statExport) addCampaignStat(summaryStat *[]StatRow, campaign ozon.Campaign) error {
	if campaign.Stat.File == "" {
		return errors.New("пропуск: нет файла статистики")
	}

	file := s.storage.Downloads().AbsolutePath(campaign.Stat.File)
	all := Stat{}
	err := utils.JsonFileRead(file, &all, "{}")
	if err != nil {
		return errors.New("пропуск: ошибка парсинга " + err.Error())
	}

	stat, ok := all[campaign.ID]
	if !ok {
		return errors.New("пропуск: в файле отсутствует кампания " + campaign.ID)
	}

	for _, row := range stat.Report.Rows {
		row.CampaignId = campaign.ID
		row.CampaignType = campaign.AdvObjectType
		row.CampaignTitle = campaign.Title
		*summaryStat = append(*summaryStat, row)
	}

	return nil
}

func logCampaign(c ozon.Campaign, msg ...any) {
	fmt.Printf("[%s] %s\n", c.ID, fmt.Sprint(msg...))
}
