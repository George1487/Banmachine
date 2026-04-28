namespace Infrastructure.Entities;

public class IngestJobEntity
{
    public Guid Id { get; set; }

    public Guid SubmissionId { get; set; }

    public string Status { get; set; } = null!;

    public DateTimeOffset CreatedAt { get; set; }

    public DateTimeOffset? FinishedAt { get; set; }

    public string? ErrorMessage { get; set; }
}
