using Domain.Labs;

namespace Presentation.ResponseDTO.Labs;

public sealed record LabDetailedInfoResponse(
    
    Guid LabId,
    
    string Title,
    
    string Description,
    
    LabStatus LabStatus,
    
    DateTimeOffset DeadlineAt);