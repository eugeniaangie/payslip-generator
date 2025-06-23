package model

type Salary struct {
	ID        string `db:"id" json:"id"`
	UserID    string `db:"user_id" json:"user_id"`
	Amount    int    `db:"amount" json:"amount"`
	IsActive  bool   `db:"is_active" json:"is_active"`
	CreatedBy string `db:"created_by" json:"created_by"`
	UpdatedBy string `db:"updated_by" json:"updated_by"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}
