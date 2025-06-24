package main

import (
	"context"
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

	fmt.Println("Server running at http://localhost:8080/api")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func getVpcs() []string {
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
	servicesList := []string{"ec2", "vpc", "eks"}

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


func getHealthHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

    w.WriteHeader(http.StatusOK)
}