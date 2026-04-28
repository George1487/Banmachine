using Infrastructure.Mappers;
using Microsoft.EntityFrameworkCore;
using Ports.OutBound.DTO;

namespace Infrastructure.Repositories;

public class AnalysisRepository
{
    private readonly AppDbContext _context;

    public AnalysisRepository(AppDbContext context)
    {
        _context = context;
    }

    public AnalysisJobResult GetAnalysisJobById(Guid jobId)
    {
        try
        {
            var entity = _context.AnalysisJobs
                .AsNoTracking()
                .FirstOrDefault(x => x.Id == jobId);

            return entity is null
                ? new AnalysisJobResult.Failure("analysis_job_not_found")
                : new AnalysisJobResult.Success(AnalysisJobMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new AnalysisJobResult.Failure(exception.Message);
        }
    }

    public AnalysisJobResult AddAnalysisJob(Domain.Jobs.AnalysisJob analysisJob)
    {
        try
        {
            var entity = AnalysisJobMapper.ToEntity(analysisJob);
            _context.AnalysisJobs.Add(entity);
            _context.SaveChanges();

            return new AnalysisJobResult.Success(AnalysisJobMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new AnalysisJobResult.Failure(exception.Message);
        }
    }

    public AnalysisJobResult GetAnalysisJobByLabId(Guid labId)
    {
        try
        {
            var entity = _context.AnalysisJobs
                .AsNoTracking()
                .Where(x => x.LabId == labId)
                .OrderByDescending(x => x.CreatedAt)
                .FirstOrDefault();

            return entity is null
                ? new AnalysisJobResult.Failure("analysis_job_not_found")
                : new AnalysisJobResult.Success(AnalysisJobMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new AnalysisJobResult.Failure(exception.Message);
        }
    }

    public AnalysisJobResult GetAnalysisJobBySubmissionId(Guid submissionId)
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
                    (_, job) => job)
                .OrderByDescending(x => x.CreatedAt)
                .FirstOrDefault();

            if (entity is null)
            {
                var labId = _context.Submissions
                    .AsNoTracking()
                    .Where(x => x.Id == submissionId)
                    .Select(x => (Guid?)x.LabId)
                    .FirstOrDefault();

                if (!labId.HasValue)
                {
                    return new AnalysisJobResult.Failure("submission_not_found");
                }

                entity = _context.AnalysisJobs
                    .AsNoTracking()
                    .Where(x => x.LabId == labId.Value)
                    .OrderByDescending(x => x.CreatedAt)
                    .FirstOrDefault();
            }

            return entity is null
                ? new AnalysisJobResult.Failure("analysis_job_not_found")
                : new AnalysisJobResult.Success(AnalysisJobMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new AnalysisJobResult.Failure(exception.Message);
        }
    }

    public AnalysisJobResult GetAnalysisJobByLastDate()
    {
        try
        {
            var entity = _context.AnalysisJobs
                .AsNoTracking()
                .OrderByDescending(x => x.CreatedAt)
                .FirstOrDefault();

            return entity is null
                ? new AnalysisJobResult.Failure("analysis_job_not_found")
                : new AnalysisJobResult.Success(AnalysisJobMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new AnalysisJobResult.Failure(exception.Message);
        }
    }

    public AnalysisStatsResults GetAnalysisStatsByLabId(Guid labId)
    {
        try
        {
            var submissionIds = _context.Submissions
                .AsNoTracking()
                .Where(x => x.LabId == labId)
                .Select(x => x.Id)
                .ToList();

            var totalSubmissions = submissionIds.Count;
            var parsedSubmissions = _context.ParsedSubmissions
                .AsNoTracking()
                .Count(x => submissionIds.Contains(x.SubmissionId));

            var latestJobId = _context.AnalysisJobs
                .AsNoTracking()
                .Where(x => x.LabId == labId)
                .OrderByDescending(x => x.CreatedAt)
                .Select(x => (Guid?)x.Id)
                .FirstOrDefault();

            var actualSubmissions = 0;
            var highRiskCount = 0;
            var mediumRiskCount = 0;
            var lowRiskCount = 0;
            decimal maxFinalScore = 0m;

            if (latestJobId.HasValue)
            {
                actualSubmissions = _context.SubmissionAnalysisSummaries
                    .AsNoTracking()
                    .Count(x => x.AnalysisJobId == latestJobId.Value);

                highRiskCount = _context.SubmissionAnalysisSummaries
                    .AsNoTracking()
                    .Count(x => x.AnalysisJobId == latestJobId.Value && x.FinalScoreRiskLevel == "high");

                mediumRiskCount = _context.SubmissionAnalysisSummaries
                    .AsNoTracking()
                    .Count(x => x.AnalysisJobId == latestJobId.Value && x.FinalScoreRiskLevel == "medium");

                lowRiskCount = _context.SubmissionAnalysisSummaries
                    .AsNoTracking()
                    .Count(x => x.AnalysisJobId == latestJobId.Value && x.FinalScoreRiskLevel == "low");

                maxFinalScore = _context.PairwiseSimilarities
                    .AsNoTracking()
                    .Where(x => x.AnalysisJobId == latestJobId.Value)
                    .Select(x => (decimal?)x.FinalScore)
                    .Max() ?? 0m;
            }

            return new AnalysisStatsResults.Success(
                AnalysisStatsMapper.ToDomain(
                    totalSubmissions,
                    actualSubmissions,
                    parsedSubmissions,
                    highRiskCount,
                    mediumRiskCount,
                    lowRiskCount,
                    maxFinalScore));
        }
        catch (Exception exception)
        {
            return new AnalysisStatsResults.Failure(exception.Message);
        }
    }
}
