# Cline CLI Installation Guide

Complete guide for installing the Cline CLI globally on your system.

---

## Table of Contents

1. [Quick Install](#quick-install)
2. [Global Installation](#global-installation)
3. [Platform-Specific Instructions](#platform-specific-instructions)
4. [Verification](#verification)
5. [Shell Integration](#shell-integration)
6. [Troubleshooting](#troubleshooting)

---

## Quick Install

### Build from Source

```bash
cd /path/to/cline
npm run compile-cli
```

This creates:
- `cli/bin/cline` - Main CLI binary
- `cli/bin/cline-host` - Host bridge binary

---

## Global Installation

### Option 1: Symbolic Link (Recommended)

Create symbolic links so you can use `cline` from anywhere:

```bash
# Create symlink for cline
sudo ln -sf /path/to/cline/cli/bin/cline /usr/local/bin/cline

# Create symlink for cline-host
sudo ln -sf /path/to/cline/cli/bin/cline-host /usr/local/bin/cline-host

# Verify
which cline
# Output: /usr/local/bin/cline
```

**Important:** Even with global installation, you must still run `cline` from the Cline project root directory due to relative path dependencies.

### Option 2: Add to PATH

Add the CLI to your PATH:

```bash
# Add to ~/.bashrc, ~/.zshrc, or ~/.profile
export PATH="$PATH:/path/to/cline/cli/bin"

# Reload shell
source ~/.bashrc  # or ~/.zshrc
```

**Example for your system:**
```bash
# Add this line to ~/.zshrc:
export PATH="$PATH:/Users/niharturumella/projects/cline/cli/bin"

# Reload
source ~/.zshrc

# Verify
which cline
# Output: /Users/niharturumella/projects/cline/cli/bin/cline
```

### Option 3: Install Script (Advanced)

Create an install script:

```bash
#!/bin/bash
# install-cline.sh

CLINE_ROOT="/path/to/cline"
INSTALL_DIR="$HOME/.local/bin"

# Create install directory
mkdir -p "$INSTALL_DIR"

# Create symlinks
ln -sf "$CLINE_ROOT/cli/bin/cline" "$INSTALL_DIR/cline"
ln -sf "$CLINE_ROOT/cli/bin/cline-host" "$INSTALL_DIR/cline-host"

# Add to PATH if not already there
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> ~/.bashrc
    echo "✓ Added $INSTALL_DIR to PATH"
    echo "Run: source ~/.bashrc"
fi

echo "✓ Cline CLI installed successfully!"
echo ""
echo "Usage: cline --help"
echo "Remember: Must run from Cline project root!"
```

Make it executable and run:
```bash
chmod +x install-cline.sh
./install-cline.sh
```

---

## Platform-Specific Instructions

### macOS

```bash
# Using symlinks
sudo ln -sf /path/to/cline/cli/bin/cline /usr/local/bin/cline
sudo ln -sf /path/to/cline/cli/bin/cline-host /usr/local/bin/cline-host

# Or add to PATH in ~/.zshrc
echo 'export PATH="$PATH:/path/to/cline/cli/bin"' >> ~/.zshrc
source ~/.zshrc

# Verify
cline version
```

### Linux

```bash
# Using symlinks
sudo ln -sf /path/to/cline/cli/bin/cline /usr/local/bin/cline
sudo ln -sf /path/to/cline/cli/bin/cline-host /usr/local/bin/cline-host

# Or add to PATH in ~/.bashrc
echo 'export PATH="$PATH:/path/to/cline/cli/bin"' >> ~/.bashrc
source ~/.bashrc

# Verify
cline version
```

### Windows (WSL)

```bash
# Using symlinks in WSL
sudo ln -sf /mnt/c/path/to/cline/cli/bin/cline /usr/local/bin/cline
sudo ln -sf /mnt/c/path/to/cline/cli/bin/cline-host /usr/local/bin/cline-host

# Or add to PATH
echo 'export PATH="$PATH:/mnt/c/path/to/cline/cli/bin"' >> ~/.bashrc
source ~/.bashrc
```

---

## Working Directory Management

### Important Note About Project Root

**⚠️ Critical:** Even with global installation, you must run `cline` from the Cline project root directory.

The CLI looks for:
- `./cli/bin/cline-host`
- `./dist-standalone/cline-core.js`
- Other relative resources

### Selecting a Repository/Workspace

Use the `-w` or `--workdir` flag to specify which repository/directory to work in:

```bash
# From Cline project root (required)
cd /path/to/cline

# Work in a specific repository
cline task new "Fix bug in API" -w /projects/my-api

# Multiple repositories
cline task new "Sync repos" \
  -w /projects/frontend \
  -w /projects/backend

# With workflows
cline workflow run deploy -w /projects/production-app

# With OCA
cline task new "Review code" \
  -w /projects/enterprise-app \
  -s act-mode-api-provider=oca
```

### Workflow Directory Selection

Workflows automatically use the specified working directory:

```bash
# Create workflow
cline workflow create deploy-api

# Edit to include steps
cat > ~/Documents/Cline/Workflows/deploy-api.md << 'EOF'
# Deploy API

## Steps

1. Navigate to API directory (automatic)
2. Run tests: `npm test`
3. Build: `npm run build`
4. Deploy: `./deploy.sh`
EOF

# Run with specific repo
cline workflow run deploy-api -w /projects/my-api
```

### Setting Default Working Directory

You can set a default working directory per instance using environment variables or by navigating there first:

**Method 1: Navigate First (Simplest)**
```bash
# The task will use current directory context
cd /projects/my-app
cline task new "Run tests"
```

**Method 2: Always Specify -w Flag**
```bash
# Explicit working directory
cline task new "Run tests" -w /projects/my-app
```

**Method 3: Shell Alias**
```bash
# Add to ~/.bashrc or ~/.zshrc
alias cline-api='cline task new -w /projects/api'
alias cline-frontend='cline task new -w /projects/frontend'

# Usage
cline-api "Fix login bug"
cline-frontend "Update UI"
```

---

## Advanced Setup

### Create Helper Script

Create `~/bin/cline-wrapper.sh`:

```bash
#!/bin/bash
# Cline CLI wrapper - handles project root requirement

CLINE_ROOT="/path/to/cline"

# Change to cline project root
cd "$CLINE_ROOT" || exit 1

# Execute cline with all arguments
exec ./cli/bin/cline "$@"
```

Make it executable:
```bash
chmod +x ~/bin/cline-wrapper.sh
ln -sf ~/bin/cline-wrapper.sh /usr/local/bin/cline
```

Now you can run `cline` from anywhere!

### With Working Directory Selection

Enhanced wrapper:

```bash
#!/bin/bash
# Enhanced Cline CLI wrapper

CLINE_ROOT="/path/to/cline"
CURRENT_DIR="$(pwd)"

# Change to cline project root
cd "$CLINE_ROOT" || exit 1

# If no -w flag provided and not in cline root, add current directory
if [[ "$*" != *"-w"* ]] && [[ "$*" != *"--workdir"* ]] && [[ "$CURRENT_DIR" != "$CLINE_ROOT" ]]; then
    # For task new commands, add -w automatically
    if [[ "$*" == *"task new"* ]] || [[ "$*" == *"workflow run"* ]]; then
        exec ./cli/bin/cline "$@" -w "$CURRENT_DIR"
    fi
fi

# Execute normally
exec ./cli/bin/cline "$@"
```

With this wrapper:
```bash
# From anywhere
cd /projects/my-api
cline task new "Fix bug"  # Automatically uses /projects/my-api as workdir!
```

---

## Shell Integration

### Bash Completion

```bash
# Generate completion script
cline completion bash > /etc/bash_completion.d/cline

# Or for user-local:
cline completion bash > ~/.bash_completion.d/cline
```

### Zsh Completion

```bash
# Generate completion script
cline completion zsh > "${fpath[1]}/_cline"

# Reload
compinit
```

### Fish Completion

```bash
# Generate completion script
cline completion fish > ~/.config/fish/completions/cline.fish
```

---

## Verification

After installation, verify everything works:

### 1. Check Command Available

```bash
which cline
# Should show path to cline binary
```

### 2. Check Version

```bash
cline version
# Should show version info
```

### 3. Check Commands

```bash
cline --help
# Should list all commands including oca and workflow
```

### 4. Test OCA

```bash
cd /path/to/cline  # Must be in project root
cline oca --help
```

### 5. Test Workflows

```bash
cd /path/to/cline  # Must be in project root
cline workflow list
```

### 6. Test with Working Directory

```bash
cd /path/to/cline
cline task new "Test" -w /projects/test-repo
```

---

## Working Directory Examples

### Example 1: Single Repository

```bash
cd /path/to/cline  # CLI requirement

# Work in your API repo
cline task new "Add authentication" -w /projects/my-api

# Follow progress
cline task follow

# Send message
cline send "Use JWT tokens" -w /projects/my-api
```

### Example 2: Monorepo

```bash
cd /path/to/cline

# Work across multiple packages
cline task new "Update shared types" \
  -w /projects/monorepo/packages/shared \
  -w /projects/monorepo/packages/api \
  -w /projects/monorepo/packages/web
```

### Example 3: Different Projects on Different Instances

```bash
cd /path/to/cline

# Instance 1: API project
cline instance new
cline task new "API work" -w /projects/api --address localhost:50052

# Instance 2: Frontend project
cline instance new
cline task new "UI work" -w /projects/frontend --address localhost:50053

# Switch between them
cline instance use localhost:50052  # Now working on API
cline instance use localhost:50053  # Now working on frontend
```

### Example 4: Workflow with Multiple Repos

```bash
cd /path/to/cline

# Create multi-repo workflow
cline workflow create sync-repos

# Edit it
cat > ~/Documents/Cline/Workflows/sync-repos.md << 'EOF'
# Sync Multiple Repositories

## Steps

1. Pull latest from repo A
   ```bash
   cd /projects/repo-a && git pull
   ```

2. Pull latest from repo B
   ```bash
   cd /projects/repo-b && git pull
   ```

3. Check for conflicts
   - Compare shared dependencies
   - Verify API contracts match

4. Run cross-repo tests
   ```bash
   cd /projects/integration-tests && npm test
   ```
EOF

# Run it
cline workflow run sync-repos \
  -w /projects/repo-a \
  -w /projects/repo-b
```

---

## Environment Setup

### Recommended Setup

Create a permanent installation:

```bash
# 1. Add to PATH
echo 'export PATH="$PATH:/Users/niharturumella/projects/cline/cli/bin"' >> ~/.zshrc

# 2. Add Go bin (if not already)
echo 'export PATH="$PATH:/Users/niharturumella/go/bin"' >> ~/.zshrc

# 3. Create wrapper function for convenience
cat >> ~/.zshrc << 'EOF'

# Cline CLI helper function
cline-here() {
    (cd /Users/niharturumella/projects/cline && cline "$@" -w "$(pwd)")
}

# Cline with workflow
cline-wf() {
    (cd /Users/niharturumella/projects/cline && cline workflow run "$@" -w "$(pwd)")
}
EOF

# 4. Reload
source ~/.zshrc
```

Now you can use:

```bash
# From anywhere - uses current directory as workdir
cd /projects/my-app
cline-here task new "Fix bug"

# Run workflow in current directory
cd /projects/my-app
cline-wf deploy-prod --wait
```

---

## Directory Structure

After global installation:

```
System:
├── /usr/local/bin/
│   ├── cline -> /path/to/cline/cli/bin/cline
│   └── cline-host -> /path/to/cline/cli/bin/cline-host

User:
├── ~/.cline/
│   ├── instances.db              # Instance registry
│   └── default_instance.json    # Default instance

├── ~/Documents/Cline/
│   ├── Workflows/                # Your workflows
│   │   ├── example-hello.md
│   │   ├── example-git-commit.md
│   │   └── your-workflows.md
│   │
│   └── Rules/                    # Cline rules (optional)

Project:
└── /path/to/cline/               # Must exist!
    ├── cli/bin/
    │   ├── cline
    │   └── cline-host
    ├── dist-standalone/
    │   └── cline-core.js
    └── ... other files
```

---

## Usage After Global Install

### Basic Commands (from project root)

```bash
# Always start here
cd /path/to/cline

# Now use cline globally
cline version
cline instance list
cline workflow list
```

### Working in Different Repos

```bash
# The Cline CLI must run from Cline project root
cd /path/to/cline

# But your work can be in any directory using -w flag
cline task new "Your task" -w /projects/any-repo
cline workflow run deploy -w /projects/production-app
cline task new "Multi-repo task" \
  -w /projects/frontend \
  -w /projects/backend
```

### Using Shell Helpers

With the shell functions from the setup:

```bash
# Navigate to your project
cd /projects/my-awesome-app

# Use cline-here (automatically uses current directory)
cline-here task new "Add feature X"
cline-here task follow

# Or run workflow in current directory
cline-wf my-deployment --wait
```

---

## Platform-Specific Instructions

### macOS Installation

```bash
# 1. Build CLI
cd /Users/niharturumella/projects/cline
npm run compile-cli

# 2. Create symlinks
sudo ln -sf /Users/niharturumella/projects/cline/cli/bin/cline /usr/local/bin/cline
sudo ln -sf /Users/niharturumella/projects/cline/cli/bin/cline-host /usr/local/bin/cline-host

# 3. Add to PATH (alternative to symlinks)
echo 'export PATH="$PATH:/Users/niharturumella/projects/cline/cli/bin"' >> ~/.zshrc
source ~/.zshrc

# 4. Verify
cline version

# 5. Setup helper function
cat >> ~/.zshrc << 'EOF'
cline-here() {
    (cd /Users/niharturumella/projects/cline && cline "$@" -w "$(pwd)")
}
EOF
source ~/.zshrc
```

### Linux Installation

```bash
# 1. Build CLI
cd ~/projects/cline
npm run compile-cli

# 2. Create symlinks
sudo ln -sf ~/projects/cline/cli/bin/cline /usr/local/bin/cline
sudo ln -sf ~/projects/cline/cli/bin/cline-host /usr/local/bin/cline-host

# 3. Or add to PATH
echo 'export PATH="$PATH:~/projects/cline/cli/bin"' >> ~/.bashrc
source ~/.bashrc

# 4. Verify
cline version

# 5. Setup helper function
cat >> ~/.bashrc << 'EOF'
cline-here() {
    (cd ~/projects/cline && cline "$@" -w "$(pwd)")
}
EOF
source ~/.bashrc
```

### Windows (WSL)

```bash
# 1. Build CLI
cd /mnt/c/Users/YourName/projects/cline
npm run compile-cli

# 2. Create symlinks
sudo ln -sf /mnt/c/Users/YourName/projects/cline/cli/bin/cline /usr/local/bin/cline
sudo ln -sf /mnt/c/Users/YourName/projects/cline/cli/bin/cline-host /usr/local/bin/cline-host

# 3. Verify
cline version
```

---

## Repository Selection Patterns

### Pattern 1: Explicit Flag

Always specify working directory:

```bash
cd /path/to/cline
cline task new "Task" -w /projects/repo
```

### Pattern 2: Shell Function

Use helper that auto-detects:

```bash
# Add to shell config:
cline-here() {
    (cd /path/to/cline && cline "$@" -w "$(pwd)")
}

# Use from any directory:
cd /projects/my-repo
cline-here task new "Fix bug"
```

### Pattern 3: Environment Variable

Set default repo for session:

```bash
# Set for session
export CLINE_WORKDIR=/projects/my-repo

# Create wrapper that uses it
cline-work() {
    (cd /path/to/cline && cline "$@" -w "${CLINE_WORKDIR:-$(pwd)}")
}

# Use it
cline-work task new "Add feature"
```

### Pattern 4: Per-Project Aliases

Create project-specific aliases:

```bash
# Add to ~/.zshrc or ~/.bashrc
alias cline-api='(cd /path/to/cline && cline task new -w /projects/api)'
alias cline-web='(cd /path/to/cline && cline task new -w /projects/web)'
alias cline-mobile='(cd /path/to/cline && cline task new -w /projects/mobile)'

# Usage
cline-api "Fix authentication"
cline-web "Update dashboard"
cline-mobile "Add notifications"
```

### Pattern 5: Interactive Selection

Create an interactive repo selector:

```bash
#!/bin/bash
# cline-select.sh

CLINE_ROOT="/path/to/cline"

# Define your repositories
declare -A REPOS=(
    ["api"]="/projects/my-api"
    ["frontend"]="/projects/my-frontend"
    ["backend"]="/projects/my-backend"
    ["mobile"]="/projects/my-mobile"
)

# Show menu
echo "Select repository:"
select repo in "${!REPOS[@]}"; do
    if [ -n "$repo" ]; then
        workdir="${REPOS[$repo]}"
        echo "Selected: $repo ($workdir)"
        
        # Get task from remaining arguments
        shift
        task="$*"
        
        # Run cline
        cd "$CLINE_ROOT"
        ./cli/bin/cline task new "$task" -w "$workdir"
        break
    fi
done
```

Usage:
```bash
./cline-select.sh
# 1) api
# 2) frontend
# 3) backend
# 4) mobile
# Select: 1
# Enter task: Fix login bug
```

---

## Verification

### Test Global Installation

```bash
# Test from any directory
cd ~
cline version

# Test from Cline root
cd /path/to/cline
cline instance list

# Test with workdir
cd /path/to/cline
cline task new "Test" -w /tmp/test-dir
```

### Test Workflows with Repos

```bash
cd /path/to/cline

# Create test workflow
cline workflow create test-repo

# Run in specific repo
cline workflow run test-repo -w /projects/my-repo
```

### Test Helper Functions

```bash
# Test cline-here
cd /projects/any-repo
cline-here task new "Test task"

# Verify working directory was set correctly
cline task follow  # Should show working in /projects/any-repo
```

---

## Troubleshooting

### "command not found: cline"

**Solution 1:** Check PATH
```bash
echo $PATH
# Should include /path/to/cline/cli/bin or /usr/local/bin
```

**Solution 2:** Reload shell
```bash
source ~/.zshrc  # or ~/.bashrc
```

**Solution 3:** Verify symlink
```bash
ls -la /usr/local/bin/cline
# Should point to actual binary
```

### "no such file or directory: ./cli/bin/cline-host"

**Problem:** Not running from Cline project root

**Solution:**
```bash
# Must run from here:
cd /path/to/cline
cline task new "Your task"

# Or use wrapper script (see Advanced Setup)
```

### Working Directory Not Used

**Problem:** Task is working in wrong directory

**Solution:** Always specify -w flag
```bash
cline task new "Task" -w /projects/correct-repo
```

### Shell Function Not Working

**Problem:** Function not found

**Solution:**
```bash
# Check if function is defined
type cline-here

# If not, reload shell config
source ~/.zshrc

# Or redefine
cline-here() {
    (cd /path/to/cline && cline "$@" -w "$(pwd)")
}
```

---

## Complete Setup Example

Here's a complete setup for your system:

```bash
#!/bin/bash
# Complete Cline CLI Setup

CLINE_ROOT="/Users/niharturumella/projects/cline"

echo "Setting up Cline CLI..."

# 1. Build
echo "Building CLI..."
cd "$CLINE_ROOT"
npm run compile-cli

# 2. Create symlinks
echo "Creating symlinks..."
sudo ln -sf "$CLINE_ROOT/cli/bin/cline" /usr/local/bin/cline
sudo ln -sf "$CLINE_ROOT/cli/bin/cline-host" /usr/local/bin/cline-host

# 3. Setup shell helpers
echo "Adding shell helpers..."
cat >> ~/.zshrc << 'EOF'

# Cline CLI helpers
export PATH="$PATH:/Users/niharturumella/go/bin"

cline-here() {
    (cd /Users/niharturumella/projects/cline && cline "$@" -w "$(pwd)")
}

cline-wf() {
    (cd /Users/niharturumella/projects/cline && cline workflow run "$@" -w "$(pwd)")
}

alias cline-root='cd /Users/niharturumella/projects/cline'
EOF

# 4. Create workflows directory
echo "Creating workflows directory..."
mkdir -p ~/Documents/Cline/Workflows

# 5. Reload shell
echo "Reloading shell..."
source ~/.zshrc

echo "✓ Setup complete!"
echo ""
echo "Usage:"
echo "  cline version              # From anywhere"
echo "  cd /projects/my-repo"
echo "  cline-here task new \"task\" # Auto-uses current dir"
echo ""
```

---

## Quick Reference Card

Save this for quick access:

```
┌────────────────────────────────────────────────┐
│         CLINE CLI GLOBAL INSTALL               │
└────────────────────────────────────────────────┘

REQUIREMENT:
  Must run from: /path/to/cline

BASIC USAGE:
  cd /path/to/cline
  cline task new "task" -w /projects/repo

HELPER FUNCTIONS:
  cline-here task new "task"    # Uses current dir
  cline-wf workflow-name         # Runs workflow in current dir

REPOSITORY SELECTION:
  -w /path/to/repo              # Single repo
  -w /repo1 -w /repo2           # Multiple repos

EXAMPLES:
  cd /projects/my-app
  cline-here task new "Fix bug"
  
  cd /path/to/cline
  cline workflow run deploy -w /projects/prod
  
  cline task new "Multi-repo" \
    -w /projects/frontend \
    -w /projects/backend
```

---

## Uninstall

To remove global installation:

```bash
# Remove symlinks
sudo rm /usr/local/bin/cline
sudo rm /usr/local/bin/cline-host

# Remove from PATH (edit ~/.zshrc or ~/.bashrc)
# Remove the line: export PATH="$PATH:/path/to/cline/cli/bin"

# Remove shell functions (edit ~/.zshrc or ~/.bashrc)
# Remove cline-here and cline-wf functions

# Reload
source ~/.zshrc

# Verify
which cline
# Should return nothing
```

---

## Best Practices

### 1. Always Specify Working Directory

```bash
# Good
cline task new "Fix bug" -w /projects/api

# Better - use helper
cline-here task new "Fix bug"
```

### 2. Use Shell Helpers

Create helpers for your common workflows:

```bash
# In ~/.zshrc
cline-deploy() {
    (cd /path/to/cline && cline workflow run deploy -w "$1" --wait)
}

cline-review() {
    (cd /path/to/cline && cline workflow run code-review -w "$1")
}

# Usage
cline-deploy /projects/api
cline-review /projects/web
```

### 3. Per-Project Configuration

Create `.clinerc` in your project:

```bash
# /projects/my-api/.clinerc
CLINE_PROVIDER=oca
CLINE_MODEL=claude-sonnet
CLINE_MODE=plan
```

Then create wrapper that reads it:

```bash
cline-smart() {
    local project_dir="$(pwd)"
    if [ -f "$project_dir/.clinerc" ]; then
        source "$project_dir/.clinerc"
    fi
    
    (cd /path/to/cline && cline task new "$@" \
        -w "$project_dir" \
        ${CLINE_PROVIDER:+-s act-mode-api-provider=$CLINE_PROVIDER} \
        ${CLINE_MODEL:+-s act-mode-api-model-id=$CLINE_MODEL} \
        ${CLINE_MODE:+-m $CLINE_MODE})
}
```

---

## Summary

### What You Get

After following this guide:

1. ✅ **Global `cline` command** available everywhere
2. ✅ **Helper functions** for convenience
3. ✅ **Repository selection** via -w flag
4. ✅ **Shell integration** for tab completion
5. ✅ **Workflows** working with any repo

### Key Points to Remember

- ⚠️ **Must run from Cline project root** (due to relative paths)
- ✅ **Use -w flag** to specify working directory
- ✅ **Helper functions** make it easier
- ✅ **Shell completion** available
- ✅ **Multiple repos** supported with multiple -w flags

---

## Quick Start After Install

```bash
# 1. Build and install
cd /path/to/cline
npm run compile-cli
sudo ln -sf $(pwd)/cli/bin/cline /usr/local/bin/cline

# 2. Setup helper
echo 'cline-here() { (cd /path/to/cline && cline "$@" -w "$(pwd)") }' >> ~/.zshrc
source ~/.zshrc

# 3. Use from anywhere
cd /projects/my-app
cline-here task new "Your first task"
```

---

**Installation Complete! 🎉**

For usage instructions, see:
- [USER_MANUAL.md](USER_MANUAL.md) - Complete guide
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Command reference
- [README.md](README.md) - Overview
