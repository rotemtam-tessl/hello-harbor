---
name: atlas-full
description: Comprehensive Atlas database schema management. Use when working with Atlas migrations, schema inspection, diffing, linting, or applying database changes.
---

# Atlas Database Schema Management

Atlas is a language-independent tool for managing and migrating database schemas using modern DevOps principles.

## Quick Reference

```bash
# Common Atlas commands
atlas schema inspect --env <name> --url file://migrations
atlas migrate status --env <name>
atlas migrate diff --env <name>
atlas migrate lint --env <name> --latest 1
atlas migrate apply --env <name>
atlas whoami
```

## Core Concepts and Configurations

### Configuration File Structure
Atlas uses `atlas.hcl` configuration files with the following structure:

```hcl
// Basic environment configuration
env "<name>" {
  url = getenv("DATABASE_URL")
  dev = "docker://postgres/15/dev?search_path=public"
  
  migration {
    dir = "file://migrations"
  }
  
  schema {
    src = "file://schema.hcl"
  }
}
```

### Dev database
Atlas utilizes a "dev-database", which is a temporary and locally running database, usually bootstrapped by Atlas. Atlas will use the dev database to process and validate users' schemas, migrations, and more. 

Examples of dev-database configurations:

```
# When working on a single database schema
--dev-url "docker://mysql/8/dev"
--dev-url "docker://postgres/15/db_name?search_path=public"
--dev-url "sqlite://dev?mode=memory"
# When working on multiple database schemas.
--dev-url "docker://mysql/8"
--dev-url "docker://postgres/15/dev"
```

Configure the dev database using HCL:

```hcl
env "<name>" {
  dev = "docker://mysql/8"
}
```

### Environment Variables and Security

**✅ DO**: Use secure configuration patterns
```hcl
// Using environment variables (recommended)
env "<name>" {
  url = getenv("DATABASE_URL")
}
```

**❌ DON'T**: Hardcode sensitive values
```hcl
// Never do this
env "prod" {
  url = "postgres://user:password123@prod-host:5432/database"
}
```

### Schema Sources

#### External Schema (ORM Integration)
The external_schema data source enables the import of an SQL schema from an external program into Atlas' desired state.

```hcl
data "external_schema" "gorm" {
  program = [
    "go", "run", "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load", "--path", "./models", "--dialect", "sqlite"
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "sqlite://dev?mode=memory"
  migration {
    dir = "file://migrations"
  }
}
```

## Common Workflows

### 1. Schema Inspection / Visualization

1. Always start by listing tables, don't immediately try to inspect the entire schema.
2. If you see there are many tables, don't inspect the entire schema at once.

```bash
# Get table list overview
atlas schema inspect --env <name> --url file://migrations --format "{{ json . }}" | jq ".schemas[].tables[].name"

# Get full SQL schema
atlas schema inspect --env <name> --url file://migrations --format "{{ sql . }}"
```

### 2. Migration Status

```bash
# Check current migration directory status
atlas migrate status --env <name>
```

### 3. Migration Generation / Diffing

```bash
# Compare current migrations with desired schema and create a new migration file
atlas migrate diff --env <name> "add_user_table"
```

### 4. Migration Linting

```bash
# Lint last migration
atlas migrate lint --env <name> --latest 1
```

### 5. Applying Migration

```bash 
# Dry run first (always do this before applying)
atlas migrate apply --env <name> --dry-run

# Apply to configured environment
atlas migrate apply --env <name>
```

### 6. Making Changes to the Schema

**⚠️ CRITICAL: ALL schema changes MUST follow this exact workflow. NO EXCEPTIONS.**

1. Start by inspecting the schema, understand the current state
2. After making changes to the schema/models, run `atlas migrate diff --env <name> <migration_name>` to generate a migration file
3. Run `atlas migrate lint --env <name> --latest 1` to validate the migration file
4. If there are issues, fix them and run `atlas migrate hash --env <name>` to recalculate the hash
5. Run `atlas migrate validate --env <name>` to ensure migrations are valid

> **Important for applied migrations:** Never edit directly. Create a new corrective migration instead.

## Troubleshooting Commands

```bash
# Check Atlas installation
atlas version

# Repair migration integrity after manual edits
atlas migrate hash --env <name>

# Validate migrations
atlas migrate validate --env <name>
```

## Key Reminders

1. **Always read `atlas.hcl` first** before running any Atlas commands
2. **Use environment names** from the config file
3. **Never hardcode database URLs**
4. **Run `atlas migrate hash`** after manually editing migration files
5. **Use `atlas migrate lint`** to validate migrations before applying
6. **Never modify applied migrations** - create new ones instead
7. **Always use `--dry-run`** before applying migrations
8. **Prefer `atlas schema inspect`** over reading migration files directly
