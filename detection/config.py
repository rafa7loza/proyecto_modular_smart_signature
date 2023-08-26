import utils 

class DevConfig(object):
  ENV = 'development'
  DEBUG = True
  QLALCHEMY_TRACK_MODIFICATIONS = False
  SQLALCHEMY_DATABASE_URI = utils.generateDBUri()
  