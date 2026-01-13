# hello-harbor

Basic Harbor task examples for blog post.

## Structure

```
tasks/
└── hello-world/          # Simple "Hello, World!" task
    ├── task.toml         # Task configuration
    ├── instruction.md    # Agent instructions
    ├── environment/      # Docker environment
    ├── tests/            # Verification tests
    └── solution/         # Oracle solution

docs/
└── harbor.md             # Harbor framework documentation
```

## Setup

```bash
# Install dependencies with uv
uv sync
```

## Running Tasks

```bash
# System test: verify oracle solution works
uv run harbor run --agent oracle --path ./tasks/hello-world

# Run with an agent
export ANTHROPIC_API_KEY="your-api-key-here"
uv run harbor run --agent claude-code --path ./tasks/hello-world
```

## Documentation

See [docs/harbor.md](docs/harbor.md) for Harbor framework docs for agents.
