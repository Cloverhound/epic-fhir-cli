# Epic FHIR CLI

A command-line tool for querying Epic's FHIR R4 APIs using OAuth2 backend authentication.

> **Notice:** This is an open source library provided as-is with limited support. We are not responsible for how you use this tool. As it interacts with Electronic Health Record (EHR) systems, it is your responsibility to ensure that your use complies with all applicable regulations, policies, and standards governing access to and handling of health data.

## Install

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/install.ps1 | iex
```

Or download from [Releases](https://github.com/Cloverhound/epic-fhir-cli/releases).

## Update

```bash
fhir-cli update
```

## Quick Start

```bash
# Interactive setup (creates ~/.fhir-cli/config.yaml)
fhir-cli config init

# Verify authentication
fhir-cli auth token

# Search for patients
fhir-cli patient search --family "Smith" -o table

# Get a patient by ID
fhir-cli patient get abc123

# Get observations for a patient
fhir-cli observation search --patient abc123 --category vital-signs

# Active medications
fhir-cli medication active abc123 -o table

# Different output formats
fhir-cli patient search --family "Smith" -o json
fhir-cli patient search --family "Smith" -o yaml
fhir-cli patient search --family "Smith" -o table
```

## FHIR Resources

30+ FHIR R4 resource types supported:

| Command | Description |
|---------|-------------|
| `patient` | Patient demographics, search, $everything |
| `observation` | Observations, vitals, lab results |
| `condition` | Conditions, active problems, problem list |
| `medication` | MedicationRequest, Statement, Administration, Dispense |
| `encounter` | Encounters (inpatient, outpatient) |
| `procedure` | Procedures |
| `allergy` | AllergyIntolerance |
| `diagnostic` | DiagnosticReport (labs, imaging) |
| `document` | DocumentReference, Binary |
| `immunization` | Immunizations |
| `appointment` | Appointments, Schedules, Slots |
| `careplan` | CarePlan, CareTeam, Goal |
| `coverage` | Coverage, ExplanationOfBenefit, ServiceRequest |
| `practitioner` | Practitioner, PractitionerRole, Organization, Location |
| `familyhistory` | FamilyMemberHistory, RelatedPerson, Consent, Provenance, QuestionnaireResponse |
| `resource` | Generic FHIR resource operations, metadata |

## Authentication

- **OAuth2 backend auth** — RS384 JWT assertions signed with your private key
- **Token caching** — tokens stored securely in OS keyring (macOS Keychain / Linux Secret Service / Windows Credential Manager) and reused until 5 minutes before expiry
- **Multi-profile** — configure multiple environments (sandbox, staging, production) and switch between them
- **Private key sources** — file path (`/path/to/key.pem`) or environment variable (`env:FHIR_PRIVATE_KEY`)

```bash
fhir-cli config init                # Interactive setup
fhir-cli config show                # Show current config
fhir-cli config add-profile prod    # Add a new profile
fhir-cli config use prod            # Switch default profile
fhir-cli auth token                 # Get/refresh access token
fhir-cli auth status                # Check authentication status
fhir-cli auth logout                # Clear cached tokens
fhir-cli auth debug                 # Debug auth configuration
```

## Output Formats

Control output with `-o`:

| Format | Description |
|--------|-------------|
| `json` | Pretty-printed JSON (default) |
| `table` | ASCII table |
| `yaml` | YAML |

## Global Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output format: json, table, yaml (default: json) |
| `-p, --profile` | Configuration profile to use (default: default) |
| `-v, --verbose` | Verbose output (show HTTP requests/responses) |
| `--config` | Custom config file path |

## Configuration

Configuration is stored in `~/.fhir-cli/config.yaml`:

```yaml
default: sandbox
profiles:
  sandbox:
    name: sandbox
    client_id: your-client-id
    private_key: env:FHIR_PRIVATE_KEY
    token_url: https://fhir.epic.com/interconnect-fhir-oauth/oauth2/token
    fhir_base_url: https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/R4
    fhir_version: R4
    scopes: system/Patient.read system/Observation.read
    output_format: json
```

The `private_key` field supports file paths (`/path/to/key.pem`, `~/key.pem`) or environment variables (`env:FHIR_PRIVATE_KEY`).

## Coding Agent Skill

The Epic FHIR CLI includes a [skill](https://agentskills.io) that enables AI coding agents to query and manage your Epic FHIR environment. It works with any agent that supports the skills standard — [Claude Code](https://claude.com/claude-code), [OpenAI Codex](https://openai.com/codex/), [Cursor](https://cursor.com), and others.

### Setup

1. Install the Epic FHIR CLI (see [Install](#install))
2. Configure: `fhir-cli config init`
3. Download the skill into the correct folder for your coding agent:

#### Claude Code

**macOS / Linux:**
```bash
mkdir -p ~/.claude/skills/epic-fhir-cli
curl -fsSL https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/skill/SKILL.md \
  -o ~/.claude/skills/epic-fhir-cli/SKILL.md
```

**Windows (PowerShell):**
```powershell
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.claude\skills\epic-fhir-cli"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/skill/SKILL.md" `
  -OutFile "$env:USERPROFILE\.claude\skills\epic-fhir-cli\SKILL.md"
```

#### OpenAI Codex

**macOS / Linux:**
```bash
mkdir -p ~/.agents/skills/epic-fhir-cli
curl -fsSL https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/skill/SKILL.md \
  -o ~/.agents/skills/epic-fhir-cli/SKILL.md
```

**Windows (PowerShell):**
```powershell
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.agents\skills\epic-fhir-cli"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/skill/SKILL.md" `
  -OutFile "$env:USERPROFILE\.agents\skills\epic-fhir-cli\SKILL.md"
```

#### Cursor

**macOS / Linux:**
```bash
mkdir -p ~/.cursor/skills/epic-fhir-cli
curl -fsSL https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/skill/SKILL.md \
  -o ~/.cursor/skills/epic-fhir-cli/SKILL.md
```

**Windows (PowerShell):**
```powershell
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.cursor\skills\epic-fhir-cli"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/skill/SKILL.md" `
  -OutFile "$env:USERPROFILE\.cursor\skills\epic-fhir-cli\SKILL.md"
```

> These commands install the skill globally (user-level). You can also install per-project by placing the `epic-fhir-cli/SKILL.md` folder inside your project's `.claude/skills/`, `.agents/skills/`, or `.cursor/skills/` directory instead.

4. If the `fhir-cli` binary is not in your `$PATH`, ask your coding agent to update the binary path in the skill file.

### Example Prompts

```
/epic-fhir-cli search for patient Smith and show their active conditions

/epic-fhir-cli get all vital signs for patient abc123 in table format

/epic-fhir-cli check auth status and list all available resource types
```

## Development

```bash
make build      # Build binary
make test       # Run tests
make lint       # Run linters
make fmt        # Format code
make build-all  # Cross-compile for all platforms
```

## License

MIT
