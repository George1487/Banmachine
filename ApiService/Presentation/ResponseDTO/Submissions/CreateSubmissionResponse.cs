using Domain.Submissions;

namespace Presentation.ResponseDTO.Submissions;

public sealed record CreateSubmissionResponse(
    
    Guid SubmissionId,

    Guid LabId,

    Guid StudentId,

    SubmissionStatus Status,
    
    DateTimeOffset SubmittedAt
    );