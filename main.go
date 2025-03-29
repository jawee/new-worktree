package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    cursor    int
    selected  string
    input     string
    newBranch bool
    page      PageType
}

type PageType int
const (
    BRANCH_TYPE PageType = iota
    NEW_BRANCH
    NAME
)

func initialModel() model {
    return model{
        page: BRANCH_TYPE,
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

var newBranchChocies = []string{"Yes", "No"}
var branchTypeChoices = []string{"Feature", "Bug", "Fix", "Chore"}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
        input := msg.String()

        if input == "ctrl+c" {
            return model{}, tea.Quit
        }

        // type prompt
        if m.page == BRANCH_TYPE {
            switch input {
            case "up":
                if m.cursor > 0 {
                    m.cursor--
                }

            case "down":
                if m.cursor < len(branchTypeChoices)-1 {
                    m.cursor++
                }

            case "enter":
                m.selected = branchTypeChoices[m.cursor]
                m.cursor = 0
                m.page = NEW_BRANCH
            }
            return m, nil
        }

        //new branch prompt
        if m.page == NEW_BRANCH {
            switch input {
            case "up":
                if m.cursor > 0 {
                    m.cursor--
                }

            case "down":
                if m.cursor < len(newBranchChocies)-1 {
                    m.cursor++
                }

            case "enter":
                if newBranchChocies[m.cursor] == "Yes" {
                    m.newBranch = true
                } else {
                    m.newBranch = false
                }
                m.page = NAME
            }
            return m, nil
        }

        // branch name prompt
        if m.page == NAME {
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
    if m.page == BRANCH_TYPE {
        s = "What type of branch?\n\n"

        for i, choice := range branchTypeChoices {
            cursor := " "
            if m.cursor == i {
                cursor = ">"
            }

            s += fmt.Sprintf("%s %s\n", cursor, choice)
        }
    }
    if m.page == NEW_BRANCH {
        s = "New branch?\n\n"
        for i, choice := range newBranchChocies {
            cursor := " "
            if m.cursor == i {
                cursor = ">"
            }

            s += fmt.Sprintf("%s %s\n", cursor, choice)
        }
    }

    if m.page == NAME {
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
        } else if (m.selected == "Chore") {
			folderName += "chore-"
			branchName += "chore/"
		}

        if folderName == "" {
            fmt.Printf("Do nothing")
            os.Exit(1)
        }

        folderName += strings.ToLower(m.input)
        branchName += strings.ToLower(m.input)

        flag := ""
        if m.newBranch {
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
