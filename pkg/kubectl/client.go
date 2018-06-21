package kubectl

import (
	"fmt"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mad01/grafana-provisioner/pkg/kutil"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
)

func NewKubectlClient(kubeconfig string) *KubectlClient {
	config, err := K8sGetClientConfig(kubeconfig)
	if err != nil {
		panic(fmt.Sprintf("failed to get kube rest config: %v", err.Error()))
	}

	e := KubectlClient{
		waitInterval: 1 * time.Minute,
		ClientConfig: kutil.NewClientConfig(config, metav1.NamespaceAll),
	}

	return &e
}

// KubectlClient struct
type KubectlClient struct {
	waitInterval time.Duration
	ClientConfig clientcmd.ClientConfig
}

func (e *KubectlClient) Apply(path string) error {
	f := cmdutil.NewFactory(e.ClientConfig)

	options := &cmd.ApplyOptions{
		FilenameOptions: resource.FilenameOptions{
			Filenames: []string{path},
		},
		Cascade: true,
	}

	cobraCmd := &cobra.Command{
		Use: "apply",
	}
	cobraCmd.Flags().Bool("validate", true, "")
	cobraCmd.Flags().Bool("openapi-validation", true, "")
	cobraCmd.Flags().Bool("dry-run", false, "")
	cobraCmd.Flags().Bool("overwrite", true, "")
	cobraCmd.Flags().Bool("record", false, "")
	cobraCmd.Flags().String("schema-cache-dir", "", "")
	cobraCmd.Flags().String("output", "", "")

	err := cmd.RunApply(
		f,
		cobraCmd,
		os.Stdout,
		os.Stderr,
		options,
	)
	if err != nil {
		return fmt.Errorf("failed to run Apply: %v", err.Error())
	}

	return nil
}
