# User Setup

Navigate to /users/setup after signing in as an authorized school administrator. Search accounts, use quick status filters, manage advanced conditions and sorting, and configure table columns.

Create user: first name, optional middle name, last name, email, optional phone, required primary role and optional secondary roles. The server creates the account with its role assignment and queues the invitation. The form does not include password or status.

Edit user: save personal names separately from account and roles. Phone and email are read-only during edit to preserve verified contacts. The primary role cannot be duplicated in secondary roles. Role IDs are server-validated. School job roles do not themselves grant administrator access.

Lifecycle actions are separate from editing. Only permitted transitions and resend operations are displayed; backend tenant-admin authorization remains authoritative.
