package util

import (
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetPods(namespace string, context string, labels string) []string {
	var pods []string
	args := K8sCommandArgs([]string{"get", "pods"}, namespace, context, labels)
	cmdPods := RunCommand("kubectl", args...)

	for _, podInfo := range cmdPods[1:] {
		podSplit := strings.Split(podInfo, " ")
		if len(podSplit) > 0 {
			if podSplit[0] != "" {
				pods = append(pods, podSplit[0])
			}
		}
	}
	return pods
}

func RawK8sOutput(namespace string, context string, labels string, args ...string) []string {
	cmdArgs := K8sCommandArgs(args, namespace, context, labels)
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
			if strings.Contains(split[0], search) || search == "" {
				if line != "" {
					output = append(output, line)
				}
			}
		}
	}
	return output
}

func GetServicePods(service string, namespace string, context string, labels string) []string {
	var pods []string
	allPods := GetPods(namespace, context, labels)
	if service == "" {
		return allPods
	}
	for _, pod := range allPods {
		if strings.HasPrefix(pod, service) {
			pods = append(pods, pod)
		}
	}
	return pods
}

func GetMatchingPods(service string, namespace string, context string, labels string) []string {
	var pods []string
	allPods := GetPods(namespace, context, labels)
	if service == "" {
		return allPods
	}
	for _, pod := range GetPods(namespace, context, labels) {
		if strings.Contains(pod, service) {
			pods = append(pods, pod)
		}
	}
	return pods
}

func K8sCommandArgs(args []string, namespace string, context string, labels string) []string {
	if namespace != "" {
		args = append(args, fmt.Sprintf("--namespace=%v", namespace))
	}
	if context != "" {
		args = append(args, fmt.Sprintf("--context=%v", context))
	}
	if labels != "" {
		args = append(args, fmt.Sprintf("--selector=%v", labels))
	}
	return args
}

func GetClientset(kubeconfig, context string) *kubernetes.Clientset {

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.ExplicitPath = kubeconfig
	overrides := &clientcmd.ConfigOverrides{CurrentContext: context}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func KeysString(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k, v := range m {
		keys = append(keys, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(keys, ",")
}
