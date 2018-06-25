package main

import (
	"fmt"

	"os"

	"github.com/mad01/grafana-provisioner/pkg/kubectl"
	"github.com/spf13/cobra"
)

func cmdVersion() *cobra.Command {
	var command = &cobra.Command{
		Use:   "version",
		Short: "get version",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getVersion())
		},
	}
	return command
}

func cmdProvision() *cobra.Command {
	var kubeconfig, dbDNS, dbPass, dbUser, team, image, ingressDNSPrefix string
	var dbPort int
	var command = &cobra.Command{
		Use:   "provision",
		Short: "provision grafana and db",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

			dbConnStr := func(dns, pass, user string, port int) string {
				return fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, pass, dns, port)
			}

			dbGrafanaStr := func(dns, pass, user, db string, port int) string {
				return fmt.Sprintf("mysql://%s:%s@%s:%d/%s", user, pass, dns, port, db)
			}

			db := DB{
				URL: dbConnStr(dbDNS, dbPass, dbUser, dbPort),
			}
			err := db.connect()
			defer db.conn.Close()
			errCheck(err)

			err = db.createDB(team)
			errCheck(err)

			name := fmt.Sprintf("grafana-%s", team)
			values := manifestValues{
				databaseURL: dbGrafanaStr(dbDNS, dbPass, dbUser, team, dbPort),

				ingressClass: "nginx",
				ingressHost:  fmt.Sprintf("%s.%s", team, ingressDNSPrefix),
				ingressName:  name,

				serviceName: name,
				namespace:   team,
				image:       image,

				deploymentName:       name,
				deploymentLabelKey:   "app",
				deploymentLabelValue: name,
			}

			manifest := manifestRender(values)

			kubectl := kubectl.NewKubectlClient(kubeconfig)
			err = kubectl.Apply(manifest)
			errCheck(err)
		},
	}

	command.Flags().StringVarP(&kubeconfig, "kube.config", "k", "", "outside cluster path to kube config")
	command.Flags().StringVarP(&dbDNS, "db.dns", "d", "localhost", "mysql database dns")
	command.Flags().IntVarP(&dbPort, "db.port", "P", 3306, "mysql database port")
	command.Flags().StringVarP(&dbPass, "db.pass", "p", "", "mysql database password")
	command.Flags().StringVarP(&dbUser, "db.user", "u", "", "mysql database username")
	command.Flags().StringVarP(&team, "team", "t", "foo", "team name")
	command.Flags().StringVarP(&image, "image", "i", "grafana/grafana:5.1.0", "grafana official container image")
	command.Flags().StringVarP(&ingressDNSPrefix, "ingress.prefix", "I", "grafana.example.com", "dna prefix template %s.prefix")

	return command
}

func runCmd() error {
	var rootCmd = &cobra.Command{Use: "grafana-provisioner"}
	rootCmd.AddCommand(cmdVersion())
	rootCmd.AddCommand(cmdProvision())

	err := rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("%v", err.Error())
	}
	return nil
}

func errCheck(err error) {
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("failed will exit")
		os.Exit(1)
	}
}
