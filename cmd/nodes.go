package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/canopytax/ckube/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "Lists pods grouped by node",
	Long:    `Lists pods grouped by node`,
	Run: func(cmd *cobra.Command, args []string) {
		printNodeView()
	},
}

func printNodeView() {
	nodeMap := nodeMap()
	for _, nodePodInfo := range nodeMap {
		//fmt.Println(node)
		printNodeInfo(nodePodInfo.Node)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)
		headerLine := fmt.Sprintf("\t%v\t%v\t%v\t%v\t%v\t", "NAME", "READY", "STATUS", "RESTARTS", "AGE")
		fmt.Fprintln(w, headerLine)
		for _, pod := range nodePodInfo.Pods {
			ps := NewPodStatus(pod)
			age := &util.Age{Time: pod.Status.StartTime.Time}
			statusLine := fmt.Sprintf("\t%v\t%v/%v\t%v\t%v\t%v\t", pod.Name, ps.ready, ps.total, pod.Status.Phase, ps.restarts, age.Relative())
			fmt.Fprintln(w, statusLine)
		}
		w.Flush()
		fmt.Println()
	}
}

func printNodeInfo(node v1.Node) {
	nodeName := node.Name
	var nodeLabels []string
	labelColor := color.New(color.BgCyan).SprintfFunc()
	masterColor := color.New(color.BgBlue).SprintfFunc()
	for k, v := range node.Labels {
		if !strings.Contains(k, "kubernetes") {
			nodeLabels = append(nodeLabels, labelColor(k+"="+v))
		}
		if k == "kubernetes.io/role" && v == "master" {
			nodeLabels = append(nodeLabels, masterColor(v))
		}
	}
	labelString := strings.Join(nodeLabels, " ")
	fmt.Printf("%v\t%v\n", nodeName, labelString)
}

type PodStatus struct {
	total    int
	ready    int
	restarts int32
}

func NewPodStatus(pod v1.Pod) PodStatus {
	total := len(pod.Status.ContainerStatuses)
	var ready int
	var restarts int32
	for _, c := range pod.Status.ContainerStatuses {
		if c.Ready {
			ready++
		}
		restarts += c.RestartCount
	}
	return PodStatus{total: total, ready: ready, restarts: restarts}
}

func nodeMap() map[string]NodePodInfo {
	clientset := util.GetClientset(kubeconfig, context)

	podList, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(fmt.Errorf("error listing pods: %v", err))
	}

	nodeList, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(fmt.Errorf("error listing pods: %v", err))
	}

	nodeMap := make(map[string][]v1.Pod)
	for _, pod := range podList.Items {
		if _, ok := nodeMap[pod.Spec.NodeName]; ok {
			nodeMap[pod.Spec.NodeName] = append(nodeMap[pod.Spec.NodeName], pod)
		} else {
			nodeMap[pod.Spec.NodeName] = []v1.Pod{pod}
		}
	}

	nodePodMap := make(map[string]NodePodInfo)
	for _, node := range nodeList.Items {
		nodePodMap[node.Name] = NodePodInfo{Node: node, Pods: nodeMap[node.Name]}
	}
	return nodePodMap
}

func init() {
	RootCmd.AddCommand(nodesCmd)
}

type NodePodInfo struct {
	Node v1.Node
	Pods []v1.Pod
}
