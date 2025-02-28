namespace AdminAPI;

model GetLatestEntriesRow {
  Path: string;
  Title: string;
  Body: string;
  Visibility: string;
  Format: string;
  PublishedAt?: utcDateTime | null;
  LastEditedAt?: utcDateTime | null;
  CreatedAt?: utcDateTime | null;
  UpdatedAt?: utcDateTime | null;
  ImageUrl?: string | null;
}
