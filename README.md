# Strong Password Recommendation Backend

This project provides a backend API for recommending steps to strengthen passwords based on defined criteria. The API is built using Go with the Gin framework and is containerized using Docker.

## Features

- **Password Strength Check**: Determines the minimum number of steps required to make a password strong.
- **API Endpoint**: `/api/strong_password_steps`
- **Dockerized**: Uses Docker for containerization.
- **PostgreSQL Database**: Logs requests and responses into a PostgreSQL database.

## API Specification

### Endpoint

- **Base URL**: `/api/strong_password_steps`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "init_password": "your_password_here"
  }
- **Request Body**:
```json
{
  "num_of_steps": number_of_steps_required
}
