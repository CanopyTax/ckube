package cmd

import (
	"fmt"

	"github.com/canopytax/ckube/util"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls [NAME]",
	Aliases: []string{"get"},
	Short:   "Interactive list of pods",
	Long: `Interactive list of pods matching the specified value.
Use arrow keys to navigate and '/' to search the list. Selecting
a pod will print the result of 'kubectl describe [POD]'

Examples:
  # List all pods in the current namespace
  ckube ls

  # List all pods starting with 'nginx'
  ckube ls nginx`,
	Run: func(cmd *cobra.Command, args []string) {
		pods := util.RawK8sOutput(namespace, context, labels, "get", "pods")
		if len(args) > 0 {
			searchString := args[0]
			pods = util.FilterOutput(pods, searchString, false)
		}

		podInfos := util.CreatePodInfos(pods)
		prompt := util.GetPodPrompt(podInfos, "Select a pod for more info. Type '/' to search")

		i, _, err := prompt.Run()

		if err != nil {
			return
		}
		output := util.RawK8sOutput(namespace, context, labels, "describe", "pod", podInfos[i].Name)

		for _, line := range output {
			fmt.Println(line)
		}
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
