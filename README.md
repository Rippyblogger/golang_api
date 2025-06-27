# Golang-API

A Go-based API server for monitoring AWS resources and service quotas built with AWS SDK for Go v2

## Features
- Real-time monitoring of AWS resources (VPCs, EC2s, EKS clusters)
- Service quota utilization tracking
- Programmatic quota increase requests

## API Endpoints

### Required Routes
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/vpcs` | GET | Returns information about VPCs in your AWS account |
| `/ec2s` | GET | Returns information about EC2 instances in your account |
| `/eks` | GET | Returns information about EKS clusters in your account |
| `/quotas` | GET | Shows service quotas utilzation for EC2, VPC, and EKS |
| `/quota` | POST | Requests an increase for an account-level service quota |

## Getting Started

### Prerequisites
- AWS credentials with appropriate permissions
- Go 1.18+
- AWS account

### To run the API

1. **Configure AWS credentials** using one of these methods:
   - AWS CLI: `aws configure`
   - Environment variables:
     ```
     export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY
     export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
     export AWS_REGION=SET_REGION
     ```
   - `~/.aws/credentials` file

2. **Set required permissions** in IAM:
   - `ec2:Describe*`
   - `eks:ListClusters`
   - `servicequotas:ListServiceQuotas`
   - `servicequotas:RequestServiceQuotaIncrease`

### Running the API Server

1. **Clone the repository**
git clone **https://github.com/Rippyblogger/Golang-API.git**

2. Navigate to project directory

        `cd Golang-API`

3. Install dependencies

        `go mod tidy`

4. Start the API server

        `go run golang_api.go`

The server will start on **http://localhost:8080**


## Usage Examples

### Retrieve VPC Information
- curl http://localhost:8080/vpcs


### Request Quota Increase

- **Example:** `curl -X POST http://localhost:8080/quota
-H "Content-Type: application/json"
-d '{
"serviceCode": "ec2",
"quotaCode": "L-1216C47A",
"desiredValue": 50
}'`
