package cmd

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/takumakume/kubectl-credentials-broker/credentials"
	"github.com/takumakume/kubectl-credentials-broker/kubeconfig"
)

var (
	argsClientCertificatePath string
	argsClientKeyPath         string
	argsTokenPath             string
	argsBeforeExecCommand     string
)

var rootCmd = &cobra.Command{
	Use:     "credentials-broker",
	Short:   "",
	Long:    "",
	Version: "0.0.1",
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := newRunner(&arguments{
			clientCertificatePath: argsClientCertificatePath,
			clientKeyPath:         argsClientKeyPath,
			tokenPath:             argsTokenPath,
			beforeExecCommand:     argsBeforeExecCommand,
		})
		if err != nil {
			return err
		}

		return r.run()
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
	rootCmd.Flags().StringVarP(&argsClientCertificatePath, "client-certificate-path", "", "", "")
	rootCmd.Flags().StringVarP(&argsClientKeyPath, "client-key-path", "", "", "")
	rootCmd.Flags().StringVarP(&argsTokenPath, "token-path", "", "", "")
	rootCmd.Flags().StringVarP(&argsBeforeExecCommand, "before-exec-command", "", "", "")
}

type runner struct {
	args *arguments
	cred credentials.Credentials
}

type arguments struct {
	clientCertificatePath string
	clientKeyPath         string
	tokenPath             string
	beforeExecCommand     string
}

func (args *arguments) validate() error {
	switch {
	case args.clientCertificatePath != "" && args.clientKeyPath != "":
		crt, err := ioutil.ReadFile(args.clientCertificatePath)
		if err != nil {
			return err
		}

		key, err := ioutil.ReadFile(args.clientKeyPath)
		if err != nil {
			return err
		}

		_, err = tls.LoadX509KeyPair(string(crt), string(key))
		if err != nil {
			return err
		}
	case args.clientCertificatePath != "" || args.clientKeyPath != "":
		return fmt.Errorf("both client-certificate-path and client-key-path must be provided")
	case args.tokenPath != "":
		if !fileExists(args.tokenPath) {
			return fmt.Errorf("token file not found: %s", args.tokenPath)
		}
	}

	return nil
}

func newRunner(args *arguments) (*runner, error) {
	if err := args.validate(); err != nil {
		return nil, err
	}

	kconfig, err := kubeconfig.New()
	if err != nil {
		return nil, err
	}

	ct, err := kconfig.ReadCurrentContext()
	if err != nil {
		return nil, err
	}

	user, err := kconfig.ReadUser(ct.AuthInfo)
	if err != nil {
		return nil, err
	}

	r := &runner{args: args}
	switch user.Exec.APIVersion {
	case "client.authentication.k8s.io/v1beta1":
		r.cred = &credentials.V1Beta1{}
	case "client.authentication.k8s.io/v1alpha1":
		r.cred = &credentials.V1Alpha1{}
	default:
		return nil, fmt.Errorf("Unsupported API Version: %s", user.Exec.APIVersion)
	}

	return r, nil
}

func (r *runner) run() error {
	if len(r.args.beforeExecCommand) > 0 {
		if err := execCommand(r.args.beforeExecCommand); err != nil {
			return err
		}
	}

	opts := credentials.CredentialOptions{}
	if len(r.args.clientCertificatePath) > 0 {
		buf, err := ioutil.ReadFile(r.args.clientCertificatePath)
		if err != nil {
			return err
		}
		opts.ClientCertificateData = string(buf)
	}
	if len(r.args.clientKeyPath) > 0 {
		buf, err := ioutil.ReadFile(r.args.clientKeyPath)
		if err != nil {
			return err
		}
		opts.ClientKeyData = string(buf)
	}
	if len(r.args.tokenPath) > 0 {
		buf, err := ioutil.ReadFile(r.args.tokenPath)
		if err != nil {
			return err
		}
		opts.Token = string(buf)

	}

	buf, err := r.cred.ToJSON(opts)
	if err != nil {
		return err
	}

	fmt.Println(string(buf))
	return nil
}

func execCommand(cmdline string) error {
	args := strings.Fields(cmdline)
	err := exec.Command(args[0], args[1:]...).Run()
	if err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)

	if pathError, ok := err.(*os.PathError); ok {
		if pathError.Err == syscall.ENOTDIR {
			return false
		}
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}
