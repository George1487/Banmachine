using Domain.Jobs;
using Infrastructure.Repositories;
using Ports.OutBound.Contracts;
using Ports.OutBound.DTO;

namespace Infrastructure.Adapters;

public class IngestJobRepositoryAdapter : IIngestJobRepository
{
    private readonly IngestJobRepository _repository;

    public IngestJobRepositoryAdapter(IngestJobRepository repository)
    {
        _repository = repository;
    }

    public IngestJobResult AddIngestJob(IngestJob ingestJob)
    {
        return _repository.AddIngestJob(ingestJob);
    }
}
