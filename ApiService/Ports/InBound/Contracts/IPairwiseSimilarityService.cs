using Ports.OutBound.DTO;

namespace Ports.InBound.Contracts;

public interface IPairwiseSimilarityService
{
    
    PairwiseSimilaritiesResult GetPairwiseSimilarityBySubmissionId(Guid submissionId);
    
}