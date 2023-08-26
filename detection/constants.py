import os 

DIR_ENV = 'PYAPP_DIR'
APP_DIR = os.getenv(DIR_ENV)
DATASET_DIR = APP_DIR + '/datasets/'
SECRETS_PATH = APP_DIR + '/secrets.json'