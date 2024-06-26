package domain

type SearchResult interface{}

type Keyword struct {
	Q string `json:"q" gorm:"column:q"`
}

type User struct {
	Id       string `json:"id" gorm:"id"`
	Username string `json:"name" gorm:"column:name"`
	Avatar   string `json:"avatar" gorm:"column:avatar"`
	Carrier  string `json:"carrier" gorm:"column:carrier"`
}

type Group struct {
	Id       string `json:"id" gorm:"id"`
	Name     string `json:"name" gorm:"column:name"`
	Avatar   string `json:"avatar" gorm:"column:avatar"`
	Category string `json:"category" gorm:"column:category"`
}

type Page struct {
	Id       string `json:"id" gorm:"id"`
	Name     string `json:"name" gorm:"column:name"`
	Avatar   string `json:"avatar" gorm:"column:avatar"`
	Category string `json:"category" gorm:"column:category"`
}
