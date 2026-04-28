using Ports.InBound.DTO;

namespace Infrastructure.Mappers;

public static class AnalysisStatsMapper
{
    public static AnalysisStats ToDomain(
        int totalSubmissions,
        int actualSubmissions,
        int parsedSubmissions,
        int highRiskCount,
        int mediumRiskCount,
        int lowRiskCount,
        decimal maxFinalScore)
    {
        return new AnalysisStats(
            totalSubmissions,
            actualSubmissions,
            parsedSubmissions,
            highRiskCount,
            mediumRiskCount,
            lowRiskCount,
            maxFinalScore);
    }
}
