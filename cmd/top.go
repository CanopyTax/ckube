package cmd

import (
	"github.com/spf13/cobra"
	"github.com/devonmoss/ckube/util"
	"sync"
)

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "View cpu and memory usage for pods",
	Long: `View cpu and memory usage for pods
For example:
  # view cpu and memory for all pods
	ckube top

  # view cpu and memory for all pods with name containing 'nginx'
	ckube top nginx`,
	Run: func(cmd *cobra.Command, args []string) {
		var pods []string
		if len(args) > 0 {
			pods = util.GetMatchingPods(args[0], namespace, context, labels)
		} else {
			pods = util.GetPods(namespace, context, labels)
		}
		oMan := &util.OutputManager{HeaderColumns:[]string{"NAME", "CPU(cores)", "MEMORY(bytes)"}}
		if len(pods) > 0 {
			var wg sync.WaitGroup
			for _, pod := range pods {
				wg.Add(1)
				go func(p string) {
					defer wg.Done()
					lines := util.RawK8sOutput(namespace, context, "top", "pod", p)
					oMan.Append(lines[1])
				}(pod)
			}
			wg.Wait()
			oMan.Print()
		}
	},
}

func init() {
	RootCmd.AddCommand(topCmd)
}
