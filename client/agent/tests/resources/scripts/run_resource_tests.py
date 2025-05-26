#!/usr/bin/env python3
"""
OpsRamp AI Agent - Resource Management Test Runner

This script automates the execution of resource management functionality tests,
capturing real API payloads and generating comprehensive test reports.

Usage:
    python run_resource_tests.py [--provider openai|anthropic|gemini] [--complexity basic|comprehensive|ultra|all]
"""

import os
import sys
import json
import time
import argparse
import subprocess
from datetime import datetime
from pathlib import Path

# Add the agent source to Python path
sys.path.insert(0, str(Path(__file__).parent.parent.parent.parent / "src"))

class ResourceTestRunner:
    def __init__(self, provider="openai", complexity="basic"):
        self.provider = provider
        self.complexity = complexity
        self.timestamp = int(time.time())
        self.session_id = f"resource_{self.complexity}_{self.provider}_{self.timestamp}"
        
        # Set up paths
        self.base_path = Path(__file__).parent.parent
        self.test_data_path = self.base_path / "test_data"
        self.output_path = self.base_path / "output"
        self.logs_path = self.output_path / "logs"
        self.payloads_path = self.output_path / "payloads"
        self.reports_path = self.output_path / "reports"
        
        # Create output directories
        for path in [self.logs_path, self.payloads_path, self.reports_path]:
            path.mkdir(parents=True, exist_ok=True)
    
    def get_test_files(self):
        """Get list of test files based on complexity setting"""
        test_files = []
        
        if self.complexity in ["basic", "all"]:
            basic_file = self.test_data_path / "basic_resource_prompts.txt"
            if basic_file.exists():
                test_files.append(("basic", basic_file))
        
        if self.complexity in ["comprehensive", "all"]:
            comprehensive_file = self.test_data_path / "comprehensive_resource_prompts.txt"
            if comprehensive_file.exists():
                test_files.append(("comprehensive", comprehensive_file))
        
        if self.complexity in ["ultra", "all"]:
            ultra_file = self.test_data_path / "ultra_complex_resource_prompts.txt"
            if ultra_file.exists():
                test_files.append(("ultra_complex", ultra_file))
        
        return test_files
    
    def run_test_batch(self, test_type, test_file):
        """Run a batch of tests for a specific test type"""
        print(f"\n🧪 Running {test_type} resource tests with {self.provider}...")
        print(f"📁 Test file: {test_file}")
        
        # Set up output files
        log_file = self.logs_path / f"{self.session_id}_{test_type}.log"
        payload_file = self.payloads_path / f"{self.session_id}_{test_type}_payloads.jsonl"
        
        # Prepare environment variables
        env = os.environ.copy()
        if self.provider == "anthropic":
            env["LLM_PROVIDER"] = "anthropic"
            env["MODEL_NAME"] = "claude-3-haiku-20240307"
        elif self.provider == "gemini":
            env["LLM_PROVIDER"] = "gemini"
            env["MODEL_NAME"] = "gemini-1.5-flash"
        else:  # openai
            env["LLM_PROVIDER"] = "openai"
            env["MODEL_NAME"] = "gpt-3.5-turbo"
        
        # Build command
        agent_path = Path(__file__).parent.parent.parent.parent
        cmd = [
            sys.executable,
            str(agent_path / "examples" / "batch_process.py"),
            "--input", str(test_file),
            "--output", str(log_file),
            "--format", "text",
            "--verbose"
        ]
        
        # Run the test
        start_time = time.time()
        try:
            result = subprocess.run(
                cmd,
                cwd=str(agent_path),
                env=env,
                capture_output=True,
                text=True,
                timeout=2400  # 40 minute timeout for resource tests (they can be longer)
            )
            
            duration = time.time() - start_time
            success = result.returncode == 0
            
            # Count prompts processed
            prompt_count = 0
            if test_file.exists():
                with open(test_file, 'r') as f:
                    prompt_count = len([line for line in f if line.strip() and not line.strip().startswith('#')])
            
            # Log results
            test_result = {
                "session_id": self.session_id,
                "test_type": test_type,
                "provider": self.provider,
                "start_time": datetime.fromtimestamp(start_time).isoformat(),
                "duration": duration,
                "success": success,
                "return_code": result.returncode,
                "prompt_count": prompt_count,
                "stdout": result.stdout,
                "stderr": result.stderr,
                "command": " ".join(cmd)
            }
            
            # Save detailed log
            with open(log_file.with_suffix('.json'), 'w') as f:
                json.dump(test_result, f, indent=2)
            
            if success:
                print(f"✅ {test_type.title()} tests completed successfully in {duration:.1f}s")
                print(f"📊 Processed {prompt_count} prompts")
            else:
                print(f"❌ {test_type.title()} tests failed after {duration:.1f}s")
                print(f"Error: {result.stderr}")
            
            return test_result
            
        except subprocess.TimeoutExpired:
            print(f"⏰ {test_type.title()} tests timed out after 40 minutes")
            return {
                "session_id": self.session_id,
                "test_type": test_type,
                "provider": self.provider,
                "success": False,
                "error": "Timeout after 40 minutes"
            }
        except Exception as e:
            print(f"💥 {test_type.title()} tests failed with exception: {e}")
            return {
                "session_id": self.session_id,
                "test_type": test_type,
                "provider": self.provider,
                "success": False,
                "error": str(e)
            }
    
    def analyze_token_efficiency(self, test_results):
        """Analyze token efficiency across test runs"""
        analysis = {
            "total_prompts": 0,
            "total_duration": 0,
            "avg_time_per_prompt": 0,
            "efficiency_score": "N/A"
        }
        
        for result in test_results:
            if result.get("success") and result.get("prompt_count"):
                analysis["total_prompts"] += result["prompt_count"]
                analysis["total_duration"] += result.get("duration", 0)
        
        if analysis["total_prompts"] > 0:
            analysis["avg_time_per_prompt"] = analysis["total_duration"] / analysis["total_prompts"]
            
            # Calculate efficiency score based on prompts per minute
            prompts_per_minute = analysis["total_prompts"] / (analysis["total_duration"] / 60)
            if prompts_per_minute >= 10:
                analysis["efficiency_score"] = "Excellent"
            elif prompts_per_minute >= 5:
                analysis["efficiency_score"] = "Good"
            elif prompts_per_minute >= 2:
                analysis["efficiency_score"] = "Fair"
            else:
                analysis["efficiency_score"] = "Poor"
        
        return analysis
    
    def generate_summary_report(self, test_results):
        """Generate a comprehensive test summary report"""
        token_analysis = self.analyze_token_efficiency(test_results)
        
        report = {
            "session_id": self.session_id,
            "provider": self.provider,
            "complexity": self.complexity,
            "timestamp": datetime.now().isoformat(),
            "summary": {
                "total_test_batches": len(test_results),
                "successful_batches": sum(1 for r in test_results if r.get("success", False)),
                "failed_batches": sum(1 for r in test_results if not r.get("success", False)),
                "total_duration": sum(r.get("duration", 0) for r in test_results),
                "total_prompts": token_analysis["total_prompts"],
                "success_rate": 0
            },
            "token_efficiency": token_analysis,
            "test_results": test_results
        }
        
        if report["summary"]["total_test_batches"] > 0:
            report["summary"]["success_rate"] = (
                report["summary"]["successful_batches"] / report["summary"]["total_test_batches"] * 100
            )
        
        # Save report
        report_file = self.reports_path / f"{self.session_id}_summary.json"
        with open(report_file, 'w') as f:
            json.dump(report, f, indent=2)
        
        # Print summary
        print(f"\n📊 Resource Management Test Summary")
        print(f"=" * 50)
        print(f"Provider: {self.provider}")
        print(f"Complexity: {self.complexity}")
        print(f"Test Batches: {report['summary']['total_test_batches']}")
        print(f"Successful: {report['summary']['successful_batches']}")
        print(f"Failed: {report['summary']['failed_batches']}")
        print(f"Success Rate: {report['summary']['success_rate']:.1f}%")
        print(f"Total Duration: {report['summary']['total_duration']:.1f}s")
        print(f"Total Prompts: {report['summary']['total_prompts']}")
        print(f"Avg Time/Prompt: {token_analysis['avg_time_per_prompt']:.2f}s")
        print(f"Efficiency Score: {token_analysis['efficiency_score']}")
        print(f"Report saved: {report_file}")
        
        return report
    
    def run_all_tests(self):
        """Run all resource management tests"""
        print(f"🚀 Starting Resource Management Test Suite")
        print(f"Provider: {self.provider}")
        print(f"Complexity: {self.complexity}")
        print(f"Session ID: {self.session_id}")
        
        test_files = self.get_test_files()
        if not test_files:
            print("❌ No test files found!")
            return False
        
        test_results = []
        for test_type, test_file in test_files:
            result = self.run_test_batch(test_type, test_file)
            test_results.append(result)
        
        # Generate summary report
        report = self.generate_summary_report(test_results)
        
        # Return overall success
        return report["summary"]["success_rate"] == 100.0

def main():
    parser = argparse.ArgumentParser(description="Run OpsRamp AI Agent Resource Management Tests")
    parser.add_argument(
        "--provider",
        choices=["openai", "anthropic", "gemini"],
        default="openai",
        help="LLM provider to use for testing"
    )
    parser.add_argument(
        "--complexity",
        choices=["basic", "comprehensive", "ultra", "all"],
        default="basic",
        help="Test complexity level to run"
    )
    
    args = parser.parse_args()
    
    runner = ResourceTestRunner(
        provider=args.provider,
        complexity=args.complexity
    )
    
    success = runner.run_all_tests()
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main() 