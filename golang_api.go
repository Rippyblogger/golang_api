package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/vpcs", getVpcsHandler)
	http.HandleFunc("/ec2s", getEc2sHandler)
	http.HandleFunc("/eks", getEksHandler)
	http.HandleFunc("/quotas", getQuotasHandler)
	http.HandleFunc("/health", getHealthHandler)
	http.HandleFunc("/quota", quotaIncreaseHandler)

	fmt.Println("Server running at http://localhost:8080/")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func getVpcs() []string {
	// Load AWS SDK config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		// config.WithSharedConfigProfile("default"),
	)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed loading config, %v", err))
	}

	// Create service client
	ec2Client := ec2.NewFromConfig(cfg)

	// Create input variables
	vpcInput := &ec2.DescribeVpcsInput{DryRun: aws.Bool(false)}

	// Call function
	vpcOut, err := ec2Client.DescribeVpcs(context.TODO(), vpcInput)

	if err != nil {
		log.Fatalf(fmt.Sprintf("Error: %v", err))
	}

	vpcList := []string{}
	for _, v := range vpcOut.Vpcs {
		formattedOutput := fmt.Sprintf(
			"CidrBlock: %s\n"+
				"VpcId: %s\n",
			*v.CidrBlock, *v.VpcId)

		vpcList = append(vpcList, formattedOutput)
	}

	return vpcList
}

func getVpcsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	response := getVpcs()
	fmt.Fprintln(w, response)
}

func getEc2s() []string {
	// Load AWS SDK config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("default"),
	)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed loading config, %v", err))
	}

	// Create service client
	ec2Client := ec2.NewFromConfig(cfg)

	// Create input variables
	ec2Input := &ec2.DescribeInstancesInput{DryRun: aws.Bool(false)}

	// Call function
	ec2Out, err := ec2Client.DescribeInstances(context.TODO(), ec2Input)

	if err != nil {
		log.Fatalf(fmt.Sprintf("Error: %v", err))
	}

	ec2List := []string{}
	for _, v := range ec2Out.Reservations {

		for _, v := range v.Instances {
			formattedOutput := fmt.Sprintf(
				"Instance ID: %s\n"+
					"InstanceType: %s\n"+
					"PublicIpAddress: %s\n"+
					"VpcId: %s\n",
				*v.InstanceId, v.InstanceType, *v.PublicIpAddress, *v.VpcId)

			ec2List = append(ec2List, formattedOutput)
		}

	}

	return ec2List
}

func getEc2sHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	response := getEc2s()
	fmt.Fprintln(w, response)
}

func getEks() []string {
	// Load AWS SDK config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("default"),
	)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed loading config, %v", err))
	}

	// Create service client
	eksClient := eks.NewFromConfig(cfg)

	// Create input variables
	eksInput := &eks.ListClustersInput{MaxResults: aws.Int32(10)}

	// Call function
	eksOut, err := eksClient.ListClusters(context.TODO(), eksInput)

	if err != nil {
		log.Fatalf(fmt.Sprintf("Error: %v", err))

	}

	eksList := []string{}
	for _, v := range eksOut.Clusters {

		formattedOutput := fmt.Sprintf(
			"Clusters: %s\n",
			v)

		eksList = append(eksList, formattedOutput)

	}

	if len(eksList) == 0 {
		fmt.Println("There are no EKS clusters in your AWS account")
		return nil

	} else {
		return eksList
	}
}

func getEksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	response := getEks()
	fmt.Fprintln(w, response)
}

func getServiceQuotas() map[string][]string {
	// Load AWS SDK config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("default"),
	)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed loading config, %v", err))
	}

	// Create service client
	sqClient := servicequotas.NewFromConfig(cfg)

	// Create input variables=
	// servicesList := []string{"ec2", "vpc", "eks"}
	servicesList := []string{"vpc"}

	// Call function
	quotaMap := make(map[string][]string)

	for _, serviceCode := range servicesList {
		listSqInput := &servicequotas.ListServiceQuotasInput{
			ServiceCode: aws.String(serviceCode),
		}

		sqOut, err := sqClient.ListServiceQuotas(context.TODO(), listSqInput)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		quotaMap[serviceCode] = []string{}

		for _, quota := range sqOut.Quotas {
			formattedOutput := fmt.Sprintf(
				"Quota Name: %s\n"+
					"Service Name: %s\n"+
					"Value: %v\n",
				*quota.QuotaName, *quota.ServiceName, *quota.Value)

			quotaMap[serviceCode] = append(quotaMap[serviceCode], formattedOutput)
		}
	}

	return quotaMap

}

func getQuotasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	response := getServiceQuotas()
	fmt.Fprintln(w, response)
}

func getHealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func quotaIncreaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg, error := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("default"),
	)

	if error != nil {
		http.Error(w, "failed loading config", http.StatusInternalServerError)
	}

	// Create service client
	quotaIncreaseClient := servicequotas.NewFromConfig(cfg)

	type requestInput struct {
		DesiredValue float64 `json:"desired_value"`
		QuotaCode    string  `json:"quota_code"`
		ServiceCode  string  `json:"service_code"`
	}

	var input requestInput

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.DesiredValue == 0 || input.QuotaCode == "" || input.ServiceCode == "" {
		http.Error(w, "One or more fileds are invalid", http.StatusBadRequest)
		return
	}

	// Create input variables
	quotaInput := &servicequotas.RequestServiceQuotaIncreaseInput{
		DesiredValue: aws.Float64(input.DesiredValue),
		QuotaCode:    aws.String(input.QuotaCode),
		ServiceCode:  aws.String(input.ServiceCode),
	}

	// Call function
	response, err := quotaIncreaseClient.RequestServiceQuotaIncrease(context.TODO(), quotaInput)

	if err != nil {
		http.Error(w, "Failed to request quota increase: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
