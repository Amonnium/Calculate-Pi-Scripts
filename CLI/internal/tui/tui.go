package tui

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/app"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/catalog"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type service interface {
	Implementations() []catalog.Implementation
	Toolchains(context.Context) ([]app.ToolchainStatus, error)
	Prepare(context.Context, string) (app.RuntimeRequirement, error)
	Run(context.Context, string, bool, io.Writer, io.Writer) app.Result
	InstallToolchain(context.Context, string) app.Result
	UninstallToolchain(context.Context, string) app.Result
}

type screen uint8

const (
	mainScreen screen = iota
	implementationScreen
	toolchainScreen
	toolchainActionScreen
	confirmationScreen
	helpScreen
	outputScreen
)

type listItem struct {
	id          string
	title       string
	description string
}

func (item listItem) Title() string       { return item.title }
func (item listItem) Description() string { return item.description }
func (item listItem) FilterValue() string { return item.title }

type action uint8

const (
	noAction action = iota
	installForRun
	installToolchain
	uninstallToolchain
)

type operationResultMsg struct {
	result app.Result
}

type prepareResultMsg struct {
	requirement app.RuntimeRequirement
	err         error
}

type toolchainsResultMsg struct {
	statuses []app.ToolchainStatus
	err      error
}

type serviceExec struct {
	operation func(io.Writer, io.Writer) app.Result
	result    app.Result
	stdout    io.Writer
	stderr    io.Writer
}

func (command *serviceExec) Run() error {
	command.result = command.operation(command.stdout, command.stderr)
	return nil
}

func (command *serviceExec) SetStdin(io.Reader)         {}
func (command *serviceExec) SetStdout(writer io.Writer) { command.stdout = writer }
func (command *serviceExec) SetStderr(writer io.Writer) { command.stderr = writer }

type model struct {
	service           service
	screen            screen
	previousScreen    screen
	items             list.Model
	viewport          viewport.Model
	helpPage          int
	selectedImpl      string
	selectedToolchain string
	toolchainStatuses []app.ToolchainStatus
	pendingAction     action
	message           string
	busy              bool
	width             int
	height            int
}

func newModel(application service) model {
	current := model{service: application, screen: mainScreen}
	current.setList("Calculate Pi", mainMenuItems())
	return current
}

func (current model) Init() tea.Cmd { return nil }

func (current model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch value := message.(type) {
	case tea.WindowSizeMsg:
		current.width, current.height = value.Width, value.Height
		current.items.SetSize(value.Width-4, value.Height-8)
		current.viewport.Width = value.Width - 4
		current.viewport.Height = value.Height - 7
		return current, nil
	case prepareResultMsg:
		current.busy = false
		if value.err != nil {
			current.showMessage("Requirement check failed", value.err.Error())
			return current, nil
		}
		if value.requirement.Toolchain.ID != "" && !value.requirement.Status.Compatible {
			current.pendingAction = installForRun
			current.message = fmt.Sprintf("%s %s is needed to run this implementation.", value.requirement.Toolchain.Name, value.requirement.Toolchain.MinimumVersion)
			current.setScreen(confirmationScreen, "Install Required Runtime?", []listItem{
				{id: "yes", title: "Install and run", description: "Use the detected system package manager"},
				{id: "no", title: "Cancel", description: "Return without running the implementation"},
			})
			return current, nil
		}
		current.previousScreen = implementationScreen
		return current, current.runSelected(false)
	case toolchainsResultMsg:
		current.busy = false
		if value.err != nil {
			current.showMessage("Unable to inspect toolchains", value.err.Error())
			return current, nil
		}
		current.toolchainStatuses = value.statuses
		current.setScreen(toolchainScreen, "Manage Toolchains", toolchainItems(value.statuses))
		return current, nil
	case operationResultMsg:
		current.busy = false
		if value.result.NeedsRuntimeInstall {
			current.selectedToolchain = value.result.RuntimeID
			current.pendingAction = installForRun
			current.message = value.result.RuntimePrompt
			current.previousScreen = implementationScreen
			current.setScreen(confirmationScreen, "Install Required Runtime?", []listItem{
				{id: "yes", title: "Install and run", description: value.result.RuntimePrompt},
				{id: "no", title: "Cancel", description: "Return without running the implementation"},
			})
			return current, nil
		}
		current.showResult(value.result)
		return current, nil
	case tea.KeyMsg:
		return current.handleKey(value)
	}

	var command tea.Cmd
	if current.screen == outputScreen {
		current.viewport, command = current.viewport.Update(message)
		return current, command
	}
	current.items, command = current.items.Update(message)
	return current, command
}

