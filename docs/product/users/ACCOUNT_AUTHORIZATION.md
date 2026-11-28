# Temporary account authorization before RBAC

Protected account search, lookup, registration, update, lifecycle transitions, invitation resend and role catalog require a verified session and nonempty tenant ID, not the legacy administrator flag. Existing tenant predicates and ownership checks remain unchanged. Profile operations inherit the account lookup boundary. Gateway origin, CSRF, authentication, rate limiting and restricted workload policies are unchanged.

This interim policy gives every authenticated tenant member access to account management within that tenant. Do not expose it to real school users or production before permission-based RBAC is implemented. Role assignment does not itself grant gateway administrator authority.
