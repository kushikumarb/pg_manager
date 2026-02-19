# PG-Manager Backend 🏢

### _A Real-Time Smart PG Management & Automation Solution_

This is a professional Go-based backend designed to automate PG (Paying Guest) operations. The system eliminates manual tracking by handling anniversary-based billing, automated WhatsApp notifications, and secure digital onboarding.

---

## 🚀 Key Features

### 1. Multi-Property & Room Management

- **Property Hierarchy:** Owners can manage multiple buildings (properties) under a single dashboard.
- **Smart Capacity:** Define room-specific capacities. The system prevents over-onboarding by tracking live occupancy.

### 2. Secure Tenant Onboarding (OTP Flow)

- **Digital Admission:** Registration triggers a secure 6-digit OTP sent to the tenant's mobile.
- **Admission Verification:** The owner enters the OTP in the platform to confirm the tenant is real. Only then is the room officially allotted and the status moved to `active`.

### 3. Automated Anniversary Billing 📅

- **The 30-Day Cycle:** Each tenant follows their own billing cycle based on their onboarding date.
- **Daily Scheduler:** A GoCron job runs every morning at 09:00 AM to identify tenants who have completed their 30-day span.
- **WhatsApp Reminders:** On the 31st day, the system automatically sends a professional rent reminder via WhatsApp.

### 4. Financials & Payment Integration

- **Razorpay Webhooks:** Real-time balance updates. When a tenant pays via the digital link, the database is updated instantly without owner intervention.
- **Cash Payments:** For manual transactions, the owner can click "Paid" to update the database and generate a digital record.
- **Expenditure Tracking:** A dedicated module for the owner to log daily/monthly PG costs (Electricity, Water, Repairs) to figure out net finance.

### 5. Maintenance & Tenant Rights (Point 7 Logic)

- **Anti-Spam Limits:** Tenants are restricted to **one balance check** and **one complaint** per 24 hours.
- **Complaint Workflow:** Complaints are routed to a private owner table. Once fixed, the owner marks it as "Resolved," which triggers an automated "Issue Fixed" notification to the tenant.

---

## 🛠️ Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin Gonic (High-performance API routing)
- **Database:** PostgreSQL / MySQL (GORM)
- **Automation:** GoCron (Daily billing scheduler)
- **Integrations:** Razorpay (Payments) & WhatsApp Business API Utility

---

## 📂 Project Structure

```text
pg-manager-backend/
├── config/             # DB & Env configurations
├── handlers/           # Controller layer (API logic)
├── middleware/         # Auth & JWT Security
├── models/             # GORM Database Structs
├── services/           # Business Logic (OTP, Billing, Webhooks)
├── utils/              # Helpers (WhatsApp API, Random Generators)
└── main.go             # Entry point & Cron Scheduler

🛣️ API Endpoints Summary:

🔑 Owner & Property Management:

    POST /auth/register — Create a new owner account
    POST /auth/login — Authenticate and receive JWT
    POST /api/properties — Add a new PG building
    POST /api/rooms — Add rooms with capacity and price

💰 Tenant & Finance:

    POST /api/tenants/onboard — Register tenant & trigger OTP
    POST /api/tenants/verify — Enter OTP to confirm admission
    POST /api/tenants/payment — Record manual cash payment
    POST /api/expenditures — Log PG expenses (Electricity, Water, etc.)
    POST /api/webhooks/razorpay — (Public) Automated payment listener

🛠️ Complaints & Maintenance:

    GET /api/complaints — Owner view of all pending issues
    PUT /api/complaints/:id/resolve — Mark issue as solved & notify tenant

🔧 Installation & Setup:

1. Clone the repository:

    git clone [https://github.com/your-username/pg-manager-backend.git](https://github.com/your-username/pg-manager-backend.git)


2. Setup Environment Variables: Create a .env file in the root directory:

    DB_URL=your_database_connection_string
    JWT_SECRET=your_super_secret_key
    RAZORPAY_KEY=your_razorpay_api_key
    RAZORPAY_SECRET=your_razorpay_webhook_secret

3. Run the Application:
    go run main.go


Developed as a high-value MVP for modern PG owners, focusing on automation, financial transparency, and effective management.
```
