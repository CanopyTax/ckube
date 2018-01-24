package util

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func CreatePodInfos(pods []string) []PodInfo {
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
	return podInfos
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

func GetPodPrompt(podInfos []PodInfo, helpMessage string) promptui.Select {
	oMan := &OutputManager{HeaderColumns: []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE"}}
	for _, pInfo := range podInfos {
		oMan.Append(pInfo.Print())
	}

	formattedOutput := oMan.FormattedStringSlice()

	templates := &promptui.SelectTemplates{
		Active:   "{{ . | yellow | underline }}",
		Inactive: "{{ . }}",
		Help:     helpMessage,
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
	return prompt
}
