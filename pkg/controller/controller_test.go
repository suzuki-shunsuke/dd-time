package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateParams(t *testing.T) {
	data := []struct {
		title  string
		isErr  bool
		params Params
	}{
		{
			title: "DATADOG_API_KEY is required",
			isErr: true,
		},
		{
			title: "executed command isn't passed to dd-time",
			isErr: true,
			params: Params{
				DataDogAPIKey: "xxx",
			},
		},
		{
			title: "success",
			params: Params{
				DataDogAPIKey: "xxx",
				Args:          []string{"sleep", "5"},
			},
		},
	}
	ctrl := &Controller{}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			err := ctrl.validateParams(d.params)
			if d.isErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
