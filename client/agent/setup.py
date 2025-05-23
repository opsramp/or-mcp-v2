from setuptools import setup, find_packages

setup(
    name="opsramp-agent",
    version="0.1.0",
    description="OpsRamp AI Agent for interacting with OpsRamp MCP",
    author="OpsRamp",
    author_email="info@opsramp.com",
    packages=find_packages(where="src"),
    package_dir={"": "src"},
    python_requires=">=3.8",
    install_requires=[
        "aiohttp>=3.8.0",
        "aiohttp-sse-client>=0.2.1",
        "requests>=2.28.0",
        "python-dotenv>=1.0.0",
    ],
    extras_require={
        "openai": ["openai>=1.0.0"],
        "anthropic": ["anthropic>=0.5.0"],
        "all": ["openai>=1.0.0", "anthropic>=0.5.0"],
    },
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
    ],
)
