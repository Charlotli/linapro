# Manifest Resources

This directory stores install and uninstall SQL resources for `linapro-oa-approval`.

## Contents

- `sql/001-linapro-oa-approval-schema.sql`: creates the flow, node, request, and record tables with indexes and dictionary seeds
- `sql/uninstall/001-linapro-oa-approval-schema.sql`: removes dictionary seeds and drops the tables on data cleanup
- `sql/mock-data/001-linapro-oa-approval-mock-data.sql`: optional demo flow and request for local demos
