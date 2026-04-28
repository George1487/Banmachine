namespace Ports.InBound.DTO;

public abstract record LabSubmissionsInfoResult
{
    private LabSubmissionsInfoResult() { }
    
    public sealed record Failure(string Reason) : LabSubmissionsInfoResult;
    
    public sealed record Success(int SubmissionCount, int ParsedSubmissionCount) 
        : LabSubmissionsInfoResult;
}