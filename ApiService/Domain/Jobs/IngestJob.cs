namespace Domain.Jobs;

public sealed record IngestJob(
    
        Guid IngestJobId,
        
        Guid SubmissionId,
        
        JobStatus Status,
        
        DateTimeOffset CreatedAt,
        
        DateTimeOffset? StartedAt,
        
        DateTimeOffset? FinishedAt,
        
        string ErrorMessage
    );