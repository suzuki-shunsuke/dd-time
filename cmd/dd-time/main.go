package main

import (
	"errors"
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

func core() error {
	helpF := pflag.BoolP("help", "h", false, "Show this help message")
	verF := pflag.BoolP("version", "v", false, "Show the version")
	metricNameF := pflag.StringP("metric-name", "m", "command-execution-time", "The name of the time series")
	metricHostF := pflag.String("host", "", "The name of the host that produced the metric")
	tagsF := pflag.StringSliceP("tag", "t", nil, "DataDog tags. The format is 'key:value'")
	pflag.Parse()
	if *helpF {
		fmt.Println(constant.Help)
		return nil
	}
	if *verF {
		fmt.Println(constant.Version)
		return nil
	}
	ddAPIKey := os.Getenv("DATADOG_API_KEY")
	if ddAPIKey == "" {
		return errors.New("The environment variable 'DATADOG_API_KEY' is required")
	}

	return cmd.Main(cmd.Params{
		DataDogAPIKey: ddAPIKey,
		Args:          pflag.Args(),
		Tags:          *tagsF,
		MetricName:    *metricNameF,
		MetricHost:    *metricHostF,
	})
}