func (current model) handleKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	pressed := strings.ToLower(key.String())
	if pressed == "q" || key.Type == tea.KeyCtrlC {
		return current, tea.Quit
	}
	if current.busy {
		return current, nil
	}
	if pressed == "b" || key.Type == tea.KeyEsc {
		if current.screen == mainScreen {
			return current, nil
		}
		if current.screen == outputScreen && current.previousScreen == toolchainActionScreen {
			current.screen = toolchainScreen
			current.busy = true
			return current, current.loadToolchains()
		}
		current.goBack()
		return current, nil
	}
	if pressed == "n" {
		current.pageDown()
		return current, nil
	}
	if pressed == "p" {
		current.pageUp()
		return current, nil
	}
	if current.screen == outputScreen {
		var command tea.Cmd
		current.viewport, command = current.viewport.Update(key)
		return current, command
	}
	if current.screen == helpScreen {
		return current, nil
	}
	if key.Type != tea.KeyEnter {
		var command tea.Cmd
		current.items, command = current.items.Update(key)
		return current, command
	}

	selected, ok := current.items.SelectedItem().(listItem)
	if !ok {
		return current, nil
	}
	switch current.screen {
	case mainScreen:
		switch selected.id {
		case "calculate":
			current.setScreen(implementationScreen, "Select Implementation", implementationItems(current.service.Implementations()))
		case "toolchains":
			current.busy = true
			return current, current.loadToolchains()
		case "help":
			current.helpPage = 0
			current.screen = helpScreen
		case "exit":
			return current, tea.Quit
		}
	case implementationScreen:
		current.selectedImpl = selected.id
		current.busy = true
		return current, func() tea.Msg {
			requirement, err := current.service.Prepare(context.Background(), selected.id)
			return prepareResultMsg{requirement: requirement, err: err}
		}
	case toolchainScreen:
		current.selectedToolchain = selected.id
		return current.openToolchainActions()
	case toolchainActionScreen:
		return current.selectToolchainAction(selected.id)
	case confirmationScreen:
		if selected.id == "yes" {
			return current.confirmPendingAction()
		}
		current.goBack()
		return current, nil
	}
	return current, nil
}

func (current model) loadToolchains() tea.Cmd {
	return func() tea.Msg {
		statuses, err := current.service.Toolchains(context.Background())
		return toolchainsResultMsg{statuses: statuses, err: err}
	}
}

func (current model) openToolchainActions() (tea.Model, tea.Cmd) {
	var selected app.ToolchainStatus
	for _, entry := range current.toolchainStatuses {
		if entry.Toolchain.ID == current.selectedToolchain {
			selected = entry
			break
		}
	}
	items := []listItem{}
	if !selected.Status.Compatible {
		items = append(items, listItem{id: "install", title: "Install", description: "Install using the system package manager"})
	}
	if selected.Status.Managed {
		items = append(items, listItem{id: "uninstall", title: "Uninstall", description: "Remove the package installed by this CLI"})
	}
	if selected.Status.Compatible && !selected.Status.Managed {
		items = append(items, listItem{id: "available", title: "Already available", description: "This system runtime is not managed by this CLI"})
	}
	items = append(items, listItem{id: "back", title: "Back", description: "Return to the toolchain list"})
	current.setScreen(toolchainActionScreen, selected.Toolchain.Name, items)
	return current, nil
}

