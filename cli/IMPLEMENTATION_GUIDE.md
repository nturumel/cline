# Cline CLI Implementation Guide

Guide for adding OCA provider support and Workflow invocation to the CLI.

## Table of Contents

1. [OCA Provider Support](#oca-provider-support)
2. [Workflow Invocation](#workflow-invocation)
3. [Implementation Checklist](#implementation-checklist)

---

## OCA Provider Support

### Overview

OCA (Oracle Cloud AI) is a provider that requires OAuth authentication with login/logout/refresh functionality. The webview already has full support via protobuf services, but the CLI needs to expose these features.

### Current State

**What exists:**
- ✅ Protobuf service definition: `proto/cline/oca_account.proto`
- ✅ gRPC service: `OcaAccountService`
- ✅ Generated Go client: `src/generated/grpc-go/client/services/ocaaccount_client.go`
- ✅ Backend implementation in cline-core
- ✅ OCA provider settings in task settings parser

**What's missing:**
- ❌ CLI commands for OCA authentication
- ❌ CLI commands for OCA model management
- ❌ Interactive auth flow handling

### Architecture Addition

```
┌─────────────────────────────────────────────────────┐
│              NEW: OCA Command                       │
│          cli/pkg/cli/oca.go                         │
│                                                     │
│  Commands:                                          │
│  • oca login                                        │
│  • oca logout                                       │
│  • oca status                                       │
│  • oca models                                       │
│  • oca models refresh                               │
└───────────────────┬─────────────────────────────────┘
                    │
                    ▼ gRPC calls
┌─────────────────────────────────────────────────────┐
│         OcaAccountService (cline-core)              │
│                                                     │
│  • ocaAccountLoginClicked()                         │
│  • ocaAccountLogoutClicked()                        │
│  • ocaSubscribeToAuthStatusUpdate()                 │
└─────────────────────────────────────────────────────┘
```

### Files to Create/Modify

#### 1. Create `cli/pkg/cli/oca.go`

New command file for OCA operations:

```go
package cli

import (
	"context"
	"fmt"
	"github.com/cline/cli/pkg/cli/global"
	"github.com/cline/grpc-go/cline"
	"github.com/spf13/cobra"
	"time"
)

func NewOcaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oca",
		Short:   "Manage OCA (Oracle Cloud AI) authentication",
		Long:    `Login, logout, and manage OCA provider authentication and models.`,
	}

	cmd.AddCommand(newOcaLoginCommand())
	cmd.AddCommand(newOcaLogoutCommand())
	cmd.AddCommand(newOcaStatusCommand())
	cmd.AddCommand(newOcaModelsCommand())

	return cmd
}

func newOcaLoginCommand() *cobra.Command {
	var address string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to OCA",
		Long:  `Initiates OAuth flow to authenticate with Oracle Cloud AI.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get client
			var client *client.ClineClient
			var err error
			if address != "" {
				client, err = global.GetClientForAddress(ctx, address)
			} else {
				client, err = global.GetDefaultClient(ctx)
			}
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}

			// Trigger login
			fmt.Println("Initiating OCA login...")
			fmt.Println("Your browser will open for authentication.")
			
			_, err = client.OcaAccount.OcaAccountLoginClicked(ctx, &cline.EmptyRequest{})
			if err != nil {
				return fmt.Errorf("login failed: %w", err)
			}

			fmt.Println("✓ Login flow started")
			fmt.Println("Complete the authentication in your browser.")
			fmt.Println("\nRun 'cline oca status' to check authentication status.")

			return nil
		},
	}

	cmd.Flags().StringVar(&address, "address", "", "specific Cline instance address")
	return cmd
}

