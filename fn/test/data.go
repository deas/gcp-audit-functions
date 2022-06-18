package test

import (
	_ "embed"
)

// https://github.com/golang/go/issues/46056
//go:embed audit-iam-service-account-compute-create.json
var ServiceAccountCreateJson string

//go:embed asset-search-start-stop.json
var AssetSearchStartStopJson string
