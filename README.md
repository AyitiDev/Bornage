# Open Land Registry (`open-land-registry` / Bornage)

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![ISO Standard](https://img.shields.io/badge/Standard-ISO_19152_LADM-green.svg)](https://www.iso.org/standard/51206.html)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8.svg)](https://go.dev/)
[![PostGIS](https://img.shields.io/badge/PostGIS-3.4+-336791.svg)](https://postgis.net/)

**`open-land-registry`** is an open source, **Fit-for-Purpose (FFP)** cadastral and land registration engine

---

## Technical Architecture

```mermaid
flowchart TD
    subgraph Clients["Frontend Clients (Offline-First PWAs)"]
        FieldApp["Field Surveyor PWA<br/>IndexedDB + PMTiles + GPS"]
        CitizenApp["Citizen Portal PWA<br/>Location Lookup and Registration Request"]
    end

    subgraph API["Backend API (Go + Fiber)"]
        FieldAPI["/api/v1/field/*<br/>(Offline Batch Sync)"]
        InternalAPI["/api/v1/internal/*<br/>(Review, Objection, Registration)"]
        PublicAPI["/api/v1/public/*<br/>(Coordinates / Geo Only)"]
        InteropAPI["/api/v1/interop/*<br/>"]
    end

    subgraph Data["Data Infrastructure"]
        DB[(PostgreSQL + PostGIS)]
        S3[(MinIO Storage - S3 Evidence)]
    end

    FieldApp -->|Idempotent Sync| FieldAPI
    CitizenApp -->|Geo Query| PublicAPI
    CitizenApp -->|Submit Claim| FieldAPI
    
    FieldAPI --> DB
    InternalAPI --> DB
    PublicAPI --> DB
    InteropAPI --> DB

    FieldAPI --> S3
    InternalAPI --> S3
```

---

## 🔄 Claim State Machine

```mermaid
stateDiagram-v2
    [*] --> DRAFT
    DRAFT --> SUBMITTED
    SUBMITTED --> UNDER_REVIEW
    UNDER_REVIEW --> RETURNED_TO_FIELD: Missing evidence / witness data
    UNDER_REVIEW --> REJECTED: Invalid / fraudulent claim
    UNDER_REVIEW --> PUBLISHED_FOR_OBJECTION: Preliminary approval
    
    PUBLISHED_FOR_OBJECTION --> DISPUTED: Public objection registered
    DISPUTED --> UNDER_REVIEW: Returned to review after mediation
    
    PUBLISHED_FOR_OBJECTION --> REGISTERED: 30-60 day window expires with zero disputes
    REGISTERED --> [*]
```

