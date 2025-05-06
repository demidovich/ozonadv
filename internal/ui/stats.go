package ui

import (
	"errors"
	"fmt"

	"time"

	"github.com/demidovich/ozonadv/internal/cabinets"
	"github.com/demidovich/ozonadv/internal/models"
	"github.com/demidovich/ozonadv/internal/stats"
	"github.com/demidovich/ozonadv/internal/ui/colors"
	"github.com/demidovich/ozonadv/internal/ui/forms"
	"github.com/demidovich/ozonadv/internal/ui/forms/validators"
	"github.com/demidovich/ozonadv/internal/ui/helpers"

	"github.com/charmbracelet/huh"
	"github.com/jedib0t/go-pretty/v6/table"
)

type statsPage struct {
	cabsService  *cabinets.Service
	statsService *stats.Service
}

func newStats(cabsService *cabinets.Service, statsService *stats.Service) statsPage {
	return statsPage{
		cabsService:  cabsService,
		statsService: statsService,
	}
}

func (c statsPage) Home() error {
	options := []helpers.ListOption{}
	for _, stat := range c.statsService.All() {
		options = append(options, helpers.ListOption{
			Key:   stat.Options.Name + " " + colors.Gray().Sprintf("(%s)", stat.Options.CabinetName),
			Value: stat.UUID,
		})
	}

	options = append(
		options,
		helpers.ListOption{Key: "Добавить отчет", Value: "create_stat"},
		helpers.ListOption{Key: "Назад", Value: "back"},
	)

	fmt.Println("")
	action, err := helpers.List("Отчеты", options...)
	if err != nil {
		return err
	}

	var stat *models.Stat
	var ok bool

	switch action {
	case "create_stat":
		cabinet, err := c.chooseCabinetForm()
		if isFormCanceled(err) {
			return c.Home()
		}
		if err != nil {
			printError(err)
			return c.Home()
		}
		fmt.Println("")
		stat, err = c.createStat(*cabinet)
		if err == nil || isFormCanceled(err) {
			return c.Home()
		}
	case "back":
		return ErrGoBack
	default:
		statUUID := action
		if stat, ok = c.statsService.Find(statUUID); !ok {
			err = errors.New("отчет не найден")
		}
	}

	if err != nil {
		printError(err)
		return c.Home()
	}

	err = c.stat(stat)
	if isGoBack(err) {
		return c.Home()
	}

	return err
}

func (c statsPage) stat(stat *models.Stat) error {
	fmt.Println("Параметры отчета")
	c.printStatTable(stat)
	fmt.Println("")

	options := []helpers.ListOption{
		{Key: "Кампании", Value: "campaigns"},
		{Key: "Загрузка", Value: "download"},
		{Key: "Экспорт", Value: "export"},
		{Key: "Удалить", Value: "remove"},
		{Key: "Назад", Value: "back"},
	}

	action, err := helpers.List("Отчеты > "+stat.Options.Name, options...)
	if err != nil {
		return err
	}

	switch action {
	case "campaigns":
		c.statCampaigns(stat)
		err = c.stat(stat)
	case "download":
		if helpers.Confirm("Запустить загрузку отчета?") {
			c.statsService.Download(stat)
		}
		return c.stat(stat)
	case "export":
		var file string
		file, err = c.statExportForm()
		if err == nil {
			err = c.statsService.ExportToFile(stat, file)
		}
	case "remove":
		if helpers.Confirm("Удалить отчет \"" + stat.Options.Name + "\"?") {
			c.statsService.Remove(stat)
			err = ErrGoBack
			fmt.Println("Отчет \"" + stat.Options.Name + "\" удален.")
		} else {
			err = c.stat(stat)
		}
	case "back":
		err = ErrGoBack
	}

	return err
}

func (c statsPage) statCampaigns(stat *models.Stat) {
	fmt.Println("Кампании отчета")
	c.printStatCampaignsTable(stat)
	fmt.Println("")

	helpers.WaitButton("Назад")
}

