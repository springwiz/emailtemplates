package sender

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	t.Parallel()

	customers := []Customer{
		{`TITLE`, `FIRST_NAME`, `LAST_NAME`, `EMAIL`, ``},
		{`Mr`, `John`, `Smith`, `john.smith@example.com`, ``},
		{`Mrs`, `Michelle`, `Smith`, `michelle.smith@example.com`, ``},
		{`Mr`, `Sumit`, `Narayan`, ``, ``},
	}
	bytes, err := ioutil.ReadFile("../test-data/template")
	require.Nil(t, err)
	err = NewFile().SendEmail(customers, string(bytes), "../test-data/output_email.json", "../test-data/error.csv")
	require.Nil(t, err)
}

func TestSendEmailInvalidArgs(t *testing.T) {
	t.Parallel()

	customers := []Customer{
		{`TITLE`, `FIRST_NAME`, `LAST_NAME`, `EMAIL`, ``},
		{`Mr`, `John`, `Smith`, `john.smith@example.com`, ``},
		{`Mrs`, `Michelle`, `Smith`, `michelle.smith@example.com`, ``},
		{`Mr`, `Sumit`, `Narayan`, ``, ``},
	}
	bytes, err := ioutil.ReadFile("../test-data/template")
	require.Nil(t, err)
	err = NewFile().SendEmail(customers, string(bytes))
	require.NotNil(t, err)
}