func newOcaLogoutCommand() *cobra.Command {
	var address string

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from OCA",
		Long:  `Clears OCA authentication tokens and user state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get client
			var client *client.ClineClient
			var err error
			if address != "" {
				client, err = global.GetClientForAddress(ctx, address)
			} else {
				client, err = global.GetDefaultClient(ctx)
			}
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}

			// Logout
			_, err = client.OcaAccount.OcaAccountLogoutClicked(ctx, &cline.EmptyRequest{})
			if err != nil {
				return fmt.Errorf("logout failed: %w", err)
			}

			fmt.Println("✓ Successfully logged out from OCA")
			return nil
		},
	}

	cmd.Flags().StringVar(&address, "address", "", "specific Cline instance address")
	return cmd
}

func newOcaStatusCommand() *cobra.Command {
	var address string
	var watch bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check OCA authentication status",
		Long:  `Display current OCA authentication status and user information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get client
			var client *client.ClineClient
			var err error
			if address != "" {
				client, err = global.GetClientForAddress(ctx, address)
			} else {
				client, err = global.GetDefaultClient(ctx)
			}
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}

			if watch {
				// Watch mode - stream updates
				fmt.Println("Watching OCA authentication status (Press Ctrl+C to stop)...")
				return watchOcaStatus(ctx, client)
			}

			// One-time status check
			// Subscribe briefly to get current state
			statusChan := make(chan *cline.OcaAuthState, 1)
			errChan := make(chan error, 1)

			stream, err := client.OcaAccount.OcaSubscribeToAuthStatusUpdate(ctx, &cline.EmptyRequest{})
			if err != nil {
				return fmt.Errorf("failed to subscribe: %w", err)
			}

			go func() {
				state, err := stream.Recv()
				if err != nil {
					errChan <- err
					return
				}
				statusChan <- state
			}()

			select {
			case state := <-statusChan:
				displayOcaStatus(state)
				return nil
			case err := <-errChan:
				return fmt.Errorf("failed to get status: %w", err)
			case <-time.After(5 * time.Second):
				return fmt.Errorf("timeout waiting for status")
			}
		},
	}

	cmd.Flags().StringVar(&address, "address", "", "specific Cline instance address")
	cmd.Flags().BoolVarP(&watch, "watch", "w", false, "watch for status changes")
	return cmd
}

func watchOcaStatus(ctx context.Context, client *client.ClineClient) error {
	stream, err := client.OcaAccount.OcaSubscribeToAuthStatusUpdate(ctx, &cline.EmptyRequest{})
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	for {
		state, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("stream error: %w", err)
		}

		fmt.Printf("\n[%s] Status Update:\n", time.Now().Format("15:04:05"))
		displayOcaStatus(state)
	}
}

func displayOcaStatus(state *cline.OcaAuthState) {
	if state == nil || state.User == nil || state.User.Uid == "" {
		fmt.Println("Status: Not authenticated")
		fmt.Println("\nRun 'cline oca login' to authenticate.")
		return
	}

	fmt.Println("Status: Authenticated ✓")
	fmt.Printf("User ID: %s\n", state.User.Uid)
	
	if state.User.DisplayName != nil && *state.User.DisplayName != "" {
		fmt.Printf("Name: %s\n", *state.User.DisplayName)
	}
	
	if state.User.Email != nil && *state.User.Email != "" {
		fmt.Printf("Email: %s\n", *state.User.Email)
	}

	if state.ApiKey != nil && *state.ApiKey != "" {
		// Show first/last few characters
		key := *state.ApiKey
		if len(key) > 20 {
			fmt.Printf("Token: %s...%s\n", key[:8], key[len(key)-8:])
		}
	}
}

func newOcaModelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "models",
		Short: "Manage OCA models",
		Long:  `List and refresh available OCA models.`,
	}

	cmd.AddCommand(newOcaModelsListCommand())
	cmd.AddCommand(newOcaModelsRefreshCommand())

	return cmd
}

func newOcaModelsListCommand() *cobra.Command {
	var address string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available OCA models",
		Long:  `Display all available OCA models with their capabilities.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// This would need a new gRPC service method to fetch models
			// Currently models are fetched client-side in the webview
			fmt.Println("Model listing from CLI not yet implemented.")
			fmt.Println("Use the VSCode extension UI to view available models.")
			return nil
		},
	}

	cmd.Flags().StringVar(&address, "address", "", "specific Cline instance address")
	return cmd
}

func newOcaModelsRefreshCommand() *cobra.Command {
	var address string

	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh OCA models",
		Long:  `Refresh the list of available OCA models from the API.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// This would need a new gRPC service method
			fmt.Println("Model refresh from CLI not yet implemented.")
			fmt.Println("Models are automatically refreshed in the VSCode extension.")
			return nil
		},
	}

	cmd.Flags().StringVar(&address, "address", "", "specific Cline instance address")
	return cmd
}
```

#### 2. Modify `cli/cmd/cline/main.go`

Add OCA command:

```go
rootCmd.AddCommand(cli.NewTaskCommand())
rootCmd.AddCommand(cli.NewInstanceCommand())
rootCmd.AddCommand(cli.NewVersionCommand())
rootCmd.AddCommand(cli.NewAuthCommand())
rootCmd.AddCommand(cli.NewTaskSendCommand())
rootCmd.AddCommand(cli.NewOcaCommand())  // ADD THIS
```

#### 3. Update Documentation

Add to `cli/README.md`:

```markdown
### OCA Authentication

