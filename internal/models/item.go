package models

type Item struct {
	ItemUID     string `db:"-"`
	ChrtID      int    `db:"chrt_id"`
	TrackNumber string `db:"track_number"`
	RID         string `db:"rid"`
	Name        string `db:"name"`
	Brand       string `db:"brand"`
	Size        string `db:"size"`
	NmID        int    `db:"nm_id"`
	Status      int    `db:"status"`
}
