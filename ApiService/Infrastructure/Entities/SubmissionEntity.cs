namespace Infrastructure.Entities;

public class SubmissionEntity
{
    public Guid Id { get; set; }

    public Guid LabId { get; set; }

    public Guid StudentId { get; set; }

    public string Status { get; set; } = null!;

    public string SourceFileName { get; set; } = null!;

    public string MimeType { get; set; } = null!;

    public string StorageKey { get; set; } = null!;

    public DateTimeOffset SubmittedAt { get; set; }
}
