package entities

import "time"

type Item struct {
	Id          int64     `json:"item_id"`
	CampaignId  int64     `json:"campaign_id"`
	Name        string    `json:"item_name"`
	Description string    `json:"item_description"`
	Priority    int64     `json:"item_priority"`
	Removed     bool      `json:"item_removed"`
	CreatedAt   time.Time `json:"item_created_at"`
}
