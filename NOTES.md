# Notes

Contains additional notes that don't fit anywhere else.

## Pre-shipping considerations

Consider storing app instance information in the DB for backwards analytics, containing the version(commit id), timestamps, JWT public keys if applicable.

## Auditing

Consider generating a report when someone changes devices mid-session on a single token.
