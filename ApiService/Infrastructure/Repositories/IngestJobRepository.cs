using Domain.Jobs;
using Infrastructure.Mappers;
using Ports.OutBound.DTO;

namespace Infrastructure.Repositories;

public class IngestJobRepository
{
    private readonly AppDbContext _context;

    public IngestJobRepository(AppDbContext context)
    {
        _context = context;
    }

    public IngestJobResult AddIngestJob(IngestJob ingestJob)
    {
        try
        {
            var entity = IngestJobMapper.ToEntity(ingestJob);
            _context.IngestJobs.Add(entity);
            _context.SaveChanges();

            return new IngestJobResult.Success(IngestJobMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new IngestJobResult.Failure(exception.Message);
        }
    }
}
