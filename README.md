# dd-time

[![Build Status](https://cloud.drone.io/api/badges/suzuki-shunsuke/dd-time/status.svg)](https://cloud.drone.io/suzuki-shunsuke/dd-time)
[![codecov](https://codecov.io/gh/suzuki-shunsuke/dd-time/branch/master/graph/badge.svg)](https://codecov.io/gh/suzuki-shunsuke/dd-time)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/dd-time)](https://goreportcard.com/report/github.com/suzuki-shunsuke/dd-time)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/dd-time.svg)](https://github.com/suzuki-shunsuke/dd-time)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/dd-time/master/LICENSE)

CLI tool to post the command execution time as time-series data to DataDog.

`dd-time` is inspired to [circle-dd-bench](https://github.com/yuya-takeyama/circle-dd-bench).

## Install

Download the binary from [GitHub Release](https://github.com/suzuki-shunsuke/dd-time/releases).

## Getting Started

At first, please [prepare a Datadog API key](https://docs.datadoghq.com/account_management/api-app-keys/) and set the key as the environment variable `DATADOG_API_KEY`.

Let's try to use `dd-time`.

```
$ dd-time -t command:tutorial-dd-time -- sleep 5
```

## Usage

```
$ dd-time --help
dd-time - post the command execution time as time-series data to DataDog

https://github.com/suzuki-shunsuke/dd-time

USAGE:
  dd-time [options] -- command

ENVIRONMENT VARIABLE

  DATADOG_API_KEY - DataDog API Key. If DATADOG_API_KEY isn't set, the metrics can't be sent to DataDog but the command is run normally

OPTIONS:
  --help, -h                     show help
  --version, -v                  print the version
  --metric-name value, -m value  (default: "command_execution_time") The name of the time series
  --host value                   The name of the host that produced the metric
  --tag value, -t value          A tag associated with the metric
  --output value, -o value       The file path where the dd-time's standard error output is written
  --append, -a                   Write the dd-time's standard error output by the appended mode
```

## Tags for CircleCI

In CircleCI, there are many built-in environment variables.

https://circleci.com/docs/2.0/env-vars/#built-in-environment-variables

`dd-time` sets the following environment variables as tags automatically.

* CIRCLECI
* CIRCLE_BRANCH
* CIRCLE_BUILD_NUM
* CIRCLE_BUILD_URL
* CIRCLE_JOB
* CIRCLE_NODE_INDEX
* CIRCLE_NODE_TOTAL
* CIRCLE_PR_NUMBER
* CIRCLE_PR_REPONAME
* CIRCLE_PR_USERNAME
* CIRCLE_PROJECT_REPONAME
* CIRCLE_PROJECT_USERNAME
* CIRCLE_REPOSITORY_URL
* CIRCLE_SHA1
* CIRCLE_TAG
* CIRCLE_USERNAME
* CIRCLE_WORKFLOW_ID

## License

[MIT](LICENSE)
