#!/usr/bin/env python3
"""
Test Data Cleanup Utility

Cleans up old test data while preserving recent results and evidence.
"""

import os
import sys
import argparse
from datetime import datetime, timedelta
from pathlib import Path

def cleanup_old_files(directory, days_to_keep=30):
    """Remove files older than specified days"""
    cutoff_date = datetime.now() - timedelta(days=days_to_keep)
    removed_count = 0
    
    for file_path in Path(directory).rglob("*"):
        if file_path.is_file():
            file_time = datetime.fromtimestamp(file_path.stat().st_mtime)
            if file_time < cutoff_date:
                print(f"Removing old file: {file_path}")
                file_path.unlink()
                removed_count += 1
    
    return removed_count

def main():
    parser = argparse.ArgumentParser(description="Cleanup old test data")
    parser.add_argument("--days", type=int, default=30, 
                       help="Days of data to keep (default: 30)")
    parser.add_argument("--dry-run", action="store_true",
                       help="Show what would be deleted without deleting")
    
    args = parser.parse_args()
    
    tests_path = Path(__file__).parent.parent
    output_dirs = [
        tests_path / "integration/output",
        tests_path / "resources/output", 
        tests_path / "multi_provider/output"
    ]
    
    total_removed = 0
    for output_dir in output_dirs:
        if output_dir.exists():
            print(f"Cleaning {output_dir}...")
            if not args.dry_run:
                removed = cleanup_old_files(output_dir, args.days)
                total_removed += removed
            else:
                print(f"DRY RUN: Would clean files older than {args.days} days")
    
    print(f"Cleanup complete. Removed {total_removed} files.")

if __name__ == "__main__":
    main()
