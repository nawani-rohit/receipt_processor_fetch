# Receipt Processor Challenge

This repository contains my Python-Flask solution for a take-home exam from Fetch Rewards. The application processes receipts, calculates points based on a set of rules, and exposes two API endpoints. The project is containerized using Docker for easy deployment and testing.

## API Endpoints

### POST `/receipts/process`

- **Description:**  
  Accepts a JSON receipt object, processes it, applies the point calculation rules, and returns a unique ID that represents the processed receipt.

- **Request Payload:**  
  A JSON object representing the receipt. For example:
  ```json
  {
      "retailer": "Walgreens",
      "purchaseDate": "2022-01-02",
      "purchaseTime": "08:13",
      "total": "2.65",
      "items": [
          {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
          {"shortDescription": "Dasani", "price": "1.40"}
      ]
  }
  ```

- **Response:**  
  A JSON object containing a unique ID. For example:
  ```json
  { "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
  ```

### GET `/receipts/{id}/points`

- **Description:**  
  Returns the total points associated with the specified receipt ID.

- **Response:**  
  A JSON object with the calculated points. For example:
  ```json
  { "points": 32 }
  ```

## Points Calculation Rules

The application calculates points using the following rules:

- **Retailer Name:**  
  1 point for every alphanumeric character in the retailer name.

- **Total Amount:**  
  - 50 points if the total is a round dollar amount (no cents).  
  - 25 points if the total is a multiple of 0.25.  
  - 5 bonus points if the total is greater than 10.00 (if applicable).

- **Item Count:**  
  5 points for every two items on the receipt.

- **Item Description:**  
  For each item, if the trimmed length of the item’s short description is a multiple of 3, multiply the item’s price by 0.2 and round up to the nearest integer. The result is the number of points earned for that item.

- **Purchase Date:**  
  6 points if the day in the purchase date is odd.

- **Purchase Time:**  
  10 points if the purchase time is after 2:00 PM and before 4:00 PM.

## Project Structure

```
receipt-processor/
├── app.py              # Main Flask application; contains the endpoints and points calculation logic.
├── requirements.txt    # Lists Python dependencies (Flask, etc.)
├── Dockerfile          # Dockerfile to build and run the application.
└── README.md           # This file.
```

## Docker Instructions

### Build the Docker Image

From the root directory (where the Dockerfile is located), run:

```bash
docker build -f Dockerfile -t receipt-processor .
```

### Run the Container

Start the container by mapping the container's port 5000 to your host's port 5000:

```bash
docker run -d -p 5000:5000 --name my-receipt-app receipt-processor
```

### Verify the Container is Running

Use the following command to check that the container is up and running:

```bash
docker ps
```

### Stop the Container

To stop the running container, use:

```bash
docker stop my-receipt-app
```

## Testing with Postman

Follow these steps to test the API endpoints using Postman:

1. **Test the POST `/receipts/process` Endpoint:**

   - Open Postman and create a new **POST** request to:
     ```
     http://localhost:5000/receipts/process
     ```
   - In the **Body** tab, select **raw** and choose **JSON** as the type.
   - Paste a sample JSON receipt:
     ```json
     {
         "retailer": "Walgreens",
         "purchaseDate": "2022-01-02",
         "purchaseTime": "08:13",
         "total": "2.65",
         "items": [
             {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
             {"shortDescription": "Dasani", "price": "1.40"}
         ]
     }
     ```
   - Click **Send**. You should receive a `200 OK` response with a JSON body containing a unique ID, for example:
     ```json
     { "id": "some-uuid-value" }
     ```

2. **Test the GET `/receipts/{id}/points` Endpoint:**

   - Copy the `id` from the POST response.
   - Create a new **GET** request to:
     ```
     http://localhost:5000/receipts/<copied-id>/points
     ```
     Replace `<copied-id>` with the actual ID.
   - Click **Send**. You should receive a JSON response with the calculated points, for example:
     ```json
     { "points": 32 }
     ```

## Notes

- The application uses in-memory storage for receipts, so data will not persist across restarts.
- The code is implemented using Python and Flask.
- The Docker container simplifies deployment and testing.
- For any issues or questions, please refer to this README or contact me.

## License

This project is licensed under the [MIT License](LICENSE).

## Contact

For further inquiries, please contact [Your Name] at [your.email@example.com].
```

---

### How to Customize

- **Project Title & Description:**  
  Modify the title, description, and contact information to match your details.

- **API Details:**  
  Update the sample JSON and responses as needed to reflect any changes in your implementation.

- **Docker Commands:**  
  Adjust the Docker instructions if your container configuration changes (e.g., if you decide to use a different port).

This README file now incorporates all the provided information—from endpoints and points calculation rules to Docker build/run instructions—ensuring that reviewers and users have clear, step-by-step guidance on building, running, and testing your Receipt Processor Challenge application.
