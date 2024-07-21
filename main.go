package main

import (
	"WebpageArchiver/api"
	"WebpageArchiver/common"
)

func main() {
	common.ParseFlag()
	api.WebStarter(common.DEBUG)
}
