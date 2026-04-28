using Domain.Jobs;

namespace Ports.InBound.DTO;

public abstract record FullLabAnalysisResult
{
    private FullLabAnalysisResult() { }
    
    public sealed record Failure(string Reason) : FullLabAnalysisResult;
    
    public sealed record Success(
        AnalysisJob Job,
        AnalysisStats Stats,
        List<SubItem> SubItems) 
        : FullLabAnalysisResult;
}