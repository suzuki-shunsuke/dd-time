package constant

const (
	Help = `dd-time - post the command execution time as time-series data to DataDog

https://github.com/suzuki-shunsuke/dd-time

USAGE:
  dd-time [options] -- command

ENVIRONMENT VARIABLE

  DATADOG_API_KEY - DataDog API Key. If DATADOG_API_KEY isn't set, the metrics can't be sent to DataDog but the command is run normally

OPTIONS:
  --help, -h                     show help
  --version, -v                  print the version
  --metric-name value, -m value  (default: "command-execution-time") The name of the time series
  --host value                   The name of the host that produced the metric
  --tag value, -t value          A tag associated with the metric
  --output value, -o value       The file path where the dd-time's standard error output is written
  --append, -a                   Write the dd-time's standard error output by the appended mode`
)
