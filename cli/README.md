# Cline CLI

```
/_____/\ /_/\      /_______/\/__/\ /__/\ /_____/\     
\:::__\/ \:\ \     \__.::._\/\::\_\\  \ \\::::_\/_    
 \:\ \  __\:\ \       \::\ \  \:. `-\  \ \\:\/___/\   
  \:\ \/_/\\:\ \____  _\::\ \__\:. _    \ \\::___\/_  
   \:\_\ \ \\:\/___/\/__\::\__/\\. \`-\  \ \\:\____/\ 
    \_____\/ \_____\/\________\/ \__\/ \__\/ \_____\/ 
```

A powerful command-line interface for interacting with the Cline AI coding assistant. The CLI provides access to Cline's task management, instance control, and monitoring capabilities directly from your terminal.

## Table of Contents

- [Features](#features)
- [Building the CLI](#building-the-cli)
- [Quick Start](#quick-start)
- [Documentation](#documentation)
- [Command Reference](#command-reference)
  - [Version](#version)
  - [Instance Management](#instance-management)
  - [Task Management](#task-management)
  - [Authentication](#authentication)
  - [OCA Provider](#oca-provider)
  - [Workflow Management](#workflow-management)
- [Common Workflows](#common-workflows)
- [Global Flags](#global-flags)
- [Tips & Troubleshooting](#tips--troubleshooting)

## Features

- 🚀 **Task Management**: Create, monitor, and control Cline AI tasks
- 🔄 **Instance Management**: Run multiple Cline instances simultaneously (similar to `kubectl` contexts)
- 📊 **Real-time Streaming**: Follow task conversations in real-time with rich terminal output
- 🎯 **Flexible Modes**: Support for interactive (act) and autonomous (plan) modes
- ⚡ **Quick Operations**: Oneshot command for fire-and-forget tasks
- 📎 **File Attachments**: Attach files and images to your tasks
- 🎨 **Multiple Output Formats**: Rich terminal UI, JSON, or plain text output

## Building the CLI

From the project root:

```bash
npm run compile-cli
```

This will:
1. Generate Protocol Buffer code
2. Build the Go CLI binary to `cli/bin/cline`
3. Build the host bridge binary to `cli/bin/cline-host`

You can also build manually from the `cli` directory:

```bash
cd cli
go build -o bin/cline ./cmd/cline
```

## Quick Start

> **⚠️ Important:** The CLI must be run from the Cline project root directory (where this README's parent directory is located), not from within the `cli/` directory itself. The CLI expects to find binaries at `./cli/bin/cline-host` and other resources relative to the project root.

### Try the Examples Script

Run the included examples script to see the CLI in action:

```bash
cd /path/to/cline
./cli/examples.sh
```

This will demonstrate various CLI commands without requiring a running Cline instance.

## Documentation

Comprehensive documentation is available:

- **[USER_MANUAL.md](USER_MANUAL.md)** - Complete usage guide with examples for:
  - Starting and following tasks
  - Switching providers and models (40+ providers supported!)
  - Working with modes (act/plan/yolo)
  - Managing multiple instances
  - Advanced task settings
  - Auto-approval configuration
  - Real-world workflows
  - Complete settings reference

- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Technical architecture documentation:
  - High-level architecture diagram
  - Component details and interactions
  - Data flow diagrams
  - Design patterns used
  - Process architecture
  - Extension points

- **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - Command cheat sheet:
  - All commands with examples
  - Shortcuts and aliases
  - Common patterns
  - Quick troubleshooting

### 1. Check Installation

```bash
# From the cline project root
./cli/bin/cline version
```

### 2. Create Your First Task

The simplest way to get started is to create a task. This will automatically start a Cline instance if none exists:

```bash
./cli/bin/cline task new "Create a hello world function in Python"
```

### 3. List Running Instances

```bash
./cli/bin/cline instance list
```

Output:
```
ADDRESS            STATUS   VERSION  LAST SEEN  PID    DEFAULT
localhost:50052    SERVING  dev      14:23:45   12345  *
```

### 4. Follow Task Progress

```bash
./cli/bin/cline task follow
```

## Command Reference

### Version

Display version information about the CLI.

```bash
# Full version info
./cli/bin/cline version

# Short version only
./cli/bin/cline version --short
```

**Example Output:**
```
Cline Go Host
Version:    dev
Commit:     unknown
Built:      unknown
Built by:   unknown
Go version: go1.24.7
OS/Arch:    darwin/amd64
```

---

### Instance Management

Manage multiple Cline instances (similar to kubectl contexts). Each instance runs independently with its own port.

#### List Instances

```bash
./cli/bin/cline instance list
# Shorthand: cline i l
```

Shows all running instances with their status, version, PID, and which one is default.

#### Create New Instance

```bash
./cli/bin/cline instance new
# Shorthand: cline i n
```

**Example Output:**
```
Starting new Cline instance...
Successfully started new instance:
  Address: localhost:50053
  Core Port: 50053
  Host Bridge Port: 50103
  Status: Default instance
```

#### Switch Default Instance

```bash
./cli/bin/cline instance use localhost:50053
# Shorthand: cline i u localhost:50053
```

All subsequent task commands will use this instance unless you specify `--address`.

#### Kill Instance

```bash
# Kill specific instance
./cli/bin/cline instance kill localhost:50053
# Shorthand: cline i k localhost:50053

# Kill all instances
./cli/bin/cline instance kill --all
```

---

### Task Management

Create and manage AI-powered coding tasks.

#### Create New Task

Create a new task and optionally wait for completion:

```bash
# Basic task
./cli/bin/cline task new "Fix the bug in auth.go"
# Shorthand: cline t n "Fix the bug in auth.go"

# With attached files
./cli/bin/cline task new "Review this code" -f auth.go -f utils.go

# With images (for UI/design tasks)
./cli/bin/cline task new "Implement this design" -i mockup.png

# With specific working directory
./cli/bin/cline task new "Set up the project" -w /path/to/project

# In plan mode (autonomous)
./cli/bin/cline task new "Refactor the database layer" -m plan

# In yolo mode (non-interactive)
./cli/bin/cline task new "Update dependencies" -y

# Wait for completion
./cli/bin/cline task new "Run tests and fix failures" --wait

# With custom settings
./cli/bin/cline task new "Deploy to AWS" -s aws-region=us-west-2 -s environment=production

# To specific instance
./cli/bin/cline task new "Debug issue" --address localhost:50053
```

**Flags:**
- `-f, --file`: Attach files to the task
- `-i, --image`: Attach images
- `-w, --workdir`: Specify working directory paths
- `-m, --mode`: Set mode (act or plan)
- `-y, --yolo`: Enable non-interactive mode
- `-s, --setting`: Set task-specific settings (key=value format)
- `--wait`: Wait for task to complete
- `--address`: Use specific instance

#### Oneshot Task

Create a task in yolo+plan mode (fully autonomous) and stream until completion. Perfect for CI/CD or scripting:

```bash
./cli/bin/cline task oneshot "Run all tests and generate report"
# Shorthand: cline t o "Run all tests and generate report"

# With files
./cli/bin/cline task oneshot "Analyze code quality" -f src/**/*.go

# With custom settings
./cli/bin/cline task oneshot "Deploy application" -s environment=staging
```

This is equivalent to:
```bash
cline task new "..." -y -m plan --wait
```

#### Follow Task

Stream the current task conversation in real-time:

```bash
./cli/bin/cline task follow
# Shorthand: cline t f

# Follow task on specific instance
./cli/bin/cline task follow --address localhost:50053
```

**Use Case:** Connect to an ongoing task and watch progress in real-time.

#### Send Message

Send a followup message to the current task:

```bash
# Send text message
./cli/bin/cline send "That looks good, please continue"
# Shorthand: cline s "That looks good, please continue"

# With attached files
./cli/bin/cline send "Review these changes" -f updated.go

# Change mode
./cli/bin/cline send --mode plan

# Approve pending request
./cli/bin/cline send --approve true

# Deny pending request
./cli/bin/cline send --approve false

# Read from stdin
echo "Here's more context..." | ./cli/bin/cline send
```

**Flags:**
- `-m, --mode`: Change task mode (act or plan)
- `-a, --approve`: Approve (true) or deny (false) pending request
- `-f, --file`: Attach files
- `-i, --image`: Attach images

#### View Task

View task conversation with various options:

```bash
# Follow until completion
./cli/bin/cline task view
# Shorthand: cline t v

# View current state without following
./cli/bin/cline task view --current

# Get just the summary
./cli/bin/cline task view --summary
```

#### List Tasks

View recent task history:

```bash
./cli/bin/cline task list
# Shorthand: cline t l
```

#### Resume Task

Resume a previous task by ID:

```bash
./cli/bin/cline task resume <task-id>
# Shorthand: cline t r <task-id>
```

#### Cancel Task

Cancel the currently running task:

```bash
./cli/bin/cline task cancel
# Shorthand: cline t c
```

#### Restore Checkpoint

Restore task to a specific checkpoint:

```bash
# Restore task state only
./cli/bin/cline task restore <checkpoint-id> --type task

# Restore workspace only
./cli/bin/cline task restore <checkpoint-id> --type workspace

# Restore both task and workspace
./cli/bin/cline task restore <checkpoint-id> --type taskAndWorkspace
```

**Note:** Checkpoint IDs are timestamps. Use `task list` or `task view` to see available checkpoints.

---

### Authentication

Sign in to Cline (opens browser for OAuth flow):

```bash
./cli/bin/cline auth
```

---

### OCA Provider

Manage Oracle Cloud AI (OCA) authentication:

```bash
# Login to OCA
./cli/bin/cline oca login

# Check authentication status
./cli/bin/cline oca status

# Watch for status changes
./cli/bin/cline oca status --watch

# Logout from OCA
./cli/bin/cline oca logout

# Create task with OCA provider
./cli/bin/cline task new "Write tests" \
  -s act-mode-api-provider=oca \
  -s act-mode-oca-model-id=your-model-id \
  -s oca-base-url=https://your-oca-endpoint
```

**OCA Models:**
```bash
# List available models (coming soon)
./cli/bin/cline oca models list

# Refresh model cache (coming soon)
./cli/bin/cline oca models refresh
```

---

### Workflow Management

Create and run Cline workflows:

```bash
# List available workflows
./cli/bin/cline workflow list

# Create a new workflow
./cli/bin/cline workflow create my-deploy

# Edit a workflow
./cli/bin/cline workflow edit my-deploy
# Opens: ~/Documents/Cline/Workflows/my-deploy.md

# Run a workflow
./cli/bin/cline workflow run my-deploy

# Run and follow progress
./cli/bin/cline workflow run my-deploy --wait

# Run with settings
./cli/bin/cline workflow run my-deploy \
  -w /projects/myapp \
  -s yolo-mode-toggled=true

# Delete a workflow
./cli/bin/cline workflow delete my-deploy

# Toggle workflow on/off
./cli/bin/cline workflow toggle my-deploy --enable
./cli/bin/cline workflow toggle my-deploy --disable
```

**Workflow Structure:**

Workflows are markdown files stored in `~/Documents/Cline/Workflows/`:

```markdown
# My Deployment Workflow

## Steps

1. Run tests
   ```bash
   npm test
   ```

2. Build the project
   ```bash
   npm run build
   ```

3. Deploy to production
   ```bash
   ./deploy.sh
   ```

## Notes
Remember to check environment variables before deploying.
```

**Invoking Workflows:**

When you run `cline workflow run my-deploy.md`, it creates a task with `/my-deploy.md` as the prompt, which tells Cline to follow the workflow instructions.

---

## Common Workflows

### Workflow 1: Quick One-off Task

Perfect for CI/CD pipelines or automated scripts:

```bash
./cli/bin/cline task oneshot "Run the test suite and fix any failing tests" \
  -w /path/to/project \
  -s max-iterations=10
```

### Workflow 2: Interactive Development Session

```bash
# Start a task
./cli/bin/cline task new "Implement user authentication" -w ./src

# In another terminal, follow progress
./cli/bin/cline task follow

# Send followup messages as needed
./cli/bin/cline send "Use JWT tokens instead of sessions"

# Approve or adjust as it works
./cli/bin/cline send --approve true
```

### Workflow 3: Multiple Projects

```bash
# Start instance for project A
./cli/bin/cline instance new
./cli/bin/cline task new "Refactor API endpoints" -w /projects/api

# Start instance for project B
./cli/bin/cline instance new
./cli/bin/cline task new "Update documentation" -w /projects/docs

# List all instances
./cli/bin/cline instance list

# Switch between them
./cli/bin/cline instance use localhost:50052  # Switch to project A
./cli/bin/cline task follow

./cli/bin/cline instance use localhost:50053  # Switch to project B
./cli/bin/cline task follow
```

### Workflow 4: Code Review with Context

```bash
# Submit code for review with full context
./cli/bin/cline task new "Review these changes for security issues" \
  -f auth/handler.go \
  -f auth/middleware.go \
  -f tests/auth_test.go \
  -m plan
```

### Workflow 5: Piping Input

```bash
# Pipe command output
git diff | ./cli/bin/cline task new "Explain these changes"

# Pipe from file
cat requirements.txt | ./cli/bin/cline task new "Update these dependencies"

# Chain commands
echo "Add comprehensive error handling" | ./cli/bin/cline send
```

## Global Flags

These flags work with any command:

- `--address string`: Specify Cline Core gRPC address (default: `localhost:50052`)
- `-o, --output-format string`: Output format - `rich`, `json`, or `plain` (default: `rich`)
- `-v, --verbose`: Enable verbose output for debugging

**Examples:**

```bash
# JSON output (useful for scripting)
./cli/bin/cline instance list -o json

# Plain text output
./cli/bin/cline task view -o plain

# Verbose mode
./cli/bin/cline task new "Debug this issue" -v

# Custom instance address
./cli/bin/cline task new "Test" --address localhost:50055
```

## Tips & Troubleshooting

### Tips

1. **Aliases**: Use the shorthand versions for faster typing:
   - `cline t n` = `cline task new`
   - `cline t f` = `cline task follow`
   - `cline i l` = `cline instance list`
   - `cline s` = `cline send`

2. **Stdin Support**: Most commands that accept text input also support stdin, enabling powerful shell pipelines.

3. **Multiple Instances**: Run separate instances for different projects or contexts. Each instance maintains its own task history and state.

4. **Yolo Mode**: Use `--yolo` or `-y` for fully autonomous execution without interruptions for approval.

5. **Task Settings**: Pass custom settings with `-s key=value`. Examples:
   - `-s model=claude-3-opus` - Use specific AI model
   - `-s max-iterations=20` - Set iteration limit
   - `-s aws-region=us-west-2` - Cloud provider settings

### Troubleshooting

#### No instances found

```bash
# Just create a new one
./cli/bin/cline instance new

# Or let task auto-create one
./cli/bin/cline task new "Your task"
```

#### Instance won't start

Check that the ports aren't already in use:
```bash
lsof -i :50052
lsof -i :50102
```

#### Task appears stuck

Try canceling and restarting:
```bash
./cli/bin/cline task cancel
./cli/bin/cline task new "Your task again"
```

#### Connection refused errors

Ensure the Cline instance is running:
```bash
./cli/bin/cline instance list
```

If no instances are listed, create one:
```bash
./cli/bin/cline instance new
```

#### Clean up all instances

```bash
./cli/bin/cline instance kill --all
```

## Integration with Shell

### Add to PATH

Add the CLI to your PATH for easy access from anywhere:

```bash
# In ~/.bashrc, ~/.zshrc, or equivalent:
export PATH="$PATH:/path/to/cline/cli/bin"

# Reload your shell
source ~/.bashrc  # or ~/.zshrc

# Then you can use from anywhere, but remember to cd to the project root first:
cd /path/to/cline
cline task new "Your task"
```

> **Note:** Even with the CLI in your PATH, you must run it from the Cline project root directory due to relative path dependencies.

### Bash/Zsh Completion

Generate shell completion script:

```bash
# Bash
./cli/bin/cline completion bash > /etc/bash_completion.d/cline

# Zsh
./cli/bin/cline completion zsh > "${fpath[1]}/_cline"

# Fish
./cli/bin/cline completion fish > ~/.config/fish/completions/cline.fish
```

## Architecture

The Cline CLI consists of two main components:

1. **cline**: The CLI tool that communicates with Cline Core via gRPC
2. **cline-host**: The host bridge server that provides filesystem and environment access

The CLI automatically manages starting and stopping host bridge servers as needed for each instance.

## Development

### Running Tests

```bash
cd cli
go test ./...

# Run e2e tests
go test ./e2e/... -v
```

### Building

```bash
# Build both binaries
cd cli
go build -o bin/cline ./cmd/cline
go build -o bin/cline-host ./cmd/cline-host
```

---

**Happy Coding with Cline! 🚀**

For more information, visit the [main Cline documentation](../README.md).

