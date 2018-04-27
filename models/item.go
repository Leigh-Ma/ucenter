package models

type Item struct {
	TCom
	PlayerId int64
	Name     string
	Category string
	Amount   int
}

func NewItem(playerId int64, sn string) *Item {
	return &Item{}
}

func (t *Item) TableName() string {
	return "items"
}
