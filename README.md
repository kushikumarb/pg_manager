# PG Manager - Full-Stack Dockerized Accommodation Solution

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.x-4FC08D?style=flat&logo=vuedotjs)](https://vuejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker)](https://www.docker.com/)

A comprehensive, production-ready management system for Paying Guest (PG) accommodations. This project streamlines the entire tenant lifecycle—from KYC and onboarding to automated billing and secure offboarding.


---

## 🚀 Technical Architecture

The application is built using a modern decoupled architecture, fully containerized for consistent deployment across development and production environments.

* **Backend:** High-performance REST API built with **Golang** (Gin Framework) and **GORM**.
* **Frontend:** Responsive User Interface built with **Vue.js 3** and **Tailwind CSS**.
* **Database:** Relational data management using **PostgreSQL 15**.
* **DevOps:** Multi-container orchestration via **Docker Compose**.



---

## ✨ Key Features

### 🏨 Room & Property Management
* Real-time tracking of room availability and capacity.
* Property-specific dashboards for owners to manage multiple buildings.

### 👤 Tenant Lifecycle Management
* **Automated Onboarding:** KYC registration with built-in OTP-based verification.
* **Smart Offboarding:** Secure "Soft-Archive" logic that backs up tenant history to an `Archived` table before performing hard deletions to free up resources.
* **PostgreSQL Optimized:** Transactional integrity ensuring parent-child records are handled without foreign key conflicts.

### 💰 Finance & Payments
* **Ledger Management:** Automated tracking of security deposits and monthly rent balances.
* **Payment Integration:** Ready for **Razorpay** link generation for initial deposits and monthly dues.
* **Rate Limiting:** Controlled balance inquiries to ensure system stability.

### 🛡️ Security
* **JWT Authentication:** Secure stateless session management for owners and tenants.
* **Bcrypt Hashing:** Industry-standard password encryption for all user accounts.
* **Environment Isolation:** Zero hardcoded credentials; strictly managed via `.env` injection.

---

## 🛠️ Local Development Setup

### Prerequisites
* Docker Desktop installed
* Git

### Quick Start
1.  **Clone the Repository:**
    ```bash
    git clone [https://github.com/kushikumarb/pg_manager.git](https://github.com/kushikumarb/pg_manager.git)
    cd pg_manager
    ```

2.  **Environment Configuration:**
    Create a `.env` file in the root directory. Use the structure defined in `backend/config/config.go`:
    ```env
    PORT=8080
    DB_USER=postgres
    DB_PASS=your_secure_password
    DB_NAME=pg_management
    DB_HOST=db
    JWT_SECRET=your_secure_jwt_key
    APP_ENV=development
    ```

3.  **Spin up the Containers:**
    ```bash
    docker-compose up --build
    ```

4.  **Access the Application:**
    * **Frontend:** `http://localhost`
    * **Backend API:** `http://localhost:8080`

---

## 📂 Project Structure

```text
├── backend/            # Golang Gin API
│   ├── config/         # Database and Env configuration logic
│   ├── handlers/       # HTTP Request controllers
│   ├── models/         # GORM database schemas & associations
│   ├── services/       # Core business logic (Auth, Tenant Archiving, etc.)
│   └── main.go         # Entry point & Database Seeding logic
├── frontend/           # Vue.js SPA
├── docker-compose.yml  # Multi-container orchestration blueprint
└── .gitignore          # Secure file exclusion (prevents .env leaks)
```
## 👨‍💻 Author

### Kushikumar B 
**Final Year B.E. Student at BIET Davanagere,**<br>
**Python Full Stack Intern at Pentagon Space, Bengaluru**
