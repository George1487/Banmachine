using Domain.Labs;

namespace Presentation.ResponseDTO.Labs;

public sealed record LabDefaultInfoResponse(
        
    Guid LabId,
    
    string Title,
    
    LabStatus LabStatus,
    
    DateTimeOffset DeadlineAt);