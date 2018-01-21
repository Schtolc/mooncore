package models

import "github.com/dgrijalva/jwt-go"

// Address model
type Address struct {
	ID          int64   `json:"id"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Description string  `json:"description"`
}

// Photo model
type Photo struct {
	ID   int64  `json:"id"`
	Path string `json:"path"`
	Tags []Tag  `gorm:"many2many:photo_tags"`
}

// Tag model
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Sign  model
type Sign struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

// ManicureType model
type ManicureType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// User model
type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username";gorm:"not null;unique;"`
	Email        string `json:"email";gorm:"not null;unique;"`
	PasswordHash string `gorm:"not null"` // not serializable
	Role         int    `json:"role"`
}

// Master model
type Master struct {
	ID        int64     `json:"id"`
	UserID    int64     `sql:"type:int, FOREIGN KEY (user_id) REFERENCES users(id)" json:"user"`
	Name      string    `json:"name"`
	AddressID int64     `sql:"type:int, FOREIGN KEY (address_id) REFERENCES addresses(id)" json:"address"`
	PhotoID   int64     `sql:"type:int, FOREIGN KEY (photo_id) REFERENCES photos(id)" json:"photo"`
	Stars     int   `json:"stars"`
	Services  []Service `gorm:"ForeignKey:MasterID" json:"services"`
	Photos    []Photo   `gorm:"many2many:user_photos;" json:"photos"`
	Signs     []Sign    `gorm:"many2many:user_signs;" json:"signs"`
	User      User      //TODO load user
	Address   Address   //TODO load address
	Photo     Photo     //TODO load photo
	SalonID int64
	Salon Salon
}

//photo := &models.Photo{}
//
//if dbc := dependencies.DBInstance().First(photo, p.Source.(models.UserDetails).PhotoID); dbc.Error != nil {
//logrus.Error(dbc.Error)
//return nil, dbc.Error
//}
//if dbc := dependencies.DBInstance().Model(photo).Association("tags").Find(&photo.Tags); dbc.Error != nil {
//logrus.Error(dbc.Error)
//return nil, dbc.Error
//}
//logrus.Warn(photo)
//return *photo, nil

// Client model
type Client struct {
	ID        int64    `json:"id"`
	UserID    int64    `sql:"type:int, FOREIGN KEY (user_id) REFERENCES users(id)" json:"user"`
	Name      string   `json:"name"`
	PhotoID   int64    `sql:"type:int, FOREIGN KEY (photo_id) REFERENCES photos(id)" json:"photo"`
	Favorites []Master //TODO add table definition
	User      User
	Photo     Photo
}

// Service model
type Service struct {
	ID             int64   `json:"id"`
	MasterID       int64   `json:"master"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Description    string  `sql:"type:text" json:"description"`
	ManicureTypeID int64   `sql:"type:int, FOREIGN KEY (manicure_type_id) REFERENCES manicure_types(id)" json:"manicure_type"`
	Photos         []Photo `gorm:"many2many:service_photos;" json:"photos"`
}

type Salon struct {
	ID int64
	Name string
	AddressID int64
	Address Address
	PhotoID   int64
	Photo Photo
	Stars int
}

type Token struct {
	Token string `json:"token"`
}

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}
