# Online Photo Editor REST API

## Overview

The Online Photo Editor REST API is a robust and scalable solution for editing images via HTTP requests. This API provides various endpoints to perform common image editing operations such as uploading, cropping, resizing, converting, blurring, adjusting brightness, contrast, gamma, saturation, and sharpening.

## Features

- **Image Upload**: Upload images to the server.
- **Image Cropping**: Crop images to specified dimensions.
- **Image Resizing**: Resize images to specified dimensions.
- **Image Conversion**: Convert images between different formats.
- **Image Blurring**: Apply blur effects to images.
- **Brightness Adjustment**: Adjust the brightness of images.
- **Contrast Adjustment**: Adjust the contrast of images.
- **Gamma Correction**: Apply gamma correction to images.
- **Saturation Adjustment**: Adjust the saturation of images.
- **Sharpening**: Apply sharpening effects to images.
- **Image Processing**: Apply a sequence of image processing operations.

## Getting Started

### Prerequisites

- Go 1.23.3 or later

### Installation

1. **Clone the repository**:

   ```sh
   git clone https://github.com/yourusername/online-photo-editor.git
   cd online-photo-editor
   ```

2. **Build the project**:

   ```sh
   go build -o online-photo-editor
   ```

3. **Run the server**:
   ```sh
   ./online-photo-editor
   ```

### Configuration

The application uses a configuration file to set various parameters. You can create a `config.yaml` file with the following structure:

```yaml
env: "local" # Can be "local", "dev", or "prod"
address: ":8080"
storageImagePath: "/path/to/image/storage"
httpServer:
  timeout: 30s
  idleTimeout: 60s
```

### Environment Variables

You can also set environment variables to override the configuration:

- `ENV`: The environment (local, dev, prod)
- `ADDRESS`: The address to bind the server to
- `STORAGE_IMAGE_PATH`: The path to store images
- `HTTP_SERVER_TIMEOUT`: The HTTP server timeout
- `HTTP_SERVER_IDLE_TIMEOUT`: The HTTP server idle timeout

## API Endpoints

### Image Upload

- **URL**: `/image`
- **Method**: `POST`
- **Description**: Upload an image to the server.
- **Request Body**: Form data with the image file.
- **Response**: JSON object with the image URL.

### Image Cropping

- **URL**: `/image/crop`
- **Method**: `POST`
- **Description**: Crop an image to specified dimensions.
- **Request Body**: JSON object with crop parameters.
- **Response**: JSON object with the cropped image URL.

### Image Resizing

- **URL**: `/image/resize`
- **Method**: `POST`
- **Description**: Resize an image to specified dimensions.
- **Request Body**: JSON object with resize parameters.
- **Response**: JSON object with the resized image URL.

### Image Conversion

- **URL**: `/image/convert`
- **Method**: `POST`
- **Description**: Convert an image between different formats.
- **Request Body**: JSON object with conversion parameters.
- **Response**: JSON object with the converted image URL.

### Image Blurring

- **URL**: `/image/blur`
- **Method**: `POST`
- **Description**: Apply blur effects to an image.
- **Request Body**: JSON object with blur parameters.
- **Response**: JSON object with the blurred image URL.

### Brightness Adjustment

- **URL**: `/image/brightness`
- **Method**: `POST`
- **Description**: Adjust the brightness of an image.
- **Request Body**: JSON object with brightness parameters.
- **Response**: JSON object with the adjusted image URL.

### Contrast Adjustment

- **URL**: `/image/contrast`
- **Method**: `POST`
- **Description**: Adjust the contrast of an image.
- **Request Body**: JSON object with contrast parameters.
- **Response**: JSON object with the adjusted image URL.

### Gamma Correction

- **URL**: `/image/gamma`
- **Method**: `POST`
- **Description**: Apply gamma correction to an image.
- **Request Body**: JSON object with gamma parameters.
- **Response**: JSON object with the corrected image URL.

### Saturation Adjustment

- **URL**: `/image/saturation`
- **Method**: `POST`
- **Description**: Adjust the saturation of an image.
- **Request Body**: JSON object with saturation parameters.
- **Response**: JSON object with the adjusted image URL.

### Sharpening

- **URL**: `/image/sharpen`
- **Method**: `POST`
- **Description**: Apply sharpening effects to an image.
- **Request Body**: JSON object with sharpen parameters.
- **Response**: JSON object with the sharpened image URL.

### Image Processing

- **URL**: `/image/process`
- **Method**: `POST`
- **Description**: Apply a sequence of image processing operations.
- **Request Body**: JSON object with processing parameters.
- **Response**: JSON object with the processed image URL.

## Logging

The application uses structured logging with different handlers based on the environment:

- **Local**: Pretty-printed logs with debug level.
- **Dev**: JSON-formatted logs with debug level.
- **Prod**: Text-formatted logs with info level.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or support, please contact [mud.runner@bk.ru](mailto:mud.runner@bk.ru).

---

Thank you for using the Online Photo Editor REST API!
