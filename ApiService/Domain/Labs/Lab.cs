namespace Domain.Labs;

public sealed record Lab(
    
    Guid LabId,
    
    Guid TeacherId,
    
    string Title,
    
    string Description,
    
    LabStatus LabStatus,
    
    DateTimeOffset DeadlineAt
    );