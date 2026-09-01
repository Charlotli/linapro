# Manifest Resources

This directory stores install and uninstall SQL resources for `linapro-live-manage`.

## Contents

- `sql/001-linapro-live-manage-schema.sql`: creates the live-room and live-content tables, indexes, and dictionary seeds
- `sql/uninstall/001-linapro-live-manage-schema.sql`: removes the dictionary seeds and drops both tables when uninstall chooses data cleanup
- `sql/mock-data/001-linapro-live-manage-mock-data.sql`: optional demo rooms and live content for local demos
