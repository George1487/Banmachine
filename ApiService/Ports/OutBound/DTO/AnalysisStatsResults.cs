using Ports.InBound.DTO;

namespace Ports.OutBound.DTO;

public abstract record AnalysisStatsResults
{
    private AnalysisStatsResults() {}
    
    public sealed record Failure(string Reason) : AnalysisStatsResults;
    
    public sealed record Success(AnalysisStats Stats) : AnalysisStatsResults;
}