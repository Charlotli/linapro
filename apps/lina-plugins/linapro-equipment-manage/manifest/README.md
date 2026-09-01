# Manifest Resources

This directory stores install and uninstall SQL resources for `linapro-equipment-manage`.

## Contents

- `sql/001-linapro-equipment-manage-schema.sql`: creates the equipment and maintenance tables with indexes and dictionary seeds
- `sql/uninstall/001-linapro-equipment-manage-schema.sql`: removes dictionary seeds and drops the tables on data cleanup
- `sql/mock-data/001-linapro-equipment-manage-mock-data.sql`: optional demo equipment and maintenance record
