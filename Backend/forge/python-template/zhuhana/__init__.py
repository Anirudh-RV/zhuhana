# zhuhana/__init__.py

from .zhuhana import ZhuhanaClass

def init(api_endpoint: str, token: str) -> ZhuhanaClass:
    return ZhuhanaClass(api_endpoint=api_endpoint, token=token)
