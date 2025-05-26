#!/usr/bin/env python3
"""
OpsRamp AI Agent - Multi-Provider Test Suite

This script tests all three LLM providers (OpenAI, Anthropic, Google Gemini)
with both integration and resource management functionality.

Usage:
    python test_all_providers.py [--functionality integration|resources|all]
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

class MultiProviderTestSuite:
    def __init__(self, functionality="all"):
        self.functionality = functionality
        self.timestamp = int(time.time())
        self.session_id = f"multi_provider_{self.functionality}_{self.timestamp}"
        
        # Set up paths
        self.base_path = Path(__file__).parent.parent
        self.output_path = self.base_path / "output"
        self.output_path.mkdir(parents=True, exist_ok=True)
        
        # Test configuration
        self.providers = ["openai", "anthropic", "gemini"]
        self.test_results = {}
    
    def run_integration_tests(self, provider):
        """Run integration tests for a specific provider"""
        print(f"\n🔗 Testing Integration functionality with {provider.upper()}...")
        
        script_path = Path(__file__).parent.parent.parent / "integration" / "scripts" / "run_integration_tests.py"
        
        cmd = [
            sys.executable,
            str(script_path),
            "--provider", provider,
            "--complexity", "basic"
        ]
        
        start_time = time.time()
        try:
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                timeout=1800  # 30 minute timeout
            )
            
            duration = time.time() - start_time
            success = result.returncode == 0
            
            return {
                "provider": provider,
                "functionality": "integration",
                "success": success,
                "duration": duration,
                "return_code": result.returncode,
                "stdout": result.stdout,
                "stderr": result.stderr
            }
            
        except subprocess.TimeoutExpired:
            return {
                "provider": provider,
                "functionality": "integration",
                "success": False,
                "error": "Timeout after 30 minutes"
            }
        except Exception as e:
            return {
                "provider": provider,
                "functionality": "integration",
                "success": False,
                "error": str(e)
            }
    
    def run_resource_tests(self, provider):
        """Run resource management tests for a specific provider"""
        print(f"\n📊 Testing Resource Management functionality with {provider.upper()}...")
        
        script_path = Path(__file__).parent.parent.parent / "resources" / "scripts" / "run_resource_tests.py"
        
        cmd = [
            sys.executable,
            str(script_path),
            "--provider", provider,
            "--complexity", "basic"
        ]
        
        start_time = time.time()
        try:
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                timeout=2400  # 40 minute timeout
            )
            
            duration = time.time() - start_time
            success = result.returncode == 0
            
            return {
                "provider": provider,
                "functionality": "resources",
                "success": success,
                "duration": duration,
                "return_code": result.returncode,
                "stdout": result.stdout,
                "stderr": result.stderr
            }
            
        except subprocess.TimeoutExpired:
            return {
                "provider": provider,
                "functionality": "resources",
                "success": False,
                "error": "Timeout after 40 minutes"
            }
        except Exception as e:
            return {
                "provider": provider,
                "functionality": "resources",
                "success": False,
                "error": str(e)
            }
    
    def test_provider(self, provider):
        """Test a single provider with specified functionality"""
        print(f"\n🚀 Testing {provider.upper()} Provider")
        print("=" * 50)
        
        provider_results = []
        
        if self.functionality in ["integration", "all"]:
            result = self.run_integration_tests(provider)
            provider_results.append(result)
            
            if result["success"]:
                print(f"✅ {provider.upper()} Integration: PASS")
            else:
                print(f"❌ {provider.upper()} Integration: FAIL")
                if "error" in result:
                    print(f"   Error: {result['error']}")
        
        if self.functionality in ["resources", "all"]:
            result = self.run_resource_tests(provider)
            provider_results.append(result)
            
            if result["success"]:
                print(f"✅ {provider.upper()} Resources: PASS")
            else:
                print(f"❌ {provider.upper()} Resources: FAIL")
                if "error" in result:
                    print(f"   Error: {result['error']}")
        
        return provider_results
    
    def generate_summary_report(self):
        """Generate comprehensive multi-provider test report"""
        all_results = []
        for provider_results in self.test_results.values():
            all_results.extend(provider_results)
        
        # Calculate summary statistics
        total_tests = len(all_results)
        successful_tests = sum(1 for r in all_results if r.get("success", False))
        failed_tests = total_tests - successful_tests
        success_rate = (successful_tests / total_tests * 100) if total_tests > 0 else 0
        total_duration = sum(r.get("duration", 0) for r in all_results)
        
        # Provider-specific analysis
        provider_analysis = {}
        for provider in self.providers:
            if provider in self.test_results:
                provider_results = self.test_results[provider]
                provider_success = sum(1 for r in provider_results if r.get("success", False))
                provider_total = len(provider_results)
                provider_rate = (provider_success / provider_total * 100) if provider_total > 0 else 0
                
                provider_analysis[provider] = {
                    "total_tests": provider_total,
                    "successful_tests": provider_success,
                    "failed_tests": provider_total - provider_success,
                    "success_rate": provider_rate,
                    "duration": sum(r.get("duration", 0) for r in provider_results)
                }
        
        report = {
            "session_id": self.session_id,
            "functionality": self.functionality,
            "timestamp": datetime.now().isoformat(),
            "summary": {
                "total_tests": total_tests,
                "successful_tests": successful_tests,
                "failed_tests": failed_tests,
                "success_rate": success_rate,
                "total_duration": total_duration
            },
            "provider_analysis": provider_analysis,
            "detailed_results": self.test_results
        }
        
        # Save report
        report_file = self.output_path / f"{self.session_id}_summary.json"
        with open(report_file, 'w') as f:
            json.dump(report, f, indent=2)
        
        # Print summary
        print(f"\n📊 Multi-Provider Test Summary")
        print("=" * 60)
        print(f"Functionality: {self.functionality}")
        print(f"Total Tests: {total_tests}")
        print(f"Successful: {successful_tests}")
        print(f"Failed: {failed_tests}")
        print(f"Overall Success Rate: {success_rate:.1f}%")
        print(f"Total Duration: {total_duration:.1f}s")
        
        print(f"\n📈 Provider Performance:")
        for provider, analysis in provider_analysis.items():
            print(f"  {provider.upper()}: {analysis['success_rate']:.1f}% "
                  f"({analysis['successful_tests']}/{analysis['total_tests']}) "
                  f"in {analysis['duration']:.1f}s")
        
        print(f"\nReport saved: {report_file}")
        
        return report
    
    def run_all_tests(self):
        """Run tests for all providers"""
        print(f"🚀 Starting Multi-Provider Test Suite")
        print(f"Functionality: {self.functionality}")
        print(f"Providers: {', '.join(self.providers)}")
        print(f"Session ID: {self.session_id}")
        
        for provider in self.providers:
            try:
                provider_results = self.test_provider(provider)
                self.test_results[provider] = provider_results
            except Exception as e:
                print(f"💥 Failed to test {provider}: {e}")
                self.test_results[provider] = [{
                    "provider": provider,
                    "functionality": "unknown",
                    "success": False,
                    "error": str(e)
                }]
        
        # Generate summary report
        report = self.generate_summary_report()
        
        # Return overall success
        return report["summary"]["success_rate"] == 100.0

def main():
    parser = argparse.ArgumentParser(description="Run Multi-Provider OpsRamp AI Agent Tests")
    parser.add_argument(
        "--functionality",
        choices=["integration", "resources", "all"],
        default="all",
        help="Functionality to test across all providers"
    )
    
    args = parser.parse_args()
    
    suite = MultiProviderTestSuite(functionality=args.functionality)
    success = suite.run_all_tests()
    
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main() 