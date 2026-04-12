# Bow mobile app

Flutter client for the Bow parts cross-reference database.

## Run locally
1. `flutter pub get`
2. `flutter test`
3. `flutter run`

## Assets
- `assets/parts.db`: Bundled SQLite database shipped with the app.
- `assets/images/logo.png`: Application branding.

## Contract
- Search behavior must stay aligned with the builder and other frontends.
- The schema reference for this repo lives at `..\..\shared\schema.sql`.
