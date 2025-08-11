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
   - `~/.aws/credentials` file

2. **Ensure you have the required permissions** in IAM (Already set in [Base-AWS-Infrastructure](https://github.com/Rippyblogger/Base-AWS-Infrastructure) ):
   - `ec2:Describe*`
   - `eks:ListClusters`
   - `eks:DescribeCluster`,
   - `servicequotas:ListServiceQuotas`
   - `servicequotas:RequestServiceQuotaIncrease`,
   - `servicequotas:GetServiceQuota`,
   - `servicequotas:RequestServiceQuotaIncrease`

## Running the API Server on EKS

- AWS Account with appropriate permissions
- AWS CLI configured (for local development)
- AWS IAM role configured with the following permissions:
  ```json
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "ec2:Describe*",
          "eks:ListClusters",
          "eks:DescribeCluster",
          "vpc:Describe*",
          "servicequotas:ListServiceQuotas",
          "servicequotas:GetServiceQuota",
          "servicequotas:RequestServiceQuotaIncrease"
        ],
        "Resource": "*"
      }
    ]
  }
  ```

### Infrastructure Requirements
- EKS Cluster (deployed using the [Base-AWS-Infrastructure](https://github.com/Rippyblogger/Base-AWS-Infrastructure))

## 🚀 Deployment Options

### Option 1: Automated Deployment (Recommended)

This repository uses GitHub Actions for automated deployment. The deployment happens automatically when:

1. **On Push to Main**: Triggers build, test, and deployment
2. **On Pull Request Merge**: Deploys merged changes
3. **Manual Trigger**: To trigger the deployment manually

#### Setup for Automated Deployment:

1. **Fork this repository** to your GitHub account

2. **Set up GitHub Repository Variables** in your repository settings:
   ```
   ACCOUNT_ID: Your AWS Account ID
   AWS_REGION: Your preferred AWS region (e.g., us-west-2)
   AWS_ROLE_NAME: Name of your OIDC-enabled IAM role (e.g., oidc_role)
   ```

3. **Ensure Base Infrastructure is Deployed**: Make sure you have deployed the base infrastructure using the [Base-AWS-Infrastructure](https://github.com/Rippyblogger/Base-AWS-Infrastructure) repository first.

4. **Configure Terraform Backend**: Update the S3 backend configurations in `terraform/main.tf` to use the same S3 bucket created in Step 1 of the base infrastructure deployment:

   ```hcl
   terraform {
     backend "s3" {
       bucket  = "-terraform-state"  # Same bucket from base infrastructure
       key     = "applications/golang-api/terraform.tfstate"
       region  = "us-east-1"                          # Match your bucket region
       encrypt = true
     }
   }
 
    data "terraform_remote_state" "infra" {
      backend = "s3"
      config = {
        bucket = "s3-state-bucket73579"
        key    = "infrastructure/terraform.tfstate"
        region = "us-east-1" # Match your bucket region
      }
    }
    ```

   **⚠️ Important**: Replace `your-terraform-state-bucket-12345` with the actual S3 bucket name you created for the base infrastructure.

5. **Configure Kubernetes Deployment Parameters**

   The deployment uses a modular Terraform architecture that references the [Base-AWS-Infrastructure](https://github.com/Rippyblogger/Base-AWS-Infrastructure) . You can customize the deployment by modifying the existing sample module parameters in `terraform/main.tf`:

   ```hcl
   module "go_api_k8s" {
     source = "git::https://github.com/Rippyblogger/Base-AWS-Infrastructure.git//modules/kubernetes?ref=main"

     # Cluster connection (pulled from base infrastructure state)
     cluster_name           = data.terraform_remote_state.infra.outputs.cluster_name
     cluster_endpoint       = data.terraform_remote_state.infra.outputs.cluster_endpoint
     cluster_ca_certificate = data.terraform_remote_state.infra.outputs.cluster_ca_certificate
     
     # Application deployment settings
     deployment_name        = "go-api-deployment"        # Name of the Kubernetes deployment
     env                   = "production"                # Environment label
     replicas_count        = 1                          # Number of pod replicas
     app_name              = "go-api"                   # Application name/label
     image_name            = var.image_name             # ECR image URL (set by CI/CD)
     
     # IRSA (IAM Roles for Service Accounts) configuration
     service_account_name  = "my-go-api-irsa-sa"       # Kubernetes service account name
     oidc_arn             = data.terraform_remote_state.infra.outputs.oidc_provider_arn
     oidc_provider_url    = data.terraform_remote_state.infra.outputs.oidc_provider_url
   }
   ```

   **Key Parameters You Can Customize:**

   | Parameter | Description | Sample Value | Notes |
   |-----------|-------------|---------------|--------|
   | `deployment_name` | Name of the Kubernetes deployment | `go-api-deployment` | Must be unique within namespace |
   | `env` | Environment label (dev/staging/production) | `production` | Used for resource tagging |
   | `replicas_count` | Number of pod replicas | `1` | Increase for high availability |
   | `app_name` | Application name used in labels | `go-api` | Keep consistent across resources |
   | `service_account_name` | IRSA service account name | `my-go-api-irsa-sa` | |

   **Important Notes:**
   - The `image_name` variable is automatically set by the CI/CD pipeline with the ECR image URI
   - OIDC configuration enables the pods to assume AWS IAM roles without storing credentials
   - Cluster connection parameters are automatically retrieved from the base infrastructure Terraform state

6. **Push to Main Branch**: Any push to the main branch will trigger the deployment pipeline.

### Option 2: Manual Local Development

For local development and testing:

```bash
# Clone the repository
git clone https://github.com/Rippyblogger/golang_api.git
cd golang_api

# Install dependencies
go mod download

# Set up AWS credentials (choose one method)
# Method 1: AWS CLI
aws configure

# Method 2: Environment variables
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-west-2

# Run the application
go run main.go
```

The server will start on `http://localhost:8080`

### Option 3: Docker Deployment

```bash
# Build the Docker image
docker build -t golangapi .

# Run with AWS credentials mounted (Linux/Mac)
docker run -d \
  --name golang-api \
  --mount type=bind,source=$HOME/.aws,target=/root/.aws,readonly \
  -p 5000:8080 \
  golangapi

# Run with environment variables
docker run -d \
  --name golang-api \
  -e AWS_ACCESS_KEY_ID=your_access_key \
  -e AWS_SECRET_ACCESS_KEY=your_secret_key \
  -e AWS_REGION=us-west-2 \
  -p 5000:8080 \
  golangapi

```

## 🔧 Usage Examples

### Get VPC Information
```bash
curl -X GET http://localhost:8080/vpcs
```

### Get EC2 Instances
```bash
curl -X GET http://localhost:8080/ec2s
```

### Get EKS Clusters
```bash
curl -X GET http://localhost:8080/eks
```

### Check Service Quotas
```bash
curl -X GET http://localhost:8080/quotas
```

### Request Quota Increase
```bash
curl -X POST http://localhost:8080/quota \
  -H "Content-Type: application/json" \
  -d '{
    "serviceCode": "ec2",
    "quotaCode": "L-1216C47A",
    "desiredValue": 50
  }'
```

### Health Check
```bash
curl -X GET http://localhost:8080/health
```

## 🔄 CI/CD Pipeline

**Note**: This API is part of a larger infrastructure setup. Make sure to deploy the [Base-AWS-Infrastructure](https://github.com/Rippyblogger/Base-AWS-Infrastructure) first before deploying this application.

This repository implements a **GitOps deployment strategy** with the following workflows:

### Deploy to Production Workflow
- **Triggers**: Push to main, PR merge to main, manual dispatch
- **Process**:
  1. Builds and pushes Docker image to ECR
  2. Deploys to EKS using Terraform
  3. Updates Kubernetes deployment with new image
- **Security**: Uses OIDC authentication (no static credentials)

### Destroy Infrastructure Workflow
- **Trigger**: Manual dispatch only
- **Process**: Safely destroys all Terraform-managed resources

### Workflow Files
- `.github/workflows/deploy.yml` - Main deployment pipeline
- `.github/workflows/destroy.yml` - Container infrastructure destruction

## 🏗️ Project Structure

```
.
├── main.go                    # Main application file
├── Dockerfile                 # Container build instructions
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependency checksums
├── terraform/                 # Kubernetes deployment configuration
│   ├── main.tf               # Terraform configuration
│   ├── variables.tf          # Input variables
│   └── outputs.tf            # Output values
├── .github/
│   └── workflows/
│       ├── deploy.yml        # Deployment workflow
│       └── destroy.yml       # Destruction workflow
└── README.md                 # This file
```