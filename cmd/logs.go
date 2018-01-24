package cmd

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/canopytax/ckube/util"
	"github.com/spf13/cobra"
)

var follow bool
var tail int
var tailString string

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:     "logs (POD | PREFIX_STRING)",
	Aliases: []string{"log", "lg"},
	Short:   "get logs from a service",
	Long: `View or stream logs from a service or pod. The command will target any
pods that start with the specified value. If no pod or service is specified this
command will display an interactive list of pods. Navigate with arrow keys. Use '/'
to narrow the results by searching. Selecting a pod will show its logs.

Examples:
  # Return logs for all my-cool-service pods
  ckube logs my-cool

  # Return only the most recent 20 lines of output for all my-cool-service pods
  ckube logs my-cool --tail=20

  # Begin streaming the logs for all pods that begin with 'c'
  ckube logs -f c

  # Show interactive list of pods. Selecting a pod will print its logs.
  ckube logs

  # Show interactive list of pods. Selecting a pod will follow its logs.
  ckube logs -f
`,
	Run: func(cmd *cobra.Command, args []string) {
		tailString = strconv.Itoa(tail)

		var serviceName string
		if len(args) > 0 {
			serviceName = args[0]
		} else {
			pods := util.RawK8sOutput(namespace, context, labels, "get", "pods")
			if len(args) > 0 {
				searchString := args[0]
				pods = util.FilterOutput(pods, searchString, false)
			}

			podInfos := util.CreatePodInfos(pods)
			prompt := util.GetPodPrompt(podInfos, "Select a pod to get the logs. Type '/' to search")

			i, _, err := prompt.Run()

			if err != nil {
				return
			}
			serviceName = podInfos[i].Name
		}

		pods := util.GetServicePods(serviceName, namespace, context, labels)
		streamLogs(pods)
	},
}

func streamLogs(pods []string) {
	cm := util.ColorManager{}
	c := make(chan string)
	if len(pods) > 0 {
		if follow {
			for _, pod := range pods {
				cmdArgs := util.K8sCommandArgs([]string{"logs", "--tail", tailString, "-f", pod}, namespace, context, "")
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
					cmdArgs := util.K8sCommandArgs([]string{"logs", "--tail", tailString, p}, namespace, context, "")
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
}

func init() {
	RootCmd.AddCommand(logsCmd)
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Specify if the logs should be streamed")
	logsCmd.Flags().IntVar(&tail, "tail", -1, "lines of recent log file to display. Defaults to -1, showing all log lines")
}
