using Domain.Jobs;

namespace Presentation.ResponseDTO.Analysis;

public sealed record LastAnalysisJob(
    
    Guid Id,
    
    JobStatus Status,
    
    DateTimeOffset CreatedAt,
    
    DateTimeOffset? FinishedAt
    );