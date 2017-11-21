package util

import (
	"fmt"
	"strings"
	"sync"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"text/tabwriter"
	"os"
)

type OutputManager struct {
	sync.RWMutex
	output []string
	HeaderColumns []string
}

func (o *OutputManager) Append(s ...string) {
	o.Lock()
	defer o.Unlock()

	o.output = append(o.output, s...)
}

func (o *OutputManager) GetOutput() []string {
	return o.output
}

func (o *OutputManager) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)
	var headerLine string
	for _, s := range o.HeaderColumns {
		headerLine += s + "\t"
	}
	fmt.Fprintln(w, headerLine)
	for _, line := range o.output {
		fmt.Fprintln(w, o.tabbedString(line))
	}
	w.Flush()
}

func (o *OutputManager) tabbedString(output string) string {
	splits := strings.Split(output, " ")
	var tabbedString string
	for _, s := range splits {
		if s != "" {
			tabbedString += s + "\t"
		}
	}
	return tabbedString
}


func GetPods(namespace string, context string) []string {
	var pods []string
	args := K8sCommandArgs([]string{"get", "pods"}, namespace, context)
	cmdPods := RunCommand("kubectl", args...)

	for _, podInfo := range cmdPods[1:] {
		podSplit := strings.Split(podInfo, " ")
		if len(podSplit) > 0 {
			pods = append(pods, podSplit[0])
		}
	}
	return pods
}

func RawK8sOutput(namespace string, context string, args ...string) []string {
	cmdArgs := K8sCommandArgs(args, namespace, context)
	output := RunCommand("kubectl", cmdArgs...)
	return output
}

func FilterOutput(lines []string, search string, stripHeader bool) []string {
	var output []string
	if !stripHeader {
		output = append(output, lines[0])
	}

	for _, line := range lines[1:] {
		split := strings.Split(line, " ")
		if len(split) > 0 {
			if strings.Contains(split[0], search) || search == ""{
				if line != "" {
					output = append(output, line)
				}
			}
		}
	}
	return output
}

func GetServicePods(service string, namespace string, context string) []string {
	var pods []string
	for _, pod := range GetPods(namespace, context) {
		if strings.HasPrefix(pod, service) {
			pods = append(pods, pod)
		}
	}
	return pods
}

func GetMatchingPods(service string, namespace string, context string) []string {
	var pods []string
	for _, pod := range GetPods(namespace, context) {
		if strings.Contains(pod, service) {
			pods = append(pods, pod)
		}
	}
	return pods
}

func K8sCommandArgs(args []string, namespace string, context string) []string {
	if namespace != "" {
		args = append(args, fmt.Sprintf("--namespace=%v", namespace))
	}
	if context != "" {
		args = append(args, fmt.Sprintf("--context=%v", context))
	}
	return args
}

func GetClientset(kubeconfig string) *kubernetes.Clientset {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}