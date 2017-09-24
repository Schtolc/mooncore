package models

// Address model
type Address struct {
	ID  int
	Lat float32
	Lon float32
}

// Photo model
type Photo struct {
	ID   int
	Path string
}

// User model
type User struct {
	ID       int    `json:"-"`
	Name     string `gorm:"not null;unique; column:name"`
	Email    string `gorm:"not null;unique; column:email"`
	Password string `gorm:"not null;unique; column:password"`
}

// Client model
type Client struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Address   Address
	AddressID int `sql:"type:int, FOREIGN KEY (address_id) REFERENCES addresses(id)"`
	Photo     Photo
	PhotoID   int `sql:"type:int, FOREIGN KEY (photo_id) REFERENCES photos(id)"`
}

// Master model
type Master struct {
	ID                 int
	Name               string
	Address            Address
	AddressID          int `sql:"type:int, FOREIGN KEY (address_id) REFERENCES addresses(id)"`
	Photo              Photo
	PhotoID            int       `sql:"type:int, FOREIGN KEY (photo_id) REFERENCES photos(id)"`
	Service            []Service `gorm:"ForeignKey:MasterID"`
	WorkingPlacePhotos []Photo   `gorm:"many2many:working_place_photos;"`
}

// Service model
type Service struct {
	ID                int
	MasterID          int
	Name              string
	Description       string `sql:"type:text"`
	ManicureType      ManicureType
	ManicureTypeID    int                `sql:"type:int, FOREIGN KEY (manicure_type_id) REFERENCES manicure_types(id)"`
	ManicureMaterials []ManicureMaterial `gorm:"many2many:service_manicure_materials;"`
	Photos            []Photo            `gorm:"many2many:service_photos;"`
}

// ManicureType model
type ManicureType struct {
	ID   int
	Name string
}

// ManicureMaterial model
type ManicureMaterial struct {
	ID          int
	Firm        string
	Palette     string
	Description string    `sql:"type:text"`
	Services    []Service `gorm:"many2many:service_manicure_materials;"`
}