Authenticate with Oracle Cloud AI:

\`\`\`bash
# Login
./cli/bin/cline oca login

# Check status
./cli/bin/cline oca status

# Watch for status changes
./cli/bin/cline oca status --watch

# Logout
./cli/bin/cline oca logout
\`\`\`
```

### Usage Examples

```bash
# Login to OCA
./cli/bin/cline oca login

# Check if logged in
./cli/bin/cline oca status

# Create task with OCA provider
./cli/bin/cline task new "Write tests" \
  -s act-mode-api-provider=oca \
  -s act-mode-oca-model-id=your-model-id \
  -s oca-base-url=https://your-oca-endpoint

# Logout when done
./cli/bin/cline oca logout
```

---

## Workflow Invocation

### Overview

Workflows are markdown files that contain step-by-step instructions for Cline. They're invoked in the webview via `/workflow-name.md` slash commands. We need to add CLI support.

### Current State

**What exists:**
- ✅ Workflow storage: `~/Documents/Cline/Workflows/` (local) and global workflows
- ✅ Protobuf toggle support: `toggleWorkflow` RPC in `file.proto`
- ✅ Slash command parsing in cline-core: `src/core/slash-commands/index.ts`
- ✅ Workflow execution through regular task creation

**What's missing:**
- ❌ CLI commands to list workflows
- ❌ CLI commands to invoke workflows
- ❌ CLI commands to create/edit workflows
- ❌ CLI commands to toggle workflows on/off

### Architecture Addition

```
┌─────────────────────────────────────────────────────┐
│           NEW: Workflow Command                     │
│         cli/pkg/cli/workflow.go                     │
│                                                     │
│  Commands:                                          │
│  • workflow list                                    │
│  • workflow run <name>                              │
│  • workflow create <name>                           │
│  • workflow edit <name>                             │
│  • workflow toggle <name>                           │
│  • workflow delete <name>                           │
└───────────────────┬─────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────┐
│      Task Creation with /workflow syntax            │
│                                                     │
│  Workflow invocation happens by:                    │
│  1. Reading workflow file from disk                 │
│  2. Creating task with "/workflow-name.md" prefix   │
│  3. Cline-core processes slash command              │
│  4. Workflow content injected as instructions       │
└─────────────────────────────────────────────────────┘
```

### Files to Create/Modify

#### 1. Create `cli/pkg/cli/workflow.go`

New command file for workflow operations:

```go
package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cline/cli/pkg/cli/global"
	"github.com/cline/cli/pkg/cli/task"
	"github.com/cline/grpc-go/cline"
	"github.com/spf13/cobra"
)

func NewWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workflow",
		Aliases: []string{"wf"},
		Short:   "Manage and run Cline workflows",
		Long:    `List, create, edit, and run Cline workflows.`,
	}

	cmd.AddCommand(newWorkflowListCommand())
	cmd.AddCommand(newWorkflowRunCommand())
	cmd.AddCommand(newWorkflowCreateCommand())
	cmd.AddCommand(newWorkflowEditCommand())
	cmd.AddCommand(newWorkflowDeleteCommand())
	cmd.AddCommand(newWorkflowToggleCommand())

	return cmd
}

// Get workflows directory path
func getWorkflowsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// Default to ~/Documents/Cline/Workflows
	return filepath.Join(homeDir, "Documents", "Cline", "Workflows"), nil
}

func newWorkflowListCommand() *cobra.Command {
	var showGlobal bool
	var showLocal bool

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List available workflows",
		Long:    `List all workflows in local and global directories.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowsDir, err := getWorkflowsDir()
			if err != nil {
				return fmt.Errorf("failed to get workflows directory: %w", err)
			}

			// List workflows
			files, err := ioutil.ReadDir(workflowsDir)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("No workflows directory found.")
					fmt.Printf("Create one with: mkdir -p %s\n", workflowsDir)
					return nil
				}
				return fmt.Errorf("failed to read workflows directory: %w", err)
			}

			workflows := []string{}
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
					workflows = append(workflows, file.Name())
				}
			}

			if len(workflows) == 0 {
				fmt.Println("No workflows found.")
				fmt.Printf("Create one with: cline workflow create <name>\n")
				return nil
			}

			fmt.Printf("Available workflows (%d):\n\n", len(workflows))
			for i, wf := range workflows {
				fmt.Printf("%d. %s\n", i+1, wf)
			}

			fmt.Printf("\nRun a workflow with: cline workflow run <name>\n")
			return nil
		},
	}

	cmd.Flags().BoolVar(&showGlobal, "global", false, "show global workflows")
	cmd.Flags().BoolVar(&showLocal, "local", false, "show local workflows")

	return cmd
}

