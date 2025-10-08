# Cline CLI Documentation Index

Welcome to the Cline CLI documentation! This index will help you find the right documentation for your needs.

## 📚 Documentation Files

### 1. [README.md](README.md) - Start Here!
**Best for: First-time users, installation, quick overview**

What you'll find:
- ✅ Features overview
- ✅ Build instructions
- ✅ Quick start guide
- ✅ Command reference
- ✅ Common workflows
- ✅ Installation tips

Start with this if you're new to the CLI.

---

### 2. [USER_MANUAL.md](USER_MANUAL.md) - Complete Usage Guide
**Best for: Learning all features, provider configuration, advanced usage**

What you'll find:
- 🚀 Starting and following tasks
- 🔄 Switching between 40+ AI providers (Claude, GPT, Gemini, Ollama, and more!)
- 🎯 Working with modes (act/plan/yolo)
- 🔧 Managing multiple instances
- ⚙️ Advanced task settings (100+ configuration options)
- 🤖 Auto-approval configuration
- 💼 Real-world workflows and examples
- 📋 Complete settings reference

**Size**: ~1000 lines - Comprehensive guide with examples

Use this for:
- "How do I switch to using OpenAI instead of Claude?"
- "What's the difference between act and plan mode?"
- "How do I configure auto-approval for safe commands?"
- "What settings can I pass to -s flags?"

---

### 3. [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Command Cheat Sheet
**Best for: Quick lookups, command syntax, daily usage**

What you'll find:
- ⚡ Essential commands
- 🔤 Command aliases and shortcuts
- 📊 Quick comparison tables
- 🎨 Common usage patterns
- 🔧 Quick troubleshooting

**Size**: ~400 lines - Concise reference

Use this when:
- "What was that command to follow a task?"
- "What are the shortcuts for common commands?"
- "Quick, what's the syntax for attaching files?"

---

### 4. [ARCHITECTURE.md](ARCHITECTURE.md) - Technical Architecture
**Best for: Developers, contributors, understanding internals**

What you'll find:
- 🏗️ High-level architecture diagram
- 📦 Component details and interactions
- 🔄 Data flow diagrams
- 🎨 Design patterns used
- 🔌 Process architecture
- 🧩 Extension points for contributors

**Size**: ~800 lines - Deep technical documentation

Use this when:
- Contributing to the CLI
- Understanding how components interact
- Debugging complex issues
- Adding new features or handlers

---

### 5. [examples.sh](examples.sh) - Runnable Examples
**Best for: Seeing commands in action**

What you'll find:
- ▶️ Executable example commands
- 📺 Live demonstrations
- ✅ Verification that CLI is working

**Run it:**
```bash
cd /path/to/cline
./cli/examples.sh
```

---

### 6. [INSTALL.md](INSTALL.md) - Installation Guide
**Best for: Setting up global installation, shell integration**

What you'll find:
- 🔧 Global installation instructions
- 🖥️ Platform-specific setup (macOS, Linux, Windows/WSL)
- 📁 Repository selection patterns
- 🐚 Shell helper functions
- ✅ Verification steps

**Size**: ~500 lines - Complete installation guide

Use this when:
- Installing CLI globally
- Setting up shell helpers
- Configuring repository selection
- "How do I use cline from any directory?"

---

## 🎯 Quick Navigation Guide

### I want to...

#### Get Started
→ [README.md](README.md) - Quick Start section

#### Learn how to use a specific feature
→ [USER_MANUAL.md](USER_MANUAL.md) - Table of Contents

