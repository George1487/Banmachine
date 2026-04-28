using Domain.Jobs;

namespace Ports.OutBound.DTO;

public abstract record AnalysisJobResult()
{
    public sealed record Failure(string Reason) : AnalysisJobResult;
    
    public sealed record Success(AnalysisJob AnalysisJob) : AnalysisJobResult;
    
}