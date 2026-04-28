namespace Ports.OutBound.DTO;

public abstract record SubmissionAnalysisSummary
{
    
    private SubmissionAnalysisSummary() { }
    
    public sealed record Failure(string Reason) : SubmissionAnalysisSummary;
    
    public sealed record Success(SubmissionAnalysisSummary SubmissionAnalysisSummary) : 
        SubmissionAnalysisSummary;
}