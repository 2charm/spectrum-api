package main

import "github.com/2charm/spectrum-api/pkg/util"

func main() {
	addr := util.GetEnvironmentVariable("ADDR")
	apikey := util.GetEnvironmentVariable("APIKEY")
}
