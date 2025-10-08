# ✅ Global Installation & Repository Selection Added

**Status:** COMPLETE  
**Date:** October 8, 2025

---

## What Was Added

### 1. Global Installation Support ✅

**New File:** `cli/INSTALL.md` (500+ lines)

Complete installation guide including:
- Global installation via symlinks
- Global installation via PATH
- Platform-specific instructions (macOS, Linux, Windows/WSL)
- Shell helper functions
- Verification steps
- Troubleshooting

**New File:** `cli/setup-global.sh` (Automated installer)

One-command setup:
```bash
./cli/setup-global.sh
```

This script:
- Creates symlinks or adds to PATH
- Adds helper functions to shell config
- Creates workflows directory
- Verifies installation

### 2. Repository Selection ✅

**Helper Functions Added:**

```bash
# Auto-use current directory as workdir
cline-here() {
    (cd /path/to/cline && cline "$@" -w "$(pwd)")
}

# Run workflow in current directory
cline-wf() {
    (cd /path/to/cline && cline workflow run "$@" -w "$(pwd)")
}

# Quick navigation to Cline root
alias cline-root='cd /path/to/cline'
```

---

## How to Use

### Global Installation

**Automated (Recommended):**
```bash
cd /path/to/cline
./cli/setup-global.sh
source ~/.zshrc
```

**Manual:**
```bash
# Create symlinks
sudo ln -sf /path/to/cline/cli/bin/cline /usr/local/bin/cline
sudo ln -sf /path/to/cline/cli/bin/cline-host /usr/local/bin/cline-host

# Add helpers
cat >> ~/.zshrc << 'EOF'
cline-here() {
    (cd /path/to/cline && cline "$@" -w "$(pwd)")
}
EOF

source ~/.zshrc
```

### Repository Selection

**Method 1: Explicit -w flag**
```bash
cd /path/to/cline
cline task new "Task" -w /projects/my-repo
```

**Method 2: cline-here helper (Recommended)**
```bash
cd /projects/my-repo
cline-here task new "Task"  # Auto-uses /projects/my-repo
```

**Method 3: Multiple repositories**
```bash
cline task new "Multi-repo task" \
  -w /projects/frontend \
  -w /projects/backend
```

---

## Complete Usage Examples

### Example 1: Work from Any Directory

```bash
# Install globally
./cli/setup-global.sh
source ~/.zshrc

# Navigate to your project
cd /projects/my-awesome-app

# Use cline-here
cline-here task new "Add user authentication"
cline-here task follow
cline-here send "Use JWT tokens"
```

### Example 2: Multiple Projects

```bash
# Project 1: API
cd /projects/api
cline-here task new "Fix login endpoint"

# Project 2: Frontend  
cd /projects/frontend
cline-here task new "Update login UI"

# Both tasks run in their respective directories!
```

### Example 3: Workflows with Auto-Selection

```bash
# Navigate to your repo
cd /projects/production-app

# Run workflow in current directory
cline-wf deploy-prod --wait

# Equivalent to:
# cd /path/to/cline
# cline workflow run deploy-prod -w /projects/production-app --wait
```

### Example 4: OCA + Repository Selection

```bash
cd /projects/enterprise-app

# Login to OCA
cline oca login
cline oca status

# Create task with OCA in current repo
cline-here task new "Security audit" \
  -s act-mode-api-provider=oca \
  -s act-mode-oca-model-id=enterprise-model
```

### Example 5: Advanced Multi-Repo

```bash
# Work across multiple repositories
cd /path/to/cline

cline task new "Sync API and frontend types" \
  -w /projects/api \
  -w /projects/frontend \
  -s yolo-mode-toggled=true

# Or with helper from monorepo root
cd /projects/monorepo
cline-here task new "Update all packages" \
  -w "$(pwd)/packages/api" \
  -w "$(pwd)/packages/web" \
  -w "$(pwd)/packages/shared"
```

---

## Documentation Updates

### Updated Files

1. **README.md**
   - Added "Global Installation" section
   - Added "Repository Selection" section
   - Added reference to INSTALL.md

2. **USER_MANUAL.md**
   - Added "Global Installation (Optional)" in Getting Started
   - Added "Repository Selection" examples
   - Updated all examples to show -w flag usage

3. **QUICK_REFERENCE.md**
   - Added "Global Installation" section
   - Added "Repository Selection" section
   - Updated all workflows to show -w usage
   - Added cline-here examples

4. **DOCUMENTATION_INDEX.md**
   - Added INSTALL.md to file list
   - Added navigation for installation

### New Files

5. **INSTALL.md** (NEW)
   - Complete installation guide (500+ lines)
   - Platform-specific instructions
   - Repository selection patterns
   - Shell integration
   - Advanced setup options

6. **setup-global.sh** (NEW)
   - Automated installation script
   - Interactive setup
   - Adds helper functions
   - Creates directories

---

## Setup Script Features

The `setup-global.sh` script provides:

✅ **Automatic detection** of shell type (bash/zsh)  
✅ **Interactive prompts** for installation options  
✅ **Symlink creation** (with sudo)  
✅ **PATH addition** (alternative to symlinks)  
✅ **Helper functions** automatically added  
✅ **Workflows directory** creation  
✅ **Verification** of binaries  
✅ **Instructions** for next steps  

---

## Helper Functions Reference

### cline-here

Run cline commands using current directory as workdir:

```bash
cd /projects/my-app
cline-here task new "Your task"
cline-here task follow
cline-here send "Message"
cline-here task cancel
```

### cline-wf

Run workflows in current directory:

```bash
cd /projects/my-app
cline-wf deploy-prod --wait
cline-wf run-tests
cline-wf code-review -m plan
```

### cline-root