func (current model) selectToolchainAction(id string) (tea.Model, tea.Cmd) {
	switch id {
	case "back":
		current.screen = toolchainScreen
		current.setList("Manage Toolchains", toolchainItems(current.toolchainStatuses))
		return current, nil
	case "available":
		current.showMessage("Runtime Already Available", "This runtime is system-managed. No changes were made.")
		return current, nil
	case "install":
		current.pendingAction = installToolchain
		current.message = "Install " + current.selectedToolchain + " using the system package manager?"
	case "uninstall":
		current.pendingAction = uninstallToolchain
		current.message = "Uninstall the package recorded as installed by this CLI?"
	default:
		return current, nil
	}
	current.previousScreen = toolchainActionScreen
	current.setScreen(confirmationScreen, "Confirm Toolchain Change", []listItem{
		{id: "yes", title: "Continue", description: current.message},
		{id: "no", title: "Cancel", description: "Do not change the system"},
	})
	return current, nil
}

func (current model) confirmPendingAction() (tea.Model, tea.Cmd) {
	current.busy = true
	switch current.pendingAction {
	case installForRun:
		current.previousScreen = implementationScreen
		return current, current.runSelected(true)
	case installToolchain:
		current.previousScreen = toolchainActionScreen
		return current, current.runServiceOperation(func(_, _ io.Writer) app.Result {
			return current.service.InstallToolchain(context.Background(), current.selectedToolchain)
		})
	case uninstallToolchain:
		current.previousScreen = toolchainActionScreen
		return current, current.runServiceOperation(func(_, _ io.Writer) app.Result {
			return current.service.UninstallToolchain(context.Background(), current.selectedToolchain)
		})
	default:
		current.busy = false
		current.goBack()
		return current, nil
	}
}

func (current model) runSelected(install bool) tea.Cmd {
	return current.runServiceOperation(func(stdout, stderr io.Writer) app.Result {
		return current.service.Run(context.Background(), current.selectedImpl, install, stdout, stderr)
	})
}

func (current model) runServiceOperation(operation func(io.Writer, io.Writer) app.Result) tea.Cmd {
	command := &serviceExec{operation: operation}
	return tea.Exec(command, func(err error) tea.Msg {
		if err != nil {
			return operationResultMsg{result: app.Result{Message: "Operation failed", Err: err}}
		}
		return operationResultMsg{result: command.result}
	})
}

func (current *model) showResult(result app.Result) {
	var output strings.Builder
	if result.Message != "" {
		fmt.Fprintln(&output, result.Message)
	}
	if result.Err != nil {
		fmt.Fprintf(&output, "\nError: %v\n", result.Err)
	}
	if result.Stdout != "" {
		fmt.Fprintf(&output, "\nOutput:\n%s", result.Stdout)
	}
	if result.Stderr != "" {
		fmt.Fprintf(&output, "\nErrors:\n%s", result.Stderr)
	}
	if result.ExitCode != 0 || result.Err == nil && result.Stdout != "" {
		fmt.Fprintf(&output, "\n\nExit code: %d", result.ExitCode)
	}
	current.screen = outputScreen
	current.viewport.SetContent(output.String())
	current.viewport.GotoTop()
}

func (current *model) showMessage(title, body string) {
	current.previousScreen = current.screen
	current.showResult(app.Result{Message: title, Stderr: body})
}

func (current *model) setScreen(next screen, title string, items []listItem) {
	current.previousScreen = current.screen
	current.screen = next
	current.setList(title, items)
}

func (current *model) setList(title string, items []listItem) {
	listItems := make([]list.Item, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, item)
	}
	current.items = list.New(listItems, list.NewDefaultDelegate(), max(current.width-4, 0), max(current.height-8, 0))
	current.items.Title = title
	current.items.SetShowStatusBar(false)
	current.items.SetShowHelp(false)
	current.items.SetShowPagination(true)
	current.items.SetFilteringEnabled(false)
}

