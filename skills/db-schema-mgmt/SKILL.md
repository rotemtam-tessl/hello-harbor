---
name: db-schema-mgmt
description: Rules that MUST be followed when working on database schema related tasks with Atlas and GORM.
---

# Atlas Schema Management (Go/GORM)

This project uses Atlas for database schema management with GORM models.

## Quick Start

When the schema changes, generate a new migration:

```bash
atlas migrate diff --env gorm <migration_name>
```

## Key Commands

| Task                | Command                                |
| ------------------- | -------------------------------------- |
| Generate migration  | `atlas migrate diff --env gorm <name>` |
| Validate migrations | `atlas migrate validate --env gorm`    |
| Check status        | `atlas migrate status --env gorm`      |
| Fix checksums       | `atlas migrate hash --env gorm`        |

## Workflow for Schema Changes

1. **Modify GORM models** in `models/` directory
2. **Generate migration**: `atlas migrate diff --env gorm add_description_field`
3. **Review** the generated SQL in `migrations/`
4. **Validate**: `atlas migrate validate --env gorm`
5. **Run tests**: `go test ./...`

## Critical Rules

⚠️ **NEVER modify existing migration files** - this breaks linear history
⚠️ **ALWAYS create NEW migrations** for schema changes
⚠️ **Run `atlas migrate hash`** if you must edit an unapplied migration

## Common Issues

### Migration Validation Failed

```bash
atlas migrate validate --env gorm
# Fix issues, then:
atlas migrate hash --env gorm
```

## Project Structure

```
├── atlas.hcl          # Atlas configuration
├── models/            # GORM models (source of truth)
│   └── todos.go
├── migrations/        # Atlas versioned migrations
│   ├── atlas.sum      # Integrity checksums (don't edit)
│   └── *.sql          # Migration files
```

The `atlas.hcl` uses `atlas-provider-gorm` to read GORM model definitions and compare them against the migrations directory.
