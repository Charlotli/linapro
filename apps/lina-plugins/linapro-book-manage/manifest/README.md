# Manifest Resources

This directory stores install and uninstall SQL resources for `linapro-book-manage`.

## Contents

- `sql/001-linapro-book-manage-schema.sql`: creates the book and borrow tables with indexes and dictionary seeds
- `sql/uninstall/001-linapro-book-manage-schema.sql`: removes dictionary seeds and drops the tables on data cleanup
- `sql/mock-data/001-linapro-book-manage-mock-data.sql`: optional demo book and borrow record
