#!/usr/bin/env python3
"""
Batch processor for OpsRamp AI Agent.

This script reads prompts from a file and processes them one by one,
outputting the responses to stdout or a file.
"""

import os
import sys
import asyncio
import argparse
import json
import time
import logging
from typing import List, Dict, Union, Optional, Any

# Add parent directory to path for local development
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)

DEFAULT_SERVER_URL = 'http://localhost:8080'
DEFAULT_LLM_PROVIDER = 'openai'
DEFAULT_ENV_FILE = None

# Set to True to use mock mode (no actual MCP or LLM connections)
MOCK_MODE = False

# Import the Agent class
from src.opsramp_agent.agent import Agent, MCPError


def read_prompts_from_file(filepath: str) -> List[str]:
    """
    Read prompts from a file, one per line.
    
    Args:
        filepath: Path to the file containing prompts
        
    Returns:
        List of prompts
    """
    if not os.path.exists(filepath):
        raise FileNotFoundError(f"Prompt file not found: {filepath}")
    
    with open(filepath, 'r') as f:
        # Read lines and strip whitespace
        prompts = [line.strip() for line in f.readlines()]
        
        # Remove empty lines
        prompts = [p for p in prompts if p]
        
        # Remove comment lines (starting with #)
        prompts = [p for p in prompts if not p.startswith('#')]
        
    return prompts


def write_results_to_file(results: List[Dict[str, Any]], output_file: str = None, format_type: str = "text") -> None:
    """
    Write results to a file or stdout.
    
    Args:
        results: List of prompt processing results
        output_file: Path to output file (if None, write to stdout)
        format_type: Output format ('text' or 'json')
    """
    if output_file:
        try:
            if format_type == "json":
                with open(output_file, 'w') as f:
                    json.dump(results, f, indent=2)
            else:
                with open(output_file, 'w') as f:
                    for i, result in enumerate(results):
                        f.write(f"Prompt {i+1}: {result['prompt']}\n")
                        if result.get('success', False):
                            f.write(f"Response: {result['response']}\n\n")
                        else:
                            f.write(f"Error: {result['response']}\n\n")
            print(f"Results written to {output_file}")
        except Exception as e:
            logging.error(f"Error writing results to {output_file}: {e}")
    
    # Print summary to stdout
    total = len(results)
    successful = sum(1 for r in results if r.get('success', False))
    failed = total - successful
    
    print("\nBatch Processing Summary")
    print("-" * 30)
    print(f"Total prompts processed: {total}")
    print(f"Successful responses:    {successful}")
    print(f"Failed responses:        {failed}")
    print("-" * 30)
    
    # If verbose or no output file, print full results to stdout
    if not output_file or format_type != "json":
        for i, result in enumerate(results):
            print(f"\nPrompt {i+1}: {result['prompt']}")
            if 'error' in result:
                print(f"Error: {result['error']}")
            else:
                print(f"{result['response']}")


async def process_prompts(prompts, server_url, llm_provider, env_file=None, simple_mode=False, connection_timeout=60, model=None):
    """Process all prompts and return results."""
    results = []
    
    try:
        # Create the agent
        agent = Agent(
            server_url=server_url,
            llm_provider=llm_provider,
            model=model,
            env_file=env_file,
            simple_mode=simple_mode,
            connection_timeout=connection_timeout
        )
        
        # Connect to the server (not needed in simple mode, handled by Agent class)
        await agent.connect()
        
        # Process each prompt
        for prompt in prompts:
            if prompt.strip() and not prompt.strip().startswith('#'):
                try:
                    response = await agent.chat(prompt)
                    results.append({
                        "prompt": prompt,
                        "response": response,
                        "success": True
                    })
                except Exception as e:
                    logging.error(f"Error processing prompt '{prompt}': {e}")
                    results.append({
                        "prompt": prompt,
                        "response": f"Error: {str(e)}",
                        "success": False
                    })
    except Exception as e:
        logging.error(f"Error initializing agent: {e}")
        # Add error for all prompts
        for prompt in prompts:
            if prompt.strip() and not prompt.strip().startswith('#'):
                results.append({
                    "prompt": prompt,
                    "response": f"Agent initialization failed: {str(e)}",
                    "success": False
                })
    finally:
        # Close the agent connection if it was created
        if 'agent' in locals():
            await agent.close()
    
    return results


async def main(args):
    """Run the batch processor with the given arguments."""
    # Use simple mode if MOCK_MODE is set to True
    simple_mode = args.simple_mode or MOCK_MODE
    
    if simple_mode:
        print("WARNING: Running in simple mode without connecting to MCP server")
    
    # Override args with environment variables if they exist
    llm_provider = os.environ.get('LLM_PROVIDER', args.llm_provider)
    model = os.environ.get('MODEL_NAME', getattr(args, 'model', None))
    
    print(f"Using LLM provider: {llm_provider}")
    if model:
        print(f"Using model: {model}")
    
    try:
        # Read prompts from file
        prompts = read_prompts_from_file(args.input)
        
        # Take only the first prompt if requested
        if args.first_only and prompts:
            original_count = len(prompts)
            prompts = [prompts[0]]
            print(f"Processing only the first prompt out of {original_count} total prompts")
        else:
            print(f"Loaded {len(prompts)} prompts from {args.input}")
        
        # Process the prompts
        results = await process_prompts(
            prompts,
            args.server_url,
            llm_provider,
            args.env_file,
            simple_mode,
            args.connection_timeout,
            model
        )
        
        # Write results to file if output file is specified
        if args.output:
            write_results_to_file(results, args.output, args.format)
        
    except Exception as e:
        logger.error(f"Error in batch processing: {str(e)}", exc_info=True)
        print(f"Error: {str(e)}")
        sys.exit(1)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='OpsRamp AI Agent Batch Processor')
    
    # Server and authentication options
    parser.add_argument('--server-url', default=DEFAULT_SERVER_URL,
                        help=f'OpsRamp MCP server URL (default: {DEFAULT_SERVER_URL})')
    parser.add_argument('--llm-provider', default=DEFAULT_LLM_PROVIDER, choices=['openai', 'anthropic', 'gemini'],
                        help=f'LLM provider to use (default: {DEFAULT_LLM_PROVIDER})')
    parser.add_argument('--model', 
                        help='Model to use (e.g., gpt-4, gpt-3.5-turbo, claude-3-haiku-20240307, gemini-1.5-flash)')
    parser.add_argument('--env-file', default=DEFAULT_ENV_FILE,
                        help='Path to .env file for configuration')
    
    # Input/output options
    parser.add_argument('--input', '-i', required=True,
                        help='Path to file containing prompts (one per line)')
    parser.add_argument('--output', '-o',
                        help='Path to output file (default: print to stdout)')
    parser.add_argument('--format', '-f', choices=['text', 'json'], default='text',
                        help='Output format (default: text)')
    parser.add_argument('--verbose', '-v', action='store_true',
                        help='Print results to stdout even if output file is specified')
    
    # Mode options
    parser.add_argument('--simple-mode', action='store_true',
                        help='Run in simple mode without actually connecting to MCP')
    
    # Connection timeout option
    parser.add_argument('--connection-timeout',
                        type=int,
                        default=60,
                        help='Connection timeout in seconds (default: 60)')
    
    # First prompt only
    parser.add_argument('--first-only', action='store_true',
                        help='Process only the first prompt in the file')
    
    args = parser.parse_args()
    
    # Run the main function
    asyncio.run(main(args)) 