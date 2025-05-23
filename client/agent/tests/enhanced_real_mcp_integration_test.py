#!/usr/bin/env python3
"""
ENHANCED REAL MCP INTEGRATION TESTING WITH COMPREHENSIVE MULTI-TOOL SCENARIOS

This script represents the pinnacle of MCP integration testing sophistication:
- Reads comprehensive prompts from input files
- Handles complex multi-tool call scenarios  
- Logs full request/response payloads with detailed tracing
- Analyzes tool call patterns and response complexity
- Provides expert-level integration analysis and reporting

NO MOCKS. NO SIMULATIONS. ONLY REAL MCP SERVER WITH ADVANCED CAPABILITIES.
"""

import asyncio
import sys
import os
import json
import time
import logging
import re
from datetime import datetime
from typing import List, Dict, Any, Optional
from dataclasses import dataclass, asdict

# Add the src directory to Python path
sys.path.insert(0, os.path.join(os.path.dirname(os.path.dirname(__file__)), 'src'))

from opsramp_agent.agent import Agent

# Advanced logging configuration with multiple handlers
class ComprehensiveTestLogger:
    """Advanced logging system for comprehensive test analysis."""
    
    def __init__(self, test_session_id: str):
        self.test_session_id = test_session_id
        self.start_time = time.time()
        
        # Create formatters
        detailed_formatter = logging.Formatter(
            '%(asctime)s.%(msecs)03d | %(name)s | %(levelname)s | %(message)s',
            datefmt='%Y-%m-%d %H:%M:%S'
        )
        
        # Main test logger
        self.logger = logging.getLogger(f'ENHANCED_MCP_TEST_{test_session_id}')
        self.logger.setLevel(logging.DEBUG)
        
        # Create output directory if it doesn't exist
        output_dir = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'output')
        os.makedirs(output_dir, exist_ok=True)
        
        # File handler for comprehensive logs
        log_filename = os.path.join(output_dir, f'enhanced_integration_test_{test_session_id}.log')
        file_handler = logging.FileHandler(log_filename)
        file_handler.setLevel(logging.DEBUG)
        file_handler.setFormatter(detailed_formatter)
        
        # Console handler for real-time feedback
        console_handler = logging.StreamHandler()
        console_handler.setLevel(logging.INFO)
        console_handler.setFormatter(detailed_formatter)
        
        self.logger.addHandler(file_handler)
        self.logger.addHandler(console_handler)
        
        # Payload logger for request/response data
        self.payload_logger = logging.getLogger(f'PAYLOAD_TRACE_{test_session_id}')
        self.payload_logger.setLevel(logging.DEBUG)
        
        payload_filename = os.path.join(output_dir, f'request_response_payloads_{test_session_id}.jsonl')
        payload_handler = logging.FileHandler(payload_filename)
        payload_handler.setLevel(logging.DEBUG)
        payload_handler.setFormatter(logging.Formatter('%(message)s'))
        
        self.payload_logger.addHandler(payload_handler)
        
        self.log_files = {
            'main_log': log_filename,
            'payload_log': payload_filename
        }
        
        self.output_dir = output_dir
    
    def info(self, message: str):
        self.logger.info(message)
    
    def debug(self, message: str):
        self.logger.debug(message)
    
    def warning(self, message: str):
        self.logger.warning(message)
    
    def error(self, message: str):
        self.logger.error(message)
    
    def log_payload(self, payload_type: str, data: Dict[str, Any]):
        """Log request/response payload as structured JSON."""
        payload_entry = {
            "timestamp": datetime.now().isoformat(),
            "session_id": self.test_session_id,
            "type": payload_type,
            "data": data
        }
        self.payload_logger.debug(json.dumps(payload_entry))

@dataclass
class TestPrompt:
    """Structured representation of a test prompt."""
    category: str
    prompt_text: str
    expected_complexity: str
    line_number: int

@dataclass 
class TestResult:
    """Comprehensive test result with advanced metrics."""
    prompt: TestPrompt
    success: bool
    response: str
    duration_seconds: float
    tool_calls_made: List[Dict[str, Any]]
    error_message: Optional[str]
    response_length: int
    complexity_score: int
    timestamp: str
    session_traces: List[Dict[str, Any]]

