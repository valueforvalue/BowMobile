# parts.db — Database Schema

This document describes the SQLite schema for `parts.db`, the Canon copier parts cross-reference database used by **Bow Mobile**.

The database is built by a separate tool (the CanonBow PDF parser / database builder) and placed in `bow_mobile/assets/parts.db` before building the APK.

---

## Tables

### `manuals`
Tracks each PDF parts catalog that was ingested.

| Column         | Type    | Notes                         |
|----------------|---------|-------------------------------|
| `id`           | INTEGER | Primary key, auto-increment   |
| `filename`     | TEXT    | Source PDF filename (UNIQUE)  |
| `model_series` | TEXT    | Copier model / series name    |
| `revision`     | TEXT    | Catalog revision identifier   |

---

### `figures`
Represents figure/diagram pages within a manual.

| Column      | Type    | Notes                                              |
|-------------|---------|-----------------------------------------------------|
| `manual_id` | INTEGER | FK → `manuals.id`                                  |
| `id`        | TEXT    | Figure identifier (e.g. "FIG-3A")                  |

Primary key: `(manual_id, id)`

---

### `parts`
Core parts data — one row per part occurrence per figure.

| Column       | Type    | Notes                                          |
|--------------|---------|------------------------------------------------|
| `id`         | INTEGER | Primary key, auto-increment                    |
| `manual_id`  | INTEGER | FK → `manuals.id`                              |
| `figure_id`  | TEXT    | Figure the part appears in                     |
| `key_no`     | TEXT    | Key number on the diagram                      |
| `part_number`| TEXT    | Full part number (may include revision suffix) |
| `base_part`  | TEXT    | Normalized base part number (hyphens kept)     |
| `revision`   | TEXT    | Part revision letter/code                      |
| `qty`        | TEXT    | Quantity listed                                |
| `description`| TEXT    | Part description                               |
| `remarks`    | TEXT    | Schematic notes, connectors, locations, etc.   |

Index: `idx_base_part ON parts(base_part)`

---

### `metadata`
Single-row table storing database-level metadata.

| Column         | Type | Notes                            |
|----------------|------|----------------------------------|
| `id`           | INTEGER | Always 1                      |
| `last_updated` | TEXT    | ISO-8601 timestamp of last build |

---

## Notes

- Searches normalize part numbers by stripping hyphens so `FM2-A087-020` and `FM2A087020` both match.
- `base_part` stores the canonical form used for grouping results across multiple manuals.
- The `parts.db` file must be placed at `bow_mobile/assets/parts.db` **before** running `flutter build apk`.
