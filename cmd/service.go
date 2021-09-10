package cmd

import (
	"fmt"

	"github.com/canopytax/ckube/util"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services", "svc"},
	Short:   "Interactive view of your services",
	Long:    `Shows an interactive view of your services`,
	Run: func(cmd *cobra.Command, args []string) {
		showServiceView()
	},
}

func showServiceView() {
	serviceInfos := getServiceInfo()

	templates := &promptui.SelectTemplates{
		Active:   "{{ .Service.Name | underline | yellow }}",
		Inactive: "{{ .Service.Name }}",
		Details: `
--------- Pods ----------
{{ range $i, $pod := .PodDetails }}
{{ .Name }}	{{ .Ready }}/{{ .Total }}	{{ .Restarts }}	{{ .Status }}	{{ .Age }}{{end}}`,
	}

	prompt := promptui.Select{
		Label:     "Services",
		Items:     serviceInfos,
		Templates: templates,
		Size:      20,
	}

	_, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}

func getServiceInfo() []ServiceInfo {
	serviceList := util.GetServiceList(namespace, context, labels)

	var serviceInfos []ServiceInfo
	for _, service := range serviceList.Items {
		selector := service.Spec.Selector
		if len(selector) > 0 {
			podList := util.GetPodList(namespace, context, util.KeysString(selector))

			var podDetails []PodDetails
			for _, pod := range podList.Items {
				podDetails = append(podDetails, NewPodDetails(pod))
			}
			serviceInfo := ServiceInfo{Service: service, PodDetails: podDetails}
			serviceInfos = append(serviceInfos, serviceInfo)
		}
	}
	return serviceInfos
}

type PodDetails struct {
	Name     string
	Total    int
	Ready    int
	Restarts int32
	Status   string
	Age      string
}

func NewPodDetails(pod v1.Pod) PodDetails {
	total := len(pod.Status.ContainerStatuses)
	var ready int
	var restarts int32
	for _, c := range pod.Status.ContainerStatuses {
		if c.Ready {
			ready++
		}
		restarts += c.RestartCount
	}
	age := &util.Age{Time: pod.Status.StartTime.Time}
	return PodDetails{Name: pod.Name, Total: total, Ready: ready, Restarts: restarts, Age: age.Relative(), Status: string(pod.Status.Phase)}
}

type ServiceInfo struct {
	Service    v1.Service
	PodDetails []PodDetails
}

func init() {
	RootCmd.AddCommand(serviceCmd)
}
