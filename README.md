# Payslip Generation System

A scalable backend system to manage employee payrolls based on attendance, overtime, and reimbursement submissions.

---

## 📚 Background

In a company, employees are paid monthly based on fixed rules:

- **Working hours**: 9AM–5PM, Monday to Friday
- **Salary** is **prorated** based on attendance (per day)
- **Overtime** is paid **2× hourly rate**
- **Reimbursements** are added directly to the monthly payslip

This system allows employees to submit attendance, overtime, and reimbursements, and enables the admin to generate payroll and reports.

---

## 🎯 Objectives

- Handle salary calculation based on attendance
- Support overtime and reimbursement submissions
- Enable payroll generation per attendance period
- Generate payslips for each employee
- Generate a summary of payslips for admin

---

## 🧰 Tech Stack

- Go (Golang)
- Fiber (HTTP framework)
- PostgreSQL
- golang-migrate (for DB migration)
- Swagger (for API documentation)
- Go's built-in `testing` package

---

## ⚙️ Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/euegniaangie/payslip-generator.git
cd payslip-generator
```

### 2. Setup PostgreSQL

Make sure PostgreSQL is installed and running.

```sql
CREATE DATABASE payslip;
CREATE USER payslip_user WITH PASSWORD 'password123';
GRANT ALL PRIVILEGES ON DATABASE payslip TO payslip_user;
```

### 2.1. Create .env File

Create a `.env` file in the project root with the following content (example):

```
JWT_SECRET=super_secret_key_123
```

You can change the value to any secret string you want.

### 3. Run Database Migration

Install golang-migrate CLI if you haven't already:

```bash
brew install golang-migrate   # macOS
```

Then run the migration:

```bash
migrate -path migration -database "postgres://payslip_user:password123@localhost:5432/payslip?sslmode=disable" up
```

This will:
- Create the required tables
- Seed initial data:
  - 100 dummy employees
  - 1 dummy admin

### 4. Run the Server

```bash
make watch
```

Access the server at:
- [http://localhost:4000](http://localhost:4000)

---

## 📖 Swagger API Documentation

Swagger UI is available at:
- [http://localhost:4000/swagger/index.html](http://localhost:4000/swagger/index.html)

You can explore and test the available endpoints.

---

## ✅ Implemented Features

### 👩‍💼 Admin
- Add attendance period (start & end date)
- Run payroll (locks data for the period, can only be run once)
- Generate summary report of employee payslips (total take-home pay per employee & grand total)

### 👨‍💼 Employees
- Submit attendance (only once per day, cannot submit on weekends)
- Submit overtime (only after working hours, max 3 hours/day, can submit on any day)
- Submit reimbursement (with amount & description)
- Generate individual payslip (attendance breakdown, overtime, reimbursement, take-home pay)

---

## 🛠 Example API Usage

### ✅ Register (Admin or Employee)

```http
POST /auth/register
Content-Type: application/json

{
  "full_name": "employeexyz",
  "username": "employeexyz",
  "password": "employee123",
  "user_role": "employee"
}
```

### ✅ Login (Admin or Employee)

```http
POST /auth/login
Content-Type: application/json

{
  "username": "employeexyz",
  "password": "employee123"
}
```
> The response will include a Bearer token. Use this token in the `Authorization` header for all authenticated requests:

```
Authorization: Bearer <your_token_here>
```

### ✅ Submit Attendance

```http
POST /attendance
```
> Cannot submit on weekends.

### ✅ Submit Overtime

```http
POST /overtime
Content-Type: application/json

{
  "hours": 2
}
```
> Only allowed after working hours, max 3 hours.

### ✅ Submit Reimbursement

```http
POST /reimbursement
Content-Type: application/json

{
  "amount": 250000,
  "date" : "2025-01-20"
  "description": "Travel expense"
}
```

### ✅ Create Attendance Period (Admin)

```http
POST /admin/attendance-period
Content-Type: application/json

{
  "start_date": "2025-01-01",
  "end_date": "2025-01-31"
}
```

### ✅ Run Payroll (Admin)

```http
POST /admin/payroll/run
Content-Type: application/json

{
  "attendance_period_id": "uuid-attendance-period"
}
```

### ✅ Get Attendance Period (Employee)
```http
GET /attendance-period?page=1&page_size=10
```

### ✅ Generate Payslip (Employee)

```http
POST /payslip
Content-Type: application/json
{
  "period_id": "0370f6f0-5762-4bbe-a57c-b398c2c76dd2"
}
```

### ✅ Generate Summary (Admin)

```http
GET /admin/payslip/summary?attendance_period_id=uuid-attendance-period
```

## 🗝️ Default Logins

- Admin:  
  - Username: `admin`  
  - Password: `admin123`
- Employees:  
  - Username: `employee_1`, `employee_2`, ..., `employee_100`
  - Password: `employee123`
