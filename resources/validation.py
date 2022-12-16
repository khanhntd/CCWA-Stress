import boto3

client = boto3.client('cloudwatch')

response = client.delete_anomaly_detector(
    Namespace='string',
    MetricName='string',
    Dimensions=[
        {
            'Name': 'string',
            'Value': 'string'
        },
    ],
    Stat='string',
    SingleMetricAnomalyDetector={
        'Namespace': 'string',
        'MetricName': 'string',
        'Dimensions': [
            {
                'Name': 'string',
                'Value': 'string'
            },
        ],
        'Stat': 'string'
    },
    MetricMathAnomalyDetector={
        'MetricDataQueries': [
            {
                'Id': 'string',
                'MetricStat': {
                    'Metric': {
                        'Namespace': 'string',
                        'MetricName': 'string',
                        'Dimensions': [
                            {
                                'Name': 'string',
                                'Value': 'string'
                            },
                        ]
                    },
                    'Period': 123,
                    'Stat': 'string',
                    'Unit': 'Seconds'|'Microseconds'|'Milliseconds'|'Bytes'|'Kilobytes'|'Megabytes'|'Gigabytes'|'Terabytes'|'Bits'|'Kilobits'|'Megabits'|'Gigabits'|'Terabits'|'Percent'|'Count'|'Bytes/Second'|'Kilobytes/Second'|'Megabytes/Second'|'Gigabytes/Second'|'Terabytes/Second'|'Bits/Second'|'Kilobits/Second'|'Megabits/Second'|'Gigabits/Second'|'Terabits/Second'|'Count/Second'|'None'
                },
                'Expression': 'string',
                'Label': 'string',
                'ReturnData': True|False,
                'Period': 123,
                'AccountId': 'string'
            },
        ]
    }
)