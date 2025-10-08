#!/bin/bash
# Cline CLI Examples
# Run this script from the Cline project root directory

set -e

CLINE_CLI="./cli/bin/cline"

echo "🚀 Cline CLI Examples"
echo "===================="
echo ""

# Example 1: Version
echo "1. Check version:"
echo "   Command: $CLINE_CLI version --short"
$CLINE_CLI version --short
echo ""

# Example 2: Full version info
echo "2. Full version information:"
echo "   Command: $CLINE_CLI version"
$CLINE_CLI version
echo ""

# Example 3: List instances
echo "3. List all instances:"
echo "   Command: $CLINE_CLI instance list"
$CLINE_CLI instance list
echo ""

# Example 4: Help commands
echo "4. View available commands:"
echo "   Command: $CLINE_CLI --help"
$CLINE_CLI --help | head -20
echo ""

# Example 5: Task command help
echo "5. View task commands:"
echo "   Command: $CLINE_CLI task --help"
$CLINE_CLI task --help
echo ""

# Example 6: Instance command help
echo "6. View instance commands:"
echo "   Command: $CLINE_CLI instance --help"
$CLINE_CLI instance --help
echo ""

echo "7. View OCA commands:"
echo "   Command: $CLINE_CLI oca --help"
$CLINE_CLI oca --help
echo ""

echo "8. View Workflow commands:"
echo "   Command: $CLINE_CLI workflow --help"
$CLINE_CLI workflow --help
echo ""

echo "9. List workflows:"
echo "   Command: $CLINE_CLI workflow list"
$CLINE_CLI workflow list
echo ""

echo "✅ Basic CLI commands demonstrated successfully!"
echo ""
echo "To run these commands yourself:"
echo "  cd /path/to/cline"
echo "  ./cli/bin/cline [command]"
echo ""
echo "For interactive examples (requires running instances):"
echo "  ./cli/bin/cline instance new"
echo "  ./cli/bin/cline task new \"Create a hello world function\""
echo "  ./cli/bin/cline task follow"
echo ""
echo "New features:"
echo "  ./cli/bin/cline oca login"
echo "  ./cli/bin/cline workflow run example-hello --wait"
echo ""

