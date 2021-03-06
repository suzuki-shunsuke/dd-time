package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/suzuki-shunsuke/dd-time/pkg/constant"
	"github.com/suzuki-shunsuke/dd-time/pkg/controller"
	"github.com/suzuki-shunsuke/dd-time/pkg/signal"
)

func Core() int {
	opts := parseArgs()

	if opts.Help {
		fmt.Println(constant.Help)
		return 0
	}
	if opts.Version {
		fmt.Println(constant.Version)
		return 0
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go signal.Handle(cancel)

	ctrl := controller.New(opts.DataDogAPIKey)

	return ctrl.Main(ctx, controller.Params{
		DataDogAPIKey: opts.DataDogAPIKey,
		Args:          opts.Args,
		Tags:          opts.Tags,
		MetricName:    opts.MetricName,
		MetricHost:    opts.MetricHost,
		Output:        opts.Output,
		Append:        opts.Append,
	})
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
	metricNameF := pflag.StringP("metric-name", "m", "command_execution_time", "The name of the time series")
	metricHostF := pflag.String("host", "", "The name of the host that produced the metric")
	outputDDTimeF := pflag.StringP("output", "o", "", "The file path where the dd-time's standard error output is written")
	appendF := pflag.BoolP("append", "a", false, "Write the dd-time's standard error output by the appended mode")
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
