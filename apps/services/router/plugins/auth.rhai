// Authentication Plugin for Apollo Router
// Validates JWT tokens and extracts user claims

// JWT secret (would be loaded from environment in production)
let JWT_SECRET = get_env("JWT_SECRET") ?? "your-secret-key";

// Public operations that don't require authentication
let PUBLIC_OPERATIONS = [
    "IntrospectionQuery",
    "__schema",
    "Login",
    "Register",
    "ForgotPassword",
    "ResetPassword",
    "RefreshToken",
    "GetPublicOffers",
    "GetOffer",
    "GetCategories",
    "GetPartner",
    "GetEstablishment",
    "SearchOffers",
];

// Admin-only operations
let ADMIN_OPERATIONS = [
    "ApprovePartner",
    "RejectPartner",
    "ApproveReview",
    "RejectReview",
    "DeleteUser",
    "UpdateSubscriptionPlan",
    "GetPendingReviews",
    "GetReportedReviews",
    "GetAdminLogs",
];

// Partner-only operations
let PARTNER_OPERATIONS = [
    "CreateOffer",
    "UpdateOffer",
    "DeleteOffer",
    "CreateEstablishment",
    "UpdateEstablishment",
    "InviteTeamMember",
    "PartnerOutings",
    "EstablishmentOutings",
    "ManualCheckIn",
];

fn supergraph_service(service) {
    service.map_request(|request| {
        let operation_name = request.body.operation_name ?? "";
        
        // Skip auth for public operations
        if PUBLIC_OPERATIONS.contains(operation_name) {
            return;
        }
        
        // Get Authorization header
        let auth_header = request.headers["Authorization"] ?? "";
        
        if auth_header == "" {
            throw #{
                status: 401,
                message: "Authentication required",
            };
        }
        
        // Extract Bearer token
        if !auth_header.starts_with("Bearer ") {
            throw #{
                status: 401,
                message: "Invalid authorization header format",
            };
        }
        
        let token = auth_header.sub_string(7);
        
        // Validate JWT (simplified - real implementation would verify signature)
        let claims = decode_jwt(token);
        
        if claims == () {
            throw #{
                status: 401,
                message: "Invalid or expired token",
            };
        }
        
        // Check token expiration
        let now = timestamp();
        if claims.exp < now {
            throw #{
                status: 401,
                message: "Token expired",
            };
        }
        
        // Set context for downstream services
        request.context["user_id"] = claims.sub;
        request.context["user_email"] = claims.email ?? "";
        request.context["user_roles"] = claims.roles ?? [];
        
        // Check admin operations
        if ADMIN_OPERATIONS.contains(operation_name) {
            if !claims.roles.contains("admin") {
                throw #{
                    status: 403,
                    message: "Admin access required",
                };
            }
        }
        
        // Check partner operations
        if PARTNER_OPERATIONS.contains(operation_name) {
            if !claims.roles.contains("partner") && !claims.roles.contains("admin") {
                throw #{
                    status: 403,
                    message: "Partner access required",
                };
            }
        }
        
        // Forward user info to subgraphs
        request.subgraph.headers["X-User-ID"] = claims.sub;
        request.subgraph.headers["X-User-Email"] = claims.email ?? "";
        request.subgraph.headers["X-User-Roles"] = claims.roles.join(",");
    })
}

// Decode JWT (simplified - production should use proper JWT library)
fn decode_jwt(token) {
    let parts = token.split(".");
    
    if parts.len() != 3 {
        return ();
    }
    
    // Decode payload (base64)
    let payload = base64_decode(parts[1]);
    
    if payload == () {
        return ();
    }
    
    // Parse JSON
    let claims = parse_json(payload);
    
    claims
}
