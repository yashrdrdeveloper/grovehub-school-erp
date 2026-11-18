# Shared business layer
Business/API contracts are transport-only TypeScript and never render UI. Every backend call uses `apiRequest` and the shared `{success,data,meta}` envelope. Opaque authentication cookies stay out of local storage and JavaScript. Domain pages consume business services, not fetch directly.
