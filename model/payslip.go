package model

type Payslip struct {
	ID            string `db:"id" json:"id"`
	UserID        string `db:"user_id" json:"user_id"`
	PeriodID      string `db:"period_id" json:"period_id"`
	BaseSalary    int    `db:"base_salary" json:"base_salary"`
	PresentDays   int    `db:"present_days" json:"present_days"`
	WorkingDays   int    `db:"working_days" json:"working_days"`
	OvertimeHours int    `db:"overtime_hours" json:"overtime_hours"`
	Reimbursement int    `db:"reimbursement" json:"reimbursement"`
	TakeHomePay   int    `db:"take_home_pay" json:"take_home_pay"`
	CreatedBy     string `db:"created_by" json:"created_by"`
	UpdatedBy     string `db:"updated_by" json:"updated_by"`
	IpAddress     string `db:"ip_address" json:"ip_address"`
	CreatedAt     string `db:"created_at" json:"created_at"`
	UpdatedAt     string `db:"updated_at" json:"updated_at"`
}
