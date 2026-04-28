namespace Domain.Jobs;

public sealed record PairwiseSimilarity(
    
    Guid Id,
    
    Guid AnalysisJobId,
    
    Guid LabId,
    
    Guid LeftSubmissionId,
    
    Guid RightSubmissionId,
    
    decimal TextScore,
    
    decimal CalculationScore,
    
    decimal ImagesScore,
    
    decimal FinalScore
);