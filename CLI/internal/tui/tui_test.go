package tui

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/app"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/catalog"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/toolchain"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type fakeService struct {
	implementations []catalog.Implementation
	toolchains      []app.ToolchainStatus
	toolchainCalls  int
	prepareCalls    []string
	runCalls        []string
}

func (service *fakeService) Implementations() []catalog.Implementation {
	if service.implementations != nil {
		return service.implementations
	}
	return catalog.All()
}
func (service *fakeService) Toolchains(context.Context) ([]app.ToolchainStatus, error) {
	service.toolchainCalls++
	return service.toolchains, nil
}
func (service *fakeService) Prepare(_ context.Context, implementationID string) (app.RuntimeRequirement, error) {
	service.prepareCalls = append(service.prepareCalls, implementationID)
	return app.RuntimeRequirement{}, nil
}
func (service *fakeService) Run(_ context.Context, implementationID string, _ bool, _, _ io.Writer) app.Result {
	service.runCalls = append(service.runCalls, implementationID)
	return app.Result{}
}
func (*fakeService) InstallToolchain(context.Context, string) app.Result {
	return app.Result{}
}
func (*fakeService) UninstallToolchain(context.Context, string) app.Result {
	return app.Result{}
}

func TestMainMenuHasRequestedEntries(t *testing.T) {
	current := newModel(&fakeService{})
	want := []string{"Calculate Pi", "Manage Toolchains", "Help", "Exit"}
	items := current.items.Items()
	if len(items) != len(want) {
		t.Fatalf("main menu has %d items, want %d", len(items), len(want))
	}
	for index, item := range items {
		if item.(listItem).title != want[index] {
			t.Errorf("menu item %d = %q, want %q", index, item.(listItem).title, want[index])
		}
	}
}

func TestCalculateListComesFromCatalog(t *testing.T) {
	current := newModel(&fakeService{})
	updated, _ := current.Update(tea.KeyMsg{Type: tea.KeyEnter})
	current = updated.(model)
	if current.screen != implementationScreen {
		t.Fatalf("screen = %v, want implementation list", current.screen)
	}
	items := current.items.Items()
	want := catalog.All()
	if len(items) != len(want) {
		t.Fatalf("implementation list has %d entries, want %d", len(items), len(want))
	}
	for index, implementation := range want {
		if items[index].(listItem).id != implementation.ID {
			t.Errorf("implementation %d = %q, want %q", index, items[index].(listItem).id, implementation.ID)
		}
	}
}

func TestHelpPagesRespondToNAndP(t *testing.T) {
	current := newModel(&fakeService{})
	current.screen = helpScreen
	updated, _ := current.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	current = updated.(model)
	if current.helpPage != 1 {
		t.Fatalf("help page after N = %d, want 1", current.helpPage)
	}
	updated, _ = current.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	current = updated.(model)
	if current.helpPage != 0 {
		t.Fatalf("help page after P = %d, want 0", current.helpPage)
	}
}

func TestBackReturnsToMainAndQuitQuits(t *testing.T) {
	current := newModel(&fakeService{})
	updated, _ := current.Update(tea.KeyMsg{Type: tea.KeyEnter})
	current = updated.(model)
	updated, _ = current.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	current = updated.(model)
	if current.screen != mainScreen {
		t.Fatalf("screen after B = %v, want main menu", current.screen)
	}
	_, command := current.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if command == nil {
		t.Fatal("Q did not produce a quit command")
	}
}

func TestMissingRuntimePromptsBeforeInstall(t *testing.T) {
	current := newModel(&fakeService{})
	updated, _ := current.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	current = updated.(model)
	current.screen = implementationScreen
	current.selectedImpl = "python"
	pythonRuntime, _ := catalog.FindToolchain("python")
	updated, _ = current.Update(prepareResultMsg{
		requirement: app.RuntimeRequirement{Toolchain: pythonRuntime},
	})
	result := updated.(model)
	if result.screen != confirmationScreen || result.pendingAction != installForRun {
		t.Fatalf("state after missing runtime = %#v", result)
	}
	if len(result.items.Items()) != 2 || result.items.Items()[0].(listItem).title != "Install and run" {
		t.Fatalf("confirmation choices = %#v", result.items.Items())
	}

	updated, _ = result.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	result = updated.(model)
	if result.screen != implementationScreen {
		t.Fatalf("screen after cancelling runtime install = %v", result.screen)
	}
	if len(result.items.Items()) != len(catalog.All()) {
		t.Fatalf("implementation entries after cancellation = %d, want %d", len(result.items.Items()), len(catalog.All()))
	}
}

func TestToolchainMenuContainsOnlyManagedRuntimes(t *testing.T) {
	entries := make([]app.ToolchainStatus, 0)
	for _, spec := range catalog.Toolchains() {
		entries = append(entries, app.ToolchainStatus{Toolchain: spec})
	}
	current := newModel(&fakeService{toolchains: entries})
	updated, _ := current.Update(toolchainsResultMsg{statuses: entries})
	result := updated.(model)
	if result.screen != toolchainScreen || len(result.items.Items()) != 3 {
		t.Fatalf("toolchain menu has screen %v and %d entries", result.screen, len(result.items.Items()))
	}
	for index, spec := range catalog.Toolchains() {
		if got := result.items.Items()[index].(listItem).id; got != spec.ID {
			t.Errorf("toolchain %d = %q, want %q", index, got, spec.ID)
		}
	}
}

