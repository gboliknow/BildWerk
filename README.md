# Image Processing Service (BildWerk)

## Overview
BildWerk is an image processing backend service built in Golang. It allows users to upload images, perform various transformations (resize, crop, rotate, watermark, etc.), and retrieve images efficiently. The service includes authentication, image storage, and transformation capabilities similar to Cloudinary.

## Features
- **User Authentication**
  - Sign up, log in, and secure endpoints using JWT authentication.
- **Image Management**
  - Upload images.
  - Retrieve images in different formats.
  - List all uploaded images with metadata.
- **Image Transformations**
  - Resize, crop, rotate, watermark.
  - Flip, mirror, compress, change format (JPEG, PNG, etc.).
  - Apply filters like grayscale, sepia, etc.
- **Efficient Image Retrieval**
  - Optimized storage and caching.

## Tech Stack
- **Backend:** Golang (Gin framework)
- **Database:** PostgreSQL (GORM ORM)
- **Storage:** Local storage (can be extended to S3, Cloud Storage, etc.)
- **Authentication:** JWT (JSON Web Tokens)
- **Logging:** Zerolog

## Installation
### Prerequisites
- Go 1.24+
- PostgreSQL

### Clone the Repository
```sh
git clone https://github.com/yourusername/bildwerk.git
cd bildwerk
```

### Install Dependencies
```sh
go mod tidy
```

## Configuration
Create a `.env` file and set the environment variables:
```
DB_HOST=localhost
DB_USER=youruser
DB_PASSWORD=yourpassword
DB_NAME=bildwerk
DB_PORT=5432
JWT_SECRET=your_secret_key
```

## Running the Application
```sh
go run main.go
```

## API Endpoints
### Authentication
- **POST /auth/register** – Register a new user.
- **POST /auth/login** – Login and get a JWT token.

### Image Management
- **POST /images/upload** – Upload an image.
- **GET /images/{id}** – Retrieve an image.
- **GET /images** – List all user-uploaded images.

### Image Transformations
- **GET /images/{id}/resize?width=300&height=300** – Resize image.
- **GET /images/{id}/crop?x=100&y=50&width=200&height=200** – Crop image.
- **GET /images/{id}/rotate?angle=90** – Rotate image.
- **GET /images/{id}/watermark?text=MyBrand** – Add watermark.

## Contribution
1. Fork the repository.
2. Create a new branch (`feature-branch`).
3. Commit changes and push to your fork.
4. Submit a pull request.

## License
MIT License. See `LICENSE` for details.

