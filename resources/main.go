// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/cactus/go-statsd-client/v5/statsd"
)

const (
	RetryTime = 15
)

var threshold = map[int]float64{
	100:   89000000,
	1000:  89000000,
	10000: 89000000,
}

var dimension = []types.Dimension{
	{
		Name:  aws.String("exe"),
		Value: aws.String("cloudwatch-agent"),
	},
	{
		Name:  aws.String("process_name"),
		Value: aws.String("amazon-cloudwatch-agent"),
	},
}

func main() {
	for _, tps := range []int{100, 1000, 10000} {
		err := sendStatsDStress(tps)
		if err != nil {
			log.Fatalf("Send statstd stress failed")
		}
		for currentRetry := 1; ; currentRetry++ {

			metricValues, err := GetMetricDataResults("StressTest", "procstat_memory_rss", dimension)
			if err != nil {
				log.Fatalf("Fail get metric data result because of %v", err)
			}

			for _, value := range metricValues {
				if value >= threshold[tps] {
					log.Fatalln("The metrics is past threshold")
				}
			}

			if currentRetry == 15 {
				break
			}
		}
	}
}

func sendStatsDStress(tps int) error {
	client, err := statsd.NewClient("127.0.0.1:8125", "test-client")

	// and handle any initialization errors
	if err != nil {
		return err
	}

	// make sure to close to clean up when done, to avoid leaks.
	defer client.Close()

	for time := 0; time < tps; time++ {
		client.Inc(fmt.Sprintf("statsd_%v", time), time, 1.0)
	}
	return nil
}
