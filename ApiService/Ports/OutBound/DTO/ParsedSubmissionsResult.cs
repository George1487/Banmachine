using Domain.Submissions;

namespace Ports.OutBound.DTO;

public abstract record ParsedSubmissionsResult
{
    private ParsedSubmissionsResult() {}
    
    public sealed record Failure(string Reason) : ParsedSubmissionsResult;
    
    public sealed record Success(List<ParsedSubmission> ParsedSubmission) 
        : ParsedSubmissionsResult;
}