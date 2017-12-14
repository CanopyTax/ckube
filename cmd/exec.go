package cmd

import (
	"fmt"

	"github.com/canopytax/ckube/util"
	"github.com/spf13/cobra"
	"sync"
)

var tty, stdin, all bool
var postDashArgs []string

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec (POD | PREFIX_STRING) [options] -- COMMAND [args...]",
	Short: "execute a command in a container",
	Long: `Execute a command in a container.

Examples:
  # get output from running 'date' in a nginx pod, using first container of first pod by default
  ckube exec nginx date

  # get output from running 'date' in all nginx pods, using first container of each pod by default
  ckube exec -a nginx date

  # switch to raw terminal mode, sends stdin to 'bash' in a nginx pod and sends stdout/stderr from
  # 'bash' back to the client
  ckube exec -it nginx bash

  # get output from an extended arg 'curl' command in all nginx pods
  ckube exec -a nginx -- curl https://google.com -v`,
	Run: func(cmd *cobra.Command, args []string) {
		if labels != "" {
			fmt.Println("using labels for exec is not yet supported")
			return
		}
		serviceName := args[0]
		dashLength := cmd.ArgsLenAtDash()
		if dashLength > 0 {
			postDashArgs = args[dashLength:]
		}
		pods := util.GetServicePods(serviceName, namespace, context, labels)
		cm := util.ColorManager{}
		if len(pods) == 0 {
			fmt.Println("No matching containers")
			return
		}
		if !all {
			// make pods only have 1 item
			pods = append([]string{}, pods[0])
		}
		if tty && stdin {
			pod := pods[0]
			cmdArgs := util.K8sCommandArgs([]string{"exec", "-it", pod}, namespace, context, "")
			if len(postDashArgs) > 0 {
				cmdArgs = append(cmdArgs, "--")
				cmdArgs = append(cmdArgs, postDashArgs...)
			} else {
				cmdArgs = append(cmdArgs, args[1:]...)
			}
			fmt.Println(cm.Colorize(pod))
			util.InteractiveCommand("kubectl", cmdArgs...)
		} else {
			var wg sync.WaitGroup
			for _, pod := range pods {
				wg.Add(1)
				go func(p string) {
					defer wg.Done()
					cmdArgs := util.K8sCommandArgs([]string{"exec", p}, namespace, context, "")
					if len(postDashArgs) > 0 {
						cmdArgs = append(cmdArgs, "--")
						cmdArgs = append(cmdArgs, postDashArgs...)
					} else {
						cmdArgs = append(cmdArgs, args[1:]...)
					}
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
	},
}

func init() {
	RootCmd.AddCommand(execCmd)
	execCmd.Flags().BoolVarP(&tty, "tty", "t", false, "Stdin is a TTY")
	execCmd.Flags().BoolVarP(&stdin, "stdin", "i", false, "Pass stdin to the container")
	execCmd.Flags().BoolVarP(&all, "all", "a", false, "run command in all matching containers (ignores -i, -t)")
}
