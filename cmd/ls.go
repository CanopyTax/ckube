package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/devonmoss/ckube/util"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls [NAME]",
	Aliases: []string{"get"},
	Short: "list pods in kubernetes",
	Long: `List pods in kubernetes matching the specified Value

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

		for _, pod := range pods {
			fmt.Println(pod)
		}
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
