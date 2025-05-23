from setuptools import setup, find_packages

setup(
    name="ormcp",
    version="0.1.0",
    packages=find_packages(where="src"),
    package_dir={"": "src"},
    install_requires=[
        "aiohttp>=3.8.0",
        "aiohttp-sse-client>=0.2.1",
        "requests>=2.28.0",
    ],
    python_requires=">=3.7",
) 