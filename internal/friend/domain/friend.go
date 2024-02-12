package domain

type Friend struct {
	UserId1 string `json:"userId1" gorm:"column:user_id_1"`
	UserId2 string `json:"userId2" gorm:"column:user_id_2"`
	Status  string `json:"status" gorm:"column:status"`
}
