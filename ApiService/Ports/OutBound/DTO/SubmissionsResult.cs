using Domain.Submissions;

namespace Ports.OutBound.DTO;

public abstract record SubmissionsResult
{
    private SubmissionsResult() { }
    
    public sealed record Failure(string Reason) : SubmissionsResult;
    
    public sealed record Success(List<Submission> Submissions) : SubmissionsResult;
}