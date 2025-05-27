#!/usr/bin/env python3
"""
OpsRamp AI Agent - Comprehensive Test Report Generator

This script generates comprehensive test reports by analyzing all test results,
API payloads, and evidence from Integration and Resources testing.

Usage:
    python generate_test_report.py [--period daily|weekly|monthly] [--format html|json|text]
"""

import os
import sys
import json
import argparse
from datetime import datetime, timedelta
from pathlib import Path
import glob

class TestReportGenerator:
    def __init__(self, period="daily", format="html"):
        self.period = period
        self.format = format
        self.timestamp = datetime.now()
        
        # Base paths
        self.tests_path = Path(__file__).parent.parent
        self.evidence_path = self.tests_path / "evidence"
        self.output_dirs = {
            "integration": self.tests_path / "integration" / "output",
            "resources": self.tests_path / "resources" / "output", 
            "multi_provider": self.tests_path / "multi_provider" / "output"
        }
        
        # Report data
        self.report_data = {
            "metadata": {
                "generated_at": self.timestamp.isoformat(),
                "period": self.period,
                "format": self.format
            },
            "summary": {},
            "integration_tests": {},
            "resource_tests": {},
            "multi_provider_tests": {},
            "api_evidence": {},
            "performance_metrics": {},
            "recommendations": []
        }
    
    def get_time_filter(self):
        """Get time filter based on period"""
        if self.period == "daily":
            return self.timestamp - timedelta(days=1)
        elif self.period == "weekly":
            return self.timestamp - timedelta(weeks=1)
        elif self.period == "monthly":
            return self.timestamp - timedelta(days=30)
        else:
            return self.timestamp - timedelta(days=1)
    
    def analyze_test_results(self, category):
        """Analyze test results for a specific category"""
        output_dir = self.output_dirs.get(category)
        if not output_dir or not output_dir.exists():
            return {"error": f"No output directory found for {category}"}
        
        results = {
            "total_sessions": 0,
            "successful_sessions": 0,
            "failed_sessions": 0,
            "total_duration": 0,
            "total_prompts": 0,
            "recent_sessions": []
        }
        
        time_filter = self.get_time_filter()
        
        # Analyze log files
        log_files = list((output_dir / "logs").glob("*.json")) if (output_dir / "logs").exists() else []
        
        for log_file in log_files:
            try:
                file_time = datetime.fromtimestamp(log_file.stat().st_mtime)
                if file_time >= time_filter:
                    with open(log_file, 'r') as f:
                        session_data = json.load(f)
                    
                    results["total_sessions"] += 1
                    if session_data.get("success", False):
                        results["successful_sessions"] += 1
                    else:
                        results["failed_sessions"] += 1
                    
                    results["total_duration"] += session_data.get("duration", 0)
                    results["total_prompts"] += session_data.get("prompt_count", 0)
                    
                    results["recent_sessions"].append({
                        "session_id": session_data.get("session_id"),
                        "test_type": session_data.get("test_type"),
                        "success": session_data.get("success"),
                        "duration": session_data.get("duration"),
                        "timestamp": file_time.isoformat()
                    })
            except Exception as e:
                print(f"Error analyzing {log_file}: {e}")
        
        # Calculate success rate
        if results["total_sessions"] > 0:
            results["success_rate"] = (results["successful_sessions"] / results["total_sessions"]) * 100
        else:
            results["success_rate"] = 0
        
        return results
    
    def analyze_api_evidence(self):
        """Analyze API payload evidence"""
        evidence = {
            "total_payloads": 0,
            "integration_payloads": 0,
            "resource_payloads": 0,
            "payload_sizes": [],
            "api_endpoints": set(),
            "recent_payloads": []
        }
        
        time_filter = self.get_time_filter()
        
        # Check all payload directories
        for category, output_dir in self.output_dirs.items():
            payload_dir = output_dir / "payloads"
            if payload_dir.exists():
                payload_files = list(payload_dir.glob("*.jsonl"))
                
                for payload_file in payload_files:
                    try:
                        file_time = datetime.fromtimestamp(payload_file.stat().st_mtime)
                        if file_time >= time_filter:
                            file_size = payload_file.stat().st_size
                            evidence["total_payloads"] += 1
                            evidence["payload_sizes"].append(file_size)
                            
                            if category == "integration":
                                evidence["integration_payloads"] += 1
                            elif category == "resources":
                                evidence["resource_payloads"] += 1
                            
                            # Analyze payload content for API endpoints
                            with open(payload_file, 'r') as f:
                                for line in f:
                                    try:
                                        payload = json.loads(line.strip())
                                        if "url" in payload:
                                            evidence["api_endpoints"].add(payload["url"])
                                    except:
                                        continue
                            
                            evidence["recent_payloads"].append({
                                "file": str(payload_file.name),
                                "category": category,
                                "size": file_size,
                                "timestamp": file_time.isoformat()
                            })
                    except Exception as e:
                        print(f"Error analyzing payload {payload_file}: {e}")
        
        evidence["api_endpoints"] = list(evidence["api_endpoints"])
        return evidence
    
    def calculate_performance_metrics(self):
        """Calculate performance metrics across all tests"""
        metrics = {
            "avg_response_time": 0,
            "total_test_time": 0,
            "prompts_per_minute": 0,
            "success_rate_trend": [],
            "efficiency_score": "N/A"
        }
        
        total_duration = 0
        total_prompts = 0
        
        for category in ["integration", "resources", "multi_provider"]:
            results = self.report_data.get(f"{category}_tests", {})
            total_duration += results.get("total_duration", 0)
            total_prompts += results.get("total_prompts", 0)
        
        if total_prompts > 0 and total_duration > 0:
            metrics["avg_response_time"] = total_duration / total_prompts
            metrics["total_test_time"] = total_duration
            metrics["prompts_per_minute"] = total_prompts / (total_duration / 60)
            
            # Calculate efficiency score
            if metrics["prompts_per_minute"] >= 10:
                metrics["efficiency_score"] = "Excellent"
            elif metrics["prompts_per_minute"] >= 5:
                metrics["efficiency_score"] = "Good"
            elif metrics["prompts_per_minute"] >= 2:
                metrics["efficiency_score"] = "Fair"
            else:
                metrics["efficiency_score"] = "Poor"
        
        return metrics
    
    def generate_recommendations(self):
        """Generate recommendations based on test results"""
        recommendations = []
        
        # Check success rates
        for category in ["integration", "resources"]:
            results = self.report_data.get(f"{category}_tests", {})
            success_rate = results.get("success_rate", 0)
            
            if success_rate < 80:
                recommendations.append({
                    "category": "reliability",
                    "priority": "high",
                    "message": f"{category.title()} tests have low success rate ({success_rate:.1f}%). Investigate failures and improve test stability."
                })
            elif success_rate < 95:
                recommendations.append({
                    "category": "reliability", 
                    "priority": "medium",
                    "message": f"{category.title()} tests could be more reliable ({success_rate:.1f}%). Review failed test cases."
                })
        
        # Check performance
        performance = self.report_data.get("performance_metrics", {})
        efficiency = performance.get("efficiency_score", "N/A")
        
        if efficiency == "Poor":
            recommendations.append({
                "category": "performance",
                "priority": "high", 
                "message": "Test execution is slow. Consider optimizing prompts, using faster models, or implementing parallel testing."
            })
        elif efficiency == "Fair":
            recommendations.append({
                "category": "performance",
                "priority": "medium",
                "message": "Test performance could be improved. Review prompt efficiency and model selection."
            })
        
        # Check API evidence
        evidence = self.report_data.get("api_evidence", {})
        total_payloads = evidence.get("total_payloads", 0)
        
        if total_payloads == 0:
            recommendations.append({
                "category": "evidence",
                "priority": "high",
                "message": "No API payload evidence found. Ensure test scripts are capturing real API calls."
            })
        elif total_payloads < 10:
            recommendations.append({
                "category": "evidence",
                "priority": "medium", 
                "message": "Limited API evidence collected. Consider running more comprehensive tests."
            })
        
        return recommendations
    
    def generate_html_report(self):
        """Generate HTML format report"""
        html_template = """<!DOCTYPE html>
<html>
<head>
    <title>OpsRamp AI Agent Test Report</title>
    <style>
        body {{ font-family: Arial, sans-serif; margin: 20px; }}
        .header {{ background: #2c3e50; color: white; padding: 20px; border-radius: 5px; }}
        .section {{ margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }}
        .success {{ color: #27ae60; }}
        .warning {{ color: #f39c12; }}
        .error {{ color: #e74c3c; }}
        .metric {{ display: inline-block; margin: 10px; padding: 10px; background: #ecf0f1; border-radius: 3px; }}
        table {{ width: 100%; border-collapse: collapse; }}
        th, td {{ padding: 8px; text-align: left; border-bottom: 1px solid #ddd; }}
        th {{ background-color: #f2f2f2; }}
    </style>
</head>
<body>
    <div class="header">
        <h1>🚀 OpsRamp AI Agent Test Report</h1>
        <p>Generated: {timestamp}</p>
        <p>Period: {period}</p>
    </div>
    
    <div class="section">
        <h2>📊 Executive Summary</h2>
        <div class="metric">
            <strong>Total Tests:</strong> {total_tests}
        </div>
        <div class="metric">
            <strong>Success Rate:</strong> <span class="{success_class}">{success_rate:.1f}%</span>
        </div>
        <div class="metric">
            <strong>API Evidence:</strong> {total_payloads} payloads
        </div>
        <div class="metric">
            <strong>Efficiency:</strong> {efficiency_score}
        </div>
    </div>
    
    <div class="section">
        <h2>🔗 Integration Tests</h2>
        <p><strong>Sessions:</strong> {integration_sessions} | <strong>Success Rate:</strong> {integration_success:.1f}%</p>
        <p><strong>Duration:</strong> {integration_duration:.1f}s | <strong>Prompts:</strong> {integration_prompts}</p>
    </div>
    
    <div class="section">
        <h2>📊 Resource Tests</h2>
        <p><strong>Sessions:</strong> {resource_sessions} | <strong>Success Rate:</strong> {resource_success:.1f}%</p>
        <p><strong>Duration:</strong> {resource_duration:.1f}s | <strong>Prompts:</strong> {resource_prompts}</p>
    </div>
    
    <div class="section">
        <h2>🔍 API Evidence</h2>
        <table>
            <tr><th>Category</th><th>Payloads</th><th>Endpoints</th></tr>
            <tr><td>Integration</td><td>{integration_payloads}</td><td rowspan="2">{api_endpoints}</td></tr>
            <tr><td>Resources</td><td>{resource_payloads}</td></tr>
        </table>
    </div>
    
    <div class="section">
        <h2>💡 Recommendations</h2>
        {recommendations_html}
    </div>
</body>
</html>"""
        
        # Prepare data for template
        summary = self.report_data["summary"]
        integration = self.report_data["integration_tests"]
        resources = self.report_data["resource_tests"]
        evidence = self.report_data["api_evidence"]
        performance = self.report_data["performance_metrics"]
        
        success_rate = summary.get("overall_success_rate", 0)
        success_class = "success" if success_rate >= 90 else "warning" if success_rate >= 70 else "error"
        
        # Generate recommendations HTML
        recommendations_html = ""
        for rec in self.report_data["recommendations"]:
            priority_class = "error" if rec["priority"] == "high" else "warning"
            recommendations_html += f'<p class="{priority_class}"><strong>{rec["category"].title()}:</strong> {rec["message"]}</p>'
        
        return html_template.format(
            timestamp=self.timestamp.strftime("%Y-%m-%d %H:%M:%S"),
            period=self.period,
            total_tests=summary.get("total_tests", 0),
            success_rate=success_rate,
            success_class=success_class,
            total_payloads=evidence.get("total_payloads", 0),
            efficiency_score=performance.get("efficiency_score", "N/A"),
            integration_sessions=integration.get("total_sessions", 0),
            integration_success=integration.get("success_rate", 0),
            integration_duration=integration.get("total_duration", 0),
            integration_prompts=integration.get("total_prompts", 0),
            resource_sessions=resources.get("total_sessions", 0),
            resource_success=resources.get("success_rate", 0),
            resource_duration=resources.get("total_duration", 0),
            resource_prompts=resources.get("total_prompts", 0),
            integration_payloads=evidence.get("integration_payloads", 0),
            resource_payloads=evidence.get("resource_payloads", 0),
            api_endpoints=len(evidence.get("api_endpoints", [])),
            recommendations_html=recommendations_html
        )
    
    def generate_report(self):
        """Generate the complete test report"""
        print("🚀 Generating OpsRamp AI Agent Test Report...")
        
        # Analyze test results
        print("📊 Analyzing integration tests...")
        self.report_data["integration_tests"] = self.analyze_test_results("integration")
        
        print("📊 Analyzing resource tests...")
        self.report_data["resource_tests"] = self.analyze_test_results("resources")
        
        print("📊 Analyzing multi-provider tests...")
        self.report_data["multi_provider_tests"] = self.analyze_test_results("multi_provider")
        
        # Analyze API evidence
        print("🔍 Analyzing API evidence...")
        self.report_data["api_evidence"] = self.analyze_api_evidence()
        
        # Calculate performance metrics
        print("⚡ Calculating performance metrics...")
        self.report_data["performance_metrics"] = self.calculate_performance_metrics()
        
        # Generate summary
        integration = self.report_data["integration_tests"]
        resources = self.report_data["resource_tests"]
        multi_provider = self.report_data["multi_provider_tests"]
        
        total_tests = (integration.get("total_sessions", 0) + 
                      resources.get("total_sessions", 0) + 
                      multi_provider.get("total_sessions", 0))
        
        successful_tests = (integration.get("successful_sessions", 0) + 
                           resources.get("successful_sessions", 0) + 
                           multi_provider.get("successful_sessions", 0))
        
        overall_success_rate = (successful_tests / total_tests * 100) if total_tests > 0 else 0
        
        self.report_data["summary"] = {
            "total_tests": total_tests,
            "successful_tests": successful_tests,
            "failed_tests": total_tests - successful_tests,
            "overall_success_rate": overall_success_rate
        }
        
        # Generate recommendations
        print("💡 Generating recommendations...")
        self.report_data["recommendations"] = self.generate_recommendations()
        
        # Generate output
        timestamp_str = self.timestamp.strftime("%Y%m%d_%H%M%S")
        
        if self.format == "html":
            report_content = self.generate_html_report()
            filename = f"test_report_{self.period}_{timestamp_str}.html"
        elif self.format == "json":
            report_content = json.dumps(self.report_data, indent=2)
            filename = f"test_report_{self.period}_{timestamp_str}.json"
        else:  # text
            report_content = self.generate_text_report()
            filename = f"test_report_{self.period}_{timestamp_str}.txt"
        
        # Save report
        report_path = self.evidence_path / "test_reports" / filename
        report_path.parent.mkdir(parents=True, exist_ok=True)
        
        with open(report_path, 'w') as f:
            f.write(report_content)
        
        print(f"✅ Report generated: {report_path}")
        print(f"📊 Summary: {total_tests} tests, {overall_success_rate:.1f}% success rate")
        
        return report_path
    
    def generate_text_report(self):
        """Generate text format report"""
        lines = []
        lines.append("=" * 60)
        lines.append("🚀 OPSRAMP AI AGENT TEST REPORT")
        lines.append("=" * 60)
        lines.append(f"Generated: {self.timestamp.strftime('%Y-%m-%d %H:%M:%S')}")
        lines.append(f"Period: {self.period}")
        lines.append("")
        
        # Summary
        summary = self.report_data["summary"]
        lines.append("📊 EXECUTIVE SUMMARY")
        lines.append("-" * 30)
        lines.append(f"Total Tests: {summary.get('total_tests', 0)}")
        lines.append(f"Successful: {summary.get('successful_tests', 0)}")
        lines.append(f"Failed: {summary.get('failed_tests', 0)}")
        lines.append(f"Success Rate: {summary.get('overall_success_rate', 0):.1f}%")
        lines.append("")
        
        # Integration Tests
        integration = self.report_data["integration_tests"]
        lines.append("🔗 INTEGRATION TESTS")
        lines.append("-" * 30)
        lines.append(f"Sessions: {integration.get('total_sessions', 0)}")
        lines.append(f"Success Rate: {integration.get('success_rate', 0):.1f}%")
        lines.append(f"Duration: {integration.get('total_duration', 0):.1f}s")
        lines.append(f"Prompts: {integration.get('total_prompts', 0)}")
        lines.append("")
        
        # Resource Tests
        resources = self.report_data["resource_tests"]
        lines.append("📊 RESOURCE TESTS")
        lines.append("-" * 30)
        lines.append(f"Sessions: {resources.get('total_sessions', 0)}")
        lines.append(f"Success Rate: {resources.get('success_rate', 0):.1f}%")
        lines.append(f"Duration: {resources.get('total_duration', 0):.1f}s")
        lines.append(f"Prompts: {resources.get('total_prompts', 0)}")
        lines.append("")
        
        # API Evidence
        evidence = self.report_data["api_evidence"]
        lines.append("🔍 API EVIDENCE")
        lines.append("-" * 30)
        lines.append(f"Total Payloads: {evidence.get('total_payloads', 0)}")
        lines.append(f"Integration Payloads: {evidence.get('integration_payloads', 0)}")
        lines.append(f"Resource Payloads: {evidence.get('resource_payloads', 0)}")
        lines.append(f"API Endpoints: {len(evidence.get('api_endpoints', []))}")
        lines.append("")
        
        # Performance
        performance = self.report_data["performance_metrics"]
        lines.append("⚡ PERFORMANCE METRICS")
        lines.append("-" * 30)
        lines.append(f"Avg Response Time: {performance.get('avg_response_time', 0):.2f}s")
        lines.append(f"Prompts/Minute: {performance.get('prompts_per_minute', 0):.1f}")
        lines.append(f"Efficiency Score: {performance.get('efficiency_score', 'N/A')}")
        lines.append("")
        
        # Recommendations
        recommendations = self.report_data["recommendations"]
        if recommendations:
            lines.append("💡 RECOMMENDATIONS")
            lines.append("-" * 30)
            for rec in recommendations:
                priority = rec["priority"].upper()
                lines.append(f"[{priority}] {rec['category'].title()}: {rec['message']}")
            lines.append("")
        
        return "\n".join(lines)

def main():
    parser = argparse.ArgumentParser(description="Generate OpsRamp AI Agent Test Report")
    parser.add_argument(
        "--period",
        choices=["daily", "weekly", "monthly"],
        default="daily",
        help="Time period for report analysis"
    )
    parser.add_argument(
        "--format",
        choices=["html", "json", "text"],
        default="html",
        help="Report output format"
    )
    
    args = parser.parse_args()
    
    generator = TestReportGenerator(period=args.period, format=args.format)
    report_path = generator.generate_report()
    
    print(f"\n🎉 Test report generated successfully!")
    print(f"📄 Report location: {report_path}")

if __name__ == "__main__":
    main() 