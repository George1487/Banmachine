using Domain.Jobs;
using Infrastructure.Repositories;
using Ports.OutBound.Contracts;
using Ports.OutBound.DTO;

namespace Infrastructure.Adapters;

public class AnalysisRepositoryAdapter : IAnalysisRepository
{
    private readonly AnalysisRepository _repository;

    public AnalysisRepositoryAdapter(AnalysisRepository repository)
    {
        _repository = repository;
    }

    public AnalysisJobResult GetAnalysisJobById(Guid jobId)
    {
        return _repository.GetAnalysisJobById(jobId);
    }

    public AnalysisJobResult AddAnalysisJob(AnalysisJob analysisJob)
    {
        return _repository.AddAnalysisJob(analysisJob);
    }

    public AnalysisJobResult GetAnalysisJobByLabId(Guid labId)
    {
        return _repository.GetAnalysisJobByLabId(labId);
    }

    public AnalysisJobResult GetAnalysisJobBySubmissionId(Guid submissionId)
    {
        return _repository.GetAnalysisJobBySubmissionId(submissionId);
    }

    public AnalysisJobResult GetAnalysisJobByLastDate()
    {
        return _repository.GetAnalysisJobByLastDate();
    }

    public AnalysisStatsResults GetAnalysisStatsByLabId(Guid labId)
    {
        return _repository.GetAnalysisStatsByLabId(labId);
    }
}
