package stat_processor

import (
	"fmt"
	"ozonadv/internal/ozon"
	"unicode/utf8"

	"github.com/jedib0t/go-pretty/v6/table"
)

type procresults struct {
	entries map[string]procresultsEntry
}

func newProcresults() *procresults {
	return &procresults{
		entries: make(map[string]procresultsEntry),
	}
}

func (p *procresults) RegisterCampaign(c ozon.Campaign) {
	p.entries[c.ID] = procresultsEntry{
		campaignID:    c.ID,
		campaignTitle: c.Title,
	}
}

func (p *procresults) RegisterStatRequest(s ozon.StatRequest) {
	entry := p.entryByStatRequest(s)
	entry.statRequestUUID = s.UUID
	entry.statRequestState = s.State

	p.entries[entry.campaignID] = entry
}

func (p *procresults) RegisterStatDownload(s ozon.StatRequest, fname string) {
	entry := p.entryByStatRequest(s)
	entry.statRequestUUID = s.UUID
	entry.statRequestState = s.State
	entry.downloadedFile = fname

	p.entries[entry.campaignID] = entry
}

func (p *procresults) PrintSummaryTable() {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"#", "Название", "Статус", "Файл"})
	tw.AppendRow(table.Row{"", "", "", "", ""})
	downloadFails := 0
	for _, e := range p.entries {
		tw.AppendRow(table.Row{
			e.campaignID,
			e.campaignTitleShorten(45),
			e.statRequestUUID,
			e.statRequestState,
			e.downloadedFile,
		})
		if e.downloadedFile == "" {
			downloadFails++
		}
	}

	fmt.Println("[shutdown] результаты выполнения")
	fmt.Println(tw.Render())
	fmt.Printf("Всего: %d, Не загружено: %d\n", len(p.entries), downloadFails)
	// fmt.Println("Всего:", len(p.entries))
	// fmt.Println("Не загружено:", downloadFails)
}

func (p *procresults) entryByStatRequest(s ozon.StatRequest) procresultsEntry {
	entry, ok := p.entries[s.Request.CampaignId]
	if !ok {
		entry = procresultsEntry{
			campaignID:    s.Request.CampaignId,
			campaignTitle: "undefined",
		}
	}

	return entry
}

type procresultsEntry struct {
	campaignID       string
	campaignTitle    string
	statRequestUUID  string
	statRequestState string
	downloadedFile   string
}

func (p procresultsEntry) campaignTitleShorten(maxlen int) string {
	if utf8.RuneCountInString(p.campaignTitle) <= maxlen+3 {
		return p.campaignTitle
	}

	r := []rune(p.campaignTitle)

	return string(r[:maxlen+2]) + "..."
}
