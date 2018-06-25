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
	var dryRun bool
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

			config := &Config{}
			if team != "" {
				fmt.Println("adding teams from flags")
				config.Teams = append(config.Teams, team)
			} else {
				config = GetConfig("config.yaml")
				fmt.Println("adding teams from config")
			}

			manifests := ""
			for _, teamname := range config.Teams {
				name := fmt.Sprintf("grafana-%s", teamname)
				values := manifestValues{
					databaseURL: dbGrafanaStr(dbDNS, dbPass, dbUser, teamname, dbPort),

					ingressClass: "nginx",
					ingressHost:  fmt.Sprintf("%s.%s", teamname, ingressDNSPrefix),
					ingressName:  name,

					serviceName: name,
					namespace:   teamname,
					image:       image,

					deploymentName:       name,
					deploymentLabelKey:   "app",
					deploymentLabelValue: name,
				}

				manifest := manifestRender(values)
				manifests = manifestsAppend(manifests, manifest)

				if dryRun {
					fmt.Printf("--- %s ---\n", name)
				} else {
					err = db.createDB(teamname)
					errCheck(err)

				}
			}

			if dryRun {
				fmt.Println(manifests)
			} else {
				kubectl := kubectl.NewKubectlClient(kubeconfig)
				err = kubectl.Apply(manifests)
				errCheck(err)

			}

		},
	}

	command.Flags().StringVarP(&kubeconfig, "kube.config", "k", "", "outside cluster path to kube config")
	command.Flags().StringVarP(&dbDNS, "db.dns", "d", "localhost", "mysql database dns")
	command.Flags().IntVarP(&dbPort, "db.port", "P", 3306, "mysql database port")
	command.Flags().StringVarP(&dbPass, "db.pass", "p", "", "mysql database password")
	command.Flags().StringVarP(&dbUser, "db.user", "u", "", "mysql database username")
	command.Flags().StringVarP(&team, "team", "t", "", "team name, if passed config file will be skipped ")
	command.Flags().StringVarP(&image, "image", "i", "grafana/grafana:5.1.0", "grafana official container image")
	command.Flags().StringVarP(&ingressDNSPrefix, "ingress.prefix", "I", "grafana.example.com", "dna prefix template %s.prefix")

	command.Flags().BoolVarP(&dryRun, "dry-run", "D", false, "only output data")

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
