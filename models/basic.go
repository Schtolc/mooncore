package models

// Address model
type Address struct {
	ID  int     `json:"id"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Photo model
type Photo struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
}

// User model
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Address   Address
	AddressID int `sql:"type:int, FOREIGN KEY (address_id) REFERENCES addresses(id)" json:"address"`
	Photo     Photo
	PhotoID   int `sql:"type:int, FOREIGN KEY (photo_id) REFERENCES photos(id)" json:"photo"`
}

// Master model
type Master struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Address            Address
	AddressID          int `sql:"type:int, FOREIGN KEY (address_id) REFERENCES addresses(id)" json:"address"`
	Photo              Photo
	PhotoID            int       `sql:"type:int, FOREIGN KEY (photo_id) REFERENCES photos(id)" json:"photo"`
	Service            []Service `gorm:"ForeignKey:MasterID"`
	WorkingPlacePhotos []Photo   `gorm:"many2many:working_place_photos;"`
}

// Service model
type Service struct {
	ID                int    `json:"id"`
	MasterID          int    `json:"master"`
	Name              string `json:"name"`
	Description       string `sql:"type:text" json:"description"`
	ManicureType      ManicureType
	ManicureTypeID    int                `sql:"type:int, FOREIGN KEY (manicure_type_id) REFERENCES manicure_types(id)" json:"manicure_type"`
	ManicureMaterials []ManicureMaterial `gorm:"many2many:service_manicure_materials;"`
	Photos            []Photo            `gorm:"many2many:service_photos;"`
}

// ManicureType model
type ManicureType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ManicureMaterial model
type ManicureMaterial struct {
	ID          int       `json:"id"`
	Firm        string    `json:"firm"`
	Palette     string    `json:"palette"`
	Description string    `sql:"type:text" json:"description"`
	Services    []Service `gorm:"many2many:service_manicure_materials;"`
}

// Response model
type Response struct {
	Code int         `json:"code"`
	Body interface{} `json:"body"`
}
