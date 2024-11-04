// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/gcp"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type abapEnvironmentCloneGitRepoOptions struct {
	Username          string   `json:"username,omitempty"`
	Password          string   `json:"password,omitempty"`
	ByogUsername      string   `json:"byogUsername,omitempty"`
	ByogPassword      string   `json:"byogPassword,omitempty"`
	ByogAuthMethod    string   `json:"byogAuthMethod,omitempty" validate:"possible-values=TOKEN BASIC"`
	Repositories      string   `json:"repositories,omitempty"`
	RepositoryName    string   `json:"repositoryName,omitempty"`
	BranchName        string   `json:"branchName,omitempty"`
	Host              string   `json:"host,omitempty"`
	LogOutput         string   `json:"logOutput,omitempty" validate:"possible-values=ZIP STANDARD"`
	CfAPIEndpoint     string   `json:"cfApiEndpoint,omitempty"`
	CfOrg             string   `json:"cfOrg,omitempty"`
	CfSpace           string   `json:"cfSpace,omitempty"`
	CfServiceInstance string   `json:"cfServiceInstance,omitempty"`
	CfServiceKeyName  string   `json:"cfServiceKeyName,omitempty"`
	CertificateNames  []string `json:"certificateNames,omitempty"`
}

// AbapEnvironmentCloneGitRepoCommand Clones a git repository to a SAP BTP ABAP Environment system
func AbapEnvironmentCloneGitRepoCommand() *cobra.Command {
	const STEP_NAME = "abapEnvironmentCloneGitRepo"

	metadata := abapEnvironmentCloneGitRepoMetadata()
	var stepConfig abapEnvironmentCloneGitRepoOptions
	var startTime time.Time
	var logCollector *log.CollectorHook
	var splunkClient *splunk.Splunk
	telemetryClient := &telemetry.Telemetry{}

	var createAbapEnvironmentCloneGitRepoCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Clones a git repository to a SAP BTP ABAP Environment system",
		Long: `Clones a git repository (Software Component) to a SAP BTP ABAP Environment system. If the repository is already cloned, the step will checkout the configured branch and pull the specified commit, instead.
Please provide either of the following options:

* The host and credentials the BTP ABAP Environment system itself. The credentials must be configured for the Communication Scenario [SAP_COM_0948](https://help.sap.com/docs/sap-btp-abap-environment/abap-environment/api-for-managing-software-components-61f4d47af1394b1c8ad684b71d3ad6a0?locale=en-US).
* The Cloud Foundry parameters (API endpoint, organization, space), credentials, the service instance for the ABAP service and the service key for the Communication Scenario SAP_COM_0948.
* Only provide one of those options with the respective credentials. If all values are provided, the direct communication (via host) has priority.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, err := os.Getwd()
			if err != nil {
				return err
			}
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err = PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.Username)
			log.RegisterSecret(stepConfig.Password)
			log.RegisterSecret(stepConfig.ByogUsername)
			log.RegisterSecret(stepConfig.ByogPassword)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 || len(GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint) > 0 {
				splunkClient = &splunk.Splunk{}
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			if err = log.RegisterANSHookIfConfigured(GeneralConfig.CorrelationID); err != nil {
				log.Entry().WithError(err).Warn("failed to set up SAP Alert Notification Service log hook")
			}

			validation, err := validation.New(validation.WithJSONNamesForStructFields(), validation.WithPredefinedErrorMessages())
			if err != nil {
				return err
			}
			if err = validation.ValidateStruct(stepConfig); err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			vaultClient := config.GlobalVaultClient()
			if vaultClient != nil {
				defer vaultClient.MustRevokeToken()
			}

			stepTelemetryData := telemetry.CustomData{}
			stepTelemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				stepTelemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				stepTelemetryData.ErrorCategory = log.GetErrorCategory().String()
				stepTelemetryData.PiperCommitHash = GitCommit
				telemetryClient.SetData(&stepTelemetryData)
				telemetryClient.Send()
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.Dsn,
						GeneralConfig.HookConfig.SplunkConfig.Token,
						GeneralConfig.HookConfig.SplunkConfig.Index,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
				if len(GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint) > 0 {
					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblToken,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblIndex,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
				if GeneralConfig.HookConfig.GCPPubSubConfig.Enabled {
					err := gcp.NewGcpPubsubClient(
						vaultClient,
						GeneralConfig.HookConfig.GCPPubSubConfig.ProjectNumber,
						GeneralConfig.HookConfig.GCPPubSubConfig.IdentityPool,
						GeneralConfig.HookConfig.GCPPubSubConfig.IdentityProvider,
						GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.OIDCConfig.RoleID,
					).Publish(GeneralConfig.HookConfig.GCPPubSubConfig.Topic, telemetryClient.GetDataBytes())
					if err != nil {
						log.Entry().WithError(err).Warn("event publish failed")
					}
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetryClient.Initialize(GeneralConfig.NoTelemetry, STEP_NAME, GeneralConfig.HookConfig.PendoConfig.Token)
			abapEnvironmentCloneGitRepo(stepConfig, &stepTelemetryData)
			stepTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addAbapEnvironmentCloneGitRepoFlags(createAbapEnvironmentCloneGitRepoCmd, &stepConfig)
	return createAbapEnvironmentCloneGitRepoCmd
}

func addAbapEnvironmentCloneGitRepoFlags(cmd *cobra.Command, stepConfig *abapEnvironmentCloneGitRepoOptions) {
	cmd.Flags().StringVar(&stepConfig.Username, "username", os.Getenv("PIPER_username"), "User for either the Cloud Foundry API or the Communication Arrangement for SAP_COM_0948")
	cmd.Flags().StringVar(&stepConfig.Password, "password", os.Getenv("PIPER_password"), "Password for either the Cloud Foundry API or the Communication Arrangement for SAP_COM_0948")
	cmd.Flags().StringVar(&stepConfig.ByogUsername, "byogUsername", os.Getenv("PIPER_byogUsername"), "Username for bring your own git (BYOG) authentication")
	cmd.Flags().StringVar(&stepConfig.ByogPassword, "byogPassword", os.Getenv("PIPER_byogPassword"), "Password for bring your own git (BYOG) authentication")
	cmd.Flags().StringVar(&stepConfig.ByogAuthMethod, "byogAuthMethod", `TOKEN`, "Specifies which authentication method is used for bring your own git (BYOG) repositories")
	cmd.Flags().StringVar(&stepConfig.Repositories, "repositories", os.Getenv("PIPER_repositories"), "Specifies a YAML file containing the repositories configuration")
	cmd.Flags().StringVar(&stepConfig.RepositoryName, "repositoryName", os.Getenv("PIPER_repositoryName"), "Specifies a repository (Software Components) on the SAP BTP ABAP Environment system")
	cmd.Flags().StringVar(&stepConfig.BranchName, "branchName", os.Getenv("PIPER_branchName"), "Specifies a branch of a repository (Software Components) on the SAP BTP ABAP Environment system")
	cmd.Flags().StringVar(&stepConfig.Host, "host", os.Getenv("PIPER_host"), "Specifies the host address of the SAP BTP ABAP Environment system")
	cmd.Flags().StringVar(&stepConfig.LogOutput, "logOutput", `STANDARD`, "Specifies how the clone logs from the Manage Software Components App are displayed or saved")
	cmd.Flags().StringVar(&stepConfig.CfAPIEndpoint, "cfApiEndpoint", os.Getenv("PIPER_cfApiEndpoint"), "Cloud Foundry API Enpoint")
	cmd.Flags().StringVar(&stepConfig.CfOrg, "cfOrg", os.Getenv("PIPER_cfOrg"), "Cloud Foundry target organization")
	cmd.Flags().StringVar(&stepConfig.CfSpace, "cfSpace", os.Getenv("PIPER_cfSpace"), "Cloud Foundry target space")
	cmd.Flags().StringVar(&stepConfig.CfServiceInstance, "cfServiceInstance", os.Getenv("PIPER_cfServiceInstance"), "Cloud Foundry Service Instance")
	cmd.Flags().StringVar(&stepConfig.CfServiceKeyName, "cfServiceKeyName", os.Getenv("PIPER_cfServiceKeyName"), "Cloud Foundry Service Key")
	cmd.Flags().StringSliceVar(&stepConfig.CertificateNames, "certificateNames", []string{}, "file names of trusted (self-signed) server certificates - need to be stored in .pipeline/trustStore")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("byogUsername")
	cmd.MarkFlagRequired("byogPassword")
}

// retrieve step metadata
func abapEnvironmentCloneGitRepoMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "abapEnvironmentCloneGitRepo",
			Aliases:     []config.Alias{},
			Description: "Clones a git repository to a SAP BTP ABAP Environment system",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "abapCredentialsId", Description: "Jenkins credentials ID containing user and password to authenticate to the BTP ABAP Environment system or the Cloud Foundry API", Type: "jenkins", Aliases: []config.Alias{{Name: "cfCredentialsId", Deprecated: false}, {Name: "credentialsId", Deprecated: false}}},
					{Name: "byogCredentialsId", Description: "Jenkins credentials ID containing ByogUsername and ByogPassword to authenticate to a software component which is used in a BYOG scenario. (https://help.sap.com/docs/btp/sap-business-technology-platform/cloning-software-components-to-abap-environment-system-383ce2f9e2eb40f1b8ad538ddf79e656)", Type: "jenkins"},
				},
				Parameters: []config.StepParameters{
					{
						Name: "username",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "abapCredentialsId",
								Param: "username",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_username"),
					},
					{
						Name: "password",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "abapCredentialsId",
								Param: "password",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_password"),
					},
					{
						Name: "byogUsername",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "byogCredentialsId",
								Param: "username",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_byogUsername"),
					},
					{
						Name: "byogPassword",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "byogCredentialsId",
								Param: "password",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_byogPassword"),
					},
					{
						Name:        "byogAuthMethod",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `TOKEN`,
					},
					{
						Name:        "repositories",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_repositories"),
					},
					{
						Name:        "repositoryName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_repositoryName"),
					},
					{
						Name:        "branchName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_branchName"),
					},
					{
						Name:        "host",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_host"),
					},
					{
						Name:        "logOutput",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `STANDARD`,
					},
					{
						Name:        "cfApiEndpoint",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "cloudFoundry/apiEndpoint"}},
						Default:     os.Getenv("PIPER_cfApiEndpoint"),
					},
					{
						Name:        "cfOrg",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "cloudFoundry/org"}},
						Default:     os.Getenv("PIPER_cfOrg"),
					},
					{
						Name:        "cfSpace",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "cloudFoundry/space"}},
						Default:     os.Getenv("PIPER_cfSpace"),
					},
					{
						Name:        "cfServiceInstance",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "cloudFoundry/serviceInstance"}},
						Default:     os.Getenv("PIPER_cfServiceInstance"),
					},
					{
						Name:        "cfServiceKeyName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "cloudFoundry/serviceKey"}, {Name: "cloudFoundry/serviceKeyName"}, {Name: "cfServiceKey"}},
						Default:     os.Getenv("PIPER_cfServiceKeyName"),
					},
					{
						Name:        "certificateNames",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
					},
				},
			},
			Containers: []config.Container{
				{Name: "cf", Image: "ppiper/cf-cli:v12"},
			},
		},
	}
	return theMetaData
}
