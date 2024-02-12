package domain

type Sequence struct {
	Name       string `json:"name" gorm:"column:name"`
	SequenceNo int64  `json:"sequenceNo" gorm:"column:sequence_no"`
}
