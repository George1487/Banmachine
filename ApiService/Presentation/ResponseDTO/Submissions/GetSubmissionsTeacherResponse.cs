using Domain.Submissions;

namespace Presentation.ResponseDTO.Submissions;

public sealed record GetSubmissionsTeacherResponse(
    
    Guid SubmissionId,

    Students Students,

    SubmissionStatus Status,

    DateTimeOffset SubmittedAt,
    
    decimal? TopMatchScore,

    string? FinalScoreRiskLevel
    );