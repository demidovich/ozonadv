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

func (c *campaigns) Add(item ozon.Campaign) {
	c.data[item.ID] = item
}

func (c *campaigns) Has(id string) bool {
	_, ok := c.data[id]
	return ok
}

func (c *campaigns) Remove(id string) {
	delete(c.data, id)
}

func (c *campaigns) RemoveAll() {
	for id := range c.data {
		delete(c.data, id)
	}
}

func (c *campaigns) All() map[string]ozon.Campaign {
	return c.data
}

func (c *campaigns) Size() int {
	return len(c.data)
}
