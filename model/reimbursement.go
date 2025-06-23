package model

type Reimbursement struct {
	ID          string `db:"id" json:"id"`
	UserID      string `db:"user_id" json:"user_id"`
	Date        string `db:"date" json:"date"`
	Amount      int    `db:"amount" json:"amount"`
	Description string `db:"description" json:"description"`
	CreatedBy   string `db:"created_by" json:"created_by"`
	UpdatedBy   string `db:"updated_by" json:"updated_by"`
	IpAddress   string `db:"ip_address" json:"ip_address"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}
