# Minimax Annotation System

This repository contains a Custom [Caddy](https://caddyserver.com/) Module written in Go that provides a web annotation layer. It acts as a middleware, injecting the necessary JavaScript and CSS into HTML responses to enable on-page annotations.

## 📂 Project Structure

- **`annotation/`**: The core Go package containing the Caddy module logic.
  - `api.go`: Handles the `GET`, `POST`, and `DELETE` requests for storing and retrieving annotations.
  - `handler.go`: Middleware logic to intercept HTTP responses and inject the annotation script.
  - `module.go`: Caddy module registration and configuration parsing.
- **`main.go`**: The entry point to build/run the custom Caddy binary with the annotation module.
- **`Caddyfile`**: The configuration file for the Caddy server.
- **`index.php`**: A simple example PHP application to demonstrate the annotation capability.

## 🚀 Prerequisites

- **Go**: Version 1.18 or higher.
- **PHP**: To run the example backend application.

## 🛠️ Usage Guide

### 1. Start the PHP Application (Backend Example)

The annotation module works by proxying requests to your backend application. For this example, we'll use the provided `index.php`.

Open a terminal and start the PHP built-in server on port 3000:

```bash
php -S localhost:3000
```

This serves the `index.php` file, which acts as the target page we want to annotate.

### 2. Run the Caddy Proxy

In a separate terminal, run the Caddy server using `go run`:

```bash
go run main.go run
```

This will search for the locally present `Caddyfile` and start the server on port **8080**.

### 3. Verify the Setup

Open your browser and navigate to:

👉 **http://localhost:8080**

You should see the content of `index.php` ("Test Content") along with the injected annotation toolbar (usually in the top-right corner).

## ⚙️ Configuration

The system is configured via the `Caddyfile`. Here is the example configuration:

```caddy
:8080 {
    # Enable annotation module
    annotation {
        enabled true
        api_endpoint /api/annotations
        script_path /annotation-assets
    }

    # Proxy to your local dev server (e.g., PHP app)
    reverse_proxy localhost:3000
}
```

### Directives:

- **`enabled`** (`true`|`false`): Toggles the injection mechanism.
- **`api_endpoint`**: The path usage for the internal API to save/load annotations.
- **`script_path`**: The path where the embedded JS and CSS assets will be served.

## 🔌 API Endpoints

The module provides an in-memory API to handle the annotations. Note that currently, this is **in-memory only**, meaning annotations are lost when the server restarts.

- `GET /api/annotations?url={url}`: Retrieve annotations for a specific URL.
- `POST /api/annotations`: Create a new annotation.
- `DELETE /api/annotations?id={id}&url={url}`: Delete an annotation.

## 📝 Example Flow

1. **User Request**: Browser requests `http://localhost:8080/`.
2. **Proxy**: Caddy forwards the request to the PHP server at `localhost:3000`.
3. **Response**: PHP server returns the HTML content of `index.php`.
4. **Injection**: The `annotation` module intercepts the HTML response and injects:
   - `<link rel="stylesheet" href="/annotation-assets/annotation.css">`
   - `<script src="/annotation-assets/annotation.js"></script>`
5. **Render**: The browser renders the page with the annotation UI.
