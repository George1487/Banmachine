using Ports.OutBound.DTO;

namespace Ports.InBound.Contracts;

public interface ISubmissionAnalysisSummaryService
{
    
    public SubmissionAnalysisSummaryResult GetSubmissionAnalysisSummaryBySubmissionId(Guid submissionId);
    
}