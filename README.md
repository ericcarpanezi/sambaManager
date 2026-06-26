# AG Directory Manager

Sistema Web para gerenciamento de Samba Active Directory (Samba 4 Domain Controller), com arquitetura modular, segurança elevada e interface moderna inspirada no ADUC.

## Stack

### Backend
- Go 1.24+
- Gin
- LDAP (go-ldap)
- SQLite
- JWT + Refresh Token
- WebSocket
- Clean Architecture + Repository Pattern

### Frontend
- React + TypeScript + Vite
- TailwindCSS
- React Query
- Zustand
- React Router
- Framer Motion
- Lucide Icons

## Estrutura

```text
cmd/
internal/
frontend/
docs/
scripts/
docker/
migrations/
tests/
```

## Quick Start

### Backend

```bash
export APP_JWT_SECRET='troque-este-segredo'
go run ./cmd/server
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Variáveis de ambiente

Copie `.env.example` para `.env` e ajuste as configurações de LDAP/JWT.

## Estado atual

Este repositório entrega um **MVP empresarial inicial** com:
- API REST versionada (`/api/v1`)
- Login LDAP + perfil interno da aplicação
- RBAC granular
- Restrição por OU
- Auditoria imutável
- Dashboard inicial
- Modo demonstração
- Setup de Docker Compose
- Scripts de instalação base para Linux

Os módulos estão preparados para evolução incremental até a cobertura completa do escopo Enterprise descrito no documento funcional.
