package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    choices  []string           
    cursor   int               
    selected string
    input string
}

func initialModel() model {
	return model{
		choices:  []string{"Feature", "Bug", "Fix"},
	}
}

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
        switch msg.String() {

        case "ctrl+c":
            return m, tea.Quit

        case "up":
            if m.cursor > 0 {
                m.cursor--
            }

        case "down":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }

        case "enter":

            if m.selected == "" {
                m.selected = m.choices[m.cursor]
            }

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

    return m, nil
}

func (m model) View() string {
    s := ""
    if m.selected == "" {
        s = "What type of branch?\n\n"

        for i, choice := range m.choices {
            cursor := " "
            if m.cursor == i {
                cursor = ">"
            }

            s += fmt.Sprintf("%s %s\n", cursor, choice)
        }
    }

    if m.selected != "" {
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
        cmd := exec.Command("git", "worktree", "add", folderName, "-b", branchName)

        err := cmd.Run()

        if err != nil {
            fmt.Printf("Could not create worktree. %s\n", err)
            os.Exit(1)
        }

        fmt.Printf("Worktree created\n")
    }
}
