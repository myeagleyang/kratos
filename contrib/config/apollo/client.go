package apollo

import (
	"flag"
	"gitlab.wwgame.com/wwgame/kratos/v2/config"
	"os"
)

var (
	confAppID, confCluster, confCacheDir, confEndpoint, confNamespaces, confSecret string
)

func init() {
	addApolloFlags()
	addApolloEnvs()
}

func addApolloFlags() {
	flag.StringVar(&confAppID, "apollo.appid", "", "apollo app id")
	flag.StringVar(&confCluster, "apollo.cluster", "default", "apollo cluster")
	flag.StringVar(&confCacheDir, "apollo.cachedir", "/tmp/", "apollo cache dir")
	flag.StringVar(&confEndpoint, "apollo.endpoint", "", "apollo endpoint addr, e.g. http://localhost:8080")
	flag.StringVar(&confNamespaces, "apollo.namespaces", "", "subscribed apollo namespaces, comma separated, e.g. app.yml,mysql.yml")
	flag.StringVar(&confSecret, "apollo.secret", "", "access apollo with secret")
}
func addApolloEnvs() {
	if tempEnv := os.Getenv("APOLLO_APP_ID"); tempEnv != "" {
		confAppID = tempEnv
	}
	if tempEnv := os.Getenv("APOLLO_CLUSTER"); tempEnv != "" {
		confCluster = tempEnv
	}
	if tempEnv := os.Getenv("APOLLO_CACHE_DIR"); tempEnv != "" {
		confCacheDir = tempEnv
	}
	if tempEnv := os.Getenv("APOLLO_ENDPOINT"); tempEnv != "" {
		confEndpoint = tempEnv
	}
	if tempEnv := os.Getenv("APOLLO_NAMESPACE"); tempEnv != "" {
		confNamespaces = tempEnv
	}
	if tempEnv := os.Getenv("APOLLO_SECRET"); tempEnv != "" {
		confSecret = tempEnv
	}

}

func DefaultSource() config.Source {

	return NewSource(
		WithAppID(confAppID),
		WithCluster(confCluster),
		WithEndpoint(confEndpoint),
		WithNamespace(confNamespaces),
		WithBackupPath(confCacheDir),
		WithEnableBackup(),
		WithSecret(confSecret),
	)
}
