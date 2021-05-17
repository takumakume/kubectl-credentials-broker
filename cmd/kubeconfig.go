package cmd

import (
	"fmt"

	"github.com/Songmu/prompter"
	"github.com/kballard/go-shellquote"
	"github.com/spf13/cobra"
	"github.com/takumakume/kubectl-credentials-broker/credentials"
	"github.com/takumakume/kubectl-credentials-broker/kubeconfig"
)

var defaultExecAPIVersion = (&credentials.V1Beta1{}).APIVersionString()

type kubeconfigCmdArgs struct {
	execAPIVersion string
	env            map[string]string
	force          bool
	rootCmdArgs
}

var (
	argsKubeconfigClientCertificatePath string
	argsKubeconfigClientKeyPath         string
	argsKubeconfigTokenPath             string
	argsKubeconfigBeforeExecCommand     string
	argsKubeconfigExecAPIVersion        string
	argsKubeconfigEnv                   map[string]string
	argsKubeconfigForce                 bool
)

var configCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "kubeconfig",
	Long:  "kubeconfig",
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "This command adds an exec command to the current-context user.",
	Long:  "This command adds an exec command to the current-context user.",
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
			force:          argsKubeconfigForce,
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
	configSetCmd.Flags().BoolVarP(&argsKubeconfigForce, "force", "f", false, "Do not confirm overwriting of kubeconfig (Default: false)")
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}

func (args *kubeconfigCmdArgs) validate() error {
	if err := args.rootCmdArgs.validate(); err != nil {
		return err
	}
	return nil
}

func (args *kubeconfigCmdArgs) makePluginCommand() ([]string, error) {
	c := []string{commandName}

	if len(args.beforeExecCommand) > 0 {
		quotedCmd, err := shellquote.Split(args.beforeExecCommand)
		if err != nil {
			return nil, err
		}
		c = append(c, "--before-exec-command")
		c = append(c, quotedCmd...)
	}
	if len(args.clientCertificatePath) > 0 {
		c = append(c, "--client-certificate-path", args.clientCertificatePath)
	}
	if len(args.clientKeyPath) > 0 {
		c = append(c, "--client-key-path", args.clientKeyPath)
	}
	if len(args.tokenPath) > 0 {
		c = append(c, "--token-path", args.tokenPath)
	}

	return c, nil
}

func kubeconfigSet(args *kubeconfigCmdArgs) error {
	k := kubeconfig.New()

	pluginCmd, err := args.makePluginCommand()
	if err != nil {
		return err
	}

	diff, err := k.UpdateCurrentUserExecConfigDryRun(args.execAPIVersion, "kubectl", pluginCmd, args.env)
	if err != nil {
		return err
	}

	if diff == "" {
		fmt.Println("current kubeconfig is up to date")
		return nil
	}

	fmt.Printf("---\n%s", diff)
	if args.force {
		if err := k.UpdateCurrentUserExecConfig(args.execAPIVersion, "kubectl", pluginCmd, args.env); err != nil {
			return err
		}
	} else {
		if prompter.YesNo("---\ncontinue? (y/N)", false) {
			if err := k.UpdateCurrentUserExecConfig(args.execAPIVersion, "kubectl", pluginCmd, args.env); err != nil {
				return err
			}
		} else {
			fmt.Println("---\ncanceled update kubeconfig")
			return nil
		}
	}

	fmt.Println("---\nupdate successful")
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
