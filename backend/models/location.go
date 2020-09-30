package models

type Location struct {
	ID      uint   `gorm:"primary_key" json:"-"`
	Country string `gorm:"unique"`
	State   string `gorm:"unique"`
	City    string `gorm:"unique"`
}

func GetLocation(location *Location) error {
	res := conn.FirstOrCreate(&location)
	return res.Error
}