Quickly navigate to Cline root:

```bash
cline-root
# Now in: /path/to/cline
```

---

## Verification Commands

Test that everything works:

```bash
# 1. Test global command
cline version

# 2. Test from any directory
cd ~
cline version

# 3. Test helper
cd /tmp
cline-here task new "Test" 2>&1 | head -5

# 4. Test repository selection
cd /path/to/cline
cline task new "Test" -w /tmp/test-dir

# 5. Test workflow with repo
cd /projects/my-app
cline-wf example-hello  # (if you have an instance running)
```

---

## Platform-Specific Setup

### macOS (Your System)

```bash
cd /Users/niharturumella/projects/cline
./cli/setup-global.sh
source ~/.zshrc

# Test
cline version
cd /projects/any-repo
cline-here task new "Test"
```

### Linux

```bash
cd ~/projects/cline
./cli/setup-global.sh
source ~/.bashrc

# Test
cline version
cd /projects/any-repo
cline-here task new "Test"
```

### Windows (WSL)

```bash
cd /mnt/c/projects/cline
./cli/setup-global.sh
source ~/.bashrc

# Test
cline version
```

---

## Advanced Patterns

### Pattern 1: Project-Specific Aliases

```bash
# Add to ~/.zshrc
alias work-api='cd /projects/api && cline-here'
alias work-web='cd /projects/web && cline-here'

# Usage
work-api task new "Fix bug"
work-web task new "Update UI"
```

### Pattern 2: Team Workflows

```bash
# Create team workflows directory (in git)
mkdir -p /projects/shared-workflows

# Symlink to Cline workflows
ln -sf /projects/shared-workflows ~/Documents/Cline/Workflows/team

# Now team workflows are available
cline workflow list
```

### Pattern 3: Multi-Project Helper

```bash
# Advanced helper for multiple projects
cline-multi() {
    local projects=""
    for dir in "$@"; do
        if [ -d "$dir" ]; then
            projects="$projects -w $dir"
        fi
    done
    (cd /path/to/cline && cline task new $projects)
}

# Usage
cline-multi /projects/api /projects/web /projects/shared
# Opens task with all three directories
```

---

## Quick Start Guide

### 1. Install Globally

```bash
cd /Users/niharturumella/projects/cline
./cli/setup-global.sh
```

Follow the prompts, then:

```bash
source ~/.zshrc
```

### 2. Verify Installation

```bash
cline version
cline --help | grep -E "(oca|workflow)"
```

### 3. Try Repository Selection

```bash
# Navigate to any project
cd /projects/your-app

# Use helper
cline-here task new "Test repository selection"

# Check it's using the right directory in task output
```

### 4. Create and Run Workflow

```bash
# Create workflow
cline workflow create my-first-workflow

# Edit it
$EDITOR ~/Documents/Cline/Workflows/my-first-workflow.md

# Run from your project directory
cd /projects/your-app
cline-wf my-first-workflow --wait
```

---

## Before and After

### Before

```bash
# Had to navigate to Cline root every time
cd /path/to/cline
./cli/bin/cline task new "Task" -w /projects/api

# From different directory
cd /projects/api
# Can't use cline here!

# Manual workflow invocation
cd /path/to/cline
./cli/bin/cline task new "/workflow.md" -w /projects/api
```

### After

```bash
# Can use cline from anywhere!
cline version
cline instance list

# From any directory with helper
cd /projects/api
cline-here task new "Task"  # Auto-uses /projects/api!

# Workflow helper
cd /projects/api
cline-wf deploy --wait  # Runs in /projects/api
```

---

## Complete Feature List

### Global Installation
- ✅ Automated setup script
- ✅ Symlink installation
- ✅ PATH installation
- ✅ Shell helper functions
- ✅ Platform-specific instructions
- ✅ Verification commands

### Repository Selection
- ✅ -w flag for single repo
- ✅ Multiple -w flags for multi-repo
- ✅ cline-here helper (auto-detect)
- ✅ cline-wf helper (workflows)
- ✅ Per-project patterns
- ✅ Interactive selection script

---

## File Summary

### New/Updated Files

| File | Status | Purpose |
|------|--------|---------|
| `cli/INSTALL.md` | ✅ NEW | Complete installation guide |
| `cli/setup-global.sh` | ✅ NEW | Automated setup script |
| `cli/README.md` | ✅ UPDATED | Added global install section |
| `cli/USER_MANUAL.md` | ✅ UPDATED | Added install & repo selection |
| `cli/QUICK_REFERENCE.md` | ✅ UPDATED | Added install & repo examples |
| `cli/DOCUMENTATION_INDEX.md` | ✅ UPDATED | Added INSTALL.md reference |

---

## Final Checklist

- [x] Global installation script created
- [x] Installation guide written
- [x] Repository selection documented
- [x] Helper functions provided
- [x] Platform-specific instructions
- [x] README updated
- [x] USER_MANUAL updated
- [x] QUICK_REFERENCE updated
- [x] Examples added
- [x] Verification steps provided

**Status:** ✅ **COMPLETE**

---

## Try It Now!

```bash
# 1. Run setup
cd /Users/niharturumella/projects/cline
./cli/setup-global.sh

# 2. Reload shell
source ~/.zshrc

# 3. Test from anywhere
cd /projects/any-directory
cline version

# 4. Use helper
cline-here task new "Test with auto repo selection"

# 5. Run workflow
cline-wf example-hello
```

---

**Installation Complete! 🎉**

You can now:
- ✅ Use `cline` from any directory
- ✅ Auto-select repositories with `cline-here`
- ✅ Run workflows with `cline-wf`
- ✅ Work across multiple projects easily

For more details, see [INSTALL.md](INSTALL.md).
