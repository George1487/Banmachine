using Ports.OutBound.DTO;

namespace Ports.OutBound.Contracts;

public interface ISubAnalysisSummary
{
    SubmissionAnalysisSummaryResult GetAnalysisSummaryBySubmissionId(Guid submissionId);
    
    SubmissionAnalysisSummaryResult GetAnalysisSummaryByJobId(Guid jobId);
}