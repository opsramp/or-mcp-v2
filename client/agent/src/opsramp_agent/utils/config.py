"""
Configuration utilities for the OpsRamp AI Agent.
"""

import os
import logging
from pathlib import Path
from typing import Optional, Dict, Any

try:
    from dotenv import load_dotenv
    HAS_DOTENV = True
except ImportError:
    HAS_DOTENV = False

logger = logging.getLogger(__name__)

def load_env_from_file(env_file: Optional[str] = None) -> Dict[str, str]:
    """
    Load environment variables from a .env file.
    
    Args:
        env_file: Path to the .env file (default: looks for .env in current directory and parent directories)
        
    Returns:
        Dictionary of loaded environment variables
    """
    loaded_vars = {}
    
    if not HAS_DOTENV:
        logger.warning("python-dotenv package not installed. Cannot load from .env file.")
        return loaded_vars
    
    # If no specific file is provided, look for .env in current and parent directories
    if env_file is None:
        # Start with the current directory
        current_dir = Path.cwd()
        
        # Check for .env in current directory
        if (current_dir / '.env').exists():
            env_file = str(current_dir / '.env')
        else:
            # Try parent directories (up to 3 levels)
            for i in range(1, 4):
                parent_dir = current_dir.parents[i-1] if i <= len(current_dir.parents) else None
                if parent_dir and (parent_dir / '.env').exists():
                    env_file = str(parent_dir / '.env')
                    break
    
    # If we found a .env file or one was specified, load it
    if env_file and Path(env_file).exists():
        logger.info(f"Loading environment variables from {env_file}")
        # Record the original environment variables
        original_env = os.environ.copy()
        
        # Load the .env file
        load_dotenv(env_file)
        
        # Find new or changed variables
        for key, value in os.environ.items():
            if key not in original_env or original_env[key] != value:
                loaded_vars[key] = value
        
        logger.info(f"Loaded {len(loaded_vars)} environment variables from .env file")
    else:
        logger.warning("No .env file found")
    
    return loaded_vars

def get_api_keys(openai_api_key: Optional[str] = None, anthropic_api_key: Optional[str] = None, gemini_api_key: Optional[str] = None) -> tuple:
    """
    Get API keys from provided parameters or environment variables.
    
    Args:
        openai_api_key: OpenAI API key provided directly
        anthropic_api_key: Anthropic API key provided directly
        gemini_api_key: Google Gemini API key provided directly
        
    Returns:
        Tuple of (openai_api_key, anthropic_api_key, gemini_api_key)
    """
    # Use provided keys or fall back to environment variables
    openai_key = openai_api_key or os.environ.get("OPENAI_API_KEY")
    anthropic_key = anthropic_api_key or os.environ.get("ANTHROPIC_API_KEY")
    gemini_key = gemini_api_key or os.environ.get("GEMINI_API_KEY")
    
    return openai_key, anthropic_key, gemini_key 