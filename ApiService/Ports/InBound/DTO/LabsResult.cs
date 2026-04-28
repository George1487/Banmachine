using Domain.Labs;

namespace Ports.InBound.DTO;

public record LabsResult
{
    private LabsResult() {}
   
    public sealed record Success(List<Lab> Labs) : LabsResult;
   
    public sealed record Failure(string Reason) : LabsResult;
}