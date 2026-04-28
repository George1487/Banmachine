using Domain.Submissions;

namespace Ports.OutBound.DTO;

public abstract record SubmissionResult
{
    private SubmissionResult() { }
    
    public sealed record Failure(string Reason) : SubmissionResult;
    
    public sealed record Success(Submission Submission) : SubmissionResult;
}