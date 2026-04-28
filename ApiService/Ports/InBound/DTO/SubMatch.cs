namespace Ports.InBound.DTO;

public record SubMatch(
    
    Guid OtherSubmissionId,
    
    Student Student,
    
    decimal TextScore,
    
    decimal CalculationScore,
    
    decimal ImagesScore,
    
    decimal FinalScore,
    
    string RiskLevel
);