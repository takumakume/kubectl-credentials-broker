package cmd

import (
	"fmt"

	"github.com/kballard/go-shellquote"
	"github.com/spf13/cobra"
	"github.com/takumakume/kubectl-credentials-broker/credentials"
	"github.com/takumakume/kubectl-credentials-broker/kubeconfig"
)

var defaultExecAPIVersion = (&credentials.V1Beta1{}).APIVersionString()

type kubeconfigCmdArgs struct {
	execAPIVersion string
	env            map[string]string
	rootCmdArgs
}

var (
	argsKubeconfigClientCertificatePath string
	argsKubeconfigClientKeyPath         string
	argsKubeconfigTokenPath             string
	argsKubeconfigBeforeExecCommand     string
	argsKubeconfigExecAPIVersion        string
	argsKubeconfigEnv                   map[string]string
)

var configCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "kubeconfig",
	Long:  "kubeconfig",
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	Long:  "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := &kubeconfigCmdArgs{
			rootCmdArgs: rootCmdArgs{
				clientCertificatePath: argsKubeconfigClientCertificatePath,
				clientKeyPath:         argsKubeconfigClientKeyPath,
				tokenPath:             argsKubeconfigTokenPath,
				beforeExecCommand:     argsKubeconfigBeforeExecCommand,
			},
			execAPIVersion: argsKubeconfigExecAPIVersion,
			env:            argsKubeconfigEnv,
		}

		if err := opt.validate(); err != nil {
			return err
		}

		return kubeconfigSet(opt)
	},
}

func init() {
	configSetCmd.Flags().StringVarP(&argsKubeconfigClientCertificatePath, "client-certificate-path", "", "", "PEM-encoded client certificate file path. Can contain CA certificate. If this flag is specified, --client-key-path is also required. (optional)")
	configSetCmd.Flags().StringVarP(&argsKubeconfigClientKeyPath, "client-key-path", "", "", "PEM-encoded client key file path. (optional)")
	configSetCmd.Flags().StringVarP(&argsKubeconfigTokenPath, "token-path", "", "", "Token file path. (optional)")
	configSetCmd.Flags().StringVarP(&argsKubeconfigBeforeExecCommand, "before-exec-command", "", "", "A command line to run before responding to the credential plugin. For example, it can be used to update certificate and token files. (optional)")
	configSetCmd.Flags().StringVarP(&argsKubeconfigExecAPIVersion, "exec-api-version", "", defaultExecAPIVersion, fmt.Sprintf("API version to use when decoding the ExecCredentials resource (Default: %s)", defaultExecAPIVersion))
	configSetCmd.Flags().StringToStringVarP(&argsKubeconfigEnv, "env", "", map[string]string{}, "Environment variables to set when running the plugin. (optional) ex. 'HOGE=huga,FOO=bar'")
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}

func (args *kubeconfigCmdArgs) validate() error {
	if err := args.rootCmdArgs.validate(); err != nil {
		return err
	}
	return nil
}

func kubeconfigSet(args *kubeconfigCmdArgs) error {
	k := kubeconfig.New()
	cmd, cmdArgs, err := splitCommand(args.beforeExecCommand)
	if err != nil {
		return err
	}

	if err := k.UpdateCurrentUserExecConfig(args.execAPIVersion, cmd, cmdArgs, args.env); err != nil {
		return err
	}

	return nil
}

func splitCommand(commandline string) (string, []string, error) {
	parsedCmd, err := shellquote.Split(commandline)
	if err != nil {
		return "", []string{}, err
	}

	if len(parsedCmd) == 0 {
		return "", []string{}, nil
	} else if len(parsedCmd) == 1 {
		return parsedCmd[0], []string{}, nil
	} else {
		return parsedCmd[0], parsedCmd[1:], nil
	}
}