func (current *model) goBack() {
	switch current.screen {
	case mainScreen:
		return
	case implementationScreen, toolchainScreen, helpScreen:
		current.screen = mainScreen
		current.setList("Calculate Pi", mainMenuItems())
	case toolchainActionScreen:
		current.screen = toolchainScreen
		current.setList("Manage Toolchains", toolchainItems(current.toolchainStatuses))
	case confirmationScreen:
		if current.pendingAction == installForRun {
			current.screen = implementationScreen
			current.setList("Select Implementation", implementationItems(current.service.Implementations()))
		} else {
			current.screen = current.previousScreen
		}
		current.pendingAction = noAction
	case outputScreen:
		current.screen = current.previousScreen
		if current.screen == implementationScreen {
			current.setList("Select Implementation", implementationItems(current.service.Implementations()))
		} else if current.screen == toolchainActionScreen {
			current.screen = toolchainScreen
			current.setList("Manage Toolchains", toolchainItems(current.toolchainStatuses))
		}
	}
}

func (current *model) pageDown() {
	switch current.screen {
	case helpScreen:
		current.helpPage = 1
	case outputScreen:
		current.viewport.PageDown()
	default:
		current.items.NextPage()
	}
}

func (current *model) pageUp() {
	switch current.screen {
	case helpScreen:
		current.helpPage = 0
	case outputScreen:
		current.viewport.PageUp()
	default:
		current.items.PrevPage()
	}
}

func mainMenuItems() []listItem {
	return []listItem{
		{id: "calculate", title: "Calculate Pi", description: "Choose an implementation to run"},
		{id: "toolchains", title: "Manage Toolchains", description: "View or manage required runtimes"},
		{id: "help", title: "Help", description: "View keyboard controls and notes"},
		{id: "exit", title: "Exit", description: "Close Calculate Pi Scripts"},
	}
}

func implementationItems(implementations []catalog.Implementation) []listItem {
	items := make([]listItem, 0, len(implementations))
	for _, implementation := range implementations {
		items = append(items, listItem{
			id:          implementation.ID,
			title:       implementation.Name,
			description: implementation.Description + " | " + implementation.Algorithm,
		})
	}
	return items
}

func toolchainItems(statuses []app.ToolchainStatus) []listItem {
	items := make([]listItem, 0, len(statuses))
	for _, entry := range statuses {
		description := "Not installed"
		if entry.Status.Installed && entry.Status.Compatible {
			description = "Available, " + entry.Status.Version
			if entry.Status.Managed {
				description += " | managed by this CLI"
			}
		} else if entry.Status.Installed {
			description = "Incompatible version " + entry.Status.Version + " | requires " + entry.Toolchain.MinimumVersion
		}
		items = append(items, listItem{id: entry.Toolchain.ID, title: entry.Toolchain.Name, description: description})
	}
	return items
}

func (current model) View() string {
	if current.width == 0 {
		return "Starting Calculate Pi Scripts..."
	}
	var content string
	switch current.screen {
	case helpScreen:
		content = helpText(current.helpPage)
	case outputScreen:
		content = current.viewport.View()
	default:
		content = current.items.View()
	}
	footer := "↑/↓ Navigate  Enter Select  B Back  Q Quit  N/P Pages"
	if current.busy {
		footer = "Working...  Q Quit"
	}
	return content + "\n" + footer + "\n"
}

func helpText(page int) string {
	if page == 0 {
		return "Help\n\nArrow keys: Move through menu items\nEnter: Select the highlighted item\nB: Return to the previous menu\nQ: Quit the application\nN: Next page\nP: Previous page\n\nPress N for more help."
	}
	return "Help: Running and Toolchains\n\nSelect Calculate Pi, then choose one of the seven listed implementations.\nThe CLI checks required files and runtimes before execution.\nA missing runtime is installed only after you confirm.\n\nManage Toolchains lists only Python, Java, and Julia.\nSystem packages are removed only when this CLI recorded their installation.\n\nPress P to return to the previous help page."
}

func Run(application service) error {
	program := tea.NewProgram(newModel(application), tea.WithAltScreen())
	_, err := program.Run()
	return err
}
