package storage

import (
	"ozonadv/internal/ozon"
	"ozonadv/pkg/utils"
)

type statCampaigns struct {
	data map[string]ozon.Campaign
}

func NewStatCampaigns(file string) *statCampaigns {
	instance := &statCampaigns{
		data: make(map[string]ozon.Campaign),
	}
	utils.JsonFileReadOrFail(file, &instance.data, "{}")

	return instance
}

func (c *statCampaigns) Add(item ozon.Campaign) {
	c.data[item.ID] = item
}

func (c *statCampaigns) Has(id string) bool {
	_, ok := c.data[id]
	return ok
}

func (c *statCampaigns) ByStatRequestUUID(uuid string) (ozon.Campaign, bool) {
	for _, c := range c.data {
		if c.Stat.RequestUUID == uuid {
			return c, true
		}
	}
	return ozon.Campaign{}, false
}

func (c *statCampaigns) RemoveAll() {
	for id := range c.data {
		delete(c.data, id)
	}
}

func (c *statCampaigns) All() []ozon.Campaign {
	result := make([]ozon.Campaign, 0, len(c.data))
	for _, c := range c.data {
		result = append(result, c)
	}

	return result
}

func (c *statCampaigns) Data() map[string]ozon.Campaign {
	return c.data
}

func (c *statCampaigns) Size() int {
	return len(c.data)
}
