package grafana

import (
	"github.com/linclaus/stock/pkg/gografana"
)

const (
	// (1s*2^(maxRetries-1))
	// 1s, 2s, 4s
	maxRetries = 3
)

var (
	folderId       int
	folderName     = "沪深股市"
	dataSourceName = "Prometheus"
	client         gografana.GrafanaClienter
	args           InitializationArguments
)
var GrafanaArgs = InitializationArguments{}

type InitializationArguments struct {
	GrafanaVer         string
	GrafanaAddr        string
	GrafanaAccessToken string
	GrafanaAdminUser   string
	GrafanaAdminPass   string
	GrafanaFolderId    int
	PrometheusAddr     string
	DataSourceName     string
}

func Initialize(arg InitializationArguments) error {
	args = arg
	dataSourceName = args.DataSourceName
	folderId = args.GrafanaFolderId
	var err error
	var auth gografana.Authenticator
	if args.GrafanaAccessToken != "" {
		auth = gografana.NewAPIKeyAuthenticator(args.GrafanaAccessToken)
	} else {
		auth = gografana.NewBasicAuthenticator(args.GrafanaAdminUser, args.GrafanaAdminPass)
	}
	client, err = gografana.GetClientByVersion(args.GrafanaVer, args.GrafanaAddr, auth)
	if err != nil {
		return err
	}

	err = findOrCreateFolder()
	if err != nil {
		return err
	}

	return nil
}
