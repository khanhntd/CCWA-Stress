// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/cactus/go-statsd-client/v5/statsd"
)

const (
	RetryTime = 15
)

var (
	ctx context.Context
	cwm *cloudwatch.Client
)

var threshold = map[int]float64{
	100:   67000000.0,
	1000:  90708352.0,
	10000: 120708352.0,
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
		time.Sleep(2 * time.Minute)
		log.Printf("Begin to send statsd metrics to CWA with number of metrics %d", tps)

		metricValues, err := GetMetricDataResults("StressTest", "procstat_memory_rss", dimension)
		if err != nil {
			log.Fatalf("Fail get metric data result because of %v", err)
		}

		if len(metricValues) == 0 {
			continue
		}

		for _, value := range metricValues {
			log.Printf("subtracting %v", value-threshold[tps])
			if value-threshold[tps] > 0 {
				log.Fatalln("The metrics is past threshold")
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
		client.Inc(fmt.Sprintf("statsd_%v", time), int64(time), 1.0)
	}
	return nil
}

func GetCloudWatchMetricsClient() (*cloudwatch.Client, *context.Context, error) {
	if cwm == nil {
		ctx = context.Background()
		c, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, nil, err
		}

		cwm = cloudwatch.NewFromConfig(c)
	}
	return cwm, &ctx, nil
}

func GetMetricDataResults(namespace, metricName string, metricSpecificDimensions []types.Dimension) ([]float64, error) {
	ec2InstanceId := GetInstanceId()
	instanceIdDimensions := types.Dimension{
		Name:  aws.String("InstanceId"),
		Value: aws.String(ec2InstanceId),
	}
	dimensions := append(metricSpecificDimensions, instanceIdDimensions)
	metricToFetch := types.Metric{
		Namespace:  aws.String(namespace),
		MetricName: aws.String(metricName),
		Dimensions: dimensions,
	}
	log.Printf("test %v %v %v %v", ec2InstanceId, &dimensions, namespace, metricName)
	metricQueryPeriod := int32(120)
	metricDataQueries := []types.MetricDataQuery{
		{
			MetricStat: &types.MetricStat{
				Metric: &metricToFetch,
				Period: &metricQueryPeriod,
				Stat:   aws.String("Average"),
			},
			Id: aws.String(strings.ToLower(metricName)),
		},
	}

	endTime := time.Now()
	startTime := subtractMinutes(endTime, 2)
	getMetricDataInput := cloudwatch.GetMetricDataInput{
		StartTime:         &startTime,
		EndTime:           &endTime,
		MetricDataQueries: metricDataQueries,
	}

	log.Printf("Metric data input is : %s", fmt.Sprint(getMetricDataInput))

	cwmClient, clientContext, err := GetCloudWatchMetricsClient()
	if err != nil {
		return nil, fmt.Errorf("Error occurred while creating CloudWatch client: %v", err.Error())
	}

	output, err := cwmClient.GetMetricData(*clientContext, &getMetricDataInput)
	if err != nil {
		return nil, fmt.Errorf("Error getting metric data %v", err)
	}

	result := output.MetricDataResults[0].Values
	log.Printf("Metric values are : %s", fmt.Sprint(result))
	return result, nil
}

func subtractMinutes(fromTime time.Time, minutes int) time.Time {
	tenMinutes := time.Duration(-1*minutes) * time.Minute
	return fromTime.Add(tenMinutes)
}

func GetInstanceId() string {
	ctx := context.Background()
	c, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		// fail fast so we don't continue the test
		log.Fatalf("Error occurred while creating SDK config: %v", err)
	}

	// TODO: this only works for EC2 based testing
	client := imds.NewFromConfig(c)
	metadata, err := client.GetInstanceIdentityDocument(ctx, &imds.GetInstanceIdentityDocumentInput{})
	if err != nil {
		log.Fatalf("Error occurred while retrieving EC2 instance ID: %v", err)
	}
	return metadata.InstanceID
}
