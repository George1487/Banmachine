using Infrastructure.Repositories;
using Ports.OutBound.Contracts;
using Ports.OutBound.DTO;

namespace Infrastructure.Adapters;

public class PairwiseSimilarityRepositoryAdapter : IPairwiseSimilarityRepository
{
    private readonly PairwiseSimilarityRepository _repository;

    public PairwiseSimilarityRepositoryAdapter(PairwiseSimilarityRepository repository)
    {
        _repository = repository;
    }

    public PairwiseSimilaritiesResult GetPairwiseSimilarityBySubmissionId(Guid submissionId)
    {
        return _repository.GetPairwiseSimilarityBySubmissionId(submissionId);
    }
}
