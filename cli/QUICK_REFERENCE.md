# Cline CLI Quick Reference

## Prerequisites
- Run from Cline project root: `cd /path/to/cline`
- Binaries: `./cli/bin/cline` and `./cli/bin/cline-host`

## Essential Commands

### Version & Help
```bash
cline version                    # Show version info
cline version --short            # Version number only
cline --help                     # List all commands
cline [command] --help           # Help for specific command
```

### Authentication
```bash
cline auth                       # Sign in to Cline
```

### OCA Provider
```bash
cline oca login                  # Login to OCA
cline oca logout                 # Logout from OCA
cline oca status                 # Check auth status
cline oca status --watch         # Watch status changes
cline oca models list            # List models (coming soon)
```

### Workflow Management
```bash
cline workflow list              # List workflows (wf l)
cline workflow create <name>     # Create workflow (wf create)
cline workflow run <name>        # Run workflow (wf r)
cline workflow run <name> --wait # Run and follow
cline workflow edit <name>       # Edit workflow (wf e)
cline workflow delete <name>     # Delete workflow (wf delete)
cline workflow toggle <name> --enable   # Enable workflow
cline workflow toggle <name> --disable  # Disable workflow
```

### Instance Management
```bash
cline instance list              # List all instances (i l)
cline instance new               # Start new instance (i n)
cline instance use <address>     # Set default instance (i u)
cline instance kill <address>    # Kill instance (i k)
cline instance kill --all        # Kill all instances
```

### Task Operations
```bash
cline task new "prompt"          # Create task (t n)
cline task new "prompt" --wait   # Create and wait for completion
cline task new "prompt" -y       # Create in yolo mode (non-interactive)
cline task new "prompt" -m plan  # Create in plan mode
cline task oneshot "prompt"      # Create, run autonomous, wait (t o)

cline task list                  # List task history (t l)
cline task follow                # Follow current task (t f)
cline task view                  # View until completion (t v)
cline task view --current        # View current state
cline task view --summary        # Show summary only

cline task cancel                # Cancel current task (t c)
cline task resume <task-id>      # Resume task (t r)
```

### Send Messages
```bash
cline send "message"             # Send followup (s)
cline send --mode plan           # Change to plan mode
cline send --mode act            # Change to act mode
cline send --approve true        # Approve pending request
cline send --approve false       # Deny pending request
cline send "msg" -f file.txt     # Send with file attachment
echo "msg" | cline send          # Send from stdin
```

### File Attachments
```bash
cline task new "prompt" -f file1.txt -f file2.go   # Attach files
cline task new "prompt" -i image.png               # Attach images
cline task new "prompt" -w /path/to/workdir        # Set working directory
```

### Settings & Modes
```bash
-s key=value                     # Custom settings
-s model=claude-3-opus           # Specify AI model
-s max-iterations=20             # Set iteration limit
-s aws-region=us-west-2          # Cloud provider settings
-m plan                          # Plan mode (autonomous)
-m act                           # Act mode (interactive)
-y                               # Yolo mode (non-interactive)
```

## Global Flags

```bash
--address <address>              # Target specific instance
-o json                          # JSON output format
-o plain                         # Plain text output
-o rich                          # Rich terminal output (default)
-v                               # Verbose output
```

## Common Workflows

### Quick Task
```bash
cline task oneshot "Run tests and fix failures"
```

### Interactive Session
```bash
cline task new "Add user authentication"
cline task follow                # In another terminal
cline send "Use JWT tokens"      # Send followup
cline send --approve true        # Approve actions
```

### Multiple Projects
```bash
cline instance new               # Start for project 1
cline instance new               # Start for project 2
cline instance list              # See both
cline instance use localhost:50053  # Switch to project 2
```

### Piping Content
```bash
git diff | cline task new "Review these changes"
cat error.log | cline send "Debug this error"
```

## Aliases (Shorthand)

| Long Form | Short |
|-----------|-------|
| `instance` | `i` |
| `task` | `t` |
| `send` | `s` |
| `new` | `n` |
| `list` | `l` |
| `follow` | `f` |
| `view` | `v` |
| `cancel` | `c` |
| `resume` | `r` |
| `oneshot` | `o` |
| `kill` | `k` |
| `use` | `u` |

### Examples with Aliases
```bash
cline i l                        # instance list
cline t n "prompt"               # task new
cline t f                        # task follow
cline s "message"                # send
```

## Output Formats

```bash
# Default: Rich terminal UI
cline instance list

# JSON: For scripts/parsing
cline instance list -o json

# Plain: Simple text
cline instance list -o plain
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| "no such file or directory" | Run from project root, not cli/ directory |
| "No instances found" | Run `cline instance new` or `cline task new` |
| Connection refused | Check `cline instance list`, restart if needed |
| Task stuck | `cline task cancel` then retry |
| Clean slate | `cline instance kill --all` |

## Quick Setup

```bash
# 1. Build
cd /path/to/cline
npm run compile-cli

# 2. Verify
./cli/bin/cline version

# 3. Run examples
./cli/examples.sh

# 4. Create first task
./cli/bin/cline task new "Hello Cline!"
```

## Useful Patterns

```bash
# Watch task in real-time
cline task new "prompt" --wait

# Fire and forget (fully autonomous)
cline task oneshot "prompt"

# Batch operations
cline task new "task 1"
cline task new "task 2" --address localhost:50053

# Review with context
cline task new "Review security" -f auth.go -f middleware.go

# Debug with logs
cat app.log | cline send "What's causing this error?"
```

---

**Full Documentation:** See [README.md](README.md) for detailed examples and workflows.

