package cmd

import (
	"fmt"

	"github.com/canopytax/ckube/util"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"
)

// deploymentCmd represents the deployment command
var deploymentCmd = &cobra.Command{
	Use:     "deployment",
	Aliases: []string{"dep", "deps", "deployments", "deploys"},
	Short:   "Get information about deployments",
	Long: `Shows an interactive list of deployments. Navigate with arrow keys. Use '/'
to narrow the list results by searching. Selecting a deployment will return detailed
information--the output of kubectl get deployment [NAME] -o yaml.

Examples:
  # Show interactive list of deployments. Selecting a deployment will print detailed info
  ckube deployment`,
	Run: func(cmd *cobra.Command, args []string) {
		printDeploymentView()
	},
}

func printDeploymentView() {
	deploymentInfo := deploymentInfo()

	oMan := &util.OutputManager{HeaderColumns: []string{"NAME", "DESIRED", "CURRENT", "UP-TO-DATE", "AVAILABLE", "LAST-UPDATED"}}
	for _, dInfo := range deploymentInfo {
		oMan.Append(dInfo.Print())
	}

	formattedOutput := oMan.FormattedStringSlice()

	templates := &promptui.SelectTemplates{
		Active:   "{{ . | yellow | underline }}",
		Inactive: "{{ . }}",
		Help:     "Select Deployment for more info. Use '/' to search",
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

	output := util.RawK8sOutput(namespace, context, labels, "get", "deployment", deploymentInfo[i].Name, "-oyaml")
	for _, line := range output {
		fmt.Println(line)
	}
}

func deploymentInfo() []DeploymentInfo {
	clientset := util.GetClientset(kubeconfig, context)

	depList, err := clientset.AppsV1beta2().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(fmt.Errorf("error listing deployments: %v", err))
	}

	var deploymentInfos []DeploymentInfo
	for _, deployment := range depList.Items {
		var lastUpdateTime time.Time
		conditions := deployment.Status.Conditions
		if len(conditions) > 0 {
			lastUpdateTime = conditions[0].LastUpdateTime.Time
		}
		depInfo := DeploymentInfo{
			Name:         deployment.Name,
			Desired:      *deployment.Spec.Replicas,
			Current:      deployment.Status.ReadyReplicas,
			Updated:      deployment.Status.UpdatedReplicas,
			Available:    deployment.Status.AvailableReplicas,
			LastDeployed: util.Age{Time: lastUpdateTime},
		}
		deploymentInfos = append(deploymentInfos, depInfo)
	}

	return deploymentInfos
}

type DeploymentInfo struct {
	Name         string
	Desired      int32
	Current      int32
	Updated      int32
	Available    int32
	LastDeployed util.Age
}

func (d DeploymentInfo) Print() string {
	return fmt.Sprintf("%v %v %v %v %v %v", d.Name, d.Desired, d.Current, d.Updated, d.Available, d.LastDeployed.Relative())
}

func init() {
	RootCmd.AddCommand(deploymentCmd)
}
