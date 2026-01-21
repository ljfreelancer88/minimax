package annotation

import (
    "bytes"
    "net/http"
    "strings"
    
    "github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// responseWriter wraps http.ResponseWriter to capture response
type responseWriter struct {
    http.ResponseWriter
    body       *bytes.Buffer
    statusCode int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    return rw.body.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
    rw.statusCode = statusCode
}

// ServeHTTP is the main handler - it routes to appropriate sub-handlers
func (a *AnnotationProxy) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
    // Route 1: Serve embedded assets
    if strings.HasPrefix(r.URL.Path, a.ScriptPath) {
        return a.serveAssets(w, r)
    }
    
    // Route 2: Handle API endpoints
    if strings.HasPrefix(r.URL.Path, a.APIEndpoint) {
        return a.serveAPI(w, r)
    }
    
    // Route 3: Intercept and inject for proxied content
    if !a.Enabled {
        return next.ServeHTTP(w, r)
    }
    
    return a.interceptAndInject(w, r, next)
}

// interceptAndInject captures the response and injects annotation script
func (a *AnnotationProxy) interceptAndInject(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
    // Create buffer to capture response
    buf := &bytes.Buffer{}
    wrapped := &responseWriter{
        ResponseWriter: w,
        body:          buf,
        statusCode:    http.StatusOK,
    }
    
    // Call next handler in chain
    err := next.ServeHTTP(wrapped, r)
    if err != nil {
        return err
    }
    
    // Check if response is HTML
    contentType := wrapped.Header().Get("Content-Type")
    isHTML := strings.Contains(strings.ToLower(contentType), "text/html")
    
    if isHTML && wrapped.statusCode == http.StatusOK {
        // Inject annotation script
        modified := a.injectScript(buf.String())
        
        // Write modified response
        w.WriteHeader(wrapped.statusCode)
        _, err = w.Write([]byte(modified))
        return err
    }
    
    // Write original response
    w.WriteHeader(wrapped.statusCode)
    _, err = w.Write(buf.Bytes())
    return err
}

// injectScript adds annotation JavaScript and CSS to HTML
func (a *AnnotationProxy) injectScript(html string) string {
    // Inject CSS before </head>
    cssTag := `<link rel="stylesheet" href="` + a.ScriptPath + `/annotation.css">`
    html = strings.Replace(html, "</head>", cssTag+"\n</head>", 1)
    
    // Inject JS before </body>
    scriptTag := `<script>window.ANNOTATION_API = "` + a.APIEndpoint + `";</script>
<script src="` + a.ScriptPath + `/annotation.js"></script>`
    html = strings.Replace(html, "</body>", scriptTag+"\n</body>", 1)
    
    return html
}

// serveAssets serves embedded JavaScript and CSS
func (a *AnnotationProxy) serveAssets(w http.ResponseWriter, r *http.Request) error {
    // Extract filename
    filename := strings.TrimPrefix(r.URL.Path, a.ScriptPath+"/")
    
    var contentType string
    switch filename {
    case "annotation.js":
        contentType = "application/javascript"
    case "annotation.css":
        contentType = "text/css"
    default:
        http.NotFound(w, r)
        return nil
    }
    
    // Read embedded file
    content, err := embedFS.ReadFile("assets/" + filename)
    if err != nil {
        http.NotFound(w, r)
        return nil
    }
    
    w.Header().Set("Content-Type", contentType)
    w.Header().Set("Cache-Control", "no-cache")
    w.Write(content)
    return nil
}
