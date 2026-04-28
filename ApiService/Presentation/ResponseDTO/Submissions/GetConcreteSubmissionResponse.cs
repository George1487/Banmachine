using Domain.Submissions;

namespace Presentation.ResponseDTO.Submissions;

public sealed record GetConcreteSubmissionResponse(
    
    Guid SubmissionId,

    Guid LabId,

    Students Students,

    SubmissionStatus Status,

    DateTimeOffset SubmittedAt
    
    );