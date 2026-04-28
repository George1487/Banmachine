using Ports.OutBound.DTO;

namespace Ports.OutBound.Contracts;

public interface ISubmissionAnalysisSummaryRepository
{
    
    public SubmissionAnalysisSummaryResult GetSubmissionAnalysisSummaryBySubmissionId(Guid submissionId);
}