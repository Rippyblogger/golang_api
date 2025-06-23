package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/vpcs", getVpcsHandler)
	http.HandleFunc("/ec2s", getEc2sHandler)

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
