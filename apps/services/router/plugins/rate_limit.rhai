// Rate Limiting Plugin for Apollo Router
// Implements per-user rate limiting with Redis backend

// Rate limit configuration
let RATE_LIMITS = #{
    // Per-user limits (requests per minute)
    default: 100,
    authenticated: 200,
    partner: 500,
    admin: 1000,
    
    // Operation-specific limits
    operations: #{
        "BookOffer": 10,        // 10 bookings per minute max
        "SubmitReview": 5,      // 5 reviews per minute max
        "AddToFavorites": 30,   // 30 favorites per minute max
        "SearchOffers": 60,     // 60 searches per minute max
    },
};

// Sliding window duration in seconds
let WINDOW_SECONDS = 60;

fn supergraph_service(service) {
    service.map_request(|request| {
        let user_id = request.context["user_id"] ?? "anonymous";
        let user_roles = request.context["user_roles"] ?? [];
        let operation_name = request.body.operation_name ?? "unknown";
        
        // Determine rate limit based on user role
        let limit = RATE_LIMITS.default;
        
        if user_roles.contains("admin") {
            limit = RATE_LIMITS.admin;
        } else if user_roles.contains("partner") {
            limit = RATE_LIMITS.partner;
        } else if user_id != "anonymous" {
            limit = RATE_LIMITS.authenticated;
        }
        
        // Check operation-specific limit
        if RATE_LIMITS.operations[operation_name] != () {
            let op_limit = RATE_LIMITS.operations[operation_name];
            if op_limit < limit {
                limit = op_limit;
            }
        }
        
        // Create rate limit key
        let key = `rate_limit:${user_id}:${operation_name}`;
        
        // Check rate limit (would use Redis in production)
        let current = get_rate_limit(key);
        
        if current >= limit {
            // Add rate limit headers
            throw #{
                status: 429,
                message: "Rate limit exceeded. Please try again later.",
                headers: #{
                    "X-RateLimit-Limit": `${limit}`,
                    "X-RateLimit-Remaining": "0",
                    "X-RateLimit-Reset": `${timestamp() + WINDOW_SECONDS}`,
                    "Retry-After": `${WINDOW_SECONDS}`,
                },
            };
        }
        
        // Increment counter
        increment_rate_limit(key);
        
        // Add rate limit headers to response
        request.context["rate_limit_limit"] = limit;
        request.context["rate_limit_remaining"] = limit - current - 1;
    });
    
    service.map_response(|response| {
        let limit = response.context["rate_limit_limit"] ?? 100;
        let remaining = response.context["rate_limit_remaining"] ?? 99;
        
        response.headers["X-RateLimit-Limit"] = `${limit}`;
        response.headers["X-RateLimit-Remaining"] = `${remaining}`;
    })
}

// In-memory rate limit (would be Redis in production)
let rate_limits = #{};

fn get_rate_limit(key) {
    let now = timestamp();
    let window_start = now - WINDOW_SECONDS;
    
    if rate_limits[key] == () {
        return 0;
    }
    
    let entry = rate_limits[key];
    
    // Check if window expired
    if entry.window_start < window_start {
        rate_limits.remove(key);
        return 0;
    }
    
    entry.count
}

fn increment_rate_limit(key) {
    let now = timestamp();
    let window_start = now - WINDOW_SECONDS;
    
    if rate_limits[key] == () {
        rate_limits[key] = #{
            window_start: now,
            count: 1,
        };
        return;
    }
    
    let entry = rate_limits[key];
    
    // Reset if window expired
    if entry.window_start < window_start {
        rate_limits[key] = #{
            window_start: now,
            count: 1,
        };
        return;
    }
    
    entry.count += 1;
}
