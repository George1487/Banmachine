namespace Domain.Jobs;

public sealed record SubmissionAnalysisSummary(

    Guid Id,
    
    Guid AnalysisJobId,
    
    Guid SubmissionId,
    
    Guid? TopMatchSubmissionId,
    
    decimal? TopMatchScore,
    
    string FinalScoreRiskLevel
);