func TestBackFromToolchainOperationRefreshesStatus(t *testing.T) {
	service := &fakeService{}
	current := newModel(service)
	current.screen = outputScreen
	current.previousScreen = toolchainActionScreen
	updated, command := current.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	if command == nil {
		t.Fatal("B did not start a toolchain status refresh")
	}
	message := command()
	result, ok := updated.(model)
	if !ok || !result.busy || result.screen != toolchainScreen {
		t.Fatalf("state while refreshing = %#v", updated)
	}
	statusMessage, ok := message.(toolchainsResultMsg)
	if !ok || statusMessage.err != nil || service.toolchainCalls != 1 {
		t.Fatalf("refresh message = %#v, calls = %d", message, service.toolchainCalls)
	}
}

func TestServiceExecForwardsOutputToReleasedTerminal(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder
	command := &serviceExec{
		operation: func(terminalStdout, terminalStderr io.Writer) app.Result {
			io.WriteString(terminalStdout, "calculation output")
			io.WriteString(terminalStderr, "calculation warning")
			return app.Result{Stdout: "calculation output", Stderr: "calculation warning"}
		},
	}
	command.SetStdout(&stdout)
	command.SetStderr(&stderr)
	if err := command.Run(); err != nil {
		t.Fatal(err)
	}
	if stdout.String() != "calculation output" || stderr.String() != "calculation warning" {
		t.Fatalf("terminal output = %q, terminal error output = %q", stdout.String(), stderr.String())
	}
}

func TestReturningFromCalculationAllowsSelectingAnotherImplementation(t *testing.T) {
	service := &fakeService{}
	current := newModel(service)
	updated, _ := current.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	current = updated.(model)
	updated, _ = current.Update(tea.KeyMsg{Type: tea.KeyEnter})
	current = updated.(model)

	firstSelection, prepareCommand := current.Update(tea.KeyMsg{Type: tea.KeyEnter})
	current = firstSelection.(model)
	prepareMessage := current.items.SelectedItem().(listItem)
	if prepareMessage.id != "python" {
		t.Fatalf("first implementation = %q, want python", prepareMessage.id)
	}
	prepared := prepareCommand()
	updated, _ = current.Update(prepared)
	current = updated.(model)
	if current.previousScreen != implementationScreen {
		t.Fatalf("return screen after run = %v, want implementation menu", current.previousScreen)
	}
	updated, _ = current.Update(operationResultMsg{result: app.Result{Message: "Completed", Stdout: "first output"}})
	current = updated.(model)
	updated, _ = current.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	current = updated.(model)
	if current.screen != implementationScreen {
		t.Fatalf("screen after B = %v, want implementation menu", current.screen)
	}
	updated, _ = current.Update(tea.KeyMsg{Type: tea.KeyDown})
	current = updated.(model)
	selected, _ := current.items.SelectedItem().(listItem)
	if selected.id != "c" {
		t.Fatalf("second implementation = %q, want c", selected.id)
	}
	_, secondPrepareCommand := current.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if secondPrepareCommand == nil {
		t.Fatal("selecting the second implementation did not start preparation")
	}
	secondPrepared := secondPrepareCommand()
	if _, ok := secondPrepared.(prepareResultMsg); !ok {
		t.Fatalf("second preparation result = %#v, want prepareResultMsg", secondPrepared)
	}
	if len(service.prepareCalls) != 2 || service.prepareCalls[0] != "python" || service.prepareCalls[1] != "c" {
		t.Fatalf("prepared implementations = %v, want [python c]", service.prepareCalls)
	}
}

func TestManageToolchainsOffersSafeActions(t *testing.T) {
	toolchains := catalog.Toolchains()
	statuses := []app.ToolchainStatus{
		{Toolchain: toolchains[0]},
		{Toolchain: toolchains[1], Status: toolchain.Status{Installed: true, Compatible: true}},
		{Toolchain: toolchains[2], Status: toolchain.Status{Installed: true, Compatible: true, Managed: true}},
	}
	current := newModel(&fakeService{toolchains: statuses})
	current.toolchainStatuses = statuses
	current.screen = toolchainScreen

	current.selectedToolchain = "python"
	updated, _ := current.openToolchainActions()
	python := updated.(model)
	if !hasAction(python.items, "install") || hasAction(python.items, "uninstall") {
		t.Fatalf("Python actions = %#v, want install only", python.items.Items())
	}
	updated, _ = python.selectToolchainAction("install")
	confirmation := updated.(model)
	if confirmation.screen != confirmationScreen || confirmation.pendingAction != installToolchain {
		t.Fatalf("Python install selection = %#v", confirmation)
	}

	current.selectedToolchain = "java"
	updated, _ = current.openToolchainActions()
	java := updated.(model)
	if hasAction(java.items, "install") || hasAction(java.items, "uninstall") || !hasAction(java.items, "available") {
		t.Fatalf("Java actions = %#v, want already available only", java.items.Items())
	}

	current.selectedToolchain = "julia"
	updated, _ = current.openToolchainActions()
	julia := updated.(model)
	if hasAction(julia.items, "install") || !hasAction(julia.items, "uninstall") {
		t.Fatalf("Julia actions = %#v, want managed uninstall", julia.items.Items())
	}
	updated, _ = julia.selectToolchainAction("uninstall")
	confirmation = updated.(model)
	if confirmation.screen != confirmationScreen || confirmation.pendingAction != uninstallToolchain {
		t.Fatalf("Julia uninstall selection = %#v", confirmation)
	}
}

func hasAction(items list.Model, id string) bool {
	for _, item := range items.Items() {
		if item.(listItem).id == id {
			return true
		}
	}
	return false
}
