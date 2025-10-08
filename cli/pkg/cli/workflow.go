package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cline/cli/pkg/cli/global"
	"github.com/cline/grpc-go/client"
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
			entries, err := os.ReadDir(workflowsDir)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("No workflows directory found.")
					fmt.Printf("Create one with: mkdir -p %s\n", workflowsDir)
					return nil
				}
				return fmt.Errorf("failed to read workflows directory: %w", err)
			}

			workflows := []string{}
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
					workflows = append(workflows, entry.Name())
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
	var yolo bool
	var mode string

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

			// Set mode if provided
			if mode != "" {
				if err := taskManager.SetMode(ctx, mode, nil, nil, nil); err != nil {
					return fmt.Errorf("failed to set mode: %w", err)
				}
				fmt.Printf("Mode set to: %s\n", mode)
			}

			// Inject yolo_mode_toggled setting if --yolo flag is set
			if yolo {
				settings = append(settings, "yolo_mode_toggled=true")
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
	cmd.Flags().BoolVarP(&yolo, "yolo", "y", false, "enable yolo mode (non-interactive)")
	cmd.Flags().StringVarP(&mode, "mode", "m", "", "mode (act|plan)")

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
			baseName := strings.TrimSuffix(workflowName, ".md")
			template := fmt.Sprintf(`# %s

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
`, baseName)

			// Write file
			if err := os.WriteFile(workflowPath, []byte(template), 0644); err != nil {
				return fmt.Errorf("failed to create workflow file: %w", err)
			}

			fmt.Printf("✓ Created workflow: %s\n", workflowPath)

			if edit {
				// Open in default editor
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vi"
				}
				fmt.Printf("To edit, run: %s %s\n", editor, workflowPath)
			} else {
				fmt.Printf("\nEdit with: $EDITOR %s\n", workflowPath)
				fmt.Printf("Run with: cline workflow run %s\n", baseName)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&edit, "edit", "e", false, "show editor command after creation")

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
	var isGlobal bool
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

			if !enable && !disable {
				return fmt.Errorf("must specify either --enable or --disable")
			}

			var c *client.ClineClient
			var err error
			if address != "" {
				c, err = global.GetClientForAddress(ctx, address)
			} else {
				c, err = global.GetDefaultClient(ctx)
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
			_, err = c.File.ToggleWorkflow(ctx, &cline.ToggleWorkflowRequest{
				WorkflowPath: workflowPath,
				Enabled:      enable,
				IsGlobal:     isGlobal,
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
	cmd.Flags().BoolVar(&isGlobal, "global", false, "toggle global workflow")
	cmd.Flags().StringVar(&address, "address", "", "specific Cline instance address")

	return cmd
}

