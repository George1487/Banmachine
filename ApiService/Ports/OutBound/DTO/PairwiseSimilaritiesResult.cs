using Domain.Jobs;

namespace Ports.OutBound.DTO;

public abstract record PairwiseSimilaritiesResult
{
    protected PairwiseSimilaritiesResult() {}
    
    public sealed record Failure(string Reason) : PairwiseSimilaritiesResult;
    
    public sealed record Success(List<PairwiseSimilarity> PairwiseSimilarity) 
        : PairwiseSimilaritiesResult;
}