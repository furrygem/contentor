from uuid import UUID
from sqlalchemy.orm import Session

from models import User


def get_user_by_username(username: str, db: Session) -> User | None:
    user = db.query(User).filter(User.username == username).first()
    return user


def get_user_by_id(id: UUID, db: Session) -> User | None:
    user = db.query(User).filter(User.id == id).first()
    return user


def create_user(user: User, db: Session) -> User:
    db.add(user)
    db.commit()
    db.refresh(user)
    return user
