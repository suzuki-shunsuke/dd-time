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
		os.Exit(1)
	}
}

type (
	options struct {
		Help          bool
		Version       bool
		MetricName    string
		MetricHost    string
		DataDogAPIKey string
		Tags          []string
		Args          []string
	}
)

func parseArgs() options {
	helpF := pflag.BoolP("help", "h", false, "Show this help message")
	verF := pflag.BoolP("version", "v", false, "Show the version")
	metricNameF := pflag.StringP("metric-name", "m", "command-execution-time", "The name of the time series")
	metricHostF := pflag.String("host", "", "The name of the host that produced the metric")
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
	})
}
