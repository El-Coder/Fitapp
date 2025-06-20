# Approach

I chose to use **React**, **Next.js**, and **TypeScript** on the frontend with **Golang** using **Echo** on the backend, connected to a **DynamoDB local** instance. Everything is wrapped in a Docker container. AWS provides a Docker setup and a Go SDK, which made integration straightforward.

# UI Hierarchy

There are two main UI components:

- **Gallery with linked fits and items**: Displays all fits and their currently linked items. Owns the display logic for each fit/item pair.
- **Link Fit to Item Form**: Lets the user select a fit and an item from dropdown menus and link them together. Owns input state for `selectedFit`, `selectedItem`, and the success message.

# State Management

I used **Zustand** for frontend state because it's lightweight and simpler to set up than Redux. Zustand stores arrays of fits and items, populated by API calls to the Go backend.

# API Design

I used Postman to make calls to populate data. When the user clicks the “Link” button the frontend requests `/api/fits` and `/api/items`, then stores the results in its Zustand store. User selections from the store are posted to `/api/link`, the backend persists this association in DynamoDB. After the link request returns, fetchGallery queries `/api/links/:fit_id` for every fit and updates the gallery state. If the request succeeds, it clears the selections and refreshes the gallery view.

# Persistent Data Storage

Each handler stores and retrieves records via DynamoDB’s APIs, so all fits, items, and their links are saved in these tables. The database files are mapped to the `docker/dynamodb` directory, it remains intact across sessions as long as the volume is preserved.

# Build

```bash
docker compose up --build
```

# Example API Calls
**POST**  
`http://localhost:8080/api/fits`

```json
{
  "fit_name": "Fit #1"
}
```
**POST** 
`http://localhost:8080/api/items`
```json
{
    "item_name": "Jacket #2"
}
```
# Tear Down
```bash
docker compose down
```

# Useful Links
https://dave.dev/blog/2021/07/14-07-2021-awsddb/

https://davidagood.com/dynamodb-local-go/

https://docs.aws.amazon.com/code-library/latest/ug/go_2_dynamodb_code_examples.html

https://www.postman.com/api-evangelist/amazon-web-services-aws/documentation/tuuvg4g/amazon-dynamodb

https://github.com/aws/aws-sdk-go-v2
