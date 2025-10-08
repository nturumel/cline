# Cline CLI Architecture

## Overview

The Cline CLI is a terminal-based interface for interacting with Cline AI. It consists of multiple components working together to manage instances, handle tasks, communicate with Cline Core, and provide filesystem access through a host bridge.

---

## High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         USER / TERMINAL                                  │
│                    (Executes: ./cli/bin/cline)                           │
└───────────────────────────────┬─────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         CLI ENTRY POINT                                  │
│                      cmd/cline/main.go                                   │
│                                                                           │
│  - Cobra Command Router                                                  │
│  - Global Flag Handling (--verbose, --output-format, --address)         │
│  - Command Registration (task, instance, auth, version)                 │
└───────────────────────────────┬─────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      GLOBAL INITIALIZATION                               │
│                    pkg/cli/global/global.go                              │
│                                                                           │
│  - Initialize Config (~/.cline directory)                                │
│  - Create ClineClients instance                                          │
│  - Setup ClientRegistry                                                  │
└─────┬───────────────────────────────────────────────────────────────────┘
      │
      │
      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        COMMAND LAYER                                     │
│                     pkg/cli/*.go                                         │
│                                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌────────────┐ │
│  │ task.go      │  │ instances.go │  │ auth.go      │  │ version.go │ │
│  │              │  │              │  │              │  │            │ │
│  │ • new        │  │ • list       │  │ • auth flow  │  │ • version  │ │
│  │ • oneshot    │  │ • new        │  │              │  │   info     │ │
│  │ • follow     │  │ • use        │  │              │  │            │ │
│  │ • view       │  │ • kill       │  │              │  │            │ │
│  │ • list       │  │              │  │              │  │            │ │
│  │ • send       │  │              │  │              │  │            │ │
│  │ • cancel     │  │              │  │              │  │            │ │
│  │ • resume     │  │              │  │              │  │            │ │
│  │ • restore    │  │              │  │              │  │            │ │
│  └──────┬───────┘  └──────┬───────┘  └──────────────┘  └────────────┘ │
│         │                 │                                             │
└─────────┼─────────────────┼─────────────────────────────────────────────┘
          │                 │
          ▼                 ▼
┌──────────────────────┐  ┌──────────────────────────────────────────────┐
│   TASK MANAGER       │  │      CLIENT REGISTRY                         │
│  pkg/cli/task/       │  │    pkg/cli/global/registry.go                │
│  manager.go          │  │                                              │
│                      │  │  - SQLite-backed instance tracking          │
│  • CreateTask()      │  │  - Default instance management              │
│  • FollowConversation│  │  - Client connection pooling                │
│  • SendMessage()     │  │  - Health check monitoring                  │
│  • CancelTask()      │  │  - Port allocation                          │
│  • ResumeTask()      │  │  - Instance lifecycle                       │
│  • RestoreCheckpoint │  │                                              │
│  • ListTasks()       │  │  Methods:                                    │
│  • SetMode()         │  │  • GetClient()                               │
│                      │  │  • GetDefaultClient()                        │
│  Contains:           │  │  • ListInstancesCleaned()                    │
│  • ConversationState │  │  • SetDefaultInstance()                      │
│  • Renderer          │  │  • HasInstanceAtAddress()                    │
│  • StreamingDisplay  │  │                                              │
│  • HandlerRegistry   │  └──────────────┬───────────────────────────────┘
└──────────┬───────────┘                 │
           │                             │
           │                             ▼
           │              ┌──────────────────────────────────────────────┐
           │              │      CLINE CLIENTS                           │
           │              │   pkg/cli/global/cline-clients.go            │
           │              │                                              │
           │              │  • StartNewInstance()                        │
           │              │  • EnsureInstanceAtAddress()                 │
           │              │  • GetRegistry()                             │
           │              │                                              │
           │              │  Manages:                                    │
           │              │  • cline-host process lifecycle              │
           │              │  • cline-core process lifecycle              │
           │              │  • Port allocation (core + host bridge)     │
           │              │  • Process spawning & monitoring             │
           │              └──────────────┬───────────────────────────────┘
           │                             │
           ▼                             ▼
┌─────────────────────────┐   ┌──────────────────────────────────────────┐
│  DISPLAY SYSTEM         │   │     PERSISTENCE LAYER                    │
│  pkg/cli/display/       │   │   pkg/cli/sqlite/                        │
│                         │   │                                          │
│ ┌────────────────────┐  │   │  SQLite Database: ~/.cline/instances.db │
│ │ renderer.go        │  │   │                                          │
│ │ • RenderAsk()      │  │   │  Tables:                                │
│ │ • RenderSay()      │  │   │  • instances                            │
│ │ • RenderInfo()     │  │   │    - address (primary key)              │
│ │ • RenderError()    │  │   │    - cline_core_port                    │
│ │ • RenderDebug()    │  │   │    - host_bridge_port                   │
│ │ • RenderJSON()     │  │   │    - pid                                │
│ └────────────────────┘  │   │    - started_at                         │
│                         │   │    - version                            │
│ ┌────────────────────┐  │   │    - last_seen                          │
│ │ streaming.go       │  │   │                                          │
│ │ • ProcessMessages()│  │   │  • default_instance.json                │
│ │ • RenderStream()   │  │   │                                          │
│ └────────────────────┘  │   │  Methods:                                │
│                         │   │  • NewLockManager()                      │
│ ┌────────────────────┐  │   │  • GetInstanceInfo()                    │
│ │ typewriter.go      │  │   │  • HasInstanceAtAddress()               │
│ │ • TypewriterEffect │  │   │  • ListInstances()                      │
│ └────────────────────┘  │   │  • UpsertInstance()                     │
│                         │   │  • DeleteInstance()                     │
│ ┌────────────────────┐  │   │                                          │
│ │ deduplicator.go    │  │   └──────────────────────────────────────────┘
│ │ • Dedupe Messages  │  │
│ └────────────────────┘  │
└─────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      MESSAGE HANDLERS                                    │
│                   pkg/cli/handlers/                                      │
│                                                                           │
│  ┌──────────────────────┐  ┌──────────────────────┐                    │
│  │ ask_handlers.go      │  │ say_handlers.go      │                    │
│  │                      │  │                      │                    │
│  │ Handles:             │  │ Handles:             │                    │
│  │ • tool               │  │ • text               │                    │
│  │ • command            │  │ • tool_result        │                    │
│  │ • completion_result  │  │ • error              │                    │
│  │ • api_req_started    │  │ • completion         │                    │
│  │ • api_req_finished   │  │ • api_req_started    │                    │
│  │                      │  │ • api_req_finished   │                    │
│  └──────────────────────┘  └──────────────────────┘                    │
│                                                                           │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │ handler.go - HandlerRegistry                                      │  │
│  │ • Register handlers                                               │  │
│  │ • Dispatch messages to appropriate handler                        │  │
│  └──────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    SETTINGS PARSER                                       │
│               pkg/cli/task/settings_parser.go                            │
│                                                                           │
│  Parses -s key=value flags into TaskSettings protobuf                   │
│                                                                           │
│  Supports 100+ settings including:                                       │
│  • Provider configuration (act_mode_api_provider, etc.)                  │
│  • Model selection (act_mode_api_model_id, etc.)                         │
│  • Cloud provider settings (AWS, GCP, Azure)                             │
│  • Auto-approval configuration                                           │
│  • Browser settings                                                      │
│  • Mode settings (yolo, strict_plan, etc.)                               │
│  • Performance tuning (timeouts, token budgets)                          │
└─────────────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      GRPC CLIENT LAYER                                   │
│              src/generated/grpc-go/client/                               │
│                                                                           │
│  Generated from protobuf definitions                                     │
│                                                                           │
│  ┌────────────────────┐  ┌────────────────────┐                         │
│  │ cline_client.go    │  │ connection.go      │                         │
│  │                    │  │                    │                         │
│  │ ClineClient struct │  │ • Connect()        │                         │
│  │ with services:     │  │ • Disconnect()     │                         │
│  │                    │  │ • Health checks    │                         │
│  │ • Account          │  └────────────────────┘                         │
│  │ • Browser          │                                                  │
│  │ • Checkpoints      │  ┌──────────────────────────────────────────┐  │
│  │ • Commands         │  │ services/*_client.go                     │  │
│  │ • Dictation        │  │                                          │  │
│  │ • File             │  │ 20 specialized service clients:          │  │
│  │ • MCP              │  │ • account_client.go                      │  │
│  │ • Models           │  │ • browser_client.go                      │  │
│  │ • OCAAccount       │  │ • task_client.go                         │  │
│  │ • Slash            │  │ • state_client.go                        │  │
│  │ • State            │  │ • workspace_client.go                    │  │
│  │ • Task             │  │ • window_client.go                       │  │
│  │ • UI               │  │ • ... and more                           │  │
│  │ • Web              │  │                                          │  │
│  └────────────────────┘  └──────────────────────────────────────────┘  │
└───────────────────────────────┬─────────────────────────────────────────┘
                                │
                                │ gRPC calls over localhost TCP
                                │
     ┌──────────────────────────┴──────────────────────────────┐
     │                                                           │
     ▼                                                           ▼
┌──────────────────────────────────┐    ┌────────────────────────────────┐
│     CLINE CORE INSTANCE          │    │    HOST BRIDGE INSTANCE        │
│   (Node.js Process)              │    │   (Go Process: cline-host)     │
│                                  │    │                                │
│  Port: 50052 (default)           │    │  Port: 50102 (default)         │
│                                  │    │                                │
│  Responsibilities:               │    │  cmd/cline-host/main.go        │
│  • AI model communication        │    │  pkg/hostbridge/               │
│  • Task orchestration            │    │                                │
│  • Conversation management       │    │  Implements gRPC services:     │
│  • Tool execution coordination   │    │  • DiffService                 │
│  • State persistence             │    │  • EnvService                  │
│  • UI state synchronization      │    │  • WindowService               │
│  • MCP server management         │    │  • WorkspaceService            │
│                                  │    │                                │
│  Communicates with:              │    │  Provides:                     │
│  • AI providers (Claude, GPT...) │    │  • File system operations      │
│  • Host Bridge (for file ops)    │    │  • Process execution           │
│  • Task Manager (via gRPC)       │    │  • Environment variable access │
│                                  │    │  • Diff generation             │
│  Started by: ClineClients        │    │  • Window management           │
│  Process: node cline-core.js     │    │  • Workspace queries           │
│  Config: ~/.cline/               │    │                                │
│                                  │    │  Started by: ClineClients      │
│                                  │    │  Process: ./cli/bin/cline-host │
└──────────────────────────────────┘    └────────────────────────────────┘
          │                                         │
          │                                         │
          ▼                                         ▼
┌──────────────────────────────────┐    ┌────────────────────────────────┐
│   AI PROVIDERS                   │    │    LOCAL FILESYSTEM            │
│                                  │    │                                │
│  • Anthropic (Claude)            │    │  ~/.cline/                     │
│  • OpenAI (GPT)                  │    │  ├── instances.db              │
│  • Google (Gemini, Vertex)       │    │  ├── default_instance.json    │
│  • AWS Bedrock                   │    │  ├── logs/                     │
│  • Local (Ollama, LM Studio)     │    │  └── tasks/                    │
│  • 40+ other providers           │    │                                │
│                                  │    │  /tmp/                         │
└──────────────────────────────────┘    │  └── cline-core-debug-*.log    │
                                        └────────────────────────────────┘
```

---

## Component Details

### 1. CLI Entry Point (`cmd/cline/main.go`)

**Responsibilities:**
- Parse command-line arguments
- Route commands to appropriate handlers
- Initialize global configuration
- Handle top-level flags (--verbose, --output-format, --address)

**Commands Registered:**
- `task` - Task management
- `instance` - Instance management
- `auth` - Authentication
- `version` - Version information
- `send` - Send messages (alias to task send)

---

### 2. Global Layer (`pkg/cli/global/`)

#### `global.go`
- **GlobalConfig**: Stores CLI configuration
- **InitializeGlobalConfig()**: Sets up ~/.cline directory
- **GetDefaultClient()**: Returns gRPC client for default instance
- **EnsureDefaultInstance()**: Creates instance if none exists

#### `cline-clients.go`
- **ClineClients**: Manages multiple Cline instances
- **StartNewInstance()**: Spawns cline-core and cline-host processes
- **EnsureInstanceAtAddress()**: Ensures instance exists at specific address
- **startClineHost()**: Spawns cline-host process
- **startClineCore()**: Spawns cline-core process with Node.js

#### `registry.go`
- **ClientRegistry**: Manages gRPC client connections
- **SQLite integration**: Persists instance information
- **Health monitoring**: Tracks instance status
- **Connection pooling**: Reuses gRPC connections
- **Default instance**: Manages which instance is active

---

### 3. Task Management (`pkg/cli/task/`)

#### `manager.go`
- **Manager**: Orchestrates task execution
- **CreateTask()**: Creates new AI tasks
- **FollowConversation()**: Streams task messages
- **SendMessage()**: Sends followup messages
- **CancelTask()**: Cancels running tasks
- **ResumeTask()**: Resumes previous tasks
- **RestoreCheckpoint()**: Restores from checkpoints
- **SetMode()**: Switches between act/plan modes

#### `settings_parser.go`
- **ParseTaskSettings()**: Parses -s flags
- Supports 100+ settings
- Handles nested settings (e.g., auto-approval-settings.actions.read-files)
- Validates provider names, model IDs, etc.

#### `stream_coordinator.go`
- **StreamCoordinator**: Manages message streaming
- Coordinates between gRPC streams and display
- Handles message ordering
- Manages stream lifecycle

---

### 4. Display System (`pkg/cli/display/`)

#### `renderer.go`
- **Renderer**: Formats output for terminal
- Output formats: rich, plain, JSON
- Color and styling support
- Error/warning/info/debug rendering

#### `streaming.go`
- **StreamingDisplay**: Real-time message display
- Message deduplication
- Progress indicators
- Typewriter effects

#### `typewriter.go`
- Animated text output
- Token-by-token streaming effect

#### `deduplicator.go`
- Prevents duplicate messages
- Tracks displayed content
- Efficient message filtering

---

### 5. Message Handlers (`pkg/cli/handlers/`)

#### `handler.go`
- **HandlerRegistry**: Routes messages to handlers
- **Handler interface**: Defines handler contract
- Extensible handler system

#### `ask_handlers.go`
- Handles "ask" message types
- Tool execution requests
- Command execution requests
- API request lifecycle

#### `say_handlers.go`
- Handles "say" message types
- Text responses
- Tool results
- Completion messages
- Error messages

---

### 6. Persistence Layer (`pkg/cli/sqlite/`)

#### `locks.go`
- **LockManager**: SQLite-backed instance registry
- Thread-safe operations
- Instance CRUD operations
- Health status tracking

**Database Schema:**
```sql
CREATE TABLE instances (
    address TEXT PRIMARY KEY,
    cline_core_port INTEGER,
    host_bridge_port INTEGER,
    pid INTEGER,
    started_at INTEGER,
    version TEXT,
    last_seen INTEGER
);
```

**Files:**
- `~/.cline/instances.db` - Instance registry
- `~/.cline/default_instance.json` - Default instance config

---

### 7. gRPC Client Layer (`src/generated/grpc-go/`)

**Auto-generated from protobuf definitions**

#### `client/cline_client.go`
- **ClineClient**: Main client struct
- Aggregates all service clients
- Manages connection lifecycle

#### `client/connection.go`
- Connection management
- Dial logic
- Health checks
- Reconnection handling

#### `client/services/`
20 specialized service clients:
- `account_client.go` - Account management
- `task_client.go` - Task operations
- `state_client.go` - State queries
- `workspace_client.go` - Workspace operations
- `window_client.go` - Window management
- `browser_client.go` - Browser automation
- `checkpoints_client.go` - Checkpoint management
- `commands_client.go` - Command execution
- `file_client.go` - File operations
- `mcp_client.go` - MCP server management
- And 10 more...

---

### 8. Host Bridge (`pkg/hostbridge/`)

**Provides filesystem and system access to Cline Core**

#### `grpc_server.go`
- **GrpcServer**: Main server implementation
- Health check support
- Service registration

#### `simple_workspace.go`
- **SimpleWorkspace**: Workspace operations
- File listing
- Directory operations
- Search functionality

#### `env.go`
- **EnvService**: Environment variable access
- System environment queries
- Path resolution

#### `diff.go`
- **DiffService**: Diff generation
- File comparison
- Patch application

#### `window.go`
- **WindowService**: Window management
- Terminal operations
- Process management

---

### 9. Command Layer (`pkg/cli/`)

#### `task.go`
Commands:
- `task new` - Create task
- `task oneshot` - Autonomous task
- `task follow` - Stream conversation
- `task view` - View conversation
- `task list` - List history
- `task send` - Send message
- `task cancel` - Cancel task
- `task resume` - Resume task
- `task restore` - Restore checkpoint

#### `instances.go`
Commands:
- `instance list` - List instances
- `instance new` - Create instance
- `instance use` - Switch default
- `instance kill` - Kill instance

#### `auth.go`
Commands:
- `auth` - Authenticate with Cline

#### `version.go`
Commands:
- `version` - Show version info

---

## Data Flow Diagrams

### Task Creation Flow

```
User Command
    │
    ▼
cline task new "prompt"
    │
    ▼
Global Init → Registry → Ensure Instance
    │                           │
    │                           ├→ Start cline-host
    │                           └→ Start cline-core
    ▼
Task Manager
    │
    ├→ Parse Settings (settings_parser.go)
    ├→ Create TaskSettings protobuf
    │
    ▼
gRPC Client → Task.NewTask()
    │
    ▼
Cline Core Instance
    │
    ├→ Process prompt
    ├→ Call AI provider
    ├→ Execute tools
    └→ Return task ID
```

### Message Streaming Flow

```
User Command
    │
    ▼
cline task follow
    │
    ▼
Task Manager
    │
    ├→ Get current conversation state
    │
    ▼
gRPC Client → State.StreamConversationState()
    │
    ▼
Stream Coordinator
    │
    ├→ Receive messages
    ├→ Track sequence
    ├→ Handle reconnections
    │
    ▼
Handler Registry
    │
    ├→ Route to AskHandler or SayHandler
    │
    ▼
Display System
    │
    ├→ Deduplicator (remove duplicates)
    ├→ Renderer (format output)
    ├→ StreamingDisplay (typewriter effect)
    │
    ▼
Terminal Output
```

### Multi-Instance Management Flow

```
User: cline instance new
    │
    ▼
ClineClients.StartNewInstance()
    │
    ├→ Allocate ports (find next available)
    ├→ Start cline-host process
    ├→ Start cline-core process
    ├→ Wait for health checks
    │
    ▼
SQLite Registry
    │
    ├→ Insert instance record
    ├→ Update last_seen timestamp
    │
    ▼
Set as default (if first instance)
    │
    ▼
Return instance info to user
```

---

## Key Design Patterns

### 1. **Registry Pattern**
- `ClientRegistry` manages multiple instance connections
- Lazy connection creation
- Health monitoring
- Default instance selection

### 2. **Handler Pattern**
- `HandlerRegistry` routes messages to specialized handlers
- Extensible handler system
- Clean separation of concerns

### 3. **Builder Pattern**
- `TaskSettings` parsed incrementally from flags
- Complex nested settings support
- Validation during construction

### 4. **Streaming Pattern**
- `StreamCoordinator` manages gRPC streams
- Buffering and backpressure handling
- Reconnection logic

### 5. **Repository Pattern**
- `LockManager` abstracts SQLite operations
- Clean data access layer
- Thread-safe operations

---

## Process Architecture

### Single Instance Setup

```
Terminal Process (cline)
    │
    ├→ gRPC Client Connection ──→ cline-core (Node.js) ──→ AI Providers
    │                                    │
    │                                    └──→ cline-host (Go) ──→ Filesystem
    │
    └→ SQLite Registry (~/.cline/instances.db)
```

### Multi-Instance Setup

```
Terminal Process (cline)
    │
    ├→ gRPC Client #1 ──→ cline-core:50052 ──→ cline-host:50102
    │                          (Project A)
    │
    ├→ gRPC Client #2 ──→ cline-core:50053 ──→ cline-host:50103
    │                          (Project B)
    │
    └→ SQLite Registry (tracks both instances)
```

---

## Port Allocation

Default ports:
- **Cline Core**: 50052, 50053, 50054, ... (incremental)
- **Host Bridge**: 50102, 50103, 50104, ... (corePort + 50)

Port allocation algorithm:
1. Check SQLite for existing instances
2. Find highest used port
3. Allocate next available port
4. Register in SQLite before starting process

---

## Configuration Files

### `~/.cline/instances.db`
SQLite database tracking all instances

### `~/.cline/default_instance.json`
```json
{
  "defaultInstance": "localhost:50052"
}
```

### `/tmp/cline-core-debug-localhost-50052.log`
Debug logs from cline-core processes

---

## Error Handling

### Connection Failures
1. Health check fails
2. Mark instance as unhealthy in registry
3. Attempt reconnection (exponential backoff)
4. Clean up stale instances

### Process Crashes
1. gRPC call fails
2. Check process PID
3. If dead, remove from registry
4. User can restart with `instance new`

### Port Conflicts
1. Try to bind port
2. If in use, increment and retry
3. Update registry with actual port used

---

## Performance Considerations

### Connection Pooling
- gRPC connections are reused
- Lazy connection creation
- Connection health monitoring

### Message Deduplication
- Hash-based deduplication in display layer
- Prevents UI flicker
- Reduces terminal I/O

### Streaming Efficiency
- Backpressure handling in StreamCoordinator
- Buffered message processing
- Efficient protobuf serialization

---

## Security Considerations

### Local Communication
- All gRPC over localhost TCP
- No external network exposure
- Process isolation

### Filesystem Access
- Host bridge provides controlled access
- Workspace boundaries enforced
- Path sanitization

### API Keys
- Stored by cline-core (not CLI)
- Managed through auth flow
- Environment variable support

---

## Testing Architecture

### Unit Tests
- `cli/pkg/` packages have unit tests
- Mock gRPC clients
- Table-driven tests

### E2E Tests
- `cli/e2e/` directory
- Tests full workflow
- Instance lifecycle tests
- Concurrent operation tests

---

## Build Process

```
npm run compile-cli
    │
    ├→ npm run protos (generate protobuf code)
    │   ├→ TypeScript protobuf
    │   └→ Go protobuf (src/generated/grpc-go/)
    │
    ├→ go build -o bin/cline ./cmd/cline
    └→ go build -o bin/cline-host ./cmd/cline-host
```

---

## Extension Points

### Adding New Commands
1. Create command in `pkg/cli/`
2. Register in `cmd/cline/main.go`
3. Implement command logic
4. Add tests

### Adding New Handlers
1. Implement `Handler` interface
2. Register in `HandlerRegistry`
3. Handle new message types

### Adding New Settings
1. Add to protobuf definition
2. Update `settings_parser.go`
3. Document in USER_MANUAL.md

---

## Dependencies

### Go Packages
- `github.com/spf13/cobra` - CLI framework
- `google.golang.org/grpc` - gRPC
- `google.golang.org/protobuf` - Protocol Buffers
- `github.com/mattn/go-sqlite3` - SQLite
- `github.com/charmbracelet/*` - Terminal UI

### Generated Code
- `github.com/cline/grpc-go` - Auto-generated gRPC clients
- From proto files in `proto/`

---

## Future Architecture Considerations

### Potential Improvements
1. **Remote instances**: Support non-localhost instances
2. **Distributed registry**: Replace SQLite with distributed store
3. **Load balancing**: Distribute tasks across instances
4. **Clustering**: Coordinate multiple CLI processes
5. **Plugin system**: Dynamic command loading
6. **Advanced streaming**: WebSocket or SSE support

---

## Troubleshooting Architecture

### Common Issues

**"failed to start cline-host"**
- Path issue: CLI looks for `./cli/bin/cline-host`
- Must run from project root

**"connection refused"**
- Instance not running
- Check `instance list`
- Port conflict

**"port already in use"**
- Another process using port
- CLI auto-allocates next port
- Check SQLite for stale entries

---

## Related Documentation

- [README.md](README.md) - Installation and overview
- [USER_MANUAL.md](USER_MANUAL.md) - Complete usage guide
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Command cheat sheet
- [../proto/](../proto/) - Protocol Buffer definitions

---

**Last Updated**: October 2025

