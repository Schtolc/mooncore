package models

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
)

// Address model
type Address struct {
	ID          int64          `json:"id"`
	Lat         float64        `json:"lat"`
	Lon         float64        `json:"lon"`
	Description string         `json:"description"`
	Stations    []AddressMetro `json:"stations"`
}

// AddressMetro model
type AddressMetro struct {
	ID        int64   `json:"id"`
	AddressID int64   `json:"address_id"`
	Name      string  `json:"name"`
	Line      string  `json:"line"`
	Color     string  `json:"color"`
	Distance  float64 `json:"distance"`
}

// MetroLine model
type MetroLine struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// MetroStation model
type MetroStation struct {
	ID     int       `json:"id"`
	LineID int       `json:"line_id"`
	Line   MetroLine `json:"line"`
	Name   string    `json:"name"`
	Lat    float64   `json:"lat"`
	Lon    float64   `json:"lon"`
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

// User model
type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email";gorm:"not null;unique;"`
	PasswordHash string `gorm:"not null"` // not serializable
	Role         int    `json:"role"`
	Ctime        int64
	Atime        int64
}

// Master model
type Master struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	UserID    int64         `sql:"type:bigint, FOREIGN KEY (user_id) REFERENCES users(id)"`
	User      User          `json:"user"`
	AddressID int64         `sql:"type:bigint, FOREIGN KEY (address_id) REFERENCES addresses(id)"`
	Address   Address       `json:"address"`
	PhotoID   sql.NullInt64 `sql:"type:bigint, FOREIGN KEY (photo_id) REFERENCES photos(id)"`
	Photo     Photo         `json:"photo"`
	SalonID   sql.NullInt64 `sql:"type:bigint, FOREIGN KEY (salon_id) REFERENCES salons(id)"`
	Salon     *Salon        `json:"salon"`
	Stars     int           `json:"stars"`
	Services  []Service     `json:"services"`
	Photos    []Photo       `gorm:"many2many:master_photos;" json:"photos"`
	Signs     []Sign        `gorm:"many2many:master_signs;" json:"signs"`
	Home      int           `sql:"type:int" json:"home_service"`
}

// Client model
type Client struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	UserID    int64    `sql:"type:bigint, FOREIGN KEY (user_id) REFERENCES users(id)"`
	User      User     `json:"user"`
	PhotoID   int64    `sql:"type:bigint, FOREIGN KEY (photo_id) REFERENCES photos(id)"`
	Photo     Photo    `json:"photo"`
	Home      bool     `json:"home_service"`
	Favorites []Master `gorm:"many2many:client_favorites;" json:"favorites"`
}

// Material model
type Material struct {
	ID          int64  `json:"id"`
	Firm        string `sql:"type:text" json:"firm"`
	Description string `sql:"type:text" json:"description"`
	Name        string `json:"name"`
}

// Service model
type Service struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Price       float64    `json:"price"`
	MasterID    int64      `sql:"type:bigint, FOREIGN KEY (master_id) REFERENCES masters(id)"`
	Master      Master     `json:"master"`
	Description string     `sql:"type:text" json:"description"`
	Photos      []Photo    `gorm:"many2many:service_photos;" json:"photos"`
	Materials   []Material `gorm:"many2many:service_materials;" json:"materials"`
	Ctime       int64
	TypeID      int64       `sql:"type:bigint, FOREIGN KEY (type_id) REFERENCES service_types(id)"`
	Type        ServiceType `json:"type"`
}

// ServiceType model
type ServiceType struct {
	ID       int64         `json:"id"`
	ParentID sql.NullInt64 `sql:"type:bigint, FOREIGN KEY (parent_id) REFERENCES service_types(id)"`
	Parent   *ServiceType  `json:"parent"`
	Name     string        `json:"name"`
}

// Salon model
type Salon struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	AddressID int64   `sql:"type:bigint, FOREIGN KEY (address_id) REFERENCES addresses(id)"`
	Address   Address `json:"address"`
	PhotoID   int64   `sql:"type:bigint, FOREIGN KEY (photo_id) REFERENCES photos(id)"`
	Photo     Photo   `json:"photo"`
	Stars     int     `json:"stars"`
}

// Token model
type Token struct {
	Token string `json:"token"`
}

// JwtClaims model
type JwtClaims struct {
	Email string
	jwt.StandardClaims
}
