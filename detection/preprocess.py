import matplotlib.pyplot as plt
import pandas as pd
import typing as tp
import cv2
import os
import shutil

from pathlib import Path

import constants
import utils


def fix_data(signature_ds_path: str):
    OFFSET_IDX = 2123

    for filename in ['/train.csv', '/test.csv', '/image_ids.csv']:
        csv_path = signature_ds_path + filename
        df: pd.DataFrame = pd.read_csv(csv_path)
        col_name = 'id' if filename == '/image_ids.csv' else 'image_id'
        df_len = len(df.index)
        prev = 0
        for i in range(df_len):
            image_id = df.at[i, col_name]
            if image_id < prev:
                image_id += OFFSET_IDX
            df.at[i, col_name] = image_id
            prev = image_id

        # new_file_path = signature_ds_path + filename.split('.')[0] + '_copy.csv'
        df.to_csv(csv_path, index=False)


def draw_boxes(df: pd.DataFrame, signature_ds_path: str, limit=None):
    df_len = len(df.index)
    cols = ['bbox', 'bbox_yolo', 'height', 'width', 'file_name']
    for i in range(df_len):
        if limit and i >= limit:
            break

        bbox, bbox_yolo, h, w, img_name = df.loc[i, cols]
        values = utils.get_coordinates_from_string(bbox)
        # normalized_yolo = get_coordinates_from_string(bbox_yolo)
        normalized_yolo = bbox_yolo
        print(normalized_yolo)

        a, b, c, d = values

        # Denormalize values
        x_min = int(a * w)
        y_min = int(b * h)
        x_max = int((a + c) * w)
        y_max = int((b + d) * h)
        print([values, h, w])
        print(x_max, x_min, y_max, y_min)

        x = int(normalized_yolo[0] * w)
        y = int(normalized_yolo[1] * h)
        x_width = int(normalized_yolo[2] * w)
        y_height = int(normalized_yolo[3] * h)
        print(x, y, x_width, y_height)

        img_path = f'{signature_ds_path}/images/{img_name}'
        img = cv2.imread(img_path)

        cv2.rectangle(img, pt1=(x_min, y_min), pt2=(x_max, y_max), color=(255, 0, 0), thickness=10)
        hw = int(x_width/2)
        hh = int(y_height/2)
        cv2.rectangle(img, pt1=(x-hw, y-hh), pt2=(x+hw, y+hh), color=(0, 255, 0), thickness=3)

        plt.imshow(img)
        plt.show()


def merge_datasets(
        filenames: tp.List[str],
        datasets: tp.List[pd.DataFrame],
        images_mapping: pd.DataFrame,
) -> tp.List[pd.DataFrame]:
    results = []
    for dataset, filename in zip(datasets, filenames):
        print(dataset)
        print(filename)
        # Remove id column to avoid naming conflict
        dataset.drop(columns=['id'], inplace=True)

        # Perform left join and remove
        merged_df = dataset.merge(images_mapping, left_on='image_id', right_on='id')
        merged_df.drop(columns=['id', 'category_id', 'image_id'], inplace=True)
        merged_df.sort_values(by=['file_name'], inplace=True)

        # Store the merged dataframe in a new file
        merged_df.to_csv(filename, index=False)
        results.append(merged_df)

    return results


def add_xywh_column(df: pd.DataFrame, save_path: str = None) -> pd.DataFrame:
    new_df = df.assign(
        bbox_yolo=[
            utils.bbox_to_yolo(
                utils.get_coordinates_from_string(df.loc[i, 'bbox']),
                input_normalized=True,
                output_normalized=True,
                image_size=df.loc[i, ['width', 'height']]
            )
            for i in range(len(df.index))
        ]
    )

    if save_path:
        new_df.to_csv(save_path, index=False)

    return new_df


def preprocess(signature_ds_path: str) -> tp.List[pd.DataFrame]:
    # Fix duplicated index error from original dataset
    fix_data(signature_ds_path)

    # Path to the csv that will be consumed
    training_csv = signature_ds_path + '/train.csv'
    testing_csv = signature_ds_path + '/test.csv'
    image_mapping_csv = signature_ds_path + '/image_ids.csv'

    # Load csv files
    train_df = pd.read_csv(training_csv)
    test_df = pd.read_csv(testing_csv)
    image_mapping_df = pd.read_csv(image_mapping_csv)

    # We only want the category that belongs to signature
    train_df = train_df[train_df['category_id'] == 1]
    test_df = test_df[test_df['category_id'] == 1]

    train_merged_path = signature_ds_path + '/train_merged.csv'
    test_merged_path = signature_ds_path + '/test_merged.csv'

    train_df, test_df = merge_datasets(
        [train_merged_path, test_merged_path],
        [train_df, test_df],
        image_mapping_df
    )

    # Add bbox yolo column
    train_df = add_xywh_column(train_df, train_merged_path)
    test_df = add_xywh_column(test_df, test_merged_path)

    return [train_df, test_df]


def generate_yolo_dataset(df: pd.DataFrame, signature_ds_dir: str, is_training=True, limit=None):
    dir_type = 'train' if is_training else 'valid'
    yolo_dataset_dir = f'{constants.DATASET_DIR}/yolo_data'
    images_dir = f'{yolo_dataset_dir}/images/{dir_type}'
    labels_dir = f'{yolo_dataset_dir}/labels/{dir_type}'

    # Create directories if they don't exist
    os.makedirs(images_dir, exist_ok=True)
    os.makedirs(labels_dir, exist_ok=True)

    df_len = len(df.index)
    columns = ['file_name', 'bbox_yolo']
    for i in range(df_len):
        if limit and i >= limit:
            break

        filename, bbox = df.loc[i, columns]
        coordinates = utils.get_coordinates_from_string(bbox)
        coordinates_str = ' '.join(map(str, coordinates))
        image_path = f'{signature_ds_dir}/images/{filename}'
        new_image_path = Path(f'{images_dir}/{filename}')
        label_path = Path(f"{labels_dir}/{filename.split('.')[0]}.txt")   # Change extension

        # Copy the image from other dataset
        if not new_image_path.exists():
            shutil.copy(image_path, new_image_path)

        if label_path.exists():
            with open(label_path, '+a') as f:
                f.write(f'0 {coordinates_str}\n')
        else:
            with open(label_path, '+w') as f:
                f.write(f'0 {coordinates_str}\n')


def run(read_from_file: bool = True):
    signature_ds_path = str(utils.download_data('victordibia', 'signverod'))

    if read_from_file:
        train_df = pd.read_csv(signature_ds_path + '/train_merged.csv')
        test_df = pd.read_csv(signature_ds_path + '/test_merged.csv')
    else:
        train_df, test_df = preprocess(signature_ds_path)

    print(train_df.head())
    print(test_df.head())

    generate_yolo_dataset(train_df, signature_ds_path)
    generate_yolo_dataset(test_df, signature_ds_path, is_training=False)


if __name__ == '__main__':
    run()
    # dry_run()
