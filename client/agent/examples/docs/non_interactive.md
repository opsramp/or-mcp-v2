# OpsRamp AI Agent - Non-Interactive Mode

The OpsRamp AI Agent provides the ability to process prompts in non-interactive mode, which is useful for:

1. Automation
2. Scripting
3. Integration with other tools
4. Testing and validation

## Single Prompt Processing

To process a single prompt without entering the interactive chat loop:

```bash
python examples/chat_client.py --prompt "List all integrations"
```

This will:
1. Connect to the MCP server
2. Process the prompt
3. Output the response
4. Close the connection

Example output:
```
Connecting to MCP server at http://localhost:8080...
Connected!

Processing prompt: List all integrations

Agent: I found 3 integrations:

1. Name: hpe-alletra-LabRat
   ID: INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc
   Status: Installed
   Type: SDK APP

2. Name: redfish-server-LabRat
   ID: INTG-f9e5d2aa-ee17-4e32-9251-493566ebdfca
   Status: Installed  
   Type: SDK APP

3. Name: vcenter-58.51
   ID: INTG-00ee85e2-1f84-4fc1-8cf8-5277ae6980dd
   Status: enabled
   Type: COMPUTE_INTEGRATION
```

## Batch Processing

For processing multiple prompts from a file:

```bash
python examples/batch_process.py --input-file examples/sample_prompts.txt --output-file results.txt
```

### Input File Format

The input file should contain one prompt per line:

```
List all integrations
What is the status of integration INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc?
Show me the integrations related to VMware
Tell me about hpe-alletra-LabRat
```

### Output Options

The batch processor supports two output formats:

1. Text Output (--output-file): Simple text file with prompts and responses
2. JSON Output (--json-output): Structured JSON file for programmatic processing

#### Text Output Example

```
OpsRamp Agent Batch Processing Results
Input File: examples/sample_prompts.txt

--- Prompt 1/4 ---
Prompt: List all integrations
Response: I found 3 integrations:

1. Name: hpe-alletra-LabRat
   ID: INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc
   Status: Installed
   Type: SDK APP

2. Name: redfish-server-LabRat
   ID: INTG-f9e5d2aa-ee17-4e32-9251-493566ebdfca
   Status: Installed  
   Type: SDK APP

3. Name: vcenter-58.51
   ID: INTG-00ee85e2-1f84-4fc1-8cf8-5277ae6980dd
   Status: enabled
   Type: COMPUTE_INTEGRATION
---

--- Prompt 2/4 ---
...
```

#### JSON Output Example

```json
[
  {
    "prompt_id": "1/4",
    "prompt": "List all integrations",
    "response": "I found 3 integrations...",
    "error": null
  },
  {
    "prompt_id": "2/4",
    "prompt": "What is the status of integration INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc?",
    "response": "The integration INTG-2ed93041-eb92-40e9-b6b4-f14ad13e54fc (hpe-alletra-LabRat) has a status of \"Installed\"...",
    "error": null
  },
  ...
]
```

## Integration with Other Tools

The non-interactive mode and batch processing capabilities enable integration with other tools and workflows:

1. **CI/CD Pipelines**: Automate testing of OpsRamp configurations
2. **Monitoring Scripts**: Run periodic checks on integrations
3. **Data Collection**: Gather information about your OpsRamp environment for reporting
4. **Bulk Operations**: Perform operations on multiple integrations in sequence
