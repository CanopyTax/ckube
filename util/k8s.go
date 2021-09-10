package util

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
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

func GetNodes(namespace string, context string, labels string) []string {
	var nodes []string
	args := K8sCommandArgs([]string{"get", "nodes"}, namespace, context, labels)
	cmdNodes := RunCommand("kubectl", args...)

	for _, nodeInfo := range cmdNodes[1:] {
		nodeSplit := strings.Split(nodeInfo, " ")
		if len(nodeSplit) > 0 {
			if nodeSplit[0] != "" {
				nodes = append(nodes, nodeSplit[0])
			}
		}
	}
	return nodes
}

func RawK8sOutput(namespace string, context string, labels string, args ...string) []string {
	cmdArgs := K8sCommandArgs(args, namespace, context, labels)
	output := RunCommand("kubectl", cmdArgs...)
	return output
}

func RawK8sOutputString(namespace string, context string, labels string, args ...string) string {
	output := RawK8sOutput(namespace, context, labels, args...)
	return strings.Join(output, "\n")
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

func GetDeploymentList(namespace string, context string, labels string) appsv1.DeploymentList {
	deploymentString := RawK8sOutputString(namespace, context, labels, []string{"get", "deployments", "-o", "json"}...)

	var deployList appsv1.DeploymentList
	err := json.Unmarshal([]byte(deploymentString), &deployList)
	if err != nil {
		log.Fatal(err)
	}

	return deployList
}

func GetPodListAllNamespaces(context string, labels string) corev1.PodList {
	podString := RawK8sOutputString("", context, labels, []string{"get", "pods", "--all-namespaces", "-o", "json"}...)

	var podList corev1.PodList
	err := json.Unmarshal([]byte(podString), &podList)
	if err != nil {
		log.Fatal(err)
	}

	return podList
}

func GetPodList(namespace string, context string, labels string) corev1.PodList {
	podString := RawK8sOutputString(namespace, context, labels, []string{"get", "pods", "-o", "json"}...)

	var podList corev1.PodList
	err := json.Unmarshal([]byte(podString), &podList)
	if err != nil {
		log.Fatal(err)
	}

	return podList
}

func GetServiceList(namespace string, context string, labels string) corev1.ServiceList {
	serviceString := RawK8sOutputString(namespace, context, labels, []string{"get", "services", "-o", "json"}...)

	var serviceList corev1.ServiceList
	err := json.Unmarshal([]byte(serviceString), &serviceList)
	if err != nil {
		log.Fatal(err)
	}

	return serviceList
}

func GetNodeList(context string, labels string) corev1.NodeList {
	nodesString := RawK8sOutputString("", context, labels, []string{"get", "nodes", "-o", "json"}...)

	var nodeList corev1.NodeList
	err := json.Unmarshal([]byte(nodesString), &nodeList)
	if err != nil {
		log.Fatal(err)
	}

	return nodeList
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

func KeysString(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k, v := range m {
		keys = append(keys, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(keys, ",")
}
