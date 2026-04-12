# BowMobile

## Project Overview
BowMobile is a Flutter frontend for browsing the Bow SQLite parts database on mobile devices. This repository should stay focused on the mobile app plus a shared schema reference, not the builder or desktop frontend.

## Tech Stack
- **Frontend**: Flutter
- **Database**: SQLite database bundled as an app asset
- **Language**: Dart

## Critical Files
- `frontends/mobile-app/lib/data/database_helper.dart`: Copies the bundled database into app storage and manages refreshes.
- `frontends/mobile-app/lib/data/search_service.dart`: Implements the shared search behavior against the SQLite database.
- `frontends/mobile-app/pubspec.yaml`: Declares bundled assets including `assets/parts.db`.
- `shared/schema.sql`: Shared SQLite schema reference for the frontend/builder contract.

## Search Logic Mandates
- Search must support partial part numbers (e.g. `WG8`, `5935`) and full part numbers (`WG8-5935`).
- Normalize hyphens when matching part numbers and base parts.
- Always search across `part_number`, `base_part`, `description`, and `remarks`.

## Database Schema Notes
- The `parts` table must include a `remarks` column.
- Keep `shared/schema.sql` aligned with the builder repo before updating bundled database assets.

## Development Workflows
- Work from `frontends/mobile-app`.
- Use `flutter test` for the existing test suite.
