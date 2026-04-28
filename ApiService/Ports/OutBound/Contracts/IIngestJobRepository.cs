using Domain.Jobs;
using Ports.OutBound.DTO;

namespace Ports.OutBound.Contracts;

public interface IIngestJobRepository
{
    IngestJobResult AddIngestJob(IngestJob ingestJob);
}