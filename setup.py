from setuptools import setup, find_packages

setup(
    name="tineyelogger",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "aiohttp>=3.8.0",
    ],
    author="FlatMoon",
    description="Unified logging for OpenClaw Task Router",
    url="https://github.com/FlatMoon-Eva/TineyeLogger",
)
