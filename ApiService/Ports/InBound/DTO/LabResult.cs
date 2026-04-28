using Domain.Labs;

namespace Ports.InBound.DTO;

public abstract record LabResult
{
    
   private LabResult() {}
   
   public sealed record Success(Lab Lab) : LabResult;
   
   public sealed record Failure(string Reason) : LabResult;
    
}