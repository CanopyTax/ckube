package cmd

import (
	"fmt"

	"github.com/devonmoss/ckube/util"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"strings"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls [NAME]",
	Aliases: []string{"get"},
	Short:   "Interactive list of pods",
	Long: `Interactive list of pods matching the specified value.
Use arrow keys to navigate and '/' to search the output. Selecting
a pod will output the result of 'kubectl describe [POD]'

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

		var podInfos []PodInfo
		for _, p := range pods[1:] {
			if p != "" {
				var parts []string
				splits := strings.Split(p, " ")
				for _, s := range splits {
					if s != "" {
						parts = append(parts, s)
					}
				}
				podInfos = append(podInfos, PodInfo{Name: parts[0], Ready: parts[1], Status: parts[2], Restarts: parts[3], Age: parts[4]})
			}
		}

		oMan := &util.OutputManager{HeaderColumns: []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE"}}
		for _, pInfo := range podInfos {
			oMan.Append(pInfo.Print())
		}

		formattedOutput := oMan.FormattedStringSlice()

		templates := &promptui.SelectTemplates{
			Active:   "{{ . | yellow | underline }}",
			Inactive: "{{ . }}",
			Help:     "Select a pod for more info",
		}

		searcher := func(input string, index int) bool {
			text := strings.Replace(strings.ToLower(formattedOutput[1:][index]), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(text, input)
		}

		prompt := promptui.Select{
			Label:     formattedOutput[0],
			Items:     formattedOutput[1:],
			Size:      20,
			Templates: templates,
			Searcher:  searcher,
		}

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

type PodInfo struct {
	Name     string
	Ready    string
	Status   string
	Restarts string
	Age      string
}

func (p PodInfo) Print() string {
	return fmt.Sprintf("%v %v %v %v %v", p.Name, p.Ready, p.Status, p.Restarts, p.Age)
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
