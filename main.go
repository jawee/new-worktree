package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    cursor   int
    selected string
    input    string
    newBranch   *bool
}

func initialModel() model {
	return model{
	}
}

func (m model) Init() tea.Cmd {
    return nil
}

var newBranch = []string{"Yes", "No"}
var choices = []string{"Feature", "Bug", "Fix"}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
        input := msg.String()

        if input == "ctrl+c" {
            return model{}, tea.Quit
        }

        // type prompt
        if m.selected == "" {
            switch input {
                case "up":
                    if m.cursor > 0 {
                        m.cursor--
                    }

                case "down":
                    if m.cursor < len(choices)-1 {
                        m.cursor++
                    }

                case "enter":
                    m.selected = choices[m.cursor]
                    m.cursor = 0
            }
            return m, nil
        }

        //new branch prompt
        if m.selected != "" && m.newBranch == nil {
            switch input {
                case "up":
                    if m.cursor > 0 {
                        m.cursor--
                    }

                case "down":
                    if m.cursor < len(newBranch)-1 {
                        m.cursor++
                    }

                case "enter":
                    if newBranch[m.cursor] == "Yes" {
                        m.newBranch = &[]bool{true}[0]
                    } else {
                        m.newBranch = &[]bool{false}[0]
                    }
            }
            return m, nil
        }

        // branch name prompt
        if m.selected != "" && m.newBranch != nil {
            switch input {
                case "enter":
                    if m.input != "" {
                       return m, tea.Quit
                    }
                case "backspace":
                    if len(m.input) > 0 {
                        m.input = m.input[:len(m.input)-1]
                    }
                default: 
                    m.input += msg.String()
            }
        }
    }

    return m, nil
}

func (m model) View() string {
    s := ""
    if m.selected == "" {
        s = "What type of branch?\n\n"

        for i, choice := range choices {
            cursor := " "
            if m.cursor == i {
                cursor = ">"
            }

            s += fmt.Sprintf("%s %s\n", cursor, choice)
        }
    }
    if m.selected != "" && m.newBranch == nil {
        s = "New branch??\n\n"
        for i, choice := range newBranch {
            cursor := " "
            if m.cursor == i {
                cursor = ">"
            }

            s += fmt.Sprintf("%s %s\n", cursor, choice)
        }
    }

    if m.selected != "" && m.newBranch != nil {
        s = "Name of branch?\n\n"
       s += fmt.Sprintf("> %s", m.input) 
    }

    return s
}

func main() {
    p := tea.NewProgram(initialModel())
    m, err := p.Run()
    if err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }

    fmt.Printf("%+v\n", m)
    if m, ok := m.(model); ok && m.input != "" {
        folderName := ""
        branchName := ""
        if m.selected == "Feature" {
            folderName += "feat-"
            branchName += "feature/"
        } else if m.selected == "Bug" {
            folderName += "bug-"
            branchName += "bug/"
        } else if m.selected == "Fix" {
            folderName += "fix-"
            branchName += "fix/"
        }
        if folderName == "" {
            fmt.Printf("Do nothing")
            os.Exit(1)
        }

        folderName += strings.ToLower(m.input)
        branchName += strings.ToLower(m.input)

        flag := ""
        if *m.newBranch {
            flag = "-b"
        }

        cmd := exec.Command("git", "worktree", "add", folderName, flag, branchName)

        err := cmd.Run()

        if err != nil {
            fmt.Printf("Could not create worktree. %s\n", err)
            os.Exit(1)
        }

        fmt.Printf("Worktree created\n")
    }
}
