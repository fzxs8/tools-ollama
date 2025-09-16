# Duola Desktop - Development Guide

This guide provides instructions for contributors who wish to participate in the project's development.

## 1. Tech Stack

-   **Backend**: Go
-   **Frontend**: Vue 3 (Composition API) + TypeScript
-   **UI Framework**: Element Plus
-   **Desktop App Framework**: Wails v2
-   **Data Storage**: Embedded `duolasdk` storage feature (based on SQLite)

## 2. Architecture Overview

The project backend uses a modular design, where core business logic is sunk into respective managers (`Manager`), such as `ModelManager`, `ChatManager`, etc. The main entry point, `app.go`, is responsible for dependency injection and routing, injecting core services like `store` and `logger` into each manager to achieve low coupling.

The frontend uses a component-based architecture based on Vue 3, with pages corresponding to modules. State management primarily uses Vue 3's reactive APIs (`ref`, `reactive`). Pinia is planned for future use for cross-component or global state.

The frontend and backend communicate via the Wails v2 framework. Backend Go methods are bound to the frontend, allowing direct calls. For scenarios requiring real-time data, like streaming output and logs, Wails' event system (`runtime.EventsEmit` and `runtime.EventsOn`) is used.

## 3. Development Standards

Please strictly adhere to the standards defined in `docs/设计/通用文档/开发规范.md`. Below is a summary of key points:

### Coding Standards
-   **Go**:
    -   All exported struct fields must have a `json` tag.
    -   All functions that can fail should return an `error`.
    -   Follow the general coding style of the Go community.
-   **Vue/TypeScript**:
    -   Use the `<script setup lang="ts">` syntax.
    -   Organize code within components in the following order: `imports` -> `interfaces` -> `consts` -> `reactive data` -> `computed` -> `methods` -> `lifecycle hooks`.
    -   Naming: Use `PascalCase` for components, `camelCase` for variables and methods.
    -   Strong Typing: Use explicit TypeScript types whenever possible; avoid `any`.

### UI Development Standards
-   **Component Library**: Prioritize using components from `Element Plus`.
-   **Styling**: Use Scoped CSS or CSS Modules to avoid style pollution.
-   **Responsive Design**: Layouts should adapt to different window sizes.
-   **Feedback**: All asynchronous operations must have a clear loading state (`loading`), and the result of operations should be communicated to the user via `ElMessage`.

### Error Handling
-   **Frontend**: All API calls to the backend must be wrapped in `try...catch` blocks. The `catch` block should use `ElMessage.error` to display a clear error message to the user.
-   **Backend**: Error messages returned to the frontend should be clear and specific, without exposing internal implementation details.

## 4. Reference for Key Issues and Solutions

During development, we have accumulated solutions to common problems, which are documented in the `06.待办事项.md` (TODOs) or specific fix logs for each module. For example:

-   **Wails Debug Configuration**: How to configure `wails.json` to use `dlv` for backend debugging.
-   **State Synchronization Issues**: How to ensure state consistency across pages and between the frontend and backend (e.g., model running status, default service configuration).
-   **CORS Issues**: How to correctly handle browser CORS preflight requests in the Go backend.
-   **Data Migration**: How to provide seamless migration of old user data when upgrading the storage structure.

Before starting development on a new feature or fixing a bug, please review the design documents of the relevant module, especially the "TODO" or "Fix Log" sections, to avoid common pitfalls.

## 5. Contribution Workflow

1.  Fork this repository.
2.  Create a new feature branch (`git checkout -b feature/xxx`) or fix branch (`git checkout -b fix/xxx`).
3.  Develop the code, ensuring adherence to the development standards.
4.  Write unit tests for new core logic.
5.  Commit your code and create a Pull Request.
6.  Clearly describe the content and purpose of your changes in the PR description.
7.  Wait for a code review and make changes based on the feedback.
