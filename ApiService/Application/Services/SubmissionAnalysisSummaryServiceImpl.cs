using Ports.InBound.Contracts;
using Ports.OutBound.Contracts;
using Ports.OutBound.DTO;

namespace Application.Services;

public class SubmissionAnalysisSummaryServiceImpl : ISubmissionAnalysisSummaryService
{
    
    private readonly ISubmissionAnalysisSummaryRepository _repo;

    public SubmissionAnalysisSummaryServiceImpl(
        ISubmissionAnalysisSummaryRepository submissionAnalysisSummaryRepository)
    {
        _repo = submissionAnalysisSummaryRepository;
    }
    
    public SubmissionAnalysisSummaryResult GetSubmissionAnalysisSummaryBySubmissionId(
        Guid submissionId)
    {
        return _repo.GetSubmissionAnalysisSummaryBySubmissionId(submissionId);
    }
}