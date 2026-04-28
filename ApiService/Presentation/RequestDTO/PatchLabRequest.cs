using Domain.Labs;

namespace Presentation.RequestDTO;

public sealed record PatchLabRequest(
    
    Guid LabId,
    
    string Title,
    
    string Description,
    
    LabStatus Status,
    
    DateTimeOffset DeadlineAt
    
    );