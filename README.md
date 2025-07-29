# 📣 NotifyGo - A Golang based notification service

A scalable, containerized microservice for sending notifications via Email, SMS, and WebSocket, built using [GoFr](https://gofr.dev) – a Golang-based backend framework. This project is designed with clean 3-layer architecture and extensibility in mind.

---

## 🚀 Features

- 📧 Email notifications via SMTP (MailHog)
- 📱 SMS notifications via Kannel gateway
- 🔔 WebSocket-based in-app alerts (real-time)
- 🧩 Dynamic message templating
- 🛠️ Unified API for triggering notifications
- 📦 Containerized with Docker Compose

---

## 🛠 Setup Instructions
### 1. Clone the Repo

```bash
git clone https://github.com/Kake27/notifygo
cd notification-service
```

### 2. Set Environment Variables
Some default environmental variables have already been included in the code. In case any modification is needed, the environemental variables can be set as follows and used
```env
SMTP_HOST=mailhog
SMTP_PORT=1025
SMTP_SENDER=test@example.com

KANNEL_URL=http://localhost:13013/cgi-bin/sendsms
KANNEL_USERNAME=test
KANNEL_PASSWORD=test123
KANNEL_SENDER=kannel
```

### 3. Start Services
Make sure docker is installed and then run 
```bash
docker compose up --build
```
This will start:
- notification-service (your Go app)
- MailHog for email testing at localhost:8025
- Kannel for SMS delivery on port 13013
---

## 📬 Sending Notifications
### 1. 📧 Send Email
POST `/notify`
```json
{
  "type": "email",
  "to": "user@example.com",
  "subject": "This is a test subject",
  "message": "Hello from the Notification Service!"
}
```

MailHog will capture and display the email at: http://localhost:8025

### 2. 📱 Send SMS
POST `/notify`
```json
{
  "type": "sms",
  "to": "+911234567890",
  "message": "This is a test SMS from the Notification Service!"
}
```
This uses the Kannel HTTP API behind the scenes for delivery. For now, the SMS is only being queued for delivery rather than actually being sent.

---

## 🧪 Testing Endpoints
Check if service is alive:
- `GET /health` → Returns whether the SMTP server for sending emails is online
- Basic test route: `GET /greet` → Returns "Hello from the server!" if the app is online.

---

## 🧰 Tech Stack
- **Language:** Golang
- **Framework:** GoFr
- **Email Delivery:** MailHog via SMTP
- **SMS Delivery:** Kannel Gateway
- **Real-Time Alerts:** WebSocket 
- **Containerization:** Docker, Docker Compose


