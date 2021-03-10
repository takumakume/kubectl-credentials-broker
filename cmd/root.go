package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

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
		r, err := newRunner()
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
	cred credentials.Credentials
}

func newRunner() (*runner, error) {
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

	r := &runner{}
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
	if len(argsBeforeExecCommand) > 0 {
		if err := execCommand(argsBeforeExecCommand); err != nil {
			return err
		}
	}

	opts := credentials.CredentialOptions{}
	if len(argsClientCertificatePath) > 0 {
		buf, err := ioutil.ReadFile(argsClientCertificatePath)
		if err != nil {
			return err
		}
		opts.ClientCertificateData = string(buf)
	}
	if len(argsClientKeyPath) > 0 {
		buf, err := ioutil.ReadFile(argsClientKeyPath)
		if err != nil {
			return err
		}
		opts.ClientKeyData = string(buf)
	}
	if len(argsTokenPath) > 0 {
		buf, err := ioutil.ReadFile(argsTokenPath)
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
