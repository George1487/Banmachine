using Domain.Jobs;

namespace Presentation.ResponseDTO.Analysis;

public sealed record PostAnalysisResponse(
    
    Guid AnalysisJobId,
    
    Guid LabId,
    
    JobStatus Status
    );