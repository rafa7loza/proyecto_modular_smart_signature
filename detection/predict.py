from ultralytics import YOLO
import cv2


import constants

def do_prediction(img):
    models_dir = constants.APP_DIR + '/models/'
    chosen_model = 'best.pt'
    model = YOLO(models_dir + chosen_model)
    results = model(img, task='detect')
    res = results[0]
    res_plotted = res.plot()
    return res_plotted

def main():
    models_dir = constants.APP_DIR + '/models/'
    chosen_model = 'best.pt'
    images_dir = constants.APP_DIR + '/datasets/'
    image_filenames = ['parcial_1', 'parcial_2', 'parcial_3', 'Reporte_final_ss']
    image_filenames = list(map(lambda x: x + '.jpg', image_filenames))
    model = YOLO(models_dir + chosen_model)

    for filename in image_filenames:
        image_path = images_dir + filename
        results = model(image_path, task='detect')
        print(results)
        for res in results:
            result_path = images_dir + '/results/' + filename
            print(f"Saving result in ${result_path}")
            res_plotted = res.plot()
            print(res_plotted)
            cv2.imwrite(result_path, res_plotted)
            # cv2.imshow("result", res_plotted)
            # cv2.waitKey(0)   #wait for a keyboard input
            # cv2.destroyAllWindows()


if __name__ == '__main__':
    main()

