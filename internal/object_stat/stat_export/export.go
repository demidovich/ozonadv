package stat_export

import (
	"fmt"
	"log"
	"os"
	"ozonadv/internal/ozon"
	"ozonadv/internal/storage"
	"strings"
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
	file = strings.TrimSuffix(file, ".csv")
	bannersFname := file + "-banner.csv"
	videosFname := file + "-video.csv"

	fmt.Println("")

	banners, err := os.Create(bannersFname)
	if err != nil {
		return err
	}
	defer banners.Close()

	videos, err := os.Create(videosFname)
	if err != nil {
		return err
	}
	defer videos.Close()

	campaigns := s.storage.ObjectStatCampaigns().All()

	fmt.Fprintln(videos, s.csvHeaders(&campaigns, "video"))
	fmt.Fprintln(banners, s.csvHeaders(&campaigns, "banner"))

	for _, campaign := range campaigns {
		logCampaign(campaign, " обработка кампании: "+campaign.Title)

		if campaign.Stat.File == "" {
			logCampaign(campaign, " пропуск: отсутствует скачаный файл")
			continue
		}

		rows := s.csvRows(campaign.Stat.File)
		if len(rows) == 0 {
			logCampaign(campaign, " пропуск: нет csv строк с данными")
		}

		s.appendCsvCampaignId(&rows, campaign)

		for _, r := range rows {
			if campaign.IsVideo() {
				fmt.Fprintln(videos, r)
			} else {
				fmt.Fprintln(banners, r)
			}
		}
	}

	log.Println("")
	log.Printf("Создан файл: %s\n", bannersFname)
	log.Printf("Создан файл: %s\n", videosFname)

	return nil
}

func (s statExport) csvHeaders(campaigns *[]ozon.Campaign, campaignsType string) string {
	var headers string

	for _, c := range *campaigns {
		if campaignsType == "video" && !c.IsVideo() {
			continue
		}

		if campaignsType == "banner" && c.IsVideo() {
			continue
		}

		if c.Stat.File == "" {
			continue
		}

		csvLines := s.fileLines(c.Stat.File)
		if len(csvLines) < 2 {
			continue
		}

		headers = "ID кампании;" + csvLines[1]
		break
	}

	return headers
}

func (s statExport) csvRows(file string) []string {
	lines := s.fileLines(file)
	if len(lines) < 4 {
		return []string{}
	}

	// Первая строка общие данные об отчете
	// Вторая строка заголовки
	// Последняя строка суммарная информация

	return lines[2 : len(lines)-1]
}

func (s statExport) fileLines(file string) []string {
	csv, err := s.storage.Downloads().ReadString(file)
	if err != nil {
		log.Fatal(err)
	}

	csv = strings.TrimSpace(csv)

	return strings.Split(string(csv), "\n")
}

func (s statExport) appendCsvCampaignId(rows *[]string, campaign ozon.Campaign) {
	for i := range *rows {
		(*rows)[i] = campaign.ID + ";" + (*rows)[i]
	}
}

func logCampaign(c ozon.Campaign, msg ...any) {
	fmt.Printf("[%s] %s\n", c.ID, fmt.Sprint(msg...))
}
