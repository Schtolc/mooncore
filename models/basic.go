package models

// Address model
type Address struct {
	ID          int     `json:"id"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Description string  `json:"description"`
}

// Photo model
type Photo struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
	Tags []Tag  `gorm:"many2many:photo_tags"`
}

// Tag model
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// User model
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

// UserDetails model
type UserDetails struct {
	ID        int
	UserID    int       `sql:"type:int, FOREIGN KEY (user_id) REFERENCES users(id)" json:"user"`
	Name      string    `json:"name"`
	AddressID int       `sql:"type:int, FOREIGN KEY (address_id) REFERENCES addresses(id)" json:"address"`
	PhotoID   int       `sql:"type:int, FOREIGN KEY (photo_id) REFERENCES photos(id)" json:"photo"`
	Stars     float64   `json:"stars"`
	Services  []Service `gorm:"ForeignKey:MasterID"`
	Photos    []Photo   `gorm:"many2many:user_photos;"`
	Signs     []Sign    `gorm:"many2many:user_signs;"`
}

// Service model
type Service struct {
	ID             int     `json:"id"`
	MasterID       int     `json:"master"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Description    string  `sql:"type:text" json:"description"`
	ManicureTypeID int     `sql:"type:int, FOREIGN KEY (manicure_type_id) REFERENCES manicure_types(id)" json:"manicure_type"`
	Photos         []Photo `gorm:"many2many:service_photos;"`
}

// ManicureType model
type ManicureType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Sign  model
type Sign struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhotoID     int    `sql:"type:int, FOREIGN KEY (photo_id) REFERENCES photos(id)" json:"photo"`
	Description string `json:"description"`
}

//// ManicureMaterial model
//type ManicureMaterial struct {
//	ID          int       `json:"id"`
//	Firm        string    `json:"firm"`
//	Palette     string    `json:"palette"`
//	Description string    `sql:"type:text" json:"description"`
//	Services    []Service `gorm:"many2many:service_manicure_materials;"`
//}
