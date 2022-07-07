from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from pydantic import BaseSettings


class Settings(BaseSettings):
    db_driver: str = "mysql"
    db_username: str = "root"
    db_password: str = "mysql"
    db_host: str = "127.0.0.1:3306"
    db_database: str = "auth-service"

    @property
    def db_url(self):
        return f"{self.db_driver}://{self.db_username}:{self.db_password}@{self.db_host}/{self.db_database}"

    class Config:
        fields = {
            'db_driver': {
                'env': "AUTH_SERVICE_DB_DRIVER",
            },
            'db_username': {
                'env': "AUTH_SERVICE_DB_USER",
            },
            'db_password': {
                'env': "AUTH_SERVICE_DB_PASS",
            },
            'db_host': {
                'env': "AUTH_SERVICE_DB_HOST",
            },
            'db_database': {
                'env': "AUTH_SERVICE_DB_NAME",
            },
        }


settings = Settings()
engine = create_engine(settings.db_url)
SessionLocal = sessionmaker(engine, autocommit=False, autoflush=False)

Base = declarative_base()
