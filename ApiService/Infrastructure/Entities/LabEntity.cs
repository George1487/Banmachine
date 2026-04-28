namespace Infrastructure.Entities;

public class LabEntity
{
    public Guid Id { get; set; }

    public Guid TeacherId { get; set; }

    public string Title { get; set; } = null!;

    public string? Description { get; set; }

    public string Status { get; set; } = null!;

    public DateTimeOffset DeadlineAt { get; set; }
}