func newWorkflowRunCommand() *cobra.Command {
	var wait bool
	var workspaces []string
	var address string
	var settings []string

	cmd := &cobra.Command{
		Use:     "run <workflow-name>",
		Aliases: []string{"r", "invoke"},
		Short:   "Run a workflow",
		Long:    `Invoke a workflow by creating a task with the workflow slash command.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			workflowName := args[0]

			// Add .md extension if not present
			if !strings.HasSuffix(workflowName, ".md") {
				workflowName += ".md"
			}

			// Verify workflow exists
			workflowsDir, err := getWorkflowsDir()
			if err != nil {
				return err
			}

			workflowPath := filepath.Join(workflowsDir, workflowName)
			if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
				return fmt.Errorf("workflow not found: %s\nRun 'cline workflow list' to see available workflows", workflowName)
			}

			// Create task with slash command
			prompt := fmt.Sprintf("/%s", workflowName)
			
			fmt.Printf("Running workflow: %s\n", workflowName)

			// Ensure task manager
			if err := ensureTaskManager(ctx, address); err != nil {
				return err
			}

			// Create task
			taskID, err := taskManager.CreateTask(ctx, prompt, nil, nil, workspaces, settings)
			if err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}

			fmt.Printf("Workflow task created with ID: %s\n", taskID)
			fmt.Printf("Using instance: %s\n", taskManager.GetCurrentInstance())

			// Wait for completion if requested
			if wait {
				fmt.Println("Following workflow execution...")
				return taskManager.FollowConversation(ctx)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&wait, "wait", false, "wait for workflow completion")
	cmd.Flags().StringSliceVarP(&workspaces, "workdir", "w", nil, "working directory paths")
	cmd.Flags().StringVar(&address, "address", "", "specific Cline instance address")
	cmd.Flags().StringSliceVarP(&settings, "setting", "s", nil, "task settings")

	return cmd
}

func newWorkflowCreateCommand() *cobra.Command {
	var edit bool

	cmd := &cobra.Command{
		Use:     "create <workflow-name>",
		Aliases: []string{"new"},
		Short:   "Create a new workflow",
		Long:    `Create a new workflow markdown file.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowName := args[0]

			// Add .md extension if not present
			if !strings.HasSuffix(workflowName, ".md") {
				workflowName += ".md"
			}

			workflowsDir, err := getWorkflowsDir()
			if err != nil {
				return err
			}

			// Create workflows directory if it doesn't exist
			if err := os.MkdirAll(workflowsDir, 0755); err != nil {
				return fmt.Errorf("failed to create workflows directory: %w", err)
			}

			workflowPath := filepath.Join(workflowsDir, workflowName)

			// Check if workflow already exists
			if _, err := os.Stat(workflowPath); err == nil {
				return fmt.Errorf("workflow already exists: %s", workflowName)
			}

			// Create template content
			template := `# ` + strings.TrimSuffix(workflowName, ".md") + `

## Description
Brief description of what this workflow does.

## Steps

1. First step description
   - Detail about first step
   
2. Second step description
   - Detail about second step

3. Final step description
   - Detail about final step

## Notes
Any additional notes or requirements.
`

			// Write file
			if err := ioutil.WriteFile(workflowPath, []byte(template), 0644); err != nil {
				return fmt.Errorf("failed to create workflow file: %w", err)
			}

			fmt.Printf("✓ Created workflow: %s\n", workflowPath)
			
			if edit {
				// Open in default editor
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vi"
				}
				fmt.Printf("Opening in %s...\n", editor)
				// You would exec the editor here
			} else {
				fmt.Printf("\nEdit with: $EDITOR %s\n", workflowPath)
				fmt.Printf("Run with: cline workflow run %s\n", strings.TrimSuffix(workflowName, ".md"))
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&edit, "edit", "e", false, "open in editor after creation")

	return cmd
}

func newWorkflowEditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit <workflow-name>",
		Aliases: []string{"e"},
		Short:   "Edit a workflow",
		Long:    `Open a workflow file in your default editor.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowName := args[0]

			if !strings.HasSuffix(workflowName, ".md") {
				workflowName += ".md"
			}

			workflowsDir, err := getWorkflowsDir()
			if err != nil {
				return err
			}

			workflowPath := filepath.Join(workflowsDir, workflowName)

			if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
				return fmt.Errorf("workflow not found: %s", workflowName)
			}

			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}

			fmt.Printf("Edit the file: %s\n", workflowPath)
			fmt.Printf("(Open with: %s %s)\n", editor, workflowPath)

			return nil
		},
	}

	return cmd
}

func newWorkflowDeleteCommand() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "delete <workflow-name>",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete a workflow",
		Long:    `Delete a workflow file.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowName := args[0]

			if !strings.HasSuffix(workflowName, ".md") {
				workflowName += ".md"
			}

			workflowsDir, err := getWorkflowsDir()
			if err != nil {
				return err
			}

			workflowPath := filepath.Join(workflowsDir, workflowName)

			if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
				return fmt.Errorf("workflow not found: %s", workflowName)
			}

			if !force {
				fmt.Printf("Delete workflow: %s? (y/N): ", workflowName)
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
					fmt.Println("Cancelled.")
					return nil
				}
			}

			if err := os.Remove(workflowPath); err != nil {
				return fmt.Errorf("failed to delete workflow: %w", err)
			}

			fmt.Printf("✓ Deleted workflow: %s\n", workflowName)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "skip confirmation")

	return cmd
}

