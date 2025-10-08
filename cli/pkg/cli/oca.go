package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/cline/cli/pkg/cli/global"
	"github.com/cline/grpc-go/client"
	"github.com/cline/grpc-go/cline"
	"github.com/spf13/cobra"
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

			_, err = client.Ocaaccount.OcaAccountLoginClicked(ctx, &cline.EmptyRequest{})
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
			_, err = client.Ocaaccount.OcaAccountLogoutClicked(ctx, &cline.EmptyRequest{})
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

			stream, err := client.Ocaaccount.OcaSubscribeToAuthStatusUpdate(ctx, &cline.EmptyRequest{})
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
	stream, err := client.Ocaaccount.OcaSubscribeToAuthStatusUpdate(ctx, &cline.EmptyRequest{})
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
			fmt.Println("\nTo use OCA models, specify them in task settings:")
			fmt.Println("  cline task new \"prompt\" -s act-mode-api-provider=oca -s act-mode-oca-model-id=model-name")
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

