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
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

var (
	ctx context.Context
	cwm *cloudwatch.Client
)

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
	metricToFetch := types.Metric{
		Namespace:  aws.String(namespace),
		MetricName: aws.String(metricName),
		Dimensions: metricSpecificDimensions,
	}

	metricQueryPeriod := int32(60)
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
	startTime := subtractMinutes(endTime, 10)
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