func newWorkflowToggleCommand() *cobra.Command {
	var enable bool
	var disable bool
	var global bool
	var address string

	cmd := &cobra.Command{
		Use:   "toggle <workflow-name>",
		Short: "Toggle a workflow on or off",
		Long:  `Enable or disable a workflow.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			workflowName := args[0]

			if !strings.HasSuffix(workflowName, ".md") {
				workflowName += ".md"
			}

			if enable && disable {
				return fmt.Errorf("cannot use both --enable and --disable")
			}

			var client *client.ClineClient
			var err error
			if address != "" {
				client, err = global.GetClientForAddress(ctx, address)
			} else {
				client, err = global.GetDefaultClient(ctx)
			}
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}

			// Get workflow path
			workflowsDir, err := getWorkflowsDir()
			if err != nil {
				return err
			}
			workflowPath := filepath.Join(workflowsDir, workflowName)

			// Toggle workflow
			_, err = client.File.ToggleWorkflow(ctx, &cline.ToggleWorkflowRequest{
				WorkflowPath: workflowPath,
				Enabled:      enable,
				IsGlobal:     global,
			})
			if err != nil {
				return fmt.Errorf("failed to toggle workflow: %w", err)
			}

			status := "enabled"
			if disable {
				status = "disabled"
			}
			fmt.Printf("✓ Workflow %s: %s\n", workflowName, status)

			return nil
		},
	}

	cmd.Flags().BoolVar(&enable, "enable", false, "enable the workflow")
	cmd.Flags().BoolVar(&disable, "disable", false, "disable the workflow")
	cmd.Flags().BoolVar(&global, "global", false, "toggle global workflow")
	cmd.Flags().StringVar(&address, "address", "", "specific Cline instance address")

	return cmd
}
```

#### 2. Modify `cli/cmd/cline/main.go`

Add workflow command:

```go
rootCmd.AddCommand(cli.NewTaskCommand())
rootCmd.AddCommand(cli.NewInstanceCommand())
rootCmd.AddCommand(cli.NewVersionCommand())
rootCmd.AddCommand(cli.NewAuthCommand())
rootCmd.AddCommand(cli.NewTaskSendCommand())
rootCmd.AddCommand(cli.NewOcaCommand())
rootCmd.AddCommand(cli.NewWorkflowCommand())  // ADD THIS
```

#### 3. Update Documentation

Add to `cli/README.md`:

```markdown
### Workflow Management

Manage and run Cline workflows:

\`\`\`bash
# List workflows
./cli/bin/cline workflow list

# Create new workflow
./cli/bin/cline workflow create my-deploy

# Edit workflow
./cli/bin/cline workflow edit my-deploy

# Run workflow
./cli/bin/cline workflow run my-deploy

# Run and follow
./cli/bin/cline workflow run my-deploy --wait

# Delete workflow
./cli/bin/cline workflow delete my-deploy
\`\`\`
```

### Usage Examples

```bash
# List all workflows
./cli/bin/cline workflow list

# Create a new deployment workflow
./cli/bin/cline workflow create deploy-prod

# Edit the workflow
$EDITOR ~/Documents/Cline/Workflows/deploy-prod.md

# Run the workflow
./cli/bin/cline workflow run deploy-prod --wait

# Run with custom working directory
./cli/bin/cline workflow run deploy-prod -w /projects/myapp

# Run with settings
./cli/bin/cline workflow run deploy-prod \
  -s yolo-mode-toggled=true \
  -s auto-approval-settings.actions.execute-all-commands=true
```

---

## Implementation Checklist

### Phase 1: OCA Support

- [ ] Create `cli/pkg/cli/oca.go`
- [ ] Implement `oca login` command
- [ ] Implement `oca logout` command
- [ ] Implement `oca status` command (with `--watch` flag)
- [ ] Implement `oca models` commands (list/refresh)
- [ ] Register OCA command in `main.go`
- [ ] Update `README.md` with OCA examples
- [ ] Update `USER_MANUAL.md` with OCA section
- [ ] Update `QUICK_REFERENCE.md` with OCA commands
- [ ] Test OCA login flow
- [ ] Test OCA task creation with provider setting
- [ ] Test OCA logout

### Phase 2: Workflow Support

- [ ] Create `cli/pkg/cli/workflow.go`
- [ ] Implement `workflow list` command
- [ ] Implement `workflow create` command
- [ ] Implement `workflow edit` command
- [ ] Implement `workflow run` command
- [ ] Implement `workflow delete` command
- [ ] Implement `workflow toggle` command
- [ ] Register workflow command in `main.go`
- [ ] Update `README.md` with workflow examples
- [ ] Update `USER_MANUAL.md` with workflow section
- [ ] Update `QUICK_REFERENCE.md` with workflow commands
- [ ] Test workflow creation
- [ ] Test workflow execution
- [ ] Test workflow with various settings

### Phase 3: Documentation

- [ ] Update `ARCHITECTURE.md` with new components
- [ ] Create workflow examples in `examples.sh`
- [ ] Add OCA authentication flow diagram
- [ ] Add workflow invocation flow diagram
- [ ] Document OCA-specific settings
- [ ] Document workflow file format
- [ ] Add troubleshooting section for OCA
- [ ] Add troubleshooting section for workflows

### Phase 4: Testing

- [ ] Write unit tests for OCA commands
- [ ] Write unit tests for workflow commands
- [ ] Write integration tests for OCA flow
- [ ] Write integration tests for workflow execution
- [ ] Test cross-platform (macOS, Linux, Windows)
- [ ] Test with multiple instances
- [ ] Test error handling

---

## Implementation Notes

### OCA Authentication Flow

1. User runs `cline oca login`
2. CLI calls `ocaAccountLoginClicked()`
3. Cline-core opens browser with OAuth URL
4. User completes authentication in browser
5. OAuth callback hits cline-core
6. Cline-core stores tokens and broadcasts auth state
7. CLI can subscribe to auth updates with `ocaSubscribeToAuthStatusUpdate()`
8. User runs `cline oca status` to verify

### Workflow Invocation Flow

1. User runs `cline workflow run deploy-prod`
2. CLI verifies workflow file exists in `~/Documents/Cline/Workflows/`
3. CLI creates task with prompt `/deploy-prod.md`
4. Cline-core's slash command parser detects workflow
5. Workflow content is read and injected as `<explicit_instructions>`
6. Task executes with workflow instructions
7. CLI follows task progress

### Key Considerations

1. **OCA Token Refresh**: The backend handles token refresh automatically. CLI just needs to display current state.

2. **Workflow Files**: Workflows are just markdown files. The CLI doesn't need to parse them—just pass the slash command to cline-core.

3. **Multiple Instances**: Both OCA auth and workflows are per-instance. Users can have different OCA accounts on different instances.

4. **Global vs Local Workflows**: 
   - Local: `~/Documents/Cline/Workflows/`
   - Global: Platform-specific location (needs investigation)

5. **Error Handling**: Both OCA and workflows should provide clear error messages when:
   - Not authenticated
   - Workflow file not found
   - Network errors
   - Permission errors

---

## Testing Commands

### OCA Testing

```bash
# Build CLI
cd /path/to/cline
npm run compile-cli

# Test OCA commands
./cli/bin/cline oca login
./cli/bin/cline oca status
./cli/bin/cline oca status --watch
./cli/bin/cline oca logout

# Test OCA task creation
./cli/bin/cline task new "Test OCA" \
  -s act-mode-api-provider=oca \
  -s act-mode-oca-model-id=test-model

# Test with specific instance
./cli/bin/cline instance new
./cli/bin/cline oca login --address localhost:50053
```

### Workflow Testing

```bash
# Test workflow commands
./cli/bin/cline workflow list
./cli/bin/cline workflow create test-workflow
./cli/bin/cline workflow run test-workflow
./cli/bin/cline workflow run test-workflow --wait
./cli/bin/cline workflow delete test-workflow

# Test with real workflow
echo "# Test Workflow\nRun: \`echo Hello from workflow\`" > ~/Documents/Cline/Workflows/hello.md
./cli/bin/cline workflow run hello --wait
```

---

## Additional Features to Consider

### OCA Enhancements

1. **Model Management**
   - Add gRPC service to fetch OCA models
   - Implement `cline oca models list`
   - Implement `cline oca models refresh`
   - Cache models locally

2. **Configuration**
   - `cline oca config set base-url <url>`
   - `cline oca config set mode internal|external`

### Workflow Enhancements

1. **Template System**
   - `cline workflow template list`
   - `cline workflow template use <template>`

2. **Workflow Variables**
   - Support for `${VAR}` in workflow files
   - Pass variables: `cline workflow run deploy --var env=prod`

3. **Workflow History**
   - Track workflow execution history
   - `cline workflow history`

4. **Workflow Validation**
   - `cline workflow validate <name>`
   - Check for syntax errors, tool availability

---

## Files Summary

### Files to Create

1. `cli/pkg/cli/oca.go` - OCA authentication commands
2. `cli/pkg/cli/workflow.go` - Workflow management commands

### Files to Modify

1. `cli/cmd/cline/main.go` - Register new commands
2. `cli/README.md` - Add command documentation
3. `cli/USER_MANUAL.md` - Add usage guides
4. `cli/QUICK_REFERENCE.md` - Add command reference
5. `cli/ARCHITECTURE.md` - Update architecture diagrams
6. `cli/examples.sh` - Add example commands

### Files to Reference

1. `proto/cline/oca_account.proto` - OCA service definition
2. `proto/cline/file.proto` - Workflow toggle definition
3. `src/services/auth/oca/OcaAuthService.ts` - OCA backend implementation
4. `src/core/slash-commands/index.ts` - Workflow parsing logic
5. `webview-ui/src/components/settings/providers/OcaProvider.tsx` - UI reference

---

## Timeline Estimate

### OCA Support: ~2-3 days
- Day 1: Implement basic commands (login/logout/status)
- Day 2: Add streaming status updates, error handling
- Day 3: Documentation, testing, polish

### Workflow Support: ~2-3 days
- Day 1: Implement basic commands (list/create/run)
- Day 2: Add edit/delete/toggle functionality
- Day 3: Documentation, testing, examples

### Total: ~1 week for both features

---

**Questions? Issues?**

If you encounter any issues during implementation:
1. Check the protobuf definitions for exact message formats
2. Reference the webview implementations for behavior
3. Use `--verbose` flag for debugging
4. Check cline-core debug logs in `/tmp/`

Good luck with the implementation! 🚀

