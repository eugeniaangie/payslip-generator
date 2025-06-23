package model

type Attendance struct {
	ID        string `db:"id" json:"id"`
	UserID    string `db:"user_id" json:"user_id"`
	Date      string `db:"date" json:"date"`
	IpAddress string `db:"ip_address" json:"ip_address"`
	CreatedBy string `db:"created_by" json:"created_by"`
	UpdatedBy string `db:"updated_by" json:"updated_by"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}
