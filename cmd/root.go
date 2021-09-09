package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var namespace string
var context string
var kubeconfig string
var labels string

var RootCmd = &cobra.Command{
	Use:   "ckube",
	Short: "Concurrent Kubectl",
	Long:  `A CLI to simplify working with kubectl.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// variables declared here are global for the entire application (assuming subcommands are in the `cmd` directory
	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ckube.yaml)")
	// if options are added to the cli to be passed through to kubectl they should mimiic the naming
	// used by kubectl whenever possible to provide a familiar and consistent experience
	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "the kubernetes namespace (defaults to value currently used by kubectl)")
	RootCmd.PersistentFlags().StringVar(&context, "context", "", "the kubernetes context (defaults to value currently used by kubectl)")
	RootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file to use for CLI requests (defaults to $KUBECONFIG or $HOME/.kube/kubeconfig)")
	RootCmd.PersistentFlags().StringVarP(&labels, "labels", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ckube" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ckube")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}
