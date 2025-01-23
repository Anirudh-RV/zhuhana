import os

ENV = os.environ.get("ENV")

if ENV == "dev":
    DJANGO_SECRET_KEY = os.environ.get("DJANGO_SECRET_KEY")
    DEBUG = os.environ.get("DEBUG")
    POSTGRES_DB = os.environ.get("POSTGRES_DB")
    POSTGRES_USER = os.environ.get("POSTGRES_USER")
    POSTGRES_PASSWORD = os.environ.get("POSTGRES_PASSWORD")
    POSTGRES_HOST = os.environ.get("POSTGRES_HOST")
    POSTGRES_PORT = os.environ.get("POSTGRES_PORT")
else:
    # TODO: Implement Azure Secrets Manager
    pass
