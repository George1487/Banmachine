using Infrastructure.Mappers;
using Microsoft.EntityFrameworkCore;
using Ports.OutBound.DTO;

namespace Infrastructure.Repositories;

public class PairwiseSimilarityRepository
{
    private readonly AppDbContext _context;

    public PairwiseSimilarityRepository(AppDbContext context)
    {
        _context = context;
    }

    public PairwiseSimilaritiesResult GetPairwiseSimilarityBySubmissionId(Guid submissionId)
    {
        try
        {
            var latestJobId = _context.PairwiseSimilarities
                .AsNoTracking()
                .Where(x => x.LeftSubmissionId == submissionId || x.RightSubmissionId == submissionId)
                .Join(
                    _context.AnalysisJobs.AsNoTracking(),
                    similarity => similarity.AnalysisJobId,
                    job => job.Id,
                    (similarity, job) => new { similarity.AnalysisJobId, job.CreatedAt })
                .OrderByDescending(x => x.CreatedAt)
                .Select(x => (Guid?)x.AnalysisJobId)
                .FirstOrDefault();

            if (!latestJobId.HasValue)
            {
                return new PairwiseSimilaritiesResult.Failure("pairwise_similarity_not_found");
            }

            var similarities = _context.PairwiseSimilarities
                .AsNoTracking()
                .Where(x =>
                    x.AnalysisJobId == latestJobId.Value &&
                    (x.LeftSubmissionId == submissionId || x.RightSubmissionId == submissionId))
                .ToList()
                .Select(x => PairwiseSimilarityMapper.ToDomain(x, submissionId))
                .ToList();

            return new PairwiseSimilaritiesResult.Success(similarities);
        }
        catch (Exception exception)
        {
            return new PairwiseSimilaritiesResult.Failure(exception.Message);
        }
    }
}
