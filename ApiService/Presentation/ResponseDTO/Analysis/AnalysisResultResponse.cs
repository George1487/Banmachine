using Ports.InBound.DTO;

namespace Presentation.ResponseDTO.Analysis;

public sealed record AnalysisResultResponse(
    
    Guid LabId,
    
    LastAnalysisJob LastAnalysisJob,
    
    AnalysisStats Stats,
    
    List<SubItem> Items
    );