package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/takumakume/kubectl-credentials-broker/credentials"
	"github.com/takumakume/kubectl-credentials-broker/kubeconfig"
)

var version = "dev"

var (
	argsClientCertificatePath string
	argsClientKeyPath         string
	argsTokenPath             string
	argsBeforeExecCommand     string
)

const commandName = "credentials-broker"

var rootCmd = &cobra.Command{
	Use:     "credentials-broker",
	Short:   "credentials-broker",
	Long:    "credentials-broker",
	Version: "0.0.1",
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := newRootCmdRunner(&rootCmdArgs{
			clientCertificatePath: argsClientCertificatePath,
			clientKeyPath:         argsClientKeyPath,
			tokenPath:             argsTokenPath,
			beforeExecCommand:     argsBeforeExecCommand,
		})
		if err != nil {
			return err
		}

		buf, err := r.run()
		if err != nil {
			return err
		}

		fmt.Printf("%s", string(buf))

		return nil
	},
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&argsClientCertificatePath, "client-certificate-path", "", "", "PEM-encoded client certificate file path. Can contain CA certificate. If this flag is specified, --client-key-path is also required. (optional)")
	rootCmd.Flags().StringVarP(&argsClientKeyPath, "client-key-path", "", "", "PEM-encoded client key file path. (optional)")
	rootCmd.Flags().StringVarP(&argsTokenPath, "token-path", "", "", "Token file path. (optional)")
	rootCmd.Flags().StringVarP(&argsBeforeExecCommand, "before-exec-command", "", "", "A command line to run before responding to the credential plugin. For example, it can be used to update certificate and token files. (optional)")
}

type rootCmdRunner struct {
	args *rootCmdArgs
	cred credentials.Credential
}

type rootCmdArgs struct {
	clientCertificatePath string
	clientKeyPath         string
	tokenPath             string
	beforeExecCommand     string
}

func (args *rootCmdArgs) validate() error {
	switch {
	case args.clientCertificatePath == "" && args.clientKeyPath == "" && args.tokenPath == "":
		return errors.New("requires either certificate token")
	case (args.clientCertificatePath != "" && args.clientKeyPath == "") || args.clientCertificatePath == "" && args.clientKeyPath != "":
		return fmt.Errorf("both client-certificate-path and client-key-path must be provided")
	}

	return nil
}

func newRootCmdRunner(args *rootCmdArgs) (*rootCmdRunner, error) {
	if err := args.validate(); err != nil {
		return nil, err
	}

	r := &rootCmdRunner{
		args: args,
	}

	execAPIVersion, err := kubeconfig.New().ReadCurrentUserExecVersion()
	if err != nil {
		return nil, err
	}

	switch execAPIVersion {
	case (&credentials.V1Beta1{}).APIVersionString():
		r.cred = &credentials.V1Beta1{}
	case (&credentials.V1Alpha1{}).APIVersionString():
		r.cred = &credentials.V1Alpha1{}
	default:
		return nil, fmt.Errorf("unsupported client authentication API version: %s", execAPIVersion)
	}

	return r, nil
}

func (r *rootCmdRunner) run() ([]byte, error) {
	if len(r.args.beforeExecCommand) > 0 {
		if err := execCommand(r.args.beforeExecCommand); err != nil {
			return nil, err
		}
	}

	opt, err := makeCredentialOptions(r.args)
	if err != nil {
		return nil, err
	}

	buf, err := r.cred.ToJSON(opt)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func execCommand(cmdline string) error {
	args := strings.Fields(cmdline)
	err := exec.Command(args[0], args[1:]...).Run()
	if err != nil {
		return err
	}

	return nil
}

func makeCredentialOptions(args *rootCmdArgs) (*credentials.CredentialOption, error) {
	opts := &credentials.CredentialOption{}

	if len(args.clientCertificatePath) > 0 {
		buf, err := ioutil.ReadFile(args.clientCertificatePath)
		if err != nil {
			return nil, err
		}
		opts.ClientCertificateData = string(buf)
	}

	if len(args.clientKeyPath) > 0 {
		buf, err := ioutil.ReadFile(args.clientKeyPath)
		if err != nil {
			return nil, err
		}
		opts.ClientKeyData = string(buf)
	}

	if len(args.tokenPath) > 0 {
		buf, err := ioutil.ReadFile(args.tokenPath)
		if err != nil {
			return nil, err
		}
		opts.Token = string(buf)
	}

	return opts, nil
}
