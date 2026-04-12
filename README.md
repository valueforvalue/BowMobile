# BowMobile

BowMobile is the Flutter frontend for the Bow parts cross-reference database. This repo is now mobile-focused: it contains the Android app and a checked-in database schema reference, not the PDF builder or desktop frontend.

## Structure
- `frontends/mobile-app/`: Flutter application source, assets, and tests.
- `shared/schema.sql`: Shared SQLite schema reference kept aligned with the builder repo.

## Development
1. `cd frontends/mobile-app`
2. `flutter pub get`
3. `flutter test`
4. `flutter run`

## Database contract
- The app ships with a bundled `assets/parts.db` and queries the production `manuals`, `metadata`, and `parts` tables.
- `shared/schema.sql` is the schema reference this frontend should validate against when the builder changes.
- Keep `shared/schema.sql` in sync with the builder repo's canonical schema before shipping new database assets.

## Notes
- Search supports part number, normalized part number, description, and remarks matches.
- The mobile app copies `assets/parts.db` into app storage on first launch and can refresh that local copy from the bundled asset.
