namespace Ports.InBound.DTO;

public sealed record SubItem(
    Guid SubmissionId,
    Student Student,
    Guid? TopMatchSubmissionId,
    decimal? TopMatchScore,
    string FinalScoreRiskLevel);