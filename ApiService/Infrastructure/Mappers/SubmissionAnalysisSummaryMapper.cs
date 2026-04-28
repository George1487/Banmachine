using Domain.Jobs;
using Infrastructure.Entities;

namespace Infrastructure.Mappers;

public static class SubmissionAnalysisSummaryMapper
{
    public static SubmissionAnalysisSummary ToDomain(SubmissionAnalysisSummaryEntity entity)
    {
        return new SubmissionAnalysisSummary(
            entity.Id,
            entity.AnalysisJobId,
            entity.SubmissionId,
            entity.TopMatchSubmissionId,
            entity.TopMatchScore,
            entity.FinalScoreRiskLevel);
    }
}
