namespace Ports.InBound.DTO;

public sealed record AnalysisStats(
    
    int TotalSubmissions,
    
    int ActualSubmissions,
    
    int ParsedSubmissions,
    
    int HighRiskCount,
    
    int MediumRiskCount,
    
    int LowRiskCount,
    
    decimal MaxFinalScore
);