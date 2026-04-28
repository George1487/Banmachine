namespace Infrastructure.Entities;

public class SubmissionAnalysisSummaryEntity
{
    public Guid Id { get; set; }

    public Guid AnalysisJobId { get; set; }

    public Guid SubmissionId { get; set; }

    public Guid? TopMatchSubmissionId { get; set; }

    public decimal? TopMatchScore { get; set; }

    public string FinalScoreRiskLevel { get; set; } = null!;
}
