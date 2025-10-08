# Cline CLI User Manual

**Complete Guide to Terminal-Based AI Coding with Cline**

## Table of Contents

1. [Getting Started](#getting-started)
2. [Starting Your First Task](#starting-your-first-task)
3. [Following Task Progress](#following-task-progress)
4. [Switching Providers & Models](#switching-providers--models)
5. [OCA Provider Authentication](#oca-provider-authentication)
6. [Working with Modes](#working-with-modes)
7. [Managing Multiple Instances](#managing-multiple-instances)
8. [Advanced Task Settings](#advanced-task-settings)
9. [Auto-Approval Configuration](#auto-approval-configuration)
10. [Working with Files & Images](#working-with-files--images)
11. [Task Lifecycle Management](#task-lifecycle-management)
12. [Workflow Management](#workflow-management)
13. [Real-World Workflows](#real-world-workflows)
14. [Complete Settings Reference](#complete-settings-reference)
15. [Troubleshooting](#troubleshooting)

---

## Getting Started

### Prerequisites

1. **Build the CLI** (from project root):
   ```bash
   cd /path/to/cline
   npm run compile-cli
   ```

2. **Verify installation**:
   ```bash
   ./cli/bin/cline version
   ```

3. **Important**: Always run commands from the Cline project root, not from the `cli/` directory.

### Authentication

Before using the CLI, authenticate with Cline:

```bash
./cli/bin/cline auth
```

This opens your browser for OAuth authentication.

---

## Starting Your First Task

### Basic Task

The simplest way to start:

```bash
./cli/bin/cline task new "Create a hello world function in Python"
```

This will:
- Automatically start a Cline instance if none exists
- Create a new task with your prompt
- Return immediately (task runs in background)

### Task with Wait

Wait for task completion before returning:

```bash
./cli/bin/cline task new "Run all tests" --wait
```

### One-Shot Task (Fully Autonomous)

For fire-and-forget execution (yolo + plan mode):

```bash
./cli/bin/cline task oneshot "Refactor the authentication module"
```

This is perfect for:
- CI/CD pipelines
- Automated scripts
- Tasks that don't need human interaction

---

## Following Task Progress

### Real-Time Following

Watch your task execute in real-time:

```bash
# In one terminal: create task
./cli/bin/cline task new "Implement user registration"

# In another terminal: follow progress
./cli/bin/cline task follow
```

### View Current State

See the current state without following:

```bash
./cli/bin/cline task view --current
```

### View Until Completion

Follow until the task completes, then exit:

```bash
./cli/bin/cline task view
```

### Get Summary Only

Just show the final summary:

```bash
./cli/bin/cline task view --summary
```

### Output Formats

Choose your preferred output format:

```bash
# Rich terminal UI (default)
./cli/bin/cline task follow

# Plain text
./cli/bin/cline task follow -o plain

# JSON (for parsing/scripting)
./cli/bin/cline task follow -o json
```

---

## Switching Providers & Models

The CLI supports 40+ AI providers! Configure them using task settings.

### Supported Providers

- **anthropic** - Claude models
- **openai** - GPT models
- **openrouter** - Multiple providers
- **bedrock** - AWS Bedrock
- **vertex** - Google Vertex AI
- **gemini** - Google Gemini
- **ollama** - Local models
- **lmstudio** - LM Studio
- **groq** - Groq
- **deepseek** - DeepSeek
- **xai** / **grok** - xAI Grok
- **together** - Together AI
- **fireworks** - Fireworks AI
- **cerebras** - Cerebras
- **sambanova** - SambaNova
- **mistral** - Mistral AI
- **litellm** - LiteLLM proxy
- **requesty** - Requesty
- **asksage** - AskSage
- And many more...

### Switch Provider for Act Mode

```bash
# Use OpenAI for act mode
./cli/bin/cline task new "Write unit tests" \
  -s act-mode-api-provider=openai \
  -s act-mode-api-model-id=gpt-4

# Use Anthropic (Claude)
./cli/bin/cline task new "Code review" \
  -s act-mode-api-provider=anthropic \
  -s act-mode-api-model-id=claude-3-5-sonnet-20241022

# Use local Ollama
./cli/bin/cline task new "Explain this code" \
  -s act-mode-api-provider=ollama \
  -s act-mode-ollama-model-id=llama3.2
```

### Switch Provider for Plan Mode

```bash
# Use different models for plan and act
./cli/bin/cline task new "Complex refactoring" \
  -s plan-act-separate-models-setting=true \
  -s plan-mode-api-provider=anthropic \
  -s plan-mode-api-model-id=claude-3-5-sonnet-20241022 \
  -s act-mode-api-provider=openai \
  -s act-mode-api-model-id=gpt-4
```

### Provider-Specific Configuration

#### AWS Bedrock

```bash
./cli/bin/cline task new "Deploy to AWS" \
  -s act-mode-api-provider=bedrock \
  -s aws-region=us-east-1 \
  -s aws-profile=my-profile \
  -s aws-use-profile=true
```

#### Google Vertex AI

```bash
./cli/bin/cline task new "Analyze data" \
  -s act-mode-api-provider=vertex \
  -s vertex-project-id=my-project \
  -s vertex-region=us-central1
```

#### Ollama (Local)

```bash
./cli/bin/cline task new "Quick task" \
  -s act-mode-api-provider=ollama \
  -s ollama-base-url=http://localhost:11434 \
  -s act-mode-ollama-model-id=llama3.2
```

#### OpenRouter

```bash
./cli/bin/cline task new "Use multiple models" \
  -s act-mode-api-provider=openrouter \
  -s act-mode-open-router-model-id=anthropic/claude-3-opus
```

#### Azure OpenAI

```bash
./cli/bin/cline task new "Enterprise task" \
  -s act-mode-api-provider=openai \
  -s azure-api-version=2024-02-15-preview
```

### Model-Specific Settings

#### OpenAI Reasoning Effort

```bash
./cli/bin/cline task new "Complex problem" \
  -s openai-reasoning-effort=high \
  -s act-mode-reasoning-effort=high
```

Options: `low`, `medium`, `high`

#### Thinking Budget Tokens

For models that support extended thinking:

```bash
./cli/bin/cline task new "Deep analysis" \
  -s plan-mode-thinking-budget-tokens=10000 \
  -s act-mode-thinking-budget-tokens=5000
```

---

## OCA Provider Authentication

Oracle Cloud AI (OCA) requires OAuth authentication. The CLI provides commands to manage OCA authentication flow.

### Login to OCA

```bash
# Initiate OAuth login flow
./cli/bin/cline oca login
```

This will:
1. Start the OAuth authentication flow
2. Open your default browser
3. Prompt you to authenticate with Oracle Cloud
4. Store authentication tokens securely

### Check Authentication Status

```bash
# One-time status check
./cli/bin/cline oca status
```

Output example:
```
Status: Authenticated ✓
User ID: user@oracle.com
Name: John Doe
Email: john.doe@oracle.com
Token: a1b2c3d4...x9y8z7w6
```

### Watch Authentication Status

Monitor authentication status in real-time:

```bash
# Stream status updates
./cli/bin/cline oca status --watch
```

This is useful when:
- Waiting for authentication to complete
- Monitoring token refresh
- Debugging authentication issues

### Logout from OCA

```bash
# Clear authentication tokens
./cli/bin/cline oca logout
```

### Using OCA Provider

Once authenticated, create tasks with OCA provider:

```bash
# Basic OCA task
./cli/bin/cline task new "Write unit tests" \
  -s act-mode-api-provider=oca \
  -s act-mode-oca-model-id=your-model-id

# With custom OCA endpoint
./cli/bin/cline task new "Code review" \
  -s act-mode-api-provider=oca \
  -s act-mode-oca-model-id=your-model \
  -s oca-base-url=https://your-oca-instance.com

# With separate plan/act models
./cli/bin/cline task new "Complex refactoring" \
  -s plan-act-separate-models-setting=true \
  -s plan-mode-api-provider=oca \
  -s plan-mode-oca-model-id=planning-model \
  -s act-mode-api-provider=oca \
  -s act-mode-oca-model-id=execution-model
```

### OCA Models Management

```bash
# List available models (coming soon)
./cli/bin/cline oca models list

# Refresh model cache (coming soon)
./cli/bin/cline oca models refresh
```

> **Note:** Model listing from CLI is planned for future release. Currently, use the VSCode extension UI to view available OCA models.

### OCA with Multiple Instances

Each Cline instance maintains its own OCA authentication:

```bash
# Instance 1 with OCA
./cli/bin/cline instance new
./cli/bin/cline oca login --address localhost:50052

# Instance 2 with different OCA account
./cli/bin/cline instance new
./cli/bin/cline oca login --address localhost:50053

# Check status of specific instance
./cli/bin/cline oca status --address localhost:50052
```

### OCA Authentication Flow

```
User: cline oca login
    │
    ▼
CLI calls OcaAccountService.ocaAccountLoginClicked()
    │
    ▼
Browser opens with OAuth URL
    │
    ▼
User completes authentication
    │
    ▼
Tokens stored in cline-core
    │
    ▼
User: cline oca status
    │
    ▼
CLI shows authenticated status ✓
```

### Troubleshooting OCA

**Not authenticated error:**
```bash
# Check status first
./cli/bin/cline oca status

# If not authenticated, login
./cli/bin/cline oca login
```

**Token expired:**
```bash
# Logout and login again
./cli/bin/cline oca logout
./cli/bin/cline oca login
```

**Multiple instances:**
```bash
# Make sure you're targeting the right instance
./cli/bin/cline instance list
./cli/bin/cline oca status --address localhost:XXXXX
```

---

## Working with Modes

Cline has two primary modes:

### Act Mode (Interactive)

- Asks for approval before actions
- More control over what gets executed
- Default mode

```bash
./cli/bin/cline task new "Update dependencies" -m act
```

### Plan Mode (Autonomous)

- Creates plans and executes them
- More autonomous operation
- Better for well-defined tasks

```bash
./cli/bin/cline task new "Refactor database layer" -m plan
```

### Switch Mode During Task

```bash
# Start in act mode
./cli/bin/cline task new "Complex task" -m act

# Switch to plan mode mid-task
./cli/bin/cline send --mode plan

# Switch back to act
./cli/bin/cline send --mode act
```

### Yolo Mode (Non-Interactive)

Disable all approval prompts:

```bash
# With -y flag
./cli/bin/cline task new "Run tests" -y

# Or via setting
./cli/bin/cline task new "Deploy" \
  -s yolo-mode-toggled=true
```

### Strict Plan Mode

More rigorous planning:

```bash
./cli/bin/cline task new "Architecture redesign" \
  -m plan \
  -s strict-plan-mode-enabled=true
```

---

## Managing Multiple Instances

Run multiple Cline instances simultaneously, similar to kubectl contexts.

### List Instances

```bash
./cli/bin/cline instance list
```

Output:
```
ADDRESS            STATUS   VERSION  LAST SEEN  PID    DEFAULT
localhost:50052    SERVING  dev      14:23:45   12345  *
localhost:50053    SERVING  dev      14:24:10   12346
```

### Create New Instance

```bash
./cli/bin/cline instance new
```

Output:
```
Successfully started new instance:
  Address: localhost:50053
  Core Port: 50053
  Host Bridge Port: 50103
```

### Switch Default Instance

```bash
./cli/bin/cline instance use localhost:50053
```

All subsequent commands use this instance unless you specify `--address`.

### Use Specific Instance for Task

```bash
# Create task on specific instance
./cli/bin/cline task new "Task A" --address localhost:50052

# Create task on different instance
./cli/bin/cline task new "Task B" --address localhost:50053

# Follow task on specific instance
./cli/bin/cline task follow --address localhost:50052
```

### Kill Instance

```bash
# Kill specific instance
./cli/bin/cline instance kill localhost:50053

# Kill all instances
./cli/bin/cline instance kill --all
```

### Use Cases for Multiple Instances

1. **Different Projects**:
   ```bash
   # Instance 1 for API project
   ./cli/bin/cline instance new
   ./cli/bin/cline task new "Fix API bugs" -w /projects/api
   
   # Instance 2 for frontend
   ./cli/bin/cline instance new
   ./cli/bin/cline task new "Update UI" -w /projects/frontend
   ```

2. **Different Providers**:
   ```bash
   # Instance 1 with Claude
   ./cli/bin/cline task new "Task" --address localhost:50052 \
     -s act-mode-api-provider=anthropic
   
   # Instance 2 with GPT-4
   ./cli/bin/cline task new "Task" --address localhost:50053 \
     -s act-mode-api-provider=openai
   ```

3. **Parallel Testing**:
   ```bash
   # Test different approaches simultaneously
   ./cli/bin/cline task new "Approach A" --address localhost:50052
   ./cli/bin/cline task new "Approach B" --address localhost:50053
   ```

---

## Advanced Task Settings

### Working Directory

Specify the workspace directory:

```bash
./cli/bin/cline task new "Initialize project" \
  -w /path/to/project
```

Multiple directories:

```bash
./cli/bin/cline task new "Sync repos" \
  -w /path/to/repo1 \
  -w /path/to/repo2
```

### Custom Instructions

Add custom instructions to guide the AI:

```bash
./cli/bin/cline task new "Write API endpoints" \
  -s custom-prompt="Always use TypeScript. Follow REST conventions. Include error handling."
```

### Preferred Language

Set output language:

```bash
./cli/bin/cline task new "Explique ce code" \
  -s preferred-language=French
```

### Request Timeout

Set custom timeout (in milliseconds):

```bash
./cli/bin/cline task new "Long running task" \
  -s request-timeout-ms=120000
```

### Terminal Settings

```bash
./cli/bin/cline task new "Run commands" \
  -s default-terminal-profile=bash \
  -s shell-integration-timeout=5000 \
  -s terminal-output-line-limit=1000
```

### Checkpoints

Enable/disable task checkpoints:

```bash
./cli/bin/cline task new "Important work" \
  -s enable-checkpoints-setting=true
```

### Auto-Condense

Automatically condense conversation history:

```bash
./cli/bin/cline task new "Long conversation" \
  -s use-auto-condense=true \
  -s auto-condense-threshold=0.8
```

---

## Auto-Approval Configuration

Configure what actions can run without approval.

### Enable Auto-Approval

```bash
./cli/bin/cline task new "Automated task" \
  -s auto-approval-settings.enabled=true \
  -s auto-approval-settings.max-requests=10
```

### Approve Specific Actions

```bash
./cli/bin/cline task new "Safe operations" \
  -s auto-approval-settings.enabled=true \
  -s auto-approval-settings.actions.read-files=true \
  -s auto-approval-settings.actions.execute-safe-commands=true
```

### Full Auto-Approval (Dangerous!)

```bash
./cli/bin/cline task new "Fully autonomous" \
  -s auto-approval-settings.enabled=true \
  -s auto-approval-settings.actions.read-files=true \
  -s auto-approval-settings.actions.edit-files=true \
  -s auto-approval-settings.actions.execute-all-commands=true \
  -s auto-approval-settings.actions.use-browser=true \
  -s auto-approval-settings.actions.use-mcp=true
```

### Available Auto-Approval Actions

- `read-files` - Read files from disk
- `read-files-externally` - Read files outside workspace
- `edit-files` - Edit/create files
- `edit-files-externally` - Edit files outside workspace
- `execute-safe-commands` - Run safe terminal commands
- `execute-all-commands` - Run any terminal command
- `use-browser` - Use browser automation
- `use-mcp` - Use MCP servers

---

## Working with Files & Images

### Attach Files to Task

```bash
# Single file
./cli/bin/cline task new "Review this code" \
  -f src/auth.go

# Multiple files
./cli/bin/cline task new "Analyze these modules" \
  -f src/auth.go \
  -f src/user.go \
  -f src/middleware.go
```

### Attach Images

Perfect for UI/design tasks:

```bash
./cli/bin/cline task new "Implement this design" \
  -i mockup.png \
  -i screenshot.jpg
```

### Attach Files to Followup Message

```bash
./cli/bin/cline send "Here's the updated code" \
  -f updated.go
```

### Pipe File Content

```bash
# Pipe file content as prompt
cat error.log | ./cli/bin/cline task new "Debug this error"

# Pipe to followup message
cat changes.diff | ./cli/bin/cline send "Review these changes"
```

### Complex Example

```bash
./cli/bin/cline task new "Complete code review" \
  -f src/api/*.go \
  -f tests/api_test.go \
  -i architecture-diagram.png \
  -w /projects/myapp \
  -s custom-prompt="Focus on security and performance" \
  -m plan
```

---

## Task Lifecycle Management

### Create Task

```bash
./cli/bin/cline task new "Implement feature X"
```

### List Task History

```bash
./cli/bin/cline task list
```

### Resume Previous Task

```bash
./cli/bin/cline task resume <task-id>
```

### Cancel Current Task

```bash
./cli/bin/cline task cancel
```

### Send Followup Messages

```bash
# Text message
./cli/bin/cline send "Please add error handling"

# Approve pending request
./cli/bin/cline send --approve true

# Deny pending request
./cli/bin/cline send --approve false
```

### Restore from Checkpoint

```bash
# List checkpoints (via task view)
./cli/bin/cline task view

# Restore task state only
./cli/bin/cline task restore 1699564821000 --type task

# Restore workspace only
./cli/bin/cline task restore 1699564821000 --type workspace

# Restore both
./cli/bin/cline task restore 1699564821000 --type taskAndWorkspace
```

---

## Workflow Management

Workflows are reusable step-by-step instructions stored as markdown files. They allow you to define common processes like deployments, code reviews, or testing procedures.

### What Are Workflows?

Workflows are markdown files that contain:
- Step-by-step instructions for Cline
- Commands to execute
- Questions to ask
- Processes to follow

When invoked via `/workflow-name.md`, Cline processes the workflow as explicit instructions.

### List Available Workflows

```bash
# List all workflows
./cli/bin/cline workflow list
```

Output:
```
Available workflows (3):

1. deploy-prod.md
2. pr-review.md
3. setup-project.md

Run a workflow with: cline workflow run <name>
```

### Create a New Workflow

```bash
# Create workflow
./cli/bin/cline workflow create my-workflow

# With edit flag (shows editor command)
./cli/bin/cline workflow create my-workflow --edit
```

This creates: `~/Documents/Cline/Workflows/my-workflow.md`

**Template content:**
```markdown
# My Workflow

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
```

### Edit a Workflow

```bash
# Edit workflow (shows path and editor command)
./cli/bin/cline workflow edit my-workflow
```

Output:
```
Edit the file: /Users/you/Documents/Cline/Workflows/my-workflow.md
(Open with: vi /Users/you/Documents/Cline/Workflows/my-workflow.md)
```

Then edit with your preferred editor:
```bash
$EDITOR ~/Documents/Cline/Workflows/my-workflow.md
```

### Run a Workflow

```bash
# Basic workflow execution
./cli/bin/cline workflow run my-workflow

# Run and follow progress
./cli/bin/cline workflow run my-workflow --wait

# Run with working directory
./cli/bin/cline workflow run deploy-prod \
  -w /projects/myapp

# Run with settings
./cli/bin/cline workflow run deploy-prod \
  -s yolo-mode-toggled=true \
  -s auto-approval-settings.actions.execute-all-commands=true

# Run in plan mode
./cli/bin/cline workflow run code-review -m plan

# Run in yolo mode
./cli/bin/cline workflow run quick-test -y
```

### Delete a Workflow

```bash
# Delete with confirmation
./cli/bin/cline workflow delete my-workflow

# Force delete without confirmation
./cli/bin/cline workflow delete my-workflow --force
```

### Toggle Workflow

```bash
# Enable workflow
./cli/bin/cline workflow toggle my-workflow --enable

# Disable workflow
./cli/bin/cline workflow toggle my-workflow --disable

# Toggle global workflow
./cli/bin/cline workflow toggle my-workflow --enable --global
```

### Workflow Examples

#### Example 1: Deployment Workflow

Create `~/Documents/Cline/Workflows/deploy-prod.md`:

```markdown
# Production Deployment

## Steps

1. Verify we're on main branch
   ```bash
   git branch --show-current
   ```

2. Run all tests
   ```bash
   npm test
   ```

3. Build the project
   ```bash
   npm run build
   ```

4. Tag the release
   ```bash
   git tag -a v$(date +%Y%m%d-%H%M%S) -m "Production deployment"
   ```

5. Deploy to production
   ```bash
   ./deploy.sh production
   ```

6. Verify deployment
   - Check health endpoint: https://api.example.com/health
   - Review logs for errors

## Notes
- Ensure AWS credentials are set
- Notify team in Slack after deployment
```

Run it:
```bash
./cli/bin/cline workflow run deploy-prod \
  -w /projects/api \
  --wait
```

#### Example 2: PR Review Workflow

Create `~/Documents/Cline/Workflows/pr-review.md`:

```markdown
# Pull Request Review

## Description
Comprehensive PR review using GitHub CLI.

## Steps

1. Get PR information
   ```bash
   gh pr view $PR_NUMBER --json title,body,comments
   ```

2. Get the diff
   ```bash
   gh pr diff $PR_NUMBER
   ```

3. Check which files were modified
   ```bash
   gh pr view $PR_NUMBER --json files
   ```

4. Review each modified file for:
   - Code quality
   - Security issues
   - Performance concerns
   - Test coverage

5. Leave review comment
   ```bash
   gh pr review $PR_NUMBER --comment
   ```

## Notes
Ask the user for PR_NUMBER if not provided.
```

Run it:
```bash
# Cline will ask for PR number
./cli/bin/cline workflow run pr-review -w /projects/repo
```

#### Example 3: Project Setup Workflow

Create `~/Documents/Cline/Workflows/setup-node-project.md`:

```markdown
# Node.js Project Setup

## Steps

1. Initialize npm project
   ```bash
   npm init -y
   ```

2. Install dependencies
   ```bash
   npm install express dotenv
   npm install --save-dev typescript @types/node @types/express
   ```

3. Create project structure
   ```bash
   mkdir -p src/{routes,controllers,services,models}
   mkdir -p tests
   ```

4. Create tsconfig.json
   ```json
   {
     "compilerOptions": {
       "target": "ES2020",
       "module": "commonjs",
       "outDir": "./dist",
       "rootDir": "./src",
       "strict": true,
       "esModuleInterop": true
     }
   }
   ```

5. Create .gitignore
   ```
   node_modules/
   dist/
   .env
   ```

6. Initialize git
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   ```

## Notes
Creates a basic Node.js/TypeScript project structure.
```

Run it:
```bash
./cli/bin/cline workflow run setup-node-project \
  -w /projects/new-api \
  -y
```

### Workflow Best Practices

1. **Be Specific**: Include exact commands and paths
2. **Add Context**: Explain why each step is needed
3. **Handle Errors**: Mention what to do if steps fail
4. **Ask Questions**: Use `ask_followup_question` tool for dynamic input
5. **Use Tools**: Leverage Cline's built-in tools (read_file, search_files, etc.)
6. **Test First**: Test workflows before using in production
7. **Version Control**: Consider versioning your workflows directory

### Workflow with MCP Tools

Workflows can use MCP tools:

```markdown
# Slack Notification Workflow

## Steps

1. Get deployment status
   ```bash
   ./check-deployment.sh
   ```

2. Send notification to Slack
   Use the Slack MCP tool to send a message to #deployments channel:
   - Message: "Production deployment completed successfully"
   - Include deployment timestamp

3. Create Jira ticket for post-deployment verification
   Use Jira MCP tool to create ticket
```

### Workflow with New Task Tool

Chain workflows together:

```markdown
# Master CI/CD Workflow

## Steps

1. Run tests workflow
   Use `new_task` tool to invoke `/run-tests.md`

2. If tests pass, run build workflow
   Use `new_task` tool to invoke `/build-project.md`

3. If build succeeds, run deployment workflow
   Use `new_task` tool to invoke `/deploy-prod.md`
```

### Workflow Locations

- **Local workflows**: `~/Documents/Cline/Workflows/`
- **Global workflows**: Shared across projects (location varies by OS)

### How Workflows Work Internally

1. You run: `cline workflow run my-workflow.md`
2. CLI creates task with prompt: `/my-workflow.md`
3. Cline-core detects slash command
4. Reads workflow file from disk
5. Injects content as `<explicit_instructions>`
6. Task executes following workflow steps

### Troubleshooting Workflows

**Workflow not found:**
```bash
# Check workflows directory
ls ~/Documents/Cline/Workflows/

# List available workflows
./cli/bin/cline workflow list
```

**Workflow directory doesn't exist:**
```bash
# Create it
mkdir -p ~/Documents/Cline/Workflows
```

**Workflow not executing properly:**
- Check markdown syntax
- Ensure commands are in code blocks
- Test commands individually first
- Use `--wait` flag to follow execution

---

## Real-World Workflows

### Workflow 1: Quick Bug Fix

```bash
# Start task
./cli/bin/cline task new "Fix login bug in auth.go" \
  -f src/auth.go \
  -f tests/auth_test.go \
  -m act

# Follow progress
./cli/bin/cline task follow

# Send clarification
./cli/bin/cline send "Use bcrypt for password hashing"

# Approve changes
./cli/bin/cline send --approve true
```

### Workflow 2: Full Feature Implementation

```bash
# Create autonomous task
./cli/bin/cline task oneshot "Implement user profile API endpoint with tests" \
  -w /projects/api \
  -s plan-mode-api-provider=anthropic \
  -s plan-mode-api-model-id=claude-3-5-sonnet-20241022 \
  -s auto-approval-settings.enabled=true \
  -s auto-approval-settings.actions.read-files=true \
  -s auto-approval-settings.actions.edit-files=true \
  -s auto-approval-settings.actions.execute-safe-commands=true
```

### Workflow 3: Code Review

```bash
# Review with context
git diff main..feature-branch > changes.diff

./cli/bin/cline task new "Review these changes for security issues" \
  -f changes.diff \
  -s custom-prompt="Focus on SQL injection, XSS, and authentication issues" \
  -m plan
```

### Workflow 4: Multi-Instance Development

```bash
# Terminal 1: Backend work
./cli/bin/cline instance new
./cli/bin/cline task new "Implement REST API" \
  -w /projects/backend \
  -s act-mode-api-provider=anthropic \
  --wait

# Terminal 2: Frontend work
./cli/bin/cline instance new
./cli/bin/cline task new "Create user interface" \
  -w /projects/frontend \
  -s act-mode-api-provider=openai \
  --wait

# Terminal 3: Monitor both
watch -n 5 './cli/bin/cline instance list'
```

### Workflow 5: Continuous Integration

```bash
#!/bin/bash
# ci-script.sh

# Run tests and fix failures
./cli/bin/cline task oneshot "Run npm test and fix any failures" \
  -w /projects/app \
  -s yolo-mode-toggled=true \
  -s auto-approval-settings.enabled=true \
  -s auto-approval-settings.actions.read-files=true \
  -s auto-approval-settings.actions.edit-files=true \
  -s auto-approval-settings.actions.execute-all-commands=true

# Get exit code
if [ $? -eq 0 ]; then
  echo "Tests passed!"
else
  echo "Tests failed!"
  exit 1
fi
```

### Workflow 6: Documentation Generation

```bash
# Generate docs for all Go files
find ./src -name "*.go" -type f | while read file; do
  echo "Processing $file"
done

./cli/bin/cline task oneshot "Generate comprehensive documentation" \
  -f $(find ./src -name "*.go") \
  -s custom-prompt="Generate godoc comments for all exported functions and types" \
  -m plan
```

### Workflow 7: Migration Tasks

```bash
# Complex migration
./cli/bin/cline task new "Migrate from REST to GraphQL" \
  -w /projects/api \
  -s plan-act-separate-models-setting=true \
  -s plan-mode-api-provider=anthropic \
  -s plan-mode-api-model-id=claude-3-5-sonnet-20241022 \
  -s act-mode-api-provider=openai \
  -s act-mode-api-model-id=gpt-4 \
  -s enable-checkpoints-setting=true \
  -m plan

# Follow and checkpoint regularly
./cli/bin/cline task follow
```

---

## Complete Settings Reference

### Provider Settings

| Setting | Values | Description |
|---------|--------|-------------|
| `act-mode-api-provider` | See [providers](#supported-providers) | Act mode provider |
| `plan-mode-api-provider` | See [providers](#supported-providers) | Plan mode provider |
| `plan-act-separate-models-setting` | `true`/`false` | Use different models for plan/act |

### Model Settings

| Setting | Description |
|---------|-------------|
| `act-mode-api-model-id` | Model ID for act mode |
| `plan-mode-api-model-id` | Model ID for plan mode |
| `act-mode-reasoning-effort` | `low`/`medium`/`high` |
| `plan-mode-reasoning-effort` | `low`/`medium`/`high` |
| `act-mode-thinking-budget-tokens` | Thinking token budget |
| `plan-mode-thinking-budget-tokens` | Thinking token budget |

### AWS Bedrock Settings

| Setting | Description |
|---------|-------------|
| `aws-region` | AWS region (e.g., `us-east-1`) |
| `aws-bedrock-endpoint` | Custom endpoint URL |
| `aws-profile` | AWS profile name |
| `aws-authentication` | Authentication method |
| `aws-use-profile` | Use AWS profile (`true`/`false`) |
| `aws-use-cross-region-inference` | Enable cross-region (`true`/`false`) |
| `aws-bedrock-use-prompt-cache` | Use prompt caching (`true`/`false`) |

### Google Vertex/Gemini Settings

| Setting | Description |
|---------|-------------|
| `vertex-project-id` | GCP project ID |
| `vertex-region` | GCP region |
| `gemini-base-url` | Custom Gemini endpoint |

### OpenAI Settings

| Setting | Description |
|---------|-------------|
| `open-ai-base-url` | Custom OpenAI endpoint |
| `openai-reasoning-effort` | `low`/`medium`/`high` |
| `azure-api-version` | Azure API version |

### Ollama Settings

| Setting | Description |
|---------|-------------|
| `ollama-base-url` | Ollama server URL |
| `ollama-api-options-ctx-num` | Context window size |
| `act-mode-ollama-model-id` | Model name |
| `plan-mode-ollama-model-id` | Model name |

### LM Studio Settings

| Setting | Description |
|---------|-------------|
| `lm-studio-base-url` | LM Studio server URL |
| `lm-studio-max-tokens` | Max output tokens |

### Other Provider Settings

| Setting | Provider | Description |
|---------|----------|-------------|
| `anthropic-base-url` | Anthropic | Custom endpoint |
| `open-router-provider-sorting` | OpenRouter | Provider sorting |
| `lite-llm-base-url` | LiteLLM | Proxy URL |
| `requesty-base-url` | Requesty | Server URL |
| `fireworks-model-max-completion-tokens` | Fireworks | Max tokens |
| `groq-model-id` | Groq | Model ID |
| `sap-ai-core-base-url` | SAP AI Core | API URL |
| `sap-ai-core-token-url` | SAP AI Core | Token URL |

### Mode Settings

| Setting | Values | Description |
|---------|--------|-------------|
| `mode` | `act`/`plan` | Initial mode |
| `yolo-mode-toggled` | `true`/`false` | Non-interactive mode |
| `strict-plan-mode-enabled` | `true`/`false` | Strict planning |

### Task Behavior Settings

| Setting | Description |
|---------|-------------|
| `custom-prompt` | Custom instructions |
| `preferred-language` | Output language |
| `request-timeout-ms` | Timeout in milliseconds |
| `use-auto-condense` | Auto-condense conversation |
| `auto-condense-threshold` | Condense threshold (0-1) |
| `enable-checkpoints-setting` | Enable checkpoints |

### Terminal Settings

| Setting | Description |
|---------|-------------|
| `default-terminal-profile` | Terminal profile name |
| `shell-integration-timeout` | Shell timeout (ms) |
| `terminal-output-line-limit` | Output line limit |

### Auto-Approval Settings

| Setting | Description |
|---------|-------------|
| `auto-approval-settings.enabled` | Enable auto-approval |
| `auto-approval-settings.max-requests` | Max auto-approved requests |
| `auto-approval-settings.enable-notifications` | Show notifications |
| `auto-approval-settings.actions.read-files` | Auto-approve file reads |
| `auto-approval-settings.actions.read-files-externally` | Auto-approve external reads |
| `auto-approval-settings.actions.edit-files` | Auto-approve file edits |
| `auto-approval-settings.actions.edit-files-externally` | Auto-approve external edits |
| `auto-approval-settings.actions.execute-safe-commands` | Auto-approve safe commands |
| `auto-approval-settings.actions.execute-all-commands` | Auto-approve all commands |
| `auto-approval-settings.actions.use-browser` | Auto-approve browser usage |
| `auto-approval-settings.actions.use-mcp` | Auto-approve MCP usage |

### Browser Settings

| Setting | Description |
|---------|-------------|
| `browser-settings.viewport-width` | Browser viewport width |
| `browser-settings.viewport-height` | Browser viewport height |

---

## Troubleshooting

### Instance Won't Start

**Problem**: `failed to start cline-host: no such file or directory`

**Solution**: Run from project root:
```bash
cd /path/to/cline
./cli/bin/cline instance list
```

### Connection Refused

**Problem**: `connection refused`

**Solution**: Check if instance is running:
```bash
./cli/bin/cline instance list

# If none, create one:
./cli/bin/cline instance new
```

### Task Appears Stuck

**Solution**: Cancel and restart:
```bash
./cli/bin/cline task cancel
./cli/bin/cline task new "Your task again"
```

### Port Already in Use

**Problem**: Can't start instance on default port

**Solution**: Let CLI auto-assign ports:
```bash
./cli/bin/cline instance new
# Uses next available ports automatically
```

### Clean Slate

Reset everything:
```bash
# Kill all instances
./cli/bin/cline instance kill --all

# Start fresh
./cli/bin/cline instance new
```

### Verbose Debugging

Enable verbose output:
```bash
./cli/bin/cline -v task new "Debug task"
./cli/bin/cline -v task follow
```

### Check Instance Health

```bash
./cli/bin/cline instance list
```

Look for `SERVING` status. If status is not `SERVING`, kill and restart:
```bash
./cli/bin/cline instance kill <address>
./cli/bin/cline instance new
```

---

## Tips & Best Practices

### 1. Use Plan Mode for Complex Tasks

```bash
./cli/bin/cline task new "Refactor entire module" -m plan
```

### 2. Enable Checkpoints for Important Work

```bash
./cli/bin/cline task new "Critical migration" \
  -s enable-checkpoints-setting=true
```

### 3. Use Different Models for Different Tasks

```bash
# Fast model for simple tasks
./cli/bin/cline task oneshot "Update README" \
  -s act-mode-api-provider=openai \
  -s act-mode-api-model-id=gpt-3.5-turbo

# Powerful model for complex tasks
./cli/bin/cline task oneshot "Design system architecture" \
  -s act-mode-api-provider=anthropic \
  -s act-mode-api-model-id=claude-3-opus-20240229
```

### 4. Pipe Commands for Context

```bash
git log -n 20 --oneline | ./cli/bin/cline task new "Summarize recent changes"
```

### 5. Use Multiple Instances Wisely

- One instance per project
- Different providers for comparison
- Parallel exploration of solutions

### 6. Save Common Settings

Create shell aliases:
```bash
# In ~/.bashrc or ~/.zshrc
alias cline-anthropic='cline task new -s act-mode-api-provider=anthropic -s act-mode-api-model-id=claude-3-5-sonnet-20241022'
alias cline-oneshot='cline task oneshot -s yolo-mode-toggled=true'
```

### 7. Use JSON Output for Scripts

```bash
INSTANCES=$(./cli/bin/cline instance list -o json)
echo "$INSTANCES" | jq '.[] | select(.status == "SERVING")'
```

---

## Quick Command Reference

```bash
# Version
cline version

# Auth
cline auth

# Instances
cline instance list
cline instance new
cline instance use <address>
cline instance kill <address>
cline instance kill --all

# Tasks
cline task new "prompt"
cline task new "prompt" --wait
cline task oneshot "prompt"
cline task list
cline task follow
cline task view
cline task cancel
cline task resume <id>
cline task restore <checkpoint-id> --type task

# Send
cline send "message"
cline send --mode plan
cline send --approve true

# With settings
cline task new "prompt" -s key=value
cline task new "prompt" -m plan -y -w /path
```

---

## Getting Help

```bash
# General help
./cli/bin/cline --help

# Command-specific help
./cli/bin/cline task --help
./cli/bin/cline task new --help
./cli/bin/cline instance --help

# Examples script
./cli/examples.sh
```

---

**Happy Coding with Cline! 🚀**

For additional resources:
- [README.md](README.md) - Installation and overview
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Cheat sheet
- [Main Documentation](../README.md) - Full Cline documentation

