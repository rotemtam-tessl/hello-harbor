# Harbor Framework - Complete Agent Reference Guide

**Source:** [Harbor Documentation](https://harborframework.com/docs)

Harbor is a framework for evaluating and optimizing agents and models in container environments. It provides simple, modular interfaces for environments, agents, and tasks, with integrations for cloud sandbox providers like Daytona, Modal, and E2B.

---

## Table of Contents

1. [Installation](#installation)
2. [Core Concepts](#core-concepts)
3. [CLI Commands Reference](#cli-commands-reference)
4. [Running Evaluations](#running-evaluations)
5. [Task Structure](#task-structure)
6. [Creating New Tasks & Benchmarks](#creating-new-tasks--benchmarks)
7. [Agents](#agents)
8. [Datasets](#datasets)
9. [Cloud Deployments (Daytona)](#cloud-deployments-daytona)
10. [Metrics](#metrics)
11. [Use Cases: Evals, RL, SFT](#use-cases-evals-rl-sft)
12. [Adapters](#adapters)

---

## Installation

### Using uv (recommended)

```bash
uv pip install harbor
```

### Using pip

```bash
pip install harbor
```

---

## Core Concepts

### Task

A task is a single instruction, container environment, and test script. Tasks are used to evaluate agents and models. A task is implemented as a directory of files in the Harbor task format.

### Dataset

A dataset is a collection of tasks. Datasets are used to evaluate agents and models. Usually, a dataset corresponds to a benchmark (e.g., Terminal-Bench, SWE-Bench Verified, etc.). Datasets can optionally define metrics for reward aggregation.

### Agent

An agent is a program that completes tasks. Agents are defined by implementing the `BaseAgent` or `BaseInstalledAgent` interfaces.

### Container Environment

Environments in Harbor are containers, typically defined as Docker images using a Dockerfile. The `BaseEnvironment` interface provides a unified interface for interacting with environments. Many cloud sandbox providers implement this interface.

### Trial

A trial is an agent's attempt at completing a task. Trials can be configured using the `TrialConfig` class. Essentially, a trial is a rollout that produces a reward.

### Job

A job is a collection of trials. Jobs are used to evaluate agents and models. A job can consist of multiple datasets, agents, tasks, and models. Jobs can be configured using the `JobConfig` class.

---

## CLI Commands Reference

### List Available Commands

```bash
harbor --help
```

### List Available Agents

```bash
harbor run --help
```

Check the `--agent` flag to see all available agents. Harbor includes:

- Terminus-2
- Claude Code
- Codex CLI
- Gemini CLI
- OpenHands
- Mini-SWE-Agent
- And more...

### List Available Datasets

```bash
harbor datasets list
```

### Run a Job from Config File

```bash
harbor run --config path/to/job.yaml
```

---

## Running Evaluations

### Running a Registered Dataset

```bash
harbor run --agent codex --model gpt-4o --dataset terminal-bench:v1.0
```

Harbor will automatically download registered datasets.

### Running a Local Dataset

```bash
harbor run --agent codex --model gpt-4o --dataset-path ./my-local-dataset/
```

### Running on SWE-Bench Verified

```bash
harbor run --agent codex --model gpt-4o --dataset swe-bench-verified
```

If you leave off the version, Harbor will use the latest version of the dataset.

### Running a Cloud Sandbox (Daytona)

```bash
harbor run \
  --agent codex \
  --model gpt-4o \
  --dataset terminal-bench:v1.0 \
  --environment daytona \
  --workers 100
```

Running with cloud sandboxes shifts command execution to the cloud, making trials I/O bounded rather than compute bounded. You can parallelize far above your CPU count (e.g., 100 trials in parallel on a MacBook Pro with 14 cores).

---

## Task Structure

### Creating a New Task

```bash
harbor task create my-task
```

This creates a new task directory with the following structure:

```
my-task/
├── task.toml           # Task configuration and metadata
├── instruction.md      # Task instructions (markdown)
├── environment/        # Container environment definition
│   └── Dockerfile      # Docker image definition
├── solution/           # Optional: Oracle solution
│   └── solve.sh        # Solution script
└── tests/              # Test verification
    └── test.sh         # Test script
```

### Running a Task

```bash
harbor run --agent codex --model gpt-4o --dataset-path ./my-task
```

### Task Configuration (task.toml)

```toml
version = "1.0"

[metadata]
name = "My Task"
description = "A sample task"
difficulty = "medium"
source = "https://github.com/example/repo"

[verifier]
timeout_sec = 300

[agent]
timeout_sec = 600

[environment]
build_timeout_sec = 600
docker_image = null  # Optional: use pre-built image
cpu = 2
memory_mb = 4096
storage_mb = 10240
```

**Configuration Parameters:**

- `version`: Task format version (string)
- `metadata`: Arbitrary metadata object
- `verifier.timeout_sec`: Timeout for test verification (number)
- `agent.timeout_sec`: Timeout for agent execution (number)
- `environment.build_timeout_sec`: Timeout for building Docker image (number)
- `environment.docker_image`: Pre-built Docker image or null (string|null)
- `environment.cpu`: CPU cores (integer)
- `environment.memory_mb`: Memory in MB (integer)
- `environment.storage_mb`: Storage in MB (integer)
- `source`: Source URL (string|null)

### Instruction File (instruction.md)

The instruction is a markdown file that contains the task's instruction. This is what the agent receives as input.

```markdown
# Task: Fix the Bug

The file `src/main.py` has a bug on line 42. Fix the bug so that the tests pass.

## Requirements

- Do not modify any other files
- Ensure all existing tests still pass
```

### Environment Directory

The environment definition is placed in an `environment/` folder. Harbor does not require any specific file to exist in that directory. Which file is required depends on the environment type being used.

**Special paths in the environment's filesystem:**

- `/logs/` - Downloaded to the host after the agent/verifier run; useful for debugging and analysis

### Solution Directory (Optional)

The solution folder must contain a `solution/solve.sh` script. Other dependencies are allowed. This folder is copied to `/oracle` by the Oracle agent at runtime and executed from the working directory.

If no solution is provided, the Oracle agent cannot be used to sanity check the task.

### Tests Directory

The tests folder must contain a `tests/test.sh` script. The test script should install test dependencies and verify the agent completed the instruction.

**Important:** The test script must produce a reward file in the `/logs/verifier/` directory.

**Reward file options:**

- `reward.txt` - Plain text with numeric reward
- `reward.json` - JSON format with additional metadata

Harbor reads `reward.txt` by default and falls back to `reward.json`.

**Example test.sh:**

```bash
#!/bin/bash
set -e

# Run tests
cd /workspace
python -m pytest tests/ -v

# Write reward based on exit code
if [ $? -eq 0 ]; then
    echo "1.0" > /logs/verifier/reward.txt
else
    echo "0.0" > /logs/verifier/reward.txt
fi
```

**Alternative using pytest result:**

```bash
#!/bin/bash
cd /workspace
python -m pytest tests/ -v
# Harbor can determine reward from exit code
if [ $? -eq 0 ]; then
    echo "1" > /logs/verifier/reward.txt
else
    echo "0" > /logs/verifier/reward.txt
fi
```

---

## Creating New Tasks & Benchmarks

### Step-by-Step Task Creation

1. **Create the task directory:**

   ```bash
   harbor task create my-new-task
   ```

2. **Write the instruction** in `instruction.md`

3. **Create the Dockerfile** in `environment/Dockerfile`:

   ```dockerfile
   FROM python:3.11-slim
   WORKDIR /workspace
   COPY . .
   RUN pip install -r requirements.txt
   ```

4. **Write the test script** in `tests/test.sh`:

   ```bash
   #!/bin/bash
   cd /workspace
   python -m pytest tests/ -v
   echo $? > /logs/verifier/reward.txt
   ```

5. **Optionally add a solution** in `solution/solve.sh`

6. **Test locally:**
   ```bash
   harbor run --agent oracle --dataset-path ./my-new-task
   ```

### Creating a Dataset (Collection of Tasks)

A dataset is simply a directory containing multiple task directories:

```
my-dataset/
├── task-001/
│   ├── task.toml
│   ├── instruction.md
│   ├── environment/
│   ├── solution/
│   └── tests/
├── task-002/
│   └── ...
└── task-003/
    └── ...
```

Run the entire dataset:

```bash
harbor run --agent codex --model gpt-4o --dataset-path ./my-dataset/
```

---

## Agents

### Pre-integrated Agents

Harbor comes with most popular agents pre-integrated:

- **Terminus-2**
- **Claude Code**
- **Codex CLI**
- **Gemini CLI**
- **OpenHands**
- **Mini-SWE-Agent**

### Types of Agents

1. **External agents**: Interface with the environment through the `BaseEnvironment` interface, typically by executing bash commands via the `exec` method.

2. **Installed agents**: Installed directly into the container environment and executed in headless mode. This comes with the advantage of bringing up the agent's native tooling.

### Implementing an External Agent

```python
from harbor.agents import BaseAgent
from harbor.environments import BaseEnvironment

class MyAgent(BaseAgent):
    async def run(
        self,
        instruction: str,
        environment: BaseEnvironment,
        config: dict
    ) -> dict:
        # Execute commands in the environment
        result = await environment.exec("ls -la")

        # Process and return result
        return {
            "status": "completed",
            "output": result.stdout
        }
```

### Implementing an Installed Agent

```python
from harbor.agents import BaseInstalledAgent

class MyInstalledAgent(BaseInstalledAgent):
    def get_install_script(self) -> str:
        return """
        pip install my-agent-package
        """

    def get_run_command(self, instruction: str) -> str:
        return f"my-agent run --instruction '{instruction}'"
```

### Running a Custom Agent

```bash
harbor run \
  --agent-path ./my_agent.py \
  --agent-class MyAgent \
  --model gpt-4o \
  --dataset-path ./my-dataset
```

---

## Datasets

### Local Datasets

A local dataset is a directory that contains a set of tasks:

```bash
harbor run --agent codex --model gpt-4o --dataset-path ./my-local-dataset/
```

### Harbor Registry

Harbor comes with a default registry defined in a `registry.json` file stored in the repository root.

**Using registered datasets:**

```bash
harbor run --agent codex --model gpt-4o --dataset terminal-bench:v1.0
```

**Dataset registry structure:**

```json
{
  "datasets": {
    "terminal-bench": {
      "v1.0": {
        "tasks": [
          {
            "name": "task-001",
            "git_url": "https://github.com/harbor/datasets",
            "git_commit_id": "abc123",
            "path": "datasets/terminal-bench/task-001"
          }
        ],
        "metrics": ["mean"]
      }
    }
  }
}
```

### Custom Registry

Create your own `registry.json` file and use it:

```bash
harbor run \
  --agent codex \
  --model gpt-4o \
  --dataset my-custom-dataset:v1.0 \
  --registry-path ./my-registry.json
```

Or host at a URL:

```bash
harbor run \
  --agent codex \
  --model gpt-4o \
  --dataset my-custom-dataset:v1.0 \
  --registry-url https://example.com/registry.json
```

---

## Cloud Deployments (Daytona)

### Why Use Cloud Sandboxes?

Containerized agentic tasks can be slow when performing rollouts due to:

- Container startup and teardown overhead
- Waiting for LLM API calls
- Waiting for command execution

Using a cloud sandbox provider shifts command execution to the cloud, making trials I/O bounded rather than compute bounded.

### Supported Providers

Harbor supports:

- **Daytona** (recommended) - Most flexible
- **Modal**
- **E2B**

### Running with Daytona

```bash
harbor run \
  --agent codex \
  --model gpt-4o \
  --dataset terminal-bench:v1.0 \
  --environment daytona \
  --workers 100
```

You can run up to 100 trials in parallel on a MacBook Pro with 14 cores.

### Daytona Configuration

Set your Daytona API credentials:

```bash
export DAYTONA_API_KEY="your-api-key"
export DAYTONA_API_URL="https://api.daytona.io"
```

### Limitations

All cloud sandbox providers currently do not support multi-container environments. However, the Docker environment still supports multi-container tasks - just include an `environment/docker-compose.yaml` file in your task definition.

---

## Metrics

### Available Metrics

Harbor provides built-in metrics:

- `sum`
- `min`
- `max`
- `mean` (most common - accuracy)

### Defining Metrics for Datasets

In `registry.json`:

```json
{
  "datasets": {
    "my-dataset": {
      "v1.0": {
        "tasks": [...],
        "metrics": ["mean", "sum", "custom:path/to/script.py"]
      }
    }
  }
}
```

### Custom Metrics

Create a Python script that accepts two flags:

```python
#!/usr/bin/env python3
# /// script
# dependencies = ["numpy"]
# ///

import argparse
import json

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", required=True)
    parser.add_argument("--output", required=True)
    args = parser.parse_args()

    with open(args.input) as f:
        rewards = json.load(f)

    # Calculate custom metric
    metric_value = sum(r["reward"] for r in rewards) / len(rewards)

    with open(args.output, "w") as f:
        json.dump({"accuracy": metric_value}, f)

if __name__ == "__main__":
    main()
```

Run with:

```bash
uv run custom_metric.py --input input.jsonl --output output.json
```

---

## Use Cases: Evals, RL, SFT

### Evals

Harbor is built by the creators of Terminal-Bench with evals as a core use case.

**Running an evaluation:**

```bash
harbor run --agent codex --model gpt-4o --dataset terminal-bench:v1.0
```

**Analyzing results:**

Results are stored in the `jobs/` directory:

```
jobs/
└── job-2024-01-01-123456/
    ├── config.json
    ├── results.json
    ├── trials/
    │   ├── task-001/
    │   │   ├── logs/
    │   │   ├── result.json
    │   │   └── trajectory.json
    │   └── task-002/
    │       └── ...
    └── metrics.json
```

### Reinforcement Learning (RL)

Harbor provides a simple interface for training on containerized agentic environments.

**Using TrialConfig for rollouts:**

```python
from harbor.core import TrialConfig, TaskConfig

# Configure trials for RL
trial_config = TrialConfig(
    task=TaskConfig(
        name="my-task",
        git_url="https://github.com/harbor/datasets",
        path="datasets/terminal-bench/task-001"
    ),
    agent="codex",
    model="gpt-4o"
)
```

**Example rollout interface:**

```python
from harbor.environments import DaytonaEnvironment
from harbor.agents import get_agent

async def rollout(trial_configs: list[TrialConfig]):
    results = []

    for config in trial_configs:
        env = DaytonaEnvironment(config.task)
        agent = get_agent(config.agent, config.model)

        async with env:
            result = await agent.run(
                instruction=config.task.instruction,
                environment=env,
                config=config
            )
            results.append({
                "reward": result.reward,
                "tokens": result.tokens,
                "trajectory": result.trajectory
            })

    return results
```

**Collecting tokens from agents:**

1. **Intercepting tokens from a vLLM server** - Assumes you're using a vLLM server and have a framework that handles interception
2. **Returning tokens as part of agent result metadata** - Configure your agent to return tokens in the result

### Supervised Fine-Tuning (SFT)

Harbor includes utilities for turning trials into conversational traces for SFT.

**Export traces from a job:**

```bash
harbor traces export ./jobs/job-2024-01-01-123456/ \
  --sharegpt \
  --filter success \
  --push \
  --repo my-org/my-sft-dataset
```

**Key options:**

| Flag                           | Description                                      |
| ------------------------------ | ------------------------------------------------ |
| `--episode`                    | `all` or `last` - which episodes to export       |
| `--sharegpt` / `--no-sharegpt` | Add ShareGPT-style column for instruction-tuning |
| `--filter`                     | `success`, `failure`, or `all`                   |
| `--push`                       | Upload to Hugging Face Hub                       |
| `--verbose`                    | Enable verbose output                            |

**Auto-export during run:**

```bash
harbor run \
  --agent codex \
  --model gpt-4o \
  --dataset terminal-bench:v1.0 \
  --export-traces \
  --export-filter success \
  --export-repo my-org/my-dataset
```

**Export split datasets (success/failure):**

```bash
harbor sweeps run \
  --push \
  --export-repo-success my-org/success-traces \
  --export-repo-failure my-org/failure-traces
```

**Local export to Parquet:**

```python
from harbor.traces import export_traces

export_traces(
    job_dir="./jobs/job-2024-01-01-123456/",
    output_path="./traces.parquet",
    sharegpt=True,
    filter="success"
)
```

---

## Adapters

Adapters are used to convert existing benchmarks into the Harbor task format.

### Overview

Harbor supports running various benchmarks via a unified interface. SWE-Bench, LiveCodeBench, and more benchmarks are integrated.

### Adapter Development Workflow

1. **Fork and prepare the original benchmark**

   - Fork the original benchmark's repository
   - Implement CLI agents and models that are supported by Harbor
   - Establish a baseline by running the original benchmark
   - Document the process

2. **Create the Harbor Adapter**

   - Fork the Harbor repository
   - Create adapter under `adapters/{adapter-name}/`
   - Develop the adapter code to convert tasks to Harbor format
   - Verify oracle solutions pass with 100% reward

3. **Run Parity Experiments**

   - Compare results between original benchmark and Harbor adapter

4. **Register the Dataset**

   - Add tasks to the harbor-datasets repository
   - Update registry.json in Harbor repo

5. **Document and Submit**
   - Create README.md for your adapter
   - Submit pull request

### Adapter Directory Structure

```
harbor/adapters/<adapter-name>/
├── adapter.py          # Main adapter code
├── run_adapter.py      # CLI entry point
├── template/           # Task template
│   ├── environment/
│   ├── solution/
│   └── tests/
├── parity_experiment.json
└── README.md
```

### Running Harbor Harness

**Option 1: Individual trials (testing single tasks)**

```bash
harbor run --agent oracle --dataset-path ./adapters/my-adapter/output/task-001
```

**Option 2: Jobs with local dataset path**

```bash
harbor run --agent codex --model gpt-4o --dataset-path ./adapters/my-adapter/output/
```

**Option 3: Jobs with configuration files**

```bash
harbor run --config ./examples/config/my-adapter-config.yaml
```

**Option 4: Using local registry after dataset PR is merged**

```bash
harbor run \
  --agent codex \
  --model gpt-4o \
  --dataset my-adapter:v1.0 \
  --registry-path ./registry.json
```

### Key Migration Steps from Terminal-Bench

| Aspect        | Terminal-Bench | Harbor                      |
| ------------- | -------------- | --------------------------- |
| Config format | `task.yaml`    | `task.toml`                 |
| Instruction   | Inline in YAML | Separate `instruction.md`   |
| Dockerfile    | Root level     | `environment/Dockerfile`    |
| Solution      | `solution.sh`  | `solution/solve.sh`         |
| Test script   | `run-tests.sh` | `tests/test.sh`             |
| Output        | Exit code      | `/logs/verifier/reward.txt` |

---

## Quick Reference

### Common Commands

```bash
# Install Harbor
uv pip install harbor

# List commands
harbor --help

# List datasets
harbor datasets list

# Run eval on registered dataset
harbor run --agent codex --model gpt-4o --dataset terminal-bench:v1.0

# Run eval on local dataset
harbor run --agent codex --model gpt-4o --dataset-path ./my-dataset/

# Run with Daytona cloud
harbor run --agent codex --model gpt-4o --dataset terminal-bench:v1.0 --environment daytona --workers 100

# Create new task
harbor task create my-new-task

# Export SFT traces
harbor traces export ./jobs/job-id/ --sharegpt --filter success
```

### Environment Variables

```bash
# Daytona
export DAYTONA_API_KEY="your-api-key"
export DAYTONA_API_URL="https://api.daytona.io"

# Model providers
export OPENAI_API_KEY="your-key"
export ANTHROPIC_API_KEY="your-key"
```

### File Structure Summary

```
task/
├── task.toml           # Config & metadata
├── instruction.md      # Agent instructions
├── environment/        # Container setup
│   └── Dockerfile
├── solution/           # Optional oracle
│   └── solve.sh
└── tests/              # Verification
    └── test.sh         # Must write to /logs/verifier/reward.txt
```

---

## Resources

- **Documentation**: https://harborframework.com/docs
- **GitHub**: https://github.com/harbor/harbor
- **Discord**: Join for support and discussions
- **Registry**: https://github.com/harbor/harbor-datasets

---

_This guide was compiled from the official Harbor documentation to enable coding agents to operate Harbor locally, use cloud deployments with Daytona, and create new scenarios and benchmarks._
