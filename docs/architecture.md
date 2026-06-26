# Arquitetura (MVP Inicial)

## Backend (Go + Gin)

Arquitetura orientada a camadas com princípios de Clean Architecture:

- `controllers`: entrada HTTP/REST
- `services`: regras de negócio
- `repositories`: acesso a dados SQLite
- `ldap`: integração com AD (real e modo demo)
- `middleware`: autenticação, RBAC, CSRF, rate limit e headers de segurança
- `audit`: trilha de auditoria imutável

### Fluxo de autenticação

1. Usuário envia login e senha
2. Backend valida credenciais via LDAP
3. Backend verifica perfil interno (`app_users`)
4. Carrega permissões e escopos de OU
5. Emite Access Token + Refresh Token

### RBAC + Restrição por OU

- RBAC em nível de endpoint (middleware `RequirePermission`)
- Escopo de OU carregado no token (`ouScopes`)
- Services LDAP recebem automaticamente os escopos permitidos

## Frontend (React + Vite)

- Shell administrativo responsivo com sidebar recolhível
- Tema claro/escuro
- React Query para dados assíncronos
- Zustand para estado de autenticação e preferências locais
- Rotas protegidas por token JWT

## Evolução planejada

- Implementar endpoints LDAP reais de leitura/escrita para objetos AD
- Aumentar cobertura de testes unitários/integrados
- OpenAPI gerado automaticamente
- Pipeline de CI/CD com scan de vulnerabilidades
- Plugin system e eventos em tempo real ricos (WebSocket)
