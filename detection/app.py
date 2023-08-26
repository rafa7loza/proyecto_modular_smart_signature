from flask import Flask, request
from flask.json import jsonify
from werkzeug.exceptions import BadRequest
from flask_sqlalchemy import SQLAlchemy
from pdf2image import convert_from_bytes
from PIL import Image
import io
import base64
import numpy as np
import cv2

import constants
from predict import do_prediction


app = Flask(__name__)
app.config.from_object("config.DevConfig")
db = SQLAlchemy(app)

from app import db
from sqlalchemy.dialects.mysql import LONGBLOB

class Documents(db.Model):
  id = db.Column(db.Integer, primary_key=True)
  document_content = db.Column(LONGBLOB)
  extension = db.Column(db.String(4))

class ProcessedDocument(db.Model):
  __tablename__ = 'processed_documents'

  id = db.Column(db.Integer, primary_key=True)
  document_content = db.Column(LONGBLOB)
  document_id = db.Column(db.Integer)

@app.route('/detect/<int:doc_id>', methods=['POST'])
def detect(doc_id: int = None):
  if doc_id is None:
    raise BadRequest("Missing document ID")

  print(doc_id)
  doc: Documents = Documents.query.get(doc_id)
  print(doc)

  if doc.extension == '.pdf':
    img = convert_from_bytes(
      doc.document_content, 
      output_folder=constants.DATASET_DIR,
      fmt='jpeg'
    )
  else:
    img = Image.open(io.BytesIO(doc.document_content))

  print(img)
  img_array = do_prediction(img)
  pil_image = Image.fromarray(img_array[..., ::-1])
  byte_stream = io.BytesIO()
  pil_image.save(byte_stream, format='JPEG')
  byte_stream = byte_stream.getvalue()

  processed_doc = ProcessedDocument(
    document_content = byte_stream,
    document_id = doc_id
  )
  print(processed_doc)

  try:
    db.session.add(processed_doc)
  except Exception as e:
    return str(e)
  db.session.commit()

  return "File uploaded and processed successfully!"

@app.route('/fetch/processed/<int:doc_id>', methods=['GET'])
def fetch_processed(doc_id: int = None):
  if doc_id is None:
    raise BadRequest("Missing document ID")

  doc: ProcessedDocument = ProcessedDocument.query.get(doc_id)
  if doc is None:
    raise BadRequest(f"Document with ID {doc_id} not found")
  print(doc)

  # result_path = constants.DATASET_DIR + '/results/tmp_copy.jpeg'
  bytes = io.BytesIO(doc.document_content)
  # pil_image = Image.open(bytes)
  # pil_image.save(result_path)
  img_str = base64.b64encode(bytes.getvalue())


  return jsonify(
    base64 = str(img_str)
  )

if __name__ == '__main__':
  app.run(port=7820)
