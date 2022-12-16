resource "aws_cloudwatch_metric_alarm" "foobar" {
  alarm_name                = "CWAgent-StressTest-CPUUsage"
  comparison_operator       = "GreaterThanOrEqualToThreshold"
  evaluation_periods        = "2"
  metric_name               = "CPUUtilization"
  namespace                 = "StressTest"
  period                    = "120"
  statistic                 = "Average"
  threshold                 = "80"
  alarm_description         = "Metric monitors CWAgent's CPU Usage"
  insufficient_data_actions = []
}