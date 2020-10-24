package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/springwiz/emailtemplates/sender"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Sender abstracts the email sender.
type Sender interface {
	SendEmail(customers []sender.Customer, emailParams ...string) error
}

// ClientParam houses the parameters captured from the cmd line.
type ClientParam struct {
	customers, template, successOutFile, errOutFile *string
}

func main() {
	cmdClient := kingpin.New("email", "Email command line client")
	cmdParams := configureCmdline(cmdClient)
	selectedCommand, err := cmdClient.Parse(os.Args[1:])
	if err != nil {
		log.Error(err)
		return
	}
	main2(*cmdParams, selectedCommand)
}

func main2(returnValues ClientParam, selectedCommand string) int {
	switch selectedCommand {
	case "send":
		var emailSender Sender = sender.NewFile()
		customers, err := readCsvFile(*returnValues.customers)
		if err != nil {
			log.Errorf("Error reading customer csv %s", err.Error())
		}
		tmpl, err := ioutil.ReadFile(*returnValues.template)
		if err != nil {
			log.Errorf("Error reading email template %s", err.Error())
		}
		err = emailSender.SendEmail(customers, string(tmpl), *returnValues.successOutFile,
			*returnValues.errOutFile)
		switch e := errors.Cause(err).(type) {
		case sender.CsvFileError:
			log.Errorln("Error writing csv file", e.Error())
			return 1
		case sender.TransformError:
			log.Errorln("Error generating email content", e.Error())
			return 1
		}
		return 0
	default:
		log.Error("unexpected command")
		return 1
	}
}

func readCsvFile(filePath string) ([]sender.Customer, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(sender.CsvFileError{
			Context: fmt.Sprintf("Read file: %s", filePath),
			Err:     err,
		}, "Reading csv failed")
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(sender.CsvFileError{
			Context: fmt.Sprintf("Read file: %s", filePath),
			Err:     err,
		}, "Reading csv failed")
	}
	var customers []sender.Customer
	if len(records) > 1 {
		for i := 0; i < len(records); i++ {
			customers = append(customers, sender.Customer{
				Title:     records[i][0],
				FirstName: records[i][1],
				LastName:  records[i][2],
				Email:     records[i][3],
				Today:     time.Now().Format("02 Jan 2006"),
			})
		}
	}
	return customers, nil
}

func configureCmdline(cmdClient *kingpin.Application) *ClientParam {
	sendCmd := cmdClient.Command("send", "Send emails")
	returnValues := &ClientParam{}
	returnValues.customers = sendCmd.Flag("data", "enter the path to source csv)").Short('d').
		Required().String()
	returnValues.template = sendCmd.Flag("template", "enter the path to template").Short('t').
		Required().String()
	returnValues.successOutFile = sendCmd.Flag("out", "enter the path to output file").Short('o').
		Default("output_email.json").String()
	returnValues.errOutFile = sendCmd.Flag("err", "enter the path to error file").Short('e').
		Default("errors.csv").String()
	return returnValues
}
