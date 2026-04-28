using Domain.Jobs;
using Infrastructure.Entities;

namespace Infrastructure.Mappers;

public static class PairwiseSimilarityMapper
{
    public static PairwiseSimilarity ToDomain(
        PairwiseSimilarityEntity entity,
        Guid? focusSubmissionId = null)
    {
        var leftSubmissionId = entity.LeftSubmissionId;
        var rightSubmissionId = entity.RightSubmissionId;

        if (focusSubmissionId.HasValue && entity.RightSubmissionId == focusSubmissionId.Value)
        {
            leftSubmissionId = entity.RightSubmissionId;
            rightSubmissionId = entity.LeftSubmissionId;
        }

        return new PairwiseSimilarity(
            entity.Id,
            entity.AnalysisJobId,
            entity.LabId,
            leftSubmissionId,
            rightSubmissionId,
            entity.TextScore,
            entity.CalculationScore,
            entity.ImagesScore,
            entity.FinalScore);
    }
}
