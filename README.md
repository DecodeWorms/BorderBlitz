# 💸 BorderBlitz Backend

A multi-currency payment sandbox for cross-border transfers using mock stablecoins, built for the **Operation Borderless** take-home test.

---

## 🧭 Overview

BorderBlitz is a backend service providing API endpoints for:

- Multi-currency wallet management
- Simulated deposits and swaps
- Transfers between users
- Transaction history

It's designed to support frontend integration (React) for user interfaces like deposit, transfer, and swap pages — including a wallet explorer with a pie chart view.

---

## ⚙️ Tech Stack

- **Language**: Go (Gin for routing, GORM for ORM)
- **Database**: PostgreSQL
- **Containerization**: Docker
- **Deployment**: 
  - Local: Docker Desktop with Kubernetes
  - Production: Render.com
- **Web Server**: Nginx (Ingress Controller for local setup)
- **SSL**: Let's Encrypt (Kubernetes Ingress)
- **Container Registry**: Docker Hub

---

## 🚀 Features

- 🏦 **Virtual Wallets**  
  Create and manage wallets with support for:
  - `cNGN`, `cXAF` (African currencies)
  - `USDx`, `EURx` (Stablecoins)

- 💰 **Simulated Deposits**  
  Deposit mock stablecoins into wallets.

- 🔄 **Cross-Currency Swaps**  
  Swap between stablecoins using mocked exchange rates (e.g., `500 cNGN = 1 USDx`).

- 📤 **Transfers with Auto-Swap**  
  Transfer between users, with currency auto-conversion if needed.

- 📜 **Transaction History**  
  View all transactions including deposits, swaps, and transfers.

---

## 📦 Docker Hub

- 🐳 **Backend Image**:  
  [decodeworms/borderblitz-backend](https://hub.docker.com/r/decodeworms/borderblitz-backend)

---

## 🛠️ Prerequisites

### For Local Kubernetes Deployment

- [Docker Desktop](https://www.docker.com/products/docker-desktop) (with Kubernetes enabled)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- Local DNS resolution:

  Add this to your `/etc/hosts` file (Linux/macOS):


---

## 📡 Deployment Options

### 🧪 Local Development with Kubernetes (Docker Desktop)

1. **Enable Kubernetes**:
 - Open Docker Desktop → **Settings > Kubernetes** → Enable.

2. **Install `kubectl`**:

 ```bash
 curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
 sudo install kubectl /usr/local/bin/kubectl

### 3. Deploy to Kubernetes:
Apply the necessary YAML manifests (Deployment, Service, Ingress, Secret, etc.) using:
kubectl apply -f k8s/

### 4. Access the API: Visit: http://borderblitz.local/api/v1

v1

🌍 Production Deployment on Render
The backend is deployed on Render, connected to this repo.

Render Setup
1. Create a new Render Web Service.

2 Environment:

Runtime: Go

Start Command:
./your-backend-binary

Environment Variables:
Add variables like DB_URL, PORT, etc.

3. Database: PostgreSQL (provision via Render Addons)

4. Docker Image: Optional (You can use Dockerfile or auto-deploy from repo)

5. API Base URL:
https://borderblitz.onrender.com/api/v1

🤝 Author
Abdulhameed (@DecodeWorms)
Backend Developer | Go Enthusiast | Systems Thinker
GitHub: DecodeWorms

