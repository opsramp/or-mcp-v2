# OpsRamp AI Agent - Non-Interactive Mode

This document provides information about using the OpsRamp AI Agent in non-interactive modes.

## Available Modes

The OpsRamp AI Agent can be run in several different modes:

1. **Interactive Chat** - The default mode where you have a conversation with the agent
2. **Single Prompt** - Process a single prompt and then exit
3. **Batch Processing** - Process multiple prompts from a file
4. **Simple Mode** - Run the agent without connecting to an actual MCP server (for testing/development)

## Setup

Before using any mode, make sure you've set up the agent by installing the required dependencies:

```bash
# From the client directory
make setup
```

## Running in Single Prompt Mode

To run the agent with a single prompt:

```bash
# Using the Makefile
make run-prompt PROMPT="your prompt here"

# Or directly
python examples/chat_client.py --prompt "your prompt here"
```

### Options

- `--server-url` - MCP server URL (default: http://localhost:8080)
- `--llm-provider` - LLM provider to use (default: openai, options: openai, anthropic)
- `--env-file` - Path to a specific .env file for configuration
- `--json` - Output in JSON format
- `--simple-mode` - Run without connecting to an MCP server

Example with options:

```bash
make run-prompt PROMPT="List all integrations" SIMPLE_MODE=true
```

## Running in Batch Mode

To process multiple prompts from a file:

```bash
# Using the Makefile with default input file (agent/examples/sample_prompts.txt)
make run-batch

# Using the Makefile with a specific input file
make run-batch INPUT_FILE=path/to/your/prompts.txt

# Or directly
python examples/batch_process.py --input path/to/your/prompts.txt
```

### Prompt File Format

The prompt file should contain one prompt per line. Empty lines and lines starting with `#` are ignored:

```
# This is a comment
List all integrations
Show me the integrations related to VMware
Tell me about hpe-alletra-LabRat
```

### Batch Processing Options

- `--input`, `-i` - Input file with prompts (required)
- `--output`, `-o` - Output file for results (default: print to stdout)
- `--format`, `-f` - Output format (default: text, options: text, json)
- `--verbose`, `-v` - Print results to stdout even if output file is specified
- `--server-url` - MCP server URL (default: http://localhost:8080)
- `--llm-provider` - LLM provider to use (default: openai, options: openai, anthropic)
- `--env-file` - Path to a specific .env file for configuration
- `--simple-mode` - Run without connecting to an MCP server

Example with options:

```bash
make run-batch INPUT_FILE=prompts.txt SIMPLE_MODE=true
```

## Running in Simple Mode

Simple mode lets you run the agent without connecting to an actual MCP server. This is useful for:

- Testing the client without a server
- Developing and debugging client features
- Running demos when MCP is not available
- Generating mock responses

To use simple mode, add the `SIMPLE_MODE=true` parameter to any of the commands:

```bash
# Interactive chat in simple mode
make run-agent SIMPLE_MODE=true

# Single prompt in simple mode
make run-prompt PROMPT="List all integrations" SIMPLE_MODE=true

# Batch processing in simple mode
make run-batch INPUT_FILE=prompts.txt SIMPLE_MODE=true
```

Simple mode will return a placeholder response rather than trying to connect to the MCP server. 