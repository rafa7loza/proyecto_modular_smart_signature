import typing as tp
# import kaggle
import zipfile
import os
import json

from pathlib import Path

import constants


def get_coordinates_from_string(coordinates: str) -> tp.List[float]:
    if type(coordinates) == list:
        return list(map(float, coordinates))
    return list(
        map(
            float, coordinates.removeprefix('[').removesuffix(']').split(', ')
        )
    )


def bbox_to_yolo(
        coordinates: tp.List[float],
        image_size: tp.List[float] = None,
        input_normalized=True,
        output_normalized=True
) -> tp.List[float]:
    if len(coordinates) != 4:
        raise Exception

    if input_normalized:
        a, b, c, d = coordinates
        if image_size is None or len(image_size) != 2:
            raise Exception

        w, h = image_size

        # Denormalize values
        x_min = int(a * w)
        y_min = int(b * h)
        x_max = int((a + c) * w)
        y_max = int((b + d) * h)
    else:
        x_min, y_min, x_max, y_max = coordinates

    yolo_format = [
        int((x_max + x_min) / 2),
        int((y_max + y_min) / 2),
        x_max - x_min,
        y_max - y_min
    ]

    if output_normalized:
        if image_size is None or len(image_size) != 2:
            raise Exception
        w, h = image_size
        return [
            yolo_format[i]/w if i % 2 == 0 else yolo_format[i]/h
            for i in range(len(yolo_format))
        ]

    return yolo_format


def download_data(owner: str, dataset_name: str) -> Path:
    dataset = f'{owner}/{dataset_name}'
    dataset_path = Path(constants.DATASET_DIR)/Path(dataset_name)

    if not dataset_path.exists():
        kaggle.api.dataset_download_cli(dataset)
        zip_filename = f'{dataset_name}.zip'
        zipfile.ZipFile(zip_filename).extractall(dataset_path)
        os.remove(zip_filename)

    return dataset_path


def generate_yaml(filename: str):
    content = f"""train: {constants.DATASET_DIR}/yolo_data/images/train
val: {constants.DATASET_DIR}/yolo_data/images/valid

nc: 4

names: ['Signature', 'Initials', 'Redaction', 'Date']
    """

    with open(filename, '+w') as f:
        f.write(content)
    return filename

def generateDBUri():
    with open(constants.SECRETS_PATH, 'r') as f:
        secrets = json.load(f)
    dbs = secrets['db']
    return f"mysql+pymysql://{dbs['user']}:{dbs['password']}@{dbs['ipAddress']}:{dbs['port']}/{dbs['schema']}"
