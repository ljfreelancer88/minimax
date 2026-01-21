package annotation

import (
    "embed"
    
    "github.com/caddyserver/caddy/v2"
    "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
    "github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
    "github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

//go:embed assets/*
var embedFS embed.FS

func init() {
    caddy.RegisterModule(AnnotationProxy{})
    httpcaddyfile.RegisterHandlerDirective("annotation", parseCaddyfile)
    httpcaddyfile.RegisterDirectiveOrder("annotation", "before", "reverse_proxy")
}

// AnnotationProxy is the main module
type AnnotationProxy struct {
    Enabled     bool   `json:"enabled,omitempty"`
    APIEndpoint string `json:"api_endpoint,omitempty"`
    ScriptPath  string `json:"script_path,omitempty"`
}

// CaddyModule returns module information
func (AnnotationProxy) CaddyModule() caddy.ModuleInfo {
    return caddy.ModuleInfo{
        ID:  "http.handlers.annotation",
        New: func() caddy.Module { return new(AnnotationProxy) },
    }
}

// Provision sets up the module
func (a *AnnotationProxy) Provision(ctx caddy.Context) error {
    if a.ScriptPath == "" {
        a.ScriptPath = "/annotation-assets"
    }
    if a.APIEndpoint == "" {
        a.APIEndpoint = "/api/annotations"
    }
    return nil
}

// Validate ensures configuration is valid
func (a *AnnotationProxy) Validate() error {
    return nil
}

// ServeHTTP implements the handler - routes requests appropriately
// Implementation is in handler.go

// UnmarshalCaddyfile implements caddyfile.Unmarshaler
func (a *AnnotationProxy) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
    for d.Next() {
        for d.NextBlock(0) {
            switch d.Val() {
            case "enabled":
                if !d.NextArg() {
                    return d.ArgErr()
                }
                a.Enabled = d.Val() == "true"
            case "api_endpoint":
                if !d.NextArg() {
                    return d.ArgErr()
                }
                a.APIEndpoint = d.Val()
            case "script_path":
                if !d.NextArg() {
                    return d.ArgErr()
                }
                a.ScriptPath = d.Val()
            }
        }
    }
    return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
    var a AnnotationProxy
    err := a.UnmarshalCaddyfile(h.Dispenser)
    return &a, err
}

// Interface guards
var (
    _ caddy.Provisioner           = (*AnnotationProxy)(nil)
    _ caddy.Validator             = (*AnnotationProxy)(nil)
    _ caddyhttp.MiddlewareHandler = (*AnnotationProxy)(nil)
    _ caddyfile.Unmarshaler       = (*AnnotationProxy)(nil)
)
