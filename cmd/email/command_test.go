package main

import (
	"io/ioutil"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain2(t *testing.T) {
	t.Parallel()

	returnValues := &ClientParam{}
	returnValues.customers = toStringPtr("../../test-data/customers.csv")
	returnValues.template = toStringPtr("../../test-data/template")
	returnValues.successOutFile = toStringPtr("../../test-data/output_email.json")
	returnValues.errOutFile = toStringPtr("../../test-data/error.csv")
	main2(*returnValues, "send")
	errBytes, err := ioutil.ReadFile("../../test-data/error.csv")
	if err != nil {
		log.Fatal("No error csv created")
	}
	successBytes, err := ioutil.ReadFile("../../test-data/output_email.json")
	if err != nil {
		log.Fatal("No emails created")
	}
	assert.True(t, strings.Contains(string(errBytes), "Mr,Sumit,Narayan,"))
	assert.True(t, strings.Contains(string(successBytes), "Hi Mr John Smith"))
	assert.True(t, strings.Contains(string(successBytes), "Hi Mrs Michelle Smith"))
}

func TestMain2Error(t *testing.T) {
	t.Parallel()

	returnValues := &ClientParam{}
	returnValues.customers = toStringPtr("customers.csv")
	returnValues.template = toStringPtr("template")
	returnValues.successOutFile = toStringPtr("output_email.json")
	returnValues.errOutFile = toStringPtr("error.csv")
	ret := main2(*returnValues, "send")
	assert.Equal(t, 0, ret)
	_, err := ioutil.ReadFile("error.csv")
	require.NotNil(t, err)
	_, err = ioutil.ReadFile("output_email.json")
	require.NotNil(t, err)
}

func toStringPtr(str string) *string {
	return &str
}
