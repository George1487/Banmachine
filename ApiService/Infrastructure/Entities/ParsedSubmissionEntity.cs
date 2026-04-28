namespace Infrastructure.Entities;

public class ParsedSubmissionEntity
{
    public Guid Id { get; set; }

    public Guid SubmissionId { get; set; }

    public string RawText { get; set; } = null!;

    public string StructuredData { get; set; } = null!;

    public DateTimeOffset ParsedAt { get; set; }
}
