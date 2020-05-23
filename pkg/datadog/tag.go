package datadog

import (
	"os"
	"strings"
)

func GetTags() []string {
	if os.Getenv("CIRCLECI") == "true" {
		return getCircleCITags()
	}
	return nil
}

func getCircleCITags() []string {
	// https://circleci.com/docs/2.0/env-vars/#built-in-environment-variables
	envs := []string{
		"CIRCLECI",
		"CIRCLE_BRANCH",
		"CIRCLE_BUILD_NUM",
		"CIRCLE_BUILD_URL",
		"CIRCLE_JOB",
		"CIRCLE_NODE_INDEX",
		"CIRCLE_NODE_TOTAL",
		"CIRCLE_PR_NUMBER",
		"CIRCLE_PR_REPONAME",
		"CIRCLE_PR_USERNAME",
		"CIRCLE_PROJECT_REPONAME",
		"CIRCLE_PROJECT_USERNAME",
		"CIRCLE_REPOSITORY_URL",
		"CIRCLE_SHA1",
		"CIRCLE_TAG",
		"CIRCLE_USERNAME",
		"CIRCLE_WORKFLOW_ID",
	}
	arr := make([]string, len(envs))
	for i, e := range envs {
		arr[i] = strings.ToLower(e) + ":" + os.Getenv(e)
	}
	return arr
}
