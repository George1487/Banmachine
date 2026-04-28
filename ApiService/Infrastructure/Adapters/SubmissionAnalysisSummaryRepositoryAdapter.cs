using Infrastructure.Repositories;
using Ports.OutBound.Contracts;
using Ports.OutBound.DTO;

namespace Infrastructure.Adapters;

public class SubmissionAnalysisSummaryRepositoryAdapter : ISubmissionAnalysisSummaryRepository
{
    private readonly SubmissionAnalysisSummaryRepository _repository;

    public SubmissionAnalysisSummaryRepositoryAdapter(SubmissionAnalysisSummaryRepository repository)
    {
        _repository = repository;
    }

    public SubmissionAnalysisSummaryResult GetSubmissionAnalysisSummaryBySubmissionId(Guid submissionId)
    {
        return _repository.GetSubmissionAnalysisSummaryBySubmissionId(submissionId);
    }
}
