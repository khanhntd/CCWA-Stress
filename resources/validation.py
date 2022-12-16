import boto3

client = boto3.client('cloudwatch')

response = client.describe_alarms_for_metric(
    MetricName='string',
    Namespace='string',
    Statistic='SampleCount'|'Average'|'Sum'|'Minimum'|'Maximum',
    ExtendedStatistic='string',
    Dimensions=[
        {
            'Name': 'string',
            'Value': 'string'
        },
    ],
    Period=123,
    Unit='Seconds'|'Microseconds'|'Milliseconds'|'Bytes'|'Kilobytes'|'Megabytes'|'Gigabytes'|'Terabytes'|'Bits'|'Kilobits'|'Megabits'|'Gigabits'|'Terabits'|'Percent'|'Count'|'Bytes/Second'|'Kilobytes/Second'|'Megabytes/Second'|'Gigabytes/Second'|'Terabytes/Second'|'Bits/Second'|'Kilobits/Second'|'Megabits/Second'|'Gigabits/Second'|'Terabits/Second'|'Count/Second'|'None'
)