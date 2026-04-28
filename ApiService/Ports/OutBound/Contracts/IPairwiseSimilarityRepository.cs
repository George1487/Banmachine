using Ports.OutBound.DTO;

namespace Ports.OutBound.Contracts;

public interface IPairwiseSimilarityRepository
{
    
    PairwiseSimilaritiesResult GetPairwiseSimilarityBySubmissionId(Guid submissionId);
    
}