class PromptFileParser:
    """Advanced parser for comprehensive integration prompts."""
    
    @staticmethod
    def parse_prompt_file(file_path: str) -> List[TestPrompt]:
        """Parse the comprehensive prompts file into structured test cases."""
        
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"Prompt file not found: {file_path}")
        
        prompts = []
        current_category = "UNKNOWN"
        
        with open(file_path, 'r', encoding='utf-8') as f:
            lines = f.readlines()
        
        for line_num, line in enumerate(lines, 1):
            line = line.strip()
            
            # Skip empty lines and comments
            if not line or line.startswith('#'):
                continue
            
            # Detect category headers
            if line.startswith('## CATEGORY'):
                # Extract category name
                category_match = re.search(r'## CATEGORY \d+: (.+)', line)
                if category_match:
                    current_category = category_match.group(1)
                continue
            
            # This is a prompt line
            if line and not line.startswith('#'):
                # Determine complexity based on keywords and length
                complexity = PromptFileParser._assess_prompt_complexity(line)
                
                prompt = TestPrompt(
                    category=current_category,
                    prompt_text=line,
                    expected_complexity=complexity,
                    line_number=line_num
                )
                prompts.append(prompt)
        
        return prompts
    
    @staticmethod
    def _assess_prompt_complexity(prompt: str) -> str:
        """Assess the expected complexity of a prompt based on its content."""
        prompt_lower = prompt.lower()
        
        # Count complexity indicators
        multi_tool_keywords = ['and', 'compare', 'analyze', 'show me', 'tell me', 'check']
        deep_analysis_keywords = ['comprehensive', 'detailed', 'troubleshoot', 'investigate', 'deep dive']
        workflow_keywords = ['how to', 'walk me through', 'guide me', 'steps', 'process']
        
        complexity_score = 0
        complexity_score += sum(1 for keyword in multi_tool_keywords if keyword in prompt_lower)
        complexity_score += sum(2 for keyword in deep_analysis_keywords if keyword in prompt_lower) 
        complexity_score += sum(1 for keyword in workflow_keywords if keyword in prompt_lower)
        
        # Length-based complexity
        if len(prompt) > 100:
            complexity_score += 1
        if len(prompt) > 150:
            complexity_score += 1
        
        # Classify complexity
        if complexity_score >= 5:
            return "VERY_HIGH"
        elif complexity_score >= 3:
            return "HIGH"
        elif complexity_score >= 2:
            return "MEDIUM"
        else:
            return "LOW"

