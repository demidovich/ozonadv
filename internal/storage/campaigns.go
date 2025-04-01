package storage

import (
	"ozonadv/internal/ozon"
	"ozonadv/pkg/utils"
)

type campaigns struct {
	data map[string]ozon.Campaign
}

func newCampaigns(file string) *campaigns {
	data := make(map[string]ozon.Campaign)
	utils.JsonFileReadOrFail(file, &data, "{}")

	instance := campaigns{data: data}

	return &instance
}

func (c *campaigns) Add(item ozon.Campaign) {
	c.data[item.ID] = item
}

func (c *campaigns) Has(id string) bool {
	_, ok := c.data[id]
	return ok
}

func (c *campaigns) ByStatUUID(uuid string) (ozon.Campaign, bool) {
	for _, c := range c.data {
		if c.Stat.UUID == uuid {
			return c, true
		}
	}
	return ozon.Campaign{}, false
}

func (c *campaigns) RemoveAll() {
	for id := range c.data {
		delete(c.data, id)
	}
}

func (c *campaigns) All() []ozon.Campaign {
	result := make([]ozon.Campaign, 0, len(c.data))
	for _, c := range c.data {
		result = append(result, c)
	}

	return result
}

func (c *campaigns) Data() map[string]ozon.Campaign {
	return c.data
}

func (c *campaigns) Size() int {
	return len(c.data)
}
