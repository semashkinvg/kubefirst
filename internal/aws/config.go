package aws

import (
	"fmt"
	"os"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/kubefirst/kubefirst/pkg"
	"github.com/rs/zerolog/log"
)

const (
	CloudProvider          = "aws"
	GithubHost             = "github.com"
	GitlabHost             = "gitlab.com"
	KubectlClientVersion   = "v1.25.7"
	LocalhostOS            = runtime.GOOS
	LocalhostArch          = runtime.GOARCH
	RegionUsEast1          = "us-east-1"
	TerraformClientVersion = "1.3.8"
	ArgocdHelmChartVersion = "4.10.5"

	ArgocdPortForwardURL = pkg.ArgocdPortForwardURL
	VaultPortForwardURL  = pkg.VaultPortForwardURL
)

type AwsConfig struct {
	ArgoWorkflowsDir                string
	DestinationGitopsRepoHttpsURL   string
	DestinationGitopsRepoGitURL     string
	DestinationMetaphorRepoHttpsURL string
	DestinationMetaphorRepoGitURL   string
	GitopsDir                       string
	GitProvider                     string
	K1Dir                           string
	Kubeconfig                      string
	KubectlClient                   string
	KubefirstBotSSHPrivateKey       string
	KubefirstConfig                 string
	LogsDir                         string
	MetaphorDir                     string
	RegistryAppName                 string
	RegistryYaml                    string
	SSLBackupDir                    string
	TerraformClient                 string
	ToolsDir                        string
}

// todo move shared values to pkg. or break into common shared configs across git
// GetConfig - load default values from kubefirst installer
func GetConfig(clusterName string, domainName string, gitProvider string, gitOwner string) *AwsConfig {
	config := AwsConfig{}

	// todo do we want these from envs?
	if err := env.Parse(&config); err != nil {
		log.Fatal().Msgf(fmt.Sprintf("error reading environment variables %s", err.Error()))
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	// cGitHost describes which git host to use depending on gitProvider
	var cGitHost string
	switch gitProvider {
	case "github":
		cGitHost = GithubHost
	case "gitlab":
		cGitHost = GitlabHost
	}

	config.DestinationGitopsRepoHttpsURL = fmt.Sprintf("https://%s/%s/gitops.git", cGitHost, gitOwner)
	config.DestinationGitopsRepoGitURL = fmt.Sprintf("git@%s:%s/gitops.git", cGitHost, gitOwner)
	config.DestinationMetaphorRepoHttpsURL = fmt.Sprintf("https://%s/%s/metaphor.git", cGitHost, gitOwner)
	config.DestinationMetaphorRepoGitURL = fmt.Sprintf("git@%s:%s/metaphor.git", cGitHost, gitOwner)

	config.ArgoWorkflowsDir = fmt.Sprintf("%s/.k1/argo-workflows", homeDir)
	config.GitopsDir = fmt.Sprintf("%s/.k1/gitops", homeDir)
	config.GitProvider = gitProvider
	config.Kubeconfig = fmt.Sprintf("%s/.k1/kubeconfig", homeDir)
	config.K1Dir = fmt.Sprintf("%s/.k1", homeDir)
	config.KubectlClient = fmt.Sprintf("%s/.k1/tools/kubectl", homeDir)
	config.KubefirstConfig = fmt.Sprintf("%s/.k1/%s", homeDir, ".kubefirst")
	config.LogsDir = fmt.Sprintf("%s/.k1/logs", homeDir)
	config.MetaphorDir = fmt.Sprintf("%s/.k1/metaphor", homeDir)
	config.RegistryAppName = "registry"
	config.RegistryYaml = fmt.Sprintf("%s/.k1/gitops/registry/%s/registry.yaml", homeDir, clusterName)
	config.SSLBackupDir = fmt.Sprintf("%s/.k1/ssl/%s", homeDir, domainName)
	config.TerraformClient = fmt.Sprintf("%s/.k1/tools/terraform", homeDir)
	config.ToolsDir = fmt.Sprintf("%s/.k1/tools", homeDir)

	return &config
}

// todo standardize on field names
type GitOpsDirectoryValues struct {
	AlertsEmail               string
	AtlantisAllowList         string
	CloudProvider             string
	CloudRegion               string
	ClusterId                 string
	ClusterName               string
	ClusterType               string
	DomainName                string
	Kubeconfig                string
	KubefirstArtifactsBucket  string
	KubefirstStateStoreBucket string
	KubefirstTeam             string
	KubefirstVersion          string

	ArgoCDIngressURL               string
	ArgoCDIngressNoHTTPSURL        string
	ArgoWorkflowsIngressURL        string
	ArgoWorkflowsIngressNoHTTPSURL string
	AtlantisIngressURL             string
	AtlantisIngressNoHTTPSURL      string
	ChartMuseumIngressURL          string
	VaultIngressURL                string
	VaultIngressNoHTTPSURL         string
	VouchIngressURL                string

	AtlantisWebhookURL   string
	AwsIamArnAccountRoot string
	AwsKmsKeyId          string
	AwsNodeCapacityType  string
	AwsAccountID         string

	GitDescription       string
	GitNamespace         string
	GitProvider          string
	GitRunner            string
	GitRunnerDescription string
	GitRunnerNS          string
	GitURL               string

	GitHubHost  string
	GitHubOwner string
	GitHubUser  string

	GitlabHost         string
	GitlabOwner        string
	GitlabOwnerGroupID int
	GitlabUser         string

	GitOpsRepoAtlantisWebhookURL string
	GitOpsRepoGitURL             string
	GitOpsRepoNoHTTPSURL         string

	ContainerRegistryURL string
	UseTelemetry         string
}

type MetaphorTokenValues struct {
	ClusterName                   string
	CloudRegion                   string
	ContainerRegistryURL          string
	DomainName                    string
	MetaphorDevelopmentIngressURL string
	MetaphorStagingIngressURL     string
	MetaphorProductionIngressURL  string
}
