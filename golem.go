package main

import (
	"fmt"
	"os/exec"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/alecthomas/kingpin.v2"
)

// the questions to ask
var qs = []*survey.Question{
	{
		Name:      "type",
		Validate:  survey.Required,
		Transform: GetResultString,
		Prompt: &survey.Select{
			Message: "Type of commit:",
			Options: []string{
				"feat     | new feature",
				"fix      | bug fix",
				"docs     | documentation",
				"style    | code formatting",
				"refactor | code refactor",
				"perf     | improve performance",
				"test     | add/update tests ",
				"chore    | other changes that doesn't modify src/test",
				"revert   | revert previous commit"},
			Default: "feat",
			PageSize: 10,
		},
	},
	{
		Name:     "scope",
		Prompt:   &survey.Input{Message: "Scope of this change:"},
		Validate: survey.Required,
	},
	{
		Name:     "message",
		Prompt:   &survey.Input{Message: "Commit message:"},
		Validate: survey.Required,
	},
	{
		Name:     "issue",
		Prompt:   &survey.Input{Message: "Issue(s) ID:"},
		Validate: survey.Required,
	},
}

func main() {

	// noIssue := kingpin.Flag("no-issue", "commit without an issue").Short('n').Bool()
	noIssue := true
	// the answers will be written to this struct
	answers := struct {
		Type    string `survey:"type"`
		Scope   string
		Message string
		Issue   string
	}{}

	kingpin.Parse()
	if noIssue {
		qs = RemoveIndex(qs, 3)
	}
	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	out := fmt.Sprintf("%s(%s): %s | %s", answers.Type, answers.Scope, answers.Issue, answers.Message)
	if noIssue {
		out = fmt.Sprintf("%s(%s): %s", answers.Type, answers.Scope, answers.Message)
	}
	fmt.Printf("message: %s", out)
	fmt.Println()
	output, _ := exec.Command("git", "commit", "-m", out).CombinedOutput()
	fmt.Println(string(output))
}

func GetResultString(ans interface{}) interface{} {
	transformer := survey.TransformString(SplitTrimString)
	return transformer(ans)
}

func SplitTrimString(ans string) string {
	s := strings.Split(ans, "|")
	ans = strings.TrimSpace(s[0])
	return ans
}

func RemoveIndex(s []*survey.Question, index int) []*survey.Question {
	return append(s[:index], s[index+1:]...)
}
