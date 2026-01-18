# Go Client for Tax1099.com

A very basic go client for tax1099.com's api.

## Downloading PDFs

Use `DownloadFilledForm` to download form PDFs. There are two supported request styles:

- **Single PDF:** provide `formId` and `formType`.
- **Multiple PDFs:** provide `payerTin`, `taxYear`, and `formType`.