func (c statsPage) CabinetStats(cabinet models.Cabinet) error {
	options := []helpers.ListOption{}
	for _, stat := range c.statsService.CabinetAll(cabinet) {
		options = append(options, helpers.ListOption{
			Key:   stat.Options.Name,
			Value: stat.UUID,
		})
	}

	options = append(
		options,
		helpers.ListOption{Key: "Добавить отчет", Value: "create_stat"},
		helpers.ListOption{Key: "Назад", Value: "back"},
	)

	action, err := helpers.List("Кабинеты > "+cabinet.Name+" > Отчеты", options...)
	if err != nil {
		return err
	}

	var stat *models.Stat
	var ok bool

	switch action {
	case "create_stat":
		stat, err = c.createStat(cabinet)
		if isFormCanceled(err) {
			fmt.Println("")
			return c.CabinetStats(cabinet)
		}
	case "back":
		return ErrGoBack
	default:
		statUUID := action
		if stat, ok = c.statsService.Find(statUUID); !ok {
			err = errors.New("отчет не найден")
		}
	}

	if err != nil {
		return err
	}

	err = c.stat(stat)
	if errors.Is(err, ErrGoBack) {
		return c.CabinetStats(cabinet)
	}

	return nil
}

func (c statsPage) chooseCabinetForm() (*models.Cabinet, error) {
	if len(c.cabsService.All()) == 0 {
		return nil, errors.New("Нет рекламных кабинетов.\nДля создания отчета необходимо создать рекламный кабинет")
	}

	fmt.Println("Выбор кабинета")

	options := []huh.Option[string]{}
	for _, cabinet := range c.cabsService.All() {
		options = append(options, huh.NewOption(
			cabinet.Name+colors.Gray().Sprintf(" (%s)", cabinet.ClientID),
			cabinet.UUID,
		))
	}

	var cabinetUUID string
	confirm := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Кабинет").
				Options(options...).
				Validate(validators.Required).
				Value(&cabinetUUID),
			huh.NewConfirm().
				Key("done").
				Value(&confirm).
				Inline(true).
				Affirmative("Ок").
				Negative("Отмена"),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	if !confirm {
		return nil, ErrFormCancel
	}

	cabinet, ok := c.cabsService.Find(cabinetUUID)
	if !ok {
		return nil, errors.New("кабинет не найден")
	}

	return cabinet, nil
}

func (c statsPage) createStat(cabinet models.Cabinet) (*models.Stat, error) {
	fmt.Println(cabinet.Name + " > Новый отчет")

	options := models.StatOptions{}
	if err := c.statOptionsForm(&options); err != nil {
		return nil, err
	}

	options.CabinetUUID = cabinet.UUID
	options.CabinetName = cabinet.Name
	options.CabinetClientID = cabinet.ClientID
	options.CabinetClientSecret = cabinet.ClientSecret

	// campaigns, err := chooseCampaignsForm(*c.cabsService, cabinet)
	campaigns, err := c.createStatCampaigns(options, cabinet)
	if err != nil {
		return nil, err
	}

	printCampaignsTable(campaigns)

	fmt.Println("")
	if !helpers.Confirm("Создать отчет?") {
		return nil, ErrFormCancel
	}

	stat, err := c.statsService.Create(options, campaigns)
	if err != nil {
		return stat, err
	}

	fmt.Println("Отчет создан")

	return stat, err
}

// Запрос кампаний кабинета, ограничение их по интервалу отчета
func (c statsPage) createStatCampaigns(options models.StatOptions, cabinet models.Cabinet) ([]models.Campaign, error) {
	campaigns, err := chooseCampaignsForm(*c.cabsService, cabinet)
	if err != nil {
		return campaigns, err
	}

	statDateFrom, _ := time.Parse("2006-01-02", options.DateFrom)
	statDateTo, _ := time.Parse("2006-01-02", options.DateTo)

	inInterval := []models.Campaign{}
	for _, c := range campaigns {
		campDateTo, err := time.Parse("2006-01-02", c.ToDate)
		if err == nil && statDateFrom.After(campDateTo) {
			continue
		}

		campDateFrom, err := time.Parse("2006-01-02", c.FromDate)
		if err == nil && statDateTo.Before(campDateFrom) {
			continue
		}

		inInterval = append(inInterval, c)
	}

	return inInterval, nil
}

