from app import db
from sqlalchemy.dialects.mysql import LONGBLOB

class Document(db.Model):
  id = db.Column(db.Integer, primary_key=True)
  document_content = db.Column(LONGBLOB)
  extension = db.Column(db.String(4))