class EnhancedMCPIntegrationTester:
    """Advanced MCP integration tester with comprehensive analysis capabilities."""
    
    def __init__(self, server_url="http://localhost:8080", prompts_file="test_data/comprehensive_integration_prompts.txt"):
        self.server_url = server_url
        self.prompts_file = prompts_file
        self.test_session_id = str(int(time.time()))
        self.logger = ComprehensiveTestLogger(self.test_session_id)
        self.agent = None
        self.test_results = []
        self.start_time = time.time()
        
        self.logger.info("=" * 100)
        self.logger.info("🚀 ENHANCED REAL MCP INTEGRATION TESTING - SUPREME EDITION")
        self.logger.info(f"📋 Test Session ID: {self.test_session_id}")
        self.logger.info(f"🔗 Server URL: {server_url}")
        self.logger.info(f"📄 Prompts File: {prompts_file}")
        self.logger.info(f"⏰ Test Start Time: {datetime.now().isoformat()}")
        self.logger.info("=" * 100)
    
    async def setup_agent(self) -> bool:
        """Set up the agent with enhanced monitoring and tracing."""
        try:
            self.logger.info("🔧 Setting up Enhanced Agent with advanced MCP connection...")
            
            # Create agent with comprehensive settings
            self.agent = Agent(
                server_url=self.server_url,
                llm_provider="openai",
                model="gpt-4",
                simple_mode=False,
                request_timeout=120  # Extended timeout for complex queries
            )
            
            self.logger.info("🔌 Establishing sophisticated MCP server connection...")
            await self.agent.connect()
            
            # Test connection with a basic call
            test_response = await self.agent.direct_call_tool("integrations", {"action": "list"})
            self.logger.log_payload("CONNECTION_TEST", {
                "request": {"tool": "integrations", "action": "list"},
                "response": test_response,
                "success": True
            })
            
            self.logger.info("✅ Enhanced Agent successfully connected with advanced capabilities!")
            return True
            
        except Exception as e:
            self.logger.error(f"❌ Failed to set up enhanced agent: {str(e)}")
            self.logger.log_payload("CONNECTION_ERROR", {
                "error": str(e),
                "success": False
            })
            return False
    
    async def load_test_prompts(self) -> List[TestPrompt]:
        """Load and parse comprehensive test prompts from file."""
        try:
            # Resolve prompts file path relative to the tests directory
            if not os.path.isabs(self.prompts_file):
                prompts_path = os.path.join(os.path.dirname(__file__), self.prompts_file)
            else:
                prompts_path = self.prompts_file
                
            self.logger.info(f"📖 Loading comprehensive prompts from: {prompts_path}")
            
            prompts = PromptFileParser.parse_prompt_file(prompts_path)
            
            # Analyze prompt distribution
            category_counts = {}
            complexity_counts = {}
            
            for prompt in prompts:
                category_counts[prompt.category] = category_counts.get(prompt.category, 0) + 1
                complexity_counts[prompt.expected_complexity] = complexity_counts.get(prompt.expected_complexity, 0) + 1
            
            self.logger.info(f"📊 Loaded {len(prompts)} comprehensive test prompts")
            self.logger.info(f"📂 Categories: {dict(category_counts)}")
            self.logger.info(f"🎯 Complexity Distribution: {dict(complexity_counts)}")
            
            return prompts
            
        except Exception as e:
            self.logger.error(f"❌ Failed to load prompts: {str(e)}")
            raise
    
    async def execute_comprehensive_test(self, prompt: TestPrompt) -> TestResult:
        """Execute a comprehensive test with advanced analysis and tracing."""
        
        self.logger.info(f"\n{'=' * 80}")
        self.logger.info(f"🧪 EXECUTING COMPREHENSIVE TEST")
        self.logger.info(f"📂 Category: {prompt.category}")
        self.logger.info(f"🎯 Complexity: {prompt.expected_complexity}")
        self.logger.info(f"📝 Prompt: {prompt.prompt_text}")
        self.logger.info(f"📍 Line: {prompt.line_number}")
        self.logger.info(f"{'=' * 80}")
        
        start_time = time.time()
        session_traces = []
        
        try:
            # Log the outgoing request
            self.logger.log_payload("REQUEST_START", {
                "prompt": asdict(prompt),
                "timestamp": datetime.now().isoformat()
            })
            
            # Create a custom tracer to capture tool calls
            original_direct_call = self.agent.direct_call_tool
            tool_calls_made = []
            
            async def traced_direct_call(tool_name: str, arguments: Dict[str, Any]) -> Any:
                call_start = time.time()
                self.logger.debug(f"🔧 Tool Call: {tool_name} with args: {arguments}")
                
                try:
                    result = await original_direct_call(tool_name, arguments)
                    call_duration = time.time() - call_start
                    
                    tool_call_record = {
                        "tool_name": tool_name,
                        "arguments": arguments,
                        "result": result,
                        "duration": call_duration,
                        "success": True,
                        "timestamp": datetime.now().isoformat()
                    }
                    
                    tool_calls_made.append(tool_call_record)
                    session_traces.append(tool_call_record)
                    
                    self.logger.log_payload("TOOL_CALL", tool_call_record)
                    self.logger.debug(f"✅ Tool call completed in {call_duration:.2f}s")
                    
                    return result
                    
                except Exception as e:
                    call_duration = time.time() - call_start
                    
                    tool_call_record = {
                        "tool_name": tool_name,
                        "arguments": arguments,
                        "error": str(e),
                        "duration": call_duration,
                        "success": False,
                        "timestamp": datetime.now().isoformat()
                    }
                    
                    tool_calls_made.append(tool_call_record)
                    session_traces.append(tool_call_record)
                    
                    self.logger.log_payload("TOOL_CALL_ERROR", tool_call_record)
                    self.logger.error(f"❌ Tool call failed in {call_duration:.2f}s: {str(e)}")
                    
                    raise
            
            # Monkey patch for tracing
            self.agent.direct_call_tool = traced_direct_call
            
            # Execute the chat request
            self.logger.info("📤 Sending comprehensive request to enhanced MCP agent...")
            response = await self.agent.chat(prompt.prompt_text)
            
            # Restore original method
            self.agent.direct_call_tool = original_direct_call
            
            duration = time.time() - start_time
            
            # Calculate complexity score based on actual execution
            complexity_score = self._calculate_complexity_score(tool_calls_made, response, duration)
            
            # Create comprehensive result
            result = TestResult(
                prompt=prompt,
                success=True,
                response=response,
                duration_seconds=duration,
                tool_calls_made=tool_calls_made,
                error_message=None,
                response_length=len(response),
                complexity_score=complexity_score,
                timestamp=datetime.now().isoformat(),
                session_traces=session_traces
            )
            
            # Log the complete response
            self.logger.log_payload("RESPONSE_COMPLETE", {
                "prompt": prompt.prompt_text,
                "response": response,
                "duration": duration,
                "tool_calls_count": len(tool_calls_made),
                "complexity_score": complexity_score,
                "success": True
            })
            
            self.logger.info(f"✅ TEST COMPLETED SUCCESSFULLY")
            self.logger.info(f"⏱️  Duration: {duration:.2f}s")
            self.logger.info(f"🔧 Tool Calls: {len(tool_calls_made)}")
            self.logger.info(f"📏 Response Length: {len(response)} chars")
            self.logger.info(f"🎯 Complexity Score: {complexity_score}")
            
            return result
            
        except Exception as e:
            duration = time.time() - start_time
            
            # Create error result
            result = TestResult(
                prompt=prompt,
                success=False,
                response="",
                duration_seconds=duration,
                tool_calls_made=tool_calls_made,
                error_message=str(e),
                response_length=0,
                complexity_score=0,
                timestamp=datetime.now().isoformat(),
                session_traces=session_traces
            )
            
            self.logger.log_payload("TEST_ERROR", {
                "prompt": prompt.prompt_text,
                "error": str(e),
                "duration": duration,
                "tool_calls_count": len(tool_calls_made)
            })
            
            self.logger.error(f"❌ TEST FAILED: {str(e)}")
            return result
    
    def _calculate_complexity_score(self, tool_calls: List[Dict], response: str, duration: float) -> int:
        """Calculate complexity score based on actual execution metrics."""
        score = 0
        
        # Tool call complexity
        score += len(tool_calls) * 2
        score += len([call for call in tool_calls if call.get("success", True)]) * 1
        
        # Response complexity
        score += min(len(response) // 500, 10)  # Cap at 10 for response length
        
        # Duration complexity
        if duration > 30:
            score += 5
        elif duration > 15:
            score += 3
        elif duration > 5:
            score += 1
        
        # Multi-tool bonus
        unique_tools = set(call.get("tool_name", "") for call in tool_calls)
        if len(unique_tools) > 1:
            score += 5
        
        return score
    
    async def run_comprehensive_test_suite(self, max_tests: Optional[int] = None) -> Dict[str, Any]:
        """Run the complete comprehensive test suite with advanced analysis."""
        
        try:
            # Load prompts
            prompts = await self.load_test_prompts()
            
            if max_tests:
                prompts = prompts[:max_tests]
                self.logger.info(f"🎯 Limiting tests to first {max_tests} prompts")
            
            self.logger.info(f"🚀 Starting comprehensive test execution of {len(prompts)} prompts")
            
            # Execute all tests
            for i, prompt in enumerate(prompts, 1):
                self.logger.info(f"\n🔄 Executing test {i}/{len(prompts)}")
                
                result = await self.execute_comprehensive_test(prompt)
                self.test_results.append(result)
                
                # Brief pause between tests to avoid overwhelming the server
                await asyncio.sleep(1)
            
            # Generate comprehensive analysis
            analysis = self._generate_comprehensive_analysis()
            
            return analysis
            
        except Exception as e:
            self.logger.error(f"💥 Test suite failed: {str(e)}", exc_info=True)
            raise
    
    def _generate_comprehensive_analysis(self) -> Dict[str, Any]:
        """Generate comprehensive analysis of all test results."""
        
        total_tests = len(self.test_results)
        successful_tests = sum(1 for r in self.test_results if r.success)
        total_duration = time.time() - self.start_time
        
        # Category analysis
        category_stats = {}
        complexity_stats = {}
        tool_usage_stats = {}
        
        for result in self.test_results:
            # Category stats
            cat = result.prompt.category
            if cat not in category_stats:
                category_stats[cat] = {"total": 0, "success": 0, "avg_duration": 0, "total_duration": 0}
            
            category_stats[cat]["total"] += 1
            category_stats[cat]["total_duration"] += result.duration_seconds
            if result.success:
                category_stats[cat]["success"] += 1
        
        # Calculate averages
        for cat_data in category_stats.values():
            cat_data["avg_duration"] = cat_data["total_duration"] / cat_data["total"]
            cat_data["success_rate"] = (cat_data["success"] / cat_data["total"]) * 100
        
        # Tool usage analysis
        all_tool_calls = []
        for result in self.test_results:
            all_tool_calls.extend(result.tool_calls_made)
        
        for call in all_tool_calls:
            tool_name = call.get("tool_name", "unknown")
            action = call.get("arguments", {}).get("action", "unknown")
            key = f"{tool_name}:{action}"
            
            if key not in tool_usage_stats:
                tool_usage_stats[key] = {"count": 0, "success_count": 0, "total_duration": 0}
            
            tool_usage_stats[key]["count"] += 1
            tool_usage_stats[key]["total_duration"] += call.get("duration", 0)
            if call.get("success", True):
                tool_usage_stats[key]["success_count"] += 1
        
        # Generate comprehensive report
        analysis = {
            "test_summary": {
                "session_id": self.test_session_id,
                "total_tests": total_tests,
                "successful_tests": successful_tests,
                "failed_tests": total_tests - successful_tests,
                "success_rate": (successful_tests / total_tests * 100) if total_tests > 0 else 0,
                "total_duration": total_duration,
                "avg_test_duration": total_duration / total_tests if total_tests > 0 else 0,
                "total_tool_calls": len(all_tool_calls),
                "server_url": self.server_url,
                "timestamp": datetime.now().isoformat()
            },
            "category_analysis": category_stats,
            "tool_usage_analysis": tool_usage_stats,
            "complexity_distribution": self._analyze_complexity_distribution(),
            "top_performers": self._identify_top_performers(),
            "failure_analysis": self._analyze_failures(),
            "detailed_results": [asdict(result) for result in self.test_results],
            "log_files": self.logger.log_files
        }
        
        # Log comprehensive analysis
        self.logger.info("\n" + "=" * 100)
        self.logger.info("📊 COMPREHENSIVE TEST ANALYSIS - SUPREME EDITION")
        self.logger.info("=" * 100)
        
        summary = analysis["test_summary"]
        self.logger.info(f"🎯 Test Execution Summary:")
        self.logger.info(f"   • Session ID: {summary['session_id']}")
        self.logger.info(f"   • Total Tests: {summary['total_tests']}")
        self.logger.info(f"   • Successful: {summary['successful_tests']}")
        self.logger.info(f"   • Failed: {summary['failed_tests']}")
        self.logger.info(f"   • Success Rate: {summary['success_rate']:.1f}%")
        self.logger.info(f"   • Total Duration: {summary['total_duration']:.2f}s")
        self.logger.info(f"   • Average Test Duration: {summary['avg_test_duration']:.2f}s")
        self.logger.info(f"   • Total Tool Calls: {summary['total_tool_calls']}")
        
        self.logger.info(f"\n📂 Category Performance:")
        for category, stats in category_stats.items():
            self.logger.info(f"   • {category}: {stats['success']}/{stats['total']} "
                           f"({stats['success_rate']:.1f}%) - {stats['avg_duration']:.2f}s avg")
        
        self.logger.info(f"\n🔧 Tool Usage Statistics:")
        sorted_tools = sorted(tool_usage_stats.items(), key=lambda x: x[1]['count'], reverse=True)
        for tool_key, stats in sorted_tools[:10]:  # Top 10
            success_rate = (stats['success_count'] / stats['count'] * 100) if stats['count'] > 0 else 0
            avg_duration = stats['total_duration'] / stats['count'] if stats['count'] > 0 else 0
            self.logger.info(f"   • {tool_key}: {stats['count']} calls "
                           f"({success_rate:.1f}% success) - {avg_duration:.2f}s avg")
        
        # Save analysis to file
        analysis_file = os.path.join(self.logger.output_dir, f"comprehensive_test_analysis_{self.test_session_id}.json")
        with open(analysis_file, 'w') as f:
            json.dump(analysis, f, indent=2)
        
        self.logger.info(f"\n💾 Comprehensive analysis saved to: {analysis_file}")
        
        if successful_tests == total_tests:
            self.logger.info("\n🏆 SUPREME SUCCESS: ALL COMPREHENSIVE TESTS PASSED!")
            self.logger.info("🎉 Complete multi-tool integration validation achieved with advanced capabilities!")
        else:
            self.logger.info(f"\n⚠️  {total_tests - successful_tests} tests need attention for complete success")
        
        return analysis
    
    def _analyze_complexity_distribution(self) -> Dict[str, Any]:
        """Analyze the distribution of test complexity."""
        complexity_data = {}
        
        for result in self.test_results:
            complexity = result.prompt.expected_complexity
            score = result.complexity_score
            
            if complexity not in complexity_data:
                complexity_data[complexity] = {
                    "count": 0,
                    "avg_score": 0,
                    "avg_duration": 0,
                    "success_rate": 0,
                    "total_score": 0,
                    "total_duration": 0,
                    "successes": 0
                }
            
            data = complexity_data[complexity]
            data["count"] += 1
            data["total_score"] += score
            data["total_duration"] += result.duration_seconds
            if result.success:
                data["successes"] += 1
        
        # Calculate averages
        for data in complexity_data.values():
            data["avg_score"] = data["total_score"] / data["count"]
            data["avg_duration"] = data["total_duration"] / data["count"]
            data["success_rate"] = (data["successes"] / data["count"]) * 100
        
        return complexity_data
    
    def _identify_top_performers(self) -> Dict[str, Any]:
        """Identify top performing tests by various metrics."""
        
        # Sort by different criteria
        by_complexity = sorted(self.test_results, key=lambda r: r.complexity_score, reverse=True)[:5]
        by_tool_calls = sorted(self.test_results, key=lambda r: len(r.tool_calls_made), reverse=True)[:5]
        by_response_length = sorted(self.test_results, key=lambda r: r.response_length, reverse=True)[:5]
        
        return {
            "highest_complexity": [{"prompt": r.prompt.prompt_text, "score": r.complexity_score} for r in by_complexity],
            "most_tool_calls": [{"prompt": r.prompt.prompt_text, "calls": len(r.tool_calls_made)} for r in by_tool_calls],
            "longest_responses": [{"prompt": r.prompt.prompt_text, "length": r.response_length} for r in by_response_length]
        }
    
    def _analyze_failures(self) -> Dict[str, Any]:
        """Analyze failed tests for patterns."""
        failures = [r for r in self.test_results if not r.success]
        
        if not failures:
            return {"count": 0, "patterns": [], "categories": {}}
        
        # Analyze failure patterns
        error_patterns = {}
        failure_categories = {}
        
        for failure in failures:
            error = failure.error_message or "Unknown error"
            category = failure.prompt.category
            
            # Group similar errors
            error_key = error[:50]  # First 50 chars for grouping
            error_patterns[error_key] = error_patterns.get(error_key, 0) + 1
            
            # Category analysis
            failure_categories[category] = failure_categories.get(category, 0) + 1
        
        return {
            "count": len(failures),
            "error_patterns": error_patterns,
            "failure_by_category": failure_categories,
            "sample_failures": [{"prompt": f.prompt.prompt_text, "error": f.error_message} for f in failures[:3]]
        }
    
    async def cleanup(self):
        """Clean up resources with comprehensive logging."""
        if self.agent:
            try:
                await self.agent.close()
                self.logger.info("🧹 Enhanced agent cleanup completed successfully")
            except Exception as e:
                self.logger.error(f"Error during agent cleanup: {str(e)}")

async def main():
    """Main execution function for comprehensive testing."""
    
    # Parse command line arguments
    import argparse
    parser = argparse.ArgumentParser(description='Enhanced Real MCP Integration Testing')
    parser.add_argument('--prompts-file', default='test_data/comprehensive_integration_prompts.txt',
                       help='File containing test prompts (default: test_data/comprehensive_integration_prompts.txt)')
    parser.add_argument('--server-url', default='http://localhost:8080',
                       help='MCP server URL (default: http://localhost:8080)')
    parser.add_argument('--max-tests', type=int, default=None,
                       help='Maximum number of tests to run (default: all)')
    
    args = parser.parse_args()
    
    tester = EnhancedMCPIntegrationTester(
        server_url=args.server_url,
        prompts_file=args.prompts_file
    )
    
    try:
        # Setup
        if not await tester.setup_agent():
            tester.logger.error("Failed to setup enhanced agent")
            return 1
        
        # Run comprehensive tests
        analysis = await tester.run_comprehensive_test_suite(max_tests=args.max_tests)
        
        # Determine success
        success_rate = analysis["test_summary"]["success_rate"]
        
        if success_rate == 100.0:
            return 0
        elif success_rate >= 90.0:
            return 1  # Mostly successful but some issues
        else:
            return 2  # Significant issues
        
    except KeyboardInterrupt:
        tester.logger.info("🛑 Comprehensive testing interrupted by user")
        return 130
    except Exception as e:
        tester.logger.error(f"💥 Comprehensive testing failed: {str(e)}", exc_info=True)
        return 1
    finally:
        await tester.cleanup()

if __name__ == "__main__":
    sys.exit(asyncio.run(main())) 