using Ports.InBound.DTO;

namespace Presentation.ResponseDTO.Submissions;

public sealed record GetSubmissionsMatchesResponse(
    
    Guid SubmissionId,
    
    Guid AnalysisJobId,
    
    List<SubMatch> Matches
    
    );