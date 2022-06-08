package main

import "github.com/dhis2-sre/im-inspector/pkg/di"

func main() {
	environment := di.GetEnvironment()
	environment.Inspector.Inspect()
}
