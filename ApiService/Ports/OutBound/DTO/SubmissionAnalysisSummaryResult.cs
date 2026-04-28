namespace Ports.OutBound.DTO;

public abstract record SubmissionAnalysisSummaryResult
{
    private SubmissionAnalysisSummaryResult() { }
    
    public sealed record Failure(string Reason) : SubmissionAnalysisSummaryResult;
    
    public sealed record Success(Domain.Jobs.SubmissionAnalysisSummary Summary) 
        : SubmissionAnalysisSummaryResult;
}