"""create users table

Revision ID: 161281c083f6
Revises: 
Create Date: 2022-07-06 11:21:22.094395

"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.mysql import BINARY


# revision identifiers, used by Alembic.
revision = '161281c083f6'
down_revision = None
branch_labels = None
depends_on = None


def upgrade() -> None:
    op.create_table(
        "users",
        sa.Column('id', BINARY(16), primary_key=True, index=True),
        sa.Column('username', sa.VARCHAR(50), index=True),
        sa.Column('password', sa.VARCHAR(100)),
    )


def downgrade() -> None:
    op.drop_table("users")
