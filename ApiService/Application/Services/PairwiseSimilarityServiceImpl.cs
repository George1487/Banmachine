using Ports.InBound.Contracts;
using Ports.OutBound.Contracts;
using Ports.OutBound.DTO;

namespace Application.Services;

public class PairwiseSimilarityServiceImpl : IPairwiseSimilarityService
{

    private readonly IPairwiseSimilarityRepository _repo;

    public PairwiseSimilarityServiceImpl(IPairwiseSimilarityRepository repo)
    {
        _repo = repo;
    }
    
    public PairwiseSimilaritiesResult GetPairwiseSimilarityBySubmissionId(Guid submissionId)
    {
        return _repo.GetPairwiseSimilarityBySubmissionId(submissionId);
    }
}