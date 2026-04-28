using Domain.Labs;
using Infrastructure.Repositories;
using Ports.InBound.DTO;
using Ports.OutBound.Contracts;

namespace Infrastructure.Adapters;

public class LabRepositoryAdapter : ILabRepository
{
    private readonly LabRepository _repository;

    public LabRepositoryAdapter(LabRepository repository)
    {
        _repository = repository;
    }

    public LabResult GetLab(Guid labId)
    {
        return _repository.GetLab(labId);
    }

    public LabResult GetLabByTitle(string title)
    {
        return _repository.GetLabByTitle(title);
    }

    public LabsResult GetLabs()
    {
        return _repository.GetLabs();
    }

    public LabResult AddLab(Lab lab)
    {
        return _repository.AddLab(lab);
    }

    public LabResult PatchLab(
        Guid labId,
        LabStatus newLabStatus,
        string newTitle,
        string newDescription,
        DateTimeOffset newDeadlineAt)
    {
        return _repository.PatchLab(labId, newLabStatus, newTitle, newDescription, newDeadlineAt);
    }
}
