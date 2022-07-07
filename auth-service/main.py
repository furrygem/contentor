from uuid import UUID
from sqlalchemy.orm import Session
from fastapi import Depends, FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse
from fastapi_another_jwt_auth import AuthJWT
from fastapi_another_jwt_auth.exceptions import AuthJWTException
from pydantic import BaseModel, BaseSettings
from db import SessionLocal
import models

from users_repo import create_user, get_user_by_username

app = FastAPI()


class UserCreateDTO(BaseModel):
    username: str
    password: str

    class Config:
        orm_mode = True


class UserFull(UserCreateDTO):
    id: UUID


class User(BaseModel):
    id: UUID
    username: str
    
    class Config:
        orm_mode = True


class Settings(BaseSettings):
    authjwt_algorithm: str = "EdDSA"
    authjwt_private_key: str | None = None
    authjwt_public_key: str | None = None

    authjwt_private_key_file: str = "key.pem"
    authjwt_public_key_file: str = "key.pub"

    def from_files(self):
        self.authjwt_private_key = open(self.authjwt_private_key_file, 'r')\
            .read()
        self.authjwt_public_key = open(self.authjwt_public_key_file, 'r')\
            .read()

    class Config:
        env_prefix = 'AUTH_SERVICE_'
        fields = {
            'authjwt_private_key_file': {
                'env': 'JWT_PRIVATE_KEY_FILE',
            },
            'authjwt_public_key_file': {
                'env': 'JWT_PUBLIC_KEY_FILE',
            },
        }


@AuthJWT.load_config
def get_config():
    s = Settings()
    s.from_files()
    return s


@app.exception_handler(AuthJWTException)
def authjwt_exception_handler(request: Request, exc: AuthJWTException):
    return JSONResponse(
        status_code=exc.status_code,
        content={"detail": exc.message}
    )


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


@app.post('/login')
def login(user: UserCreateDTO, Authorize: AuthJWT = Depends(), db: Session = Depends(get_db)):
    u = get_user_by_username(user.username, db)
    if u is None:
        raise HTTPException(status_code=401, detail="Bad username or password")

    if not u.check_pw(user.password):
        raise HTTPException(status_code=401, detail="Bad username or password")

    access_code = Authorize.create_access_token(subject=user.username,
                                                algorithm="EdDSA")
    return {"access_code": access_code}


@app.get('/user')
def user(Authorize: AuthJWT = Depends()):
    Authorize.jwt_required()
    current_user = Authorize.get_jwt_subject()
    return {"current_user": current_user}


@app.post('/register', response_model=User)
def register(user: UserCreateDTO, db: Session = Depends(get_db)):
    u = models.User(username=user.username)
    u.password = user.password
    u = create_user(u, db)
    user_response = User.from_orm(u)
    return user_response
