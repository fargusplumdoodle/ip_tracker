# 📡 IP Tracker

## 🚀 Description

IP Tracker is a simple Go application designed to help you keep track of your 
homelab's external IP address, using Notion as a database.

A Kubernetes cronjob is seriously overkill for this task. I am considering moving to a FaaS platform in the future.
This was whipped up real quick to work with my homelab.

## 🛠️ Features

- Fetches external IPv4 address from multiple services
- Logs IP changes to a Notion page
- Deploy through Kubernetes CronJob

## 📋 Getting Started

### 📝 Prerequisites

- Go 1.22 or later
- Docker
- Notion API token and page ID
- Kubernetes cluster (for deployment)
- Helm (for Kubernetes deployment)

### 🔧 Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/fargusplumdoodle/ip_tracker.git
   ```

2. **Navigate to the project directory**

   ```bash
   cd ip_tracker
   ```

3. **Copy the example environment file and fill in your details**

   ```bash
   cp .env_example .env
   ```

4. **Edit the `.env` file**

   - Set your `NOTION_TOKEN`
   - Set your `NOTION_PAGE_ID`

## ⚙️ Configuration

Edit the `.env` file to configure the application:

- `NOTION_TOKEN`: Your Notion integration token
- `NOTION_PAGE_ID`: The ID of the Notion page to log IP addresses
- `IP_SERVICES`: Comma-separated list of services to fetch the external IP (optional)
- `IMAGE`: Docker image repository
- `TAG`: Docker image tag
- `NAMESPACE`: Kubernetes namespace for deployment

## 📦 Usage

Run the application locally:

```bash
go run main.go notion.go get_ipv4.go
```

## 🚢 Deployment

This application is intended to run as a Kubernetes CronJob to periodically check and log your external IP address. Scripts are provided to simplify the deployment process.

### Build and Push Docker Image

Use the provided script to build and push the Docker image:

```bash
./scripts/build_and_push
```

Ensure that `IMAGE` and `TAG` are correctly set in your `.env` file.

### Deploy to Kubernetes

Use the deployment script to deploy the application to your Kubernetes cluster:

```bash
./scripts/deploy
```

This script uses Helm to deploy the application as a CronJob. Make sure `NAMESPACE`, `NOTION_TOKEN`, and `NOTION_PAGE_ID` are set in your `.env` file.


## 📄 License

This project is licensed under the MIT License.