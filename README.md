# Lilys Hidden Paradise

Lilys Hidden Paradise is a landlord-tenant CRM tool that allows landlords to manage their properties and tenants.

**Future iterations will also allows tenants to manage their leases and payments.**

*This tool is built usng Golang, JavaScript and PostgreSQL.*

## Key Features

- HTML templating with `http/template` standard library
- Web authentication using sessions & CSRF protection
- Protected dashboards for landlords and tenants (using middleware)
- Creating a database connection using `database/sql` standard library and `github.com/lib/pq` for the PostgreSQL driver
- Password hashing using `golang.org/x/crypto/bcrypt`
- Encrypting & decrypting database information 
- Message platform for landlords and tenants
- Database stubbing for testing