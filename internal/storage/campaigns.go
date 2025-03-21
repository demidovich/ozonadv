package storage

import (
	"ozonadv/internal/ozon"
	"ozonadv/pkg/utils"
)

type campaigns struct {
	data map[string]ozon.Campaign
}

func NewCampaigns(file string) *campaigns {
	instance := &campaigns{
		data: make(map[string]ozon.Campaign),
	}
	utils.JsonFileReadOrFail(file, &instance.data, "{}")

	return instance
}

func (d *campaigns) Add(item ozon.Campaign) {
	d.data[item.ID] = item
}

func (d *campaigns) Has(id string) bool {
	_, ok := d.data[id]
	return ok
}

func (d *campaigns) Next() (ozon.Campaign, bool) {
	for _, item := range d.data {
		return item, true
	}

	return ozon.Campaign{}, false
}

func (d *campaigns) Remove(id string) {
	delete(d.data, id)
}

func (d *campaigns) RemoveAll() {
	for id, _ := range d.data {
		delete(d.data, id)
	}
}

func (d *campaigns) All() map[string]ozon.Campaign {
	return d.data
}

func (d *campaigns) Size() int {
	return len(d.data)
}
