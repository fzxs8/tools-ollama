# Duola Desktop - User Guide

Welcome to Duola Desktop! This guide will help you understand and make the most of its powerful features.

## Table of Contents
1.  [Initial Setup: Configure Ollama Service](#1-initial-setup-configure-ollama-service)
2.  [Core Features](#2-core-features)
    *   [Chat Manager](#chat-manager)
    *   [Model Manager](#model-manager)
    *   [Model Market](#model-market)
    *   [Prompt Engineering](#prompt-engineering)
    *   [OpenAI Adapter](#openai-adapter)
    *   [API Debugger](#api-debugger)

---

### 1. Initial Setup: Configure Ollama Service

Before using any AI features, you need to configure your Ollama service.

1.  Navigate to the **Ollama Settings** page.
2.  The system automatically detects the local service. If your Ollama is running locally, no extra configuration is usually needed.
3.  Click **"Add Service"** to configure a remote Ollama server.
4.  Fill in the **Name** and **Service Address** (e.g., `http://192.168.1.10:11434`).
5.  You can check if the connection is successful using the **"Test"** button.
6.  You can set a service as the default for all features via the **"Set as Default"** option in the actions menu.

### 2. Core Features

#### Chat Manager

This is the main interface for your interactive conversations with AI models.

-   **Model Selection**: In the left sidebar, you can choose the server and model for the current conversation.
-   **Parameter Adjustment**: You can adjust model parameters like temperature and Top-P to control the diversity of the generated content.
-   **Streaming Chat**: The conversation content is displayed in real-time with a typewriter effect for a smooth experience.
-   **Conversation History**:
    -   All conversations are saved automatically.
    -   Click the history button in the top-right corner to load, edit, or delete previous conversations.
-   **System Prompt**: You can select a system prompt from the Prompt Engineering module to assign a role or provide context to the AI.

#### Model Manager

Manage the local models on all your Ollama services.

-   **Multi-Server Management**: Easily switch between different Ollama servers using the top dropdown menu to view their models.
-   **Model List**: A table clearly displays the model's name, size, modification time, and running status.
-   **Run/Stop Model**: In the model details, you can start or stop a model with a single click.
-   **Download New Model**:
    -   Click the "Download Model" button and enter a model name (e.g., `llama3`) to start the download.
    -   You can track the download progress in real-time via the "Download Queue".
-   **Test Model**: A test area is provided in the model details where you can input any text to quickly test the model's response.
-   **Delete Model**: You can delete local models that are no longer needed.

#### Model Market

Discover a vast number of online models from `ollama.com/library`.

-   **Browse and Search**: The page displays popular models by default, and you can also search for specific models using the search bar at the top.
-   **View Details**: Click "Details" to see detailed information about a model, including download counts and update times.
-   **Direct Download**: (This is implemented in the Model Manager) After finding a model you like in the Model Market, you can copy its name to the Model Manager page to download it.

#### Prompt Engineering

A powerful tool for creating, testing, and optimizing your prompts.

-   **Parallel Generation**: Input your core "idea," select multiple models (up to 3), and the system will generate different styles of prompts in parallel.
-   **Result Comparison**: Easily switch between and compare the results from different models using tabs in the results area on the right.
-   **Regenerate**: Not satisfied with a model's result? You can "Regenerate" it individually.
-   **Save and Manage**:
    -   Satisfactory results can be "Saved" with one click.
    -   Through the "Manage" panel, you can view, search, edit, and delete all your saved prompts.

#### OpenAI Adapter

Transforms your local Ollama service into an OpenAI API-compatible endpoint, allowing seamless integration with more third-party applications (like LobeChat).

1.  Go to the **OpenAI Adapter** page.
2.  **Configure**: Set the listening IP and port for the service, and select a target Ollama service to proxy.
3.  **Start**: Click the "Start Service" toggle.
4.  **Usage**: Once the service is running, you can use `http://<Your-IP>:<Your-Port>` as the OpenAI API address in other applications.
5.  **Real-time Logs**: Monitor all API requests and responses in real-time through the "View Logs" drawer.
6.  **API Docs & Debugger**: The "View API" drawer not only provides `curl` examples but also includes a built-in API debugger for testing the endpoint directly within the app.

#### API Debugger

A professional, Postman-like tool specifically for debugging Ollama's native APIs.

-   **Preset Endpoints**: The left-side list presets all official Ollama APIs. Clicking one automatically fills in all parameters.
-   **Full Customization**: You can modify the URL, request method, query parameters, headers, and body.
-   **Send and View**: After clicking "Send," you can clearly see the formatted response body, headers, status code, and request duration in the response area below.
