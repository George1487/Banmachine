using Domain.Jobs;

namespace Ports.OutBound.DTO;

public abstract record IngestJobResult
{
    
    private IngestJobResult() { }
    
    public sealed record Failure(string Reason) : IngestJobResult;
    
    public sealed record Success(IngestJob IngestJob) : IngestJobResult;
}