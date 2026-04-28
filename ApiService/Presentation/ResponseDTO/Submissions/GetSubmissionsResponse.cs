using Domain.Submissions;

namespace Presentation.ResponseDTO.Submissions;

public sealed record GetSubmissionsResponse(
    
    Guid SubmissionId,

    Guid LabId,

    string LabTitle,

    SubmissionStatus Status,
    
    DateTimeOffset SubmittedAt);