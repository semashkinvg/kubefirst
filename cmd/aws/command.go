package aws

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Create
	alertsEmailFlag          string
	cloudRegionFlag          string
	clusterNameFlag          string
	clusterTypeFlag          string
	dryRun                   bool
	githubOrgFlag            string
	gitlabGroupFlag          string
	gitProviderFlag          string
	gitopsTemplateURLFlag    string
	gitopsTemplateBranchFlag string
	domainNameFlag           string
	kbotPasswordFlag         string
	useTelemetryFlag         bool

	// Supported git providers
	supportedGitProviders = []string{"github", "gitlab"}
)

func NewCommand() *cobra.Command {

	awsCmd := &cobra.Command{
		Use:   "aws",
		Short: "kubefirst aws installation",
		Long:  "kubefirst aws",
	}

	// on error, doesn't show helper/usage
	awsCmd.SilenceUsage = true

	// wire up new commands
	awsCmd.AddCommand(Create(), Destroy(), Quota())

	return awsCmd
}

func Create() *cobra.Command {
	createCmd := &cobra.Command{
		Use:              "create",
		Short:            "create the kubefirst platform running in aws",
		TraverseChildren: true,
		RunE:             createAws,
	}

	// todo review defaults and update descriptions
	createCmd.Flags().StringVar(&alertsEmailFlag, "alerts-email", "", "email address for let's encrypt certificate notifications (required)")
	createCmd.MarkFlagRequired("alerts-email")
	createCmd.Flags().StringVar(&cloudRegionFlag, "cloud-region", "us-east-1", "the aws region to provision infrastructure in")
	createCmd.Flags().StringVar(&clusterNameFlag, "cluster-name", "kubefirst", "the name of the cluster to create")
	createCmd.Flags().StringVar(&clusterTypeFlag, "cluster-type", "mgmt", "the type of cluster to create (i.e. mgmt|workload)")
	createCmd.Flags().StringVar(&domainNameFlag, "domain-name", "", "the Route53 hosted zone name to use for DNS records (i.e. your-domain.com|subdomain.your-domain.com) (required)")
	createCmd.MarkFlagRequired("domain-name")
	createCmd.Flags().BoolVar(&dryRun, "dry-run", false, "don't execute the installation")
	createCmd.Flags().StringVar(&gitProviderFlag, "git-provider", "github", fmt.Sprintf("the git provider - one of: %s", supportedGitProviders))
	createCmd.Flags().StringVar(&githubOrgFlag, "github-org", "", "the GitHub organization for the new gitops and metaphor repositories - required if using github")
	createCmd.Flags().StringVar(&gitlabGroupFlag, "gitlab-group", "", "the GitLab group for the new gitops and metaphor projects - required if using gitlab")
	createCmd.Flags().StringVar(&gitopsTemplateBranchFlag, "gitops-template-branch", "main", "the branch to clone for the gitops-template repository")
	createCmd.Flags().StringVar(&gitopsTemplateURLFlag, "gitops-template-url", "https://github.com/kubefirst/gitops-template.git", "the fully qualified url to the gitops-template repository to clone")
	createCmd.Flags().StringVar(&kbotPasswordFlag, "kbot-password", "", "the default password to use for the kbot user")
	createCmd.Flags().BoolVar(&useTelemetryFlag, "use-telemetry", true, "whether to emit telemetry")

	return createCmd
}

func Destroy() *cobra.Command {
	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy the kubefirst platform",
		Long:  "deletes the GitHub resources, aws resources, and local content to re-provision",
		RunE:  destroyAws,
	}

	return destroyCmd

}

func Quota() *cobra.Command {
	quotaCmd := &cobra.Command{
		Use:   "quota",
		Short: "Check aws quota status",
		Long:  "Check aws quota status. By default, only ones close to limits will be shown.",
		RunE:  evalAwsQuota,
	}

	quotaCmd.Flags().StringVar(&cloudRegionFlag, "cloud-region", "us-east-1", "the aws region to provision infrastructure in")

	return quotaCmd
}
