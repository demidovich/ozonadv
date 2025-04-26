package stats

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/demidovich/ozonadv/internal/models"
)

type objectStatExport struct {
	storage Storage
	debug   Debug
	stat    *models.Stat
}

func newObjectStatExport(storage Storage, debug Debug, stat *models.Stat) objectStatExport {
	return objectStatExport{
		storage: storage,
		debug:   debug,
		stat:    stat,
	}
}

func (s objectStatExport) toFile(file string) error {
	file = strings.TrimSuffix(file, ".csv")

	bannersFile := filepath.Clean(file + "-banners.csv")
	videoBannersFile := filepath.Clean(file + "-video-banners.csv")
	globalPromoFile := filepath.Clean(file + "-global-promo.csv")

	s.debug.Println("")

	banners, err := os.Create(bannersFile)
	if err != nil {
		return err
	}
	defer banners.Close()

	videoBanners, err := os.Create(videoBannersFile)
	if err != nil {
		return err
	}
	defer videoBanners.Close()

	globalPromo, err := os.Create(globalPromoFile)
	if err != nil {
		return err
	}
	defer globalPromo.Close()

	fmt.Fprintln(banners, s.csvHeaders("banner"))
	fmt.Fprintln(videoBanners, s.csvHeaders("video_banner"))
	fmt.Fprintln(globalPromo, s.csvHeaders("global_promo"))

	for _, item := range s.stat.Items {
		campaign := item.Campaign

		s.debugCampaign(campaign, " обработка кампании: "+campaign.Title)

		if item.Request.File == "" {
			s.debugCampaign(campaign, " пропуск: отсутствует скачаный файл")
			continue
		}

		rows := s.csvRows(item.Request.File)
		if len(rows) == 0 {
			s.debugCampaign(campaign, " пропуск: нет csv строк с данными")
			continue
		}

		s.appendCsvCampaignID(&rows, campaign)

		for _, r := range rows {
			if campaign.IsBanner() {
				fmt.Fprintln(banners, r)
			} else if campaign.IsVideoBanner() {
				fmt.Fprintln(videoBanners, r)
			} else if campaign.IsGlobalPromo() {
				fmt.Fprintln(globalPromo, r)
			} else {
				s.debugCampaign(campaign, " пропуск csv строки: неизвестный тип кампании: ", campaign.AdvObjectType)
			}
		}
	}

	s.debug.Println("")
	s.debug.Printf("Создан файл: %s\n", bannersFile)
	s.debug.Printf("Создан файл: %s\n", videoBannersFile)
	s.debug.Printf("Создан файл: %s\n", globalPromoFile)

	return nil
}

func (s objectStatExport) csvHeaders(campaignsType string) string {
	var headers string

	for _, item := range s.stat.Items {
		if campaignsType == "banner" && !item.Campaign.IsBanner() {
			continue
		}

		if campaignsType == "video_banner" && !item.Campaign.IsVideoBanner() {
			continue
		}

		if campaignsType == "global_promo" && !item.Campaign.IsGlobalPromo() {
			continue
		}

		if item.Request.File == "" {
			continue
		}

		csvLines := s.fileLines(item.Request.File)
		if len(csvLines) < 2 {
			continue
		}

		headers = "ID кампании;" + csvLines[1]
		break
	}

	return headers
}

func (s objectStatExport) csvRows(file string) []string {
	lines := s.fileLines(file)
	if len(lines) < 4 {
		return []string{}
	}

	// Первая строка общие данные об отчете
	// Вторая строка заголовки

	return lines[2:]
}

func (s objectStatExport) fileLines(file string) []string {
	data := s.storage.ReadDownloadsFile(s.stat, file)
	data = bytes.TrimSpace(data)

	return strings.Split(string(data), "\n")
}

func (s objectStatExport) appendCsvCampaignID(rows *[]string, campaign models.Campaign) {
	for i := range *rows {
		(*rows)[i] = campaign.ID + ";" + (*rows)[i]
	}
}

func (s objectStatExport) debugCampaign(c models.Campaign, msg ...any) {
	s.debug.Printf("[%s] %s\n", c.ID, fmt.Sprint(msg...))
}
