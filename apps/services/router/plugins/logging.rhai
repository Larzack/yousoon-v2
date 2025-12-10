// Logging Plugin for Apollo Router
// Structured logging for observability

fn supergraph_service(service) {
    let start_time = timestamp();
    
    service.map_request(|request| {
        let request_id = request.headers["X-Request-ID"] ?? uuid();
        let user_id = request.context["user_id"] ?? "anonymous";
        let operation_name = request.body.operation_name ?? "unknown";
        
        // Store in context for response logging
        request.context["request_id"] = request_id;
        request.context["request_start"] = start_time;
        
        // Log request
        log_info(#{
            event: "graphql_request",
            request_id: request_id,
            user_id: user_id,
            operation: operation_name,
            timestamp: start_time,
        });
    });
    
    service.map_response(|response| {
        let end_time = timestamp();
        let start_time = response.context["request_start"] ?? end_time;
        let duration_ms = end_time - start_time;
        let request_id = response.context["request_id"] ?? "unknown";
        
        // Check for errors
        let has_errors = response.body.errors != () && response.body.errors.len() > 0;
        let error_count = if has_errors { response.body.errors.len() } else { 0 };
        
        // Log response
        log_info(#{
            event: "graphql_response",
            request_id: request_id,
            duration_ms: duration_ms,
            has_errors: has_errors,
            error_count: error_count,
            timestamp: end_time,
        });
        
        // Add timing header
        response.headers["X-Response-Time"] = `${duration_ms}ms`;
        response.headers["X-Request-ID"] = request_id;
    })
}

fn log_info(data) {
    // In production, this would send to Loki/CloudWatch
    print(`[INFO] ${json_stringify(data)}`);
}

fn log_error(data) {
    print(`[ERROR] ${json_stringify(data)}`);
}

fn log_warn(data) {
    print(`[WARN] ${json_stringify(data)}`);
}
