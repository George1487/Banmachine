namespace Infrastructure.Entities;

public class PairwiseSimilarityEntity
{
    public Guid Id { get; set; }

    public Guid AnalysisJobId { get; set; }

    public Guid LabId { get; set; }

    public Guid LeftSubmissionId { get; set; }

    public Guid RightSubmissionId { get; set; }

    public decimal TextScore { get; set; }

    public decimal CalculationScore { get; set; }

    public decimal ImagesScore { get; set; }

    public decimal FinalScore { get; set; }
}
