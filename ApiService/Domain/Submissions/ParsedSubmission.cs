namespace Domain.Submissions;

public sealed record ParsedSubmission(
    
    Guid ParsedSubmissionId,
    
    Guid SubmissionId,
    
    DateTimeOffset ParsedAt);
