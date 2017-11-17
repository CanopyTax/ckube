package cmd

import (
	"fmt"

	"github.com/devonmoss/ckube/util"
	"github.com/spf13/cobra"
	"sync"
)

var follow bool

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs (POD | PREFIX_STRING)",
	Short: "get logs from a service",
	Long: `View or stream logs from a service or pod. The command will target any
pods that start with the specified value.

Examples:
  # Return logs for all my-cool-service pods
  ck logs my-cool

  # Begin streaming the logs for all pods that begin with 'c'
  ck logs -f c`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		pods := util.GetServicePods(serviceName, namespace, context)
		cm := util.ColorManager{}

		c := make(chan string)
		if len(pods) > 0 {
			if follow {
				for _, pod := range pods {
					cmdArgs := util.K8sCommandArgs([]string{"logs", "-f", pod}, namespace, context)
					go util.StreamCommand(c, cm.GetPrefix(pod), "kubectl", cmdArgs...)
				}
				for {
					select {
					case msg := <-c:
						if msg != "" {
							fmt.Println(msg)
						}
					}
				}
			} else {
				var wg sync.WaitGroup
				for _, pod := range pods {
					wg.Add(1)
					go func(p string) {
						defer wg.Done()
						cmdArgs := util.K8sCommandArgs([]string{"logs", p}, namespace, context)
						prefix := cm.GetPrefix(p)
						logs := util.RunCommand("kubectl", cmdArgs...)
						for _, line := range logs {
							if line != "" {
								fmt.Printf("%v - %v\n", prefix, line)
							}
						}
					}(pod)
				}
				wg.Wait()
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Specify if the logs should be streamed")
}
