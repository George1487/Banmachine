namespace Infrastructure.Entities;

public class AnalysisJobEntity
{
    public Guid Id { get; set; }

    public Guid LabId { get; set; }

    public string Status { get; set; } = null!;

    public Guid CreatedBy { get; set; }

    public DateTimeOffset CreatedAt { get; set; }

    public DateTimeOffset? StartedAt { get; set; }

    public DateTimeOffset? FinishedAt { get; set; }

    public string? ErrorMessage { get; set; }
}