func (c statsPage) statOptionsForm(options *models.StatOptions) error {
	confirm := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(forms.RequiredTitle("Название отчета")).
				CharLimit(100).
				Validate(validators.Required).
				Value(&options.Name),
			huh.NewSelect[string]().
				Title(forms.RequiredTitle("Тип статистики")).
				Options(
					huh.NewOption("Рекламные объекты", "OBJECT"),
					huh.NewOption("Рекламные кампании", "TOTAL"),
				).
				Validate(validators.Required).
				Value(&options.Type),
			huh.NewInput().
				Title(forms.RequiredTitle("Начало интервала, дата")).
				Placeholder("ГГГГ-ММ-ДД").
				CharLimit(10).
				Validate(validators.DateRequiured).
				Value(&options.DateFrom),
			huh.NewInput().
				Title(forms.RequiredTitle("Конец интервала, дата")).
				Placeholder("ГГГГ-ММ-ДД").
				CharLimit(10).
				Validate(validators.DateRequiured).
				Value(&options.DateTo),
			huh.NewSelect[string]().
				Title(forms.RequiredTitle("Группировка")).
				Options(
					huh.NewOption("Не группировать", "NO_GROUP_BY"),
					huh.NewOption("День", "DATE"),
					huh.NewOption("Неделя", "START_OF_WEEK"),
					huh.NewOption("Месяц", "START_OF_MONTH"),
				).
				Validate(validators.Required).
				Value(&options.GroupBy),
			huh.NewConfirm().
				Key("done").
				Value(&confirm).
				Inline(true).
				Negative("Отмена").
				Affirmative("Далее"),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if !confirm {
		return ErrFormCancel
	}

	return nil
}

func (c statsPage) statExportForm() (string, error) {
	file := ""

	confirm := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(forms.RequiredTitle("Файл")).
				Description("Название файла или путь").
				CharLimit(1000).
				Value(&file),
			huh.NewConfirm().
				Key("done").
				Value(&confirm).
				Inline(true).
				Negative("Отмена").
				Affirmative("Экспортировать"),
		),
	)

	if err := form.Run(); err != nil {
		return file, err
	}

	if !confirm {
		return file, ErrFormCancel
	}

	return file, nil
}

func (c statsPage) printStatTable(stat *models.Stat) {
	if stat == nil {
		return
	}

	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"Отчет", stat.Options.Name})
	tw.AppendRow(table.Row{"Кабинет", stat.Options.CabinetName})
	tw.AppendRow(table.Row{"Начало интервала, дата", stat.Options.DateFrom})
	tw.AppendRow(table.Row{"Конец интервала, дата", stat.Options.DateTo})
	tw.AppendRow(table.Row{"Группировка", stat.Options.GroupBy})
	tw.AppendRow(table.Row{"Кампаний всего", len(stat.Items)})
	tw.AppendRow(table.Row{"Кампаний обработано", len(stat.CampaignsCompleted())})
	tw.AppendRow(table.Row{"Состояние", stat.StateHuman()})

	fmt.Println(tw.Render())
}

func (c statsPage) printStatCampaignsTable(stat *models.Stat) {
	if stat == nil {
		return
	}

	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.AppendRow(table.Row{"#", "Тип", "Кампания", "Запуск", "Окончание", "Состояние отчета"})
	tw.AppendRow(table.Row{"", "", "", "", "", ""})

	for _, item := range stat.Items {
		tw.AppendRow(table.Row{
			item.Campaign.ID,
			item.Campaign.AdvObjectType,
			item.Campaign.TitleTruncated(45),
			item.Campaign.FromDate,
			item.Campaign.ToDate,
			item.State(),
		})
	}

	fmt.Println(tw.Render())
	fmt.Println("Всего кампаний:", len(stat.Items))
}
