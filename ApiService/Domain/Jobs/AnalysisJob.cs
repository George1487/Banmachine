namespace Domain.Jobs;

public sealed record AnalysisJob(
    
    Guid JobId,
    
    Guid LabId,
    
    JobStatus Status,
    
    Guid UserId,
    
    DateTimeOffset CreatedAt,
    
    DateTimeOffset? StartedAt,
    
    DateTimeOffset? FinishedAt,
    
    string ErrorMessage
    
    );