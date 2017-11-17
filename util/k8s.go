package util

import (
	"fmt"
	"strings"
	"sync"
)

type OutputManager struct {
	sync.RWMutex
	output []string
	Header string
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
	fmt.Println(o.Header)
	for _, line := range o.output {
		fmt.Println(line)
	}
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
