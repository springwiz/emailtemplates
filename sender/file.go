package sender

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pingcap/errors"
)

// Customer houses the data captured from the csv file.
type Customer struct {
	Title, FirstName, LastName, Email, Today string
}

// CsvFileError represents the errors thrown while parsing csv file.
type CsvFileError struct {
	Context string
	Err     error
}

func (e CsvFileError) Error() string { return e.Context + ": " + e.Err.Error() }

// TransformError represents the errors thrown while generating email content.
type TransformError struct {
	Context string
	Err     error
}

func (e TransformError) Error() string { return e.Context + ": " + e.Err.Error() }

// InvalidArgError represents the errors thrown because of inavlid args.
type InvalidArgError struct {
	Context string
	Err     error
}

func (e InvalidArgError) Error() string { return e.Context + ": " + e.Err.Error() }

// File encapsulates the file based Email Sender.
type File func([]Customer, ...string) error

// SendEmail generates and sends the emails.
func (f File) SendEmail(customers []Customer, emailParams ...string) error {
	return f(customers, emailParams...)
}

// NewFile creates a new File based Email Sender.
func NewFile() File {
	return func(customers []Customer, emailParams ...string) error {
		if len(emailParams) < 3 {
			return errors.Wrap(InvalidArgError{
				Context: fmt.Sprintf("No of arguments passed < 3: %d", len(emailParams)),
			}, "Passed arguments invalid")
		}
		var validCustRecs, invalidCustRecs []Customer
		for i := 1; i < len(customers); i++ {
			if len(customers[i].Email) > 0 {
				validCustRecs = append(validCustRecs, customers[i])
			} else {
				invalidCustRecs = append(invalidCustRecs, customers[i])
			}
		}
		if len(invalidCustRecs) > 0 {
			err := writeCsvFile(invalidCustRecs, emailParams[2])
			if err != nil {
				return errors.Wrap(CsvFileError{
					Context: fmt.Sprintf("Write file: %s", emailParams[2]),
					Err:     err,
				}, "Writing csv failed")
			}
		}
		if len(validCustRecs) > 0 {
			tmpl := template.Must(template.New("emailTemplate").Parse(emailParams[0]))
			data := struct {
				Customers []Customer
			}{validCustRecs}

			var b bytes.Buffer
			outBytes := bufio.NewWriter(&b)
			err := tmpl.Execute(outBytes, data)
			if err != nil {
				return errors.Wrap(TransformError{
					Context: "Execute transform failed",
					Err:     err,
				}, "Error executing transform")
			}
			outBytes.Flush()
			outFile, err := os.Create(emailParams[1])
			if err != nil {
				return errors.Wrap(TransformError{
					Context: fmt.Sprintf("Create file: %s", emailParams[1]),
					Err:     err,
				}, "Writing emails failed")
			}
			content := strings.ReplaceAll(b.String(), ",]", "]")
			_, err = outFile.WriteString(content)
			if err != nil {
				return errors.Wrap(TransformError{
					Context: fmt.Sprintf("Write file: %s", emailParams[2]),
					Err:     err,
				}, "Writing emails failed")
			}
		}
		return nil
	}
}

func writeCsvFile(customers []Customer, errOutFile string) error {
	var csvWriter strings.Builder
	csvWriter.WriteString(`TITLE,FIRST_NAME,LAST_NAME,EMAIL
`)
	for _, cust := range customers {
		_, err := csvWriter.WriteString(cust.Title + `,` + cust.FirstName + `,` + cust.LastName + `,` + cust.Email)
		if err != nil {
			return err
		}
	}
	err := ioutil.WriteFile(errOutFile, []byte(csvWriter.String()), 0600)
	if err != nil {
		return err
	}
	return nil
}
