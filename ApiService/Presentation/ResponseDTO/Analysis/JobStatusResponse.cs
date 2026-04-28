using Domain.Jobs;

namespace Presentation.ResponseDTO.Analysis;

public sealed record JobStatusResponse(
    
    Guid JobId,
    
    Guid LabId,
    
    JobStatus Status,
    
    CreatedBy CreatedBy,
    
    DateTimeOffset CreatedAt,
    
    DateTimeOffset? StartedAt,
    
    DateTimeOffset? FinishedAt,
    
    string? ErrorMessage
    );