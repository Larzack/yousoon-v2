// Yousoon Apollo Router - Rhai Plugins Entry Point
// This file loads and chains all custom plugins

// Load plugin modules
import "auth" as auth;
import "rate_limit" as rate_limit;
import "logging" as logging;

// Supergraph request hook - runs before query planning
fn supergraph_service(service) {
    // Apply logging
    let logging_service = logging::supergraph_service(service);
    
    // Apply auth
    let auth_service = auth::supergraph_service(logging_service);
    
    // Apply rate limiting
    let rate_limited = rate_limit::supergraph_service(auth_service);
    
    rate_limited
}

// Subgraph request hook - runs for each subgraph call
fn subgraph_service(service, subgraph) {
    // Add subgraph-specific headers
    service.map_request(|request| {
        request.subgraph.headers["X-Subgraph-Name"] = subgraph;
        request.subgraph.headers["X-Request-Start"] = `${timestamp()}`;
    });
}

// Execution hook - runs during query execution
fn execution_service(service) {
    service
}

// Router service hook - first entry point
fn router_service(service) {
    // Add request ID if not present
    service.map_request(|request| {
        if request.headers["X-Request-ID"] == () {
            request.headers["X-Request-ID"] = `${uuid()}`;
        }
    });
}
