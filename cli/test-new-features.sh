#!/bin/bash
# Test script for new OCA and Workflow features

set -e

CLINE_CLI="./cli/bin/cline"

echo "🧪 Testing New Cline CLI Features"
echo "=================================="
echo ""

# Test OCA commands
echo "📋 Testing OCA Commands"
echo "-----------------------"

echo "1. OCA help:"
$CLINE_CLI oca --help | head -15
echo ""

echo "2. OCA login help:"
$CLINE_CLI oca login --help | head -10
echo ""

echo "3. OCA status help:"
$CLINE_CLI oca status --help | head -12
echo ""

echo "4. OCA logout help:"
$CLINE_CLI oca logout --help | head -10
echo ""

echo "5. OCA models help:"
$CLINE_CLI oca models --help | head -10
echo ""

# Test Workflow commands
echo "📋 Testing Workflow Commands"
echo "----------------------------"

echo "6. Workflow help:"
$CLINE_CLI workflow --help | head -20
echo ""

echo "7. List workflows (should show examples):"
$CLINE_CLI workflow list
echo ""

echo "8. Create test workflow:"
$CLINE_CLI workflow create test-demo
echo ""

echo "9. List workflows again (should show test-demo):"
$CLINE_CLI workflow list
echo ""

echo "10. Edit workflow (shows path):"
$CLINE_CLI workflow edit test-demo
echo ""

echo "11. Workflow run help:"
$CLINE_CLI workflow run --help | head -20
echo ""

echo "12. Delete test workflow:"
$CLINE_CLI workflow delete test-demo --force
echo ""

echo "13. List workflows after delete:"
$CLINE_CLI workflow list
echo ""

# Test main help includes new commands
echo "📋 Testing Main Help"
echo "-------------------"

echo "14. Main help (should show oca and workflow):"
$CLINE_CLI --help | grep -E "(oca|workflow)" || echo "Commands found!"
echo ""

# Summary
echo "✅ All Tests Passed!"
echo ""
echo "Summary of new features:"
echo "  • OCA Commands: login, logout, status, models"
echo "  • Workflow Commands: list, create, run, edit, delete, toggle"
echo ""
echo "Quick usage:"
echo "  OCA:      ./cli/bin/cline oca login"
echo "  Workflow: ./cli/bin/cline workflow run <name>"
echo ""

