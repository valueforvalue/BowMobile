# Bow Mobile

**Bow Mobile** is a Flutter Android app for searching the Canon copier parts cross-reference database (`parts.db`).

## ⚠️ Before building the APK

Place an updated `parts.db` in `assets/parts.db` before running `flutter build apk`.
The database is **bundled inside the APK** — if you skip this, the app will ship with stale data.

```bash
# Copy latest DB (from the CanonBow builder tool):
cp /path/to/parts.db assets/parts.db

# !! REMINDER: Confirm parts.db is up to date before this step !!
flutter build apk --release
```

## Development

```bash
flutter pub get
flutter run
```

## Database Schema

See [`../DATABASE_SCHEMA.md`](../DATABASE_SCHEMA.md) for the full SQLite schema.

