# fern-junit-client
A cli that can read junit reports and send them to a fern reporter instance in a format it understands

## Introduction
If you don't know what fern is, [check it out here!](https://github.com/Guidewire/fern-reporter) This repository is a command line tool to enable users to send their test results from other formats than Ginkgo to the Fern reporter. Currently only JUnit is supported but it could be easily extended to support other formats as well.

## Usage
Build the tool with the following command to create the `fern-junit-client` executable:
```bash
go build
```
Now send your test reports to Fern:
```bash
./fern-junit-client <test format> -u <Fern reporter URL> -n <project name> -d <test report dir>
```
For example:
```bash
./fern-junit-client junit -u "http://localhost:8080" -n "MyMicroservice" -d "/path/to/tests"
```

To see all available options, use `./fern-junit-client help`

## See Also
* [Fern UI](https://github.com/Guidewire/fern-ui)
* [Fern Reporter](https://github.com/Guidewire/fern-reporter)
* [Fern Ginkgo Client](https://github.com/guidewire-oss/fern-ginkgo-client)