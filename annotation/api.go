package annotation

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"
)

// Annotation represents a feedback/comment
type Annotation struct {
    ID        string `json:"id"`
    URL       string `json:"url"`
    Selector  string `json:"selector"`
    Comment   string `json:"comment"`
    Author    string `json:"author"`
    Timestamp int64  `json:"timestamp"`
    XPath     string `json:"xpath,omitempty"`
    Position  struct {
        X int `json:"x"`
        Y int `json:"y"`
    } `json:"position"`
}

// In-memory storage (replace with database in production)
var (
    annotations = make(map[string][]Annotation)
    mu          sync.RWMutex
    idCounter   int64
)

// serveAPI processes annotation API requests
func (a *AnnotationProxy) serveAPI(w http.ResponseWriter, r *http.Request) error {
    // Set CORS headers
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return nil
    }
    
    // Route to appropriate handler
    switch r.Method {
    case "GET":
        return a.getAnnotations(w, r)
    case "POST":
        return a.createAnnotation(w, r)
    case "DELETE":
        return a.deleteAnnotation(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return nil
    }
}

func (a *AnnotationProxy) getAnnotations(w http.ResponseWriter, r *http.Request) error {
    url := r.URL.Query().Get("url")
    if url == "" {
        http.Error(w, `{"error":"URL parameter required"}`, http.StatusBadRequest)
        return nil
    }
    
    mu.RLock()
    items := annotations[url]
    mu.RUnlock()
    
    if items == nil {
        items = []Annotation{}
    }
    
    response := map[string]interface{}{
        "annotations": items,
        "count":       len(items),
    }
    
    return json.NewEncoder(w).Encode(response)
}

func (a *AnnotationProxy) createAnnotation(w http.ResponseWriter, r *http.Request) error {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, `{"error":"Failed to read body"}`, http.StatusBadRequest)
        return nil
    }
    defer r.Body.Close()
    
    var annotation Annotation
    if err := json.Unmarshal(body, &annotation); err != nil {
        http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
        return nil
    }
    
    // Validate required fields
    if annotation.URL == "" || annotation.Comment == "" {
        http.Error(w, `{"error":"URL and comment are required"}`, http.StatusBadRequest)
        return nil
    }
    
    // Set metadata
    annotation.ID = generateID()
    annotation.Timestamp = time.Now().Unix()
    
    // Store annotation
    mu.Lock()
    annotations[annotation.URL] = append(annotations[annotation.URL], annotation)
    mu.Unlock()
    
    w.WriteHeader(http.StatusCreated)
    return json.NewEncoder(w).Encode(annotation)
}

func (a *AnnotationProxy) deleteAnnotation(w http.ResponseWriter, r *http.Request) error {
    id := r.URL.Query().Get("id")
    url := r.URL.Query().Get("url")
    
    if id == "" || url == "" {
        http.Error(w, `{"error":"ID and URL parameters required"}`, http.StatusBadRequest)
        return nil
    }
    
    mu.Lock()
    defer mu.Unlock()
    
    items := annotations[url]
    for i, ann := range items {
        if ann.ID == id {
            annotations[url] = append(items[:i], items[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return nil
        }
    }
    
    http.Error(w, `{"error":"Annotation not found"}`, http.StatusNotFound)
    return nil
}

// Helper functions
func generateID() string {
    mu.Lock()
    idCounter++
    id := idCounter
    mu.Unlock()
    return fmt.Sprintf("ann_%d_%d", time.Now().Unix(), id)
}