#### Find a command quickly
→ [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Essential Commands

#### Switch AI providers
→ [USER_MANUAL.md](USER_MANUAL.md) - "Switching Providers & Models" section

#### Understand how it works
→ [ARCHITECTURE.md](ARCHITECTURE.md) - High-Level Architecture Diagram

#### See working examples
→ [examples.sh](examples.sh) - Run the script

#### Troubleshoot an issue
→ [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Troubleshooting section
→ [USER_MANUAL.md](USER_MANUAL.md) - Troubleshooting section

#### Contribute code
→ [ARCHITECTURE.md](ARCHITECTURE.md) - Complete technical docs

---

## 📖 Reading Order for Different Users

### For End Users (Using the CLI)
1. [README.md](README.md) - Get started
2. [examples.sh](examples.sh) - See it in action
3. [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Bookmark for daily use
4. [USER_MANUAL.md](USER_MANUAL.md) - Read sections as needed

### For Power Users (Advanced Features)
1. [README.md](README.md) - Quick overview
2. [USER_MANUAL.md](USER_MANUAL.md) - Read fully
3. [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Keep open while working
4. [ARCHITECTURE.md](ARCHITECTURE.md) - Understand internals (optional)

### For Developers (Contributing)
1. [README.md](README.md) - Overview
2. [ARCHITECTURE.md](ARCHITECTURE.md) - Read fully
3. [USER_MANUAL.md](USER_MANUAL.md) - Understand features
4. Code exploration in `cli/pkg/`

---

## 🔍 Search by Topic

### Authentication
- [README.md](README.md) - Authentication section
- [USER_MANUAL.md](USER_MANUAL.md) - Getting Started > Authentication

### Task Management
- [README.md](README.md) - Task Management section
- [USER_MANUAL.md](USER_MANUAL.md) - Task Lifecycle Management
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Task Operations

### Instance Management
- [README.md](README.md) - Instance Management section
- [USER_MANUAL.md](USER_MANUAL.md) - Managing Multiple Instances
- [ARCHITECTURE.md](ARCHITECTURE.md) - Multi-Instance Management Flow

### Provider Configuration
- [USER_MANUAL.md](USER_MANUAL.md) - Switching Providers & Models
- [USER_MANUAL.md](USER_MANUAL.md) - Complete Settings Reference

### Modes (Act/Plan/Yolo)
- [README.md](README.md) - Common Workflows
- [USER_MANUAL.md](USER_MANUAL.md) - Working with Modes

### Auto-Approval
- [USER_MANUAL.md](USER_MANUAL.md) - Auto-Approval Configuration

### Settings (-s flags)
- [USER_MANUAL.md](USER_MANUAL.md) - Complete Settings Reference
- [USER_MANUAL.md](USER_MANUAL.md) - Advanced Task Settings

### Troubleshooting
- [README.md](README.md) - Tips & Troubleshooting
- [USER_MANUAL.md](USER_MANUAL.md) - Troubleshooting
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Troubleshooting

---

## 📏 Documentation Stats

| File | Lines | Purpose | Audience |
|------|-------|---------|----------|
| README.md | ~570 | Overview & Getting Started | Everyone |
| USER_MANUAL.md | ~1000 | Complete Usage Guide | End Users |
| QUICK_REFERENCE.md | ~400 | Command Cheat Sheet | Daily Users |
| ARCHITECTURE.md | ~800 | Technical Architecture | Developers |
| examples.sh | ~80 | Runnable Examples | Everyone |

**Total**: ~2,850 lines of documentation

---

## 🤝 Contributing to Documentation

Found an error or want to improve the docs?

1. **Typos/Errors**: Open an issue or submit a PR
2. **Missing Examples**: Add to USER_MANUAL.md or examples.sh
3. **Architecture Changes**: Update ARCHITECTURE.md
4. **New Commands**: Update all relevant docs

---

## 💡 Tips for Using These Docs

1. **Bookmark QUICK_REFERENCE.md** for daily use
2. **Search with Ctrl+F** - all docs have detailed keywords
3. **Follow links** - docs are interconnected
4. **Try examples** - all code snippets are tested
5. **Check Table of Contents** - quickly jump to sections

---

## 🆘 Still Need Help?

If you can't find what you need:

1. **Search the docs**: Use Ctrl+F across all files
2. **Run examples**: `./cli/examples.sh`
3. **Check help**: `./cli/bin/cline --help`
4. **Command help**: `./cli/bin/cline [command] --help`
5. **Verbose mode**: Add `-v` flag for debugging

---

**Happy Coding with Cline! 🚀**

Last Updated: October 2025

