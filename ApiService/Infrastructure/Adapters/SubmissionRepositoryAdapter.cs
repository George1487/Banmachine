using Domain.Submissions;
using Infrastructure.Repositories;
using Ports.OutBound.Contracts;
using Ports.OutBound.DTO;

namespace Infrastructure.Adapters;

public class SubmissionRepositoryAdapter : ISubmissionRepository
{
    private readonly SubmissionRepository _repository;

    public SubmissionRepositoryAdapter(SubmissionRepository repository)
    {
        _repository = repository;
    }

    public SubmissionResult GetSubmission(Guid submissionId)
    {
        return _repository.GetSubmission(submissionId);
    }

    public SubmissionsResult GetSubmissionsByLabId(Guid labId)
    {
        return _repository.GetSubmissionsByLabId(labId);
    }

    public SubmissionsResult GetSubmissions()
    {
        return _repository.GetSubmissions();
    }

    public SubmissionsResult GetSubmissionsByUserId(Guid userId)
    {
        return _repository.GetSubmissionsByUserId(userId);
    }

    public ParsedSubmissionsResult GetParsedSubmissions()
    {
        return _repository.GetParsedSubmissions();
    }

    public SubmissionResult AddSubmission(Submission submission)
    {
        return _repository.AddSubmission(submission);
    }

    public ParsedSubmissionsResult GetParsedSubmissionByLabId(Guid labId)
    {
        return _repository.GetParsedSubmissionByLabId(labId);
    }
}
