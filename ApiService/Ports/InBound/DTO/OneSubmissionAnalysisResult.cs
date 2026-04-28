namespace Ports.InBound.DTO;

public abstract record OneSubmissionAnalysisResult
{
    
    private OneSubmissionAnalysisResult() { }
    
    public sealed record Failure(string Reason) : OneSubmissionAnalysisResult;
    
    public sealed record Success(Guid SubmissionId, 
        Guid JobId,
        List<SubMatch> SubMatches) : OneSubmissionAnalysisResult;
}