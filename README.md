# BorderBlitz

A multi-currency payment sandbox for cross-border transfers using mock stablecoins, built for the **Operation Borderless** take-home test.

## Overview
BorderBlitz is a backend-only application (at this stage) that provides a multi-currency payment sandbox. It supports virtual wallets in mock stablecoins (cNGN, cXAF, USDx, EURx), simulated deposits, cross-currency swaps, transfers, and transaction history. The frontend (React) integration, including UI components for signup, deposit, swap, transfer, transaction history, and a pie chart for wallet balances.

## Tech Stack
- **Backend**: Go (Gin for routing, GORM for ORM)
- **Database**: PostgreSQL
- **Web Server**: Nginx (Ingress controller), HTTPS via Let’s Encrypt
- **Deployment**: Docker Desktop Kubernetes
- **Container Registry**: Docker Hub

## Features
- **Virtual Wallets**: Create and manage wallets in multiple currencies (cNGN, cXAF, USDx, EURx).
- **Simulated Deposits**: Deposit funds into a wallet.
- **Cross-Currency Swaps**: Swap between currencies using mock exchange rates (e.g., 500 cNGN = 1 USDx).
- **Transfers**: Transfer funds between users, with auto-swap if currencies differ.
- **Transaction History**: View a user’s transaction history (deposits, swaps, transfers).

## Docker Hub Repositories
- Backend: [decodeworms/borderblitz-backend](https://hub.docker.com/r/decodeworms/borderblitz-backend)

## Prerequisites
- **Docker Desktop**: Installed with Kubernetes enabled.
- **kubectl**: Kubernetes command-line tool.
- **Local Domain Setup**: Add `127.0.0.1 borderblitz.local` to your `/etc/hosts` (or `C:\Windows\System32\drivers\etc\hosts` on Windows).

## Deployment Instructions (Docker Desktop Kubernetes - Backend Only)
1. **Install Docker Desktop**:
   - Download and install Docker Desktop from [docker.com](https://www.docker.com/products/docker-desktop).
   - Enable Kubernetes in **Settings > Kubernetes**.

2. **Install kubectl**:
   ```bash
   curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
   sudo install kubectl /usr/local/bin/kubectl