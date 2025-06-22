package model

type AttendancePeriod struct {
	ID        string `db:"id" json:"id"`
	StartDate string `db:"start_date" json:"start_date"`
	EndDate   string `db:"end_date" json:"end_date"`
	CreatedBy string `db:"created_by" json:"created_by"`
	UpdatedBy string `db:"updated_by" json:"updated_by"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}
