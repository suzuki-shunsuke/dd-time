package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/suzuki-shunsuke/dd-time/pkg/cmd"
	"github.com/suzuki-shunsuke/dd-time/pkg/constant"
)

func main() {
	if err := core(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cmd.GetExitCode(err))
	}
}

type (
	options struct {
		Help          bool
		Version       bool
		Append        bool
		MetricName    string
		MetricHost    string
		DataDogAPIKey string
		Output        string
		Tags          []string
		Args          []string
	}
)

func parseArgs() options {
	helpF := pflag.BoolP("help", "h", false, "Show this help message")
	verF := pflag.BoolP("version", "v", false, "Show the version")
	metricNameF := pflag.StringP("metric-name", "m", "command-execution-time", "The name of the time series")
	metricHostF := pflag.String("host", "", "The name of the host that produced the metric")
	outputDDTimeF := pflag.StringP("output", "o", "", "The file path where the dd-time's standard error output is written")
	appendF := pflag.BoolP("append", "a", false, "Write the error by the append mode")
	tagsF := pflag.StringSliceP("tag", "t", nil, "DataDog tags. The format is 'key:value'")
	pflag.Parse()
	return options{
		Help:          *helpF,
		Version:       *verF,
		MetricName:    *metricNameF,
		MetricHost:    *metricHostF,
		DataDogAPIKey: os.Getenv("DATADOG_API_KEY"),
		Tags:          *tagsF,
		Args:          pflag.Args(),
		Output:        *outputDDTimeF,
		Append:        *appendF,
	}
}

func core() error {
	opts := parseArgs()

	if opts.Help {
		fmt.Println(constant.Help)
		return nil
	}
	if opts.Version {
		fmt.Println(constant.Version)
		return nil
	}

	return cmd.Main(cmd.Params{
		DataDogAPIKey: opts.DataDogAPIKey,
		Args:          opts.Args,
		Tags:          opts.Tags,
		MetricName:    opts.MetricName,
		MetricHost:    opts.MetricHost,
		Output:        opts.Output,
		Append:        opts.Append,
	})
}
