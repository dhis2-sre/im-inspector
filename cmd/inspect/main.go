package main

import "github.com/dhis2-sre/instance-inspector/pgk/di"

func main() {
	environment := di.GetEnvironment()
	environment.Inspector.Inspect()
}
