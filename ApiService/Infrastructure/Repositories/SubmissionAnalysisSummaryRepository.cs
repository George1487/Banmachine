using Infrastructure.Mappers;
using Microsoft.EntityFrameworkCore;
using Ports.OutBound.DTO;

namespace Infrastructure.Repositories;

public class SubmissionAnalysisSummaryRepository
{
    private readonly AppDbContext _context;

    public SubmissionAnalysisSummaryRepository(AppDbContext context)
    {
        _context = context;
    }

    public SubmissionAnalysisSummaryResult GetSubmissionAnalysisSummaryBySubmissionId(Guid submissionId)
    {
        try
        {
            var entity = _context.SubmissionAnalysisSummaries
                .AsNoTracking()
                .Where(x => x.SubmissionId == submissionId)
                .Join(
                    _context.AnalysisJobs.AsNoTracking(),
                    summary => summary.AnalysisJobId,
                    job => job.Id,
                    (summary, job) => new { Summary = summary, job.CreatedAt })
                .OrderByDescending(x => x.CreatedAt)
                .Select(x => x.Summary)
                .FirstOrDefault();

            return entity is null
                ? new SubmissionAnalysisSummaryResult.Failure("submission_analysis_summary_not_found")
                : new SubmissionAnalysisSummaryResult.Success(SubmissionAnalysisSummaryMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new SubmissionAnalysisSummaryResult.Failure(exception.Message);
        }
    }
}
