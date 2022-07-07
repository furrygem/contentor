from uuid import UUID, uuid4
import bcrypt
from sqlalchemy import VARCHAR, Column
from sqlalchemy.types import TypeDecorator
from sqlalchemy.dialects.mysql import BINARY
from db import Base


class BinUUID(TypeDecorator):
    impl = BINARY(16)

    def process_bind_param(self, value, dialect) -> bytes:
        try:
            return value.bytes
        except AttributeError:
            try:
                return UUID(value).bytes
            except TypeError:
                return value

    def process_result_value(self, value, dialect) -> UUID:
        return UUID(bytes=value)


class User(Base):
    __tablename__ = "users"
    id = Column(BinUUID, default=uuid4(), primary_key=True)
    username = Column(VARCHAR(50), index=True)
    _password = Column(VARCHAR(100), name="password")

    @property
    def password(self):
        return self._password

    @password.setter
    def password(self, value: str):
        salt = bcrypt.gensalt()
        self._password = bcrypt.hashpw(value.encode(), salt)

    def check_pw(self, password: str) -> bool:
        return bcrypt.checkpw(password.encode(), bytes(self.password, "utf8"))
