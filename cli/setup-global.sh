#!/bin/bash
# Global Installation Script for Cline CLI

set -e

CLINE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SHELL_RC="${HOME}/.zshrc"

# Detect shell
if [ -n "$BASH_VERSION" ]; then
    SHELL_RC="${HOME}/.bashrc"
elif [ -n "$ZSH_VERSION" ]; then
    SHELL_RC="${HOME}/.zshrc"
fi

echo "🚀 Cline CLI Global Setup"
echo "========================="
echo ""
echo "Cline root: $CLINE_ROOT"
echo "Shell config: $SHELL_RC"
echo ""

# Check if binaries exist
if [ ! -f "$CLINE_ROOT/cli/bin/cline" ]; then
    echo "❌ Error: cline binary not found!"
    echo "   Run 'npm run compile-cli' first"
    exit 1
fi

# Option 1: Symlinks
echo "Option 1: Create symlinks in /usr/local/bin (requires sudo)"
read -p "Create symlinks? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    sudo ln -sf "$CLINE_ROOT/cli/bin/cline" /usr/local/bin/cline
    sudo ln -sf "$CLINE_ROOT/cli/bin/cline-host" /usr/local/bin/cline-host
    echo "✓ Symlinks created"
    echo "  - /usr/local/bin/cline -> $CLINE_ROOT/cli/bin/cline"
    echo "  - /usr/local/bin/cline-host -> $CLINE_ROOT/cli/bin/cline-host"
else
    # Option 2: Add to PATH
    echo ""
    echo "Option 2: Add to PATH"
    if ! grep -q "cline/cli/bin" "$SHELL_RC"; then
        echo "export PATH=\"\$PATH:$CLINE_ROOT/cli/bin\"" >> "$SHELL_RC"
        echo "✓ Added to PATH in $SHELL_RC"
    else
        echo "✓ Already in PATH"
    fi
fi

# Add helper functions
echo ""
echo "Adding helper functions..."

# Check if functions already exist
if ! grep -q "cline-here()" "$SHELL_RC"; then
    cat >> "$SHELL_RC" << EOFUNC

# Cline CLI helper functions
cline-here() {
    (cd "$CLINE_ROOT" && cline "\$@" -w "\$(pwd)")
}

cline-wf() {
    (cd "$CLINE_ROOT" && cline workflow run "\$@" -w "\$(pwd)")
}

alias cline-root='cd $CLINE_ROOT'
EOFUNC
    echo "✓ Added helper functions to $SHELL_RC"
else
    echo "✓ Helper functions already exist"
fi

# Create workflows directory
mkdir -p ~/Documents/Cline/Workflows
echo "✓ Created workflows directory"

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Reload your shell: source $SHELL_RC"
echo "  2. Verify: cline version"
echo "  3. Try: cd /your/project && cline-here task new \"Your task\""
echo ""
echo "Helper functions available:"
echo "  cline-here     - Run cline with current directory as workdir"
echo "  cline-wf       - Run workflow in current directory"
echo "  cline-root     - Navigate to Cline project root"
echo ""
