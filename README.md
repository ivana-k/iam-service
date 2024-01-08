# iam-service

Identity management service providing:
- user registration
- user organization management
- user sign in with Vault
- token verification
- granting permissions for authorization purposes

Integrated with:
- Vault
- Oort policy engine

## POST /user/register route

### Description

Provides registration to user. This endpoint also creates new organization, new user on Vault and creates org - user relationship with default org permissions on Oort service. Each user is owner of his own organization and receives default permissions upon registration.

|parameter| type  |                    description                      |
|---------|-------|-----------------------------------------------------|
| email    | string  | **Required.**  |
| username    | string  | **Required.** Used later for login |
| name    | string  | **Required.**  |
| surname    | string  | **Required.**  |
| password    | string  | **Required.** Stored securely on Vault server. |
| org    | string  | **Required.** Name of the organization. Should be unique. If not provided, it will be created as username_default_org |
