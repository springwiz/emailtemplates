Promotion emails sender

Package implements a console based application which sends promotion emails to a client database passed to it in the form of csv file. It relies on OOB support in golang for html templating i.e. html/template to transform the given customer data into email content. Other templating engines like Ace could be chosen depending upon need.

The main features of the Promotion email sender are:
1. Uses the OOB html templating in golang to generate the email content.
2. Implements a command line api using kingpin having features like help, sensible defaults and validation.
3. The design is intentionally flexible to accomodate other REST/SMTP based implementations. The implementation
   currently supports a FILE based implementation. 
4. Writes the failed records to an errors.csv file.
5. Supports a robust error handling mechanism by defining new error implemenations and errors package to wrap them.
6. It uses go mod to build the solution.

Installation
1. Install the Go lang runtime on the local machine.
2. Create the path on the local machine at github.com/springwiz.
3. Change into springwiz and Run git clone https://github.com/springwiz/emailtemplates.git.
4. Change into emailtemplates and run go build ./cmd/email.
5. Run the test-data/test-commands. 

Assumptions
1. The Golang runtime is installed.
2. The Git cmd line is available.

Implementation details

1. "gopkg.in/alecthomas/kingpin.v2" is used to implement the command line client. The implementation is housed in cmd/email/command.go.2. The package relies on OOB Go Unit Testing for testing the solution.
3. The design is intentionally flexible and could accomodate REST/SMTP based email senders later. It only supports a FILE based seender for now(sender/file.go).

The GoLang stack was chosen from the following:
1. Python
2. Java
3. Golang

The following are the considered factors: The Python stack is more suitable for machine learning processes. The concurrency in python is severly limited by GIL. The Java needs a JVM for running and thus not ideal for building command line tools. Further concurrency in Java is generally more expensive than golang. The Java's model of 1 thread per request fails for applications as thread context switches take most of the time/memory. The Golang apps are light weight and offer rich support for concurrency via go routines and channels. The go routines are light weight and offer rich support for parellelism. Golang compiles to native directly and also cross compiles thus offers best support for building command line apps.

Improvements
1. Support REST/SMTP based email senders.
2. Chunk customers.csv file into smaller chunks and improve performance by employing goroutines.
3. Improve unit test code coverage.
4. Performance testing.
