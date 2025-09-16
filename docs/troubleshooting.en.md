# Duola Desktop - Troubleshooting

This document lists common issues and their solutions.

### 1. Model List Shows "Invalid Date"

-   **Problem**: On the "Model Management" page, the "Modified At" column for models displays "Invalid Date".
-   **Cause**: Mismatch between frontend and backend field names. The frontend expected `modifiedAt`, but the backend returned `modified_at`.
-   **Solution**: This has been fixed in the code. Ensure the frontend data structure and the backend Go struct's `json` tag match completely, standardizing on `modified_at`.

### 2. Running Status of a Model is Lost After Refresh

-   **Problem**: On the "Model Management" page, a running model's status changes back to "Stopped" after manually refreshing the list.
-   **Cause**: The running state was previously maintained only temporarily on the frontend and was overwritten by the raw data from the backend (which did not include running status) upon refresh.
-   **Solution**: The Go backend has been made the single source of truth for model running status. The backend maintains a `runningModels` list, and every time the frontend requests the model list, the backend appends the accurate running status.

### 3. Application Crashes When Clicking "Run" or "Test"

-   **Problem**: The application crashes when clicking "Run Model" on the "Model Management" page.
-   **Cause**: The core backend component `httpClient` was not properly initialized, leading to a `nil pointer dereference` when used.
-   **Solution**: It has been ensured that all core components are correctly initialized in the backend's constructor.

### 4. OpenAI Adapter or API Debugger Request Fails with "CORS" or "Load failed" Error

-   **Problem**: When using the API Debugger or an external application to call the OpenAI Adapter endpoint, the browser console reports a CORS (Cross-Origin Resource Sharing) error.
-   **Cause**: The backend HTTP service was not correctly handling the browser's `OPTIONS` preflight request and was not adding the necessary `Access-Control-Allow-*` headers to the response.
-   **Solution**: A CORS middleware has been added to the Go backend's HTTP service. This middleware:
    1.  Correctly responds to `OPTIONS` requests (returns 200 OK).
    2.  Adds headers like `Access-Control-Allow-Origin: *`, `Access-Control-Allow-Methods`, and `Access-Control-Allow-Headers` to all responses.

### 5. Model Download Reports "Success" Even When an Error Occurs

-   **Problem**: When downloading an incompatible or non-existent model, the backend has actually failed, but the frontend UI still shows the download as successful.
-   **Cause**: The backend's stream processing logic failed to correctly capture and handle error messages embedded in the download stream, causing it to send a "done" event regardless of success or failure.
-   **Solution**: The backend's error detection logic has been enhanced. Now, when an error appears in the download stream, it immediately interrupts the process and sends a dedicated error event to the frontend instead of a completion event.

### 6. After Deleting All Service Configurations, a "Local Service" is Automatically Re-created on Restart

-   **Problem**: Even if the user clears all servers in "Service Settings," a "Local Service" configuration reappears the next time the application starts.
-   **Cause**: The backend had "clever" logic that automatically created a default local service when no configuration was detected.
-   **Solution**: This logic has been removed. The application now accurately reflects the user's configuration state. If the list is empty, the frontend will prompt the user to add a service in the settings page.
