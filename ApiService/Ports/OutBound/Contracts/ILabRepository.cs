using Domain.Labs;
using Ports.InBound.DTO;

namespace Ports.OutBound.Contracts;

public interface ILabRepository
{
    LabResult GetLab(Guid labId);
    
    LabResult GetLabByTitle(string title);
    
    LabsResult GetLabs();

    LabResult AddLab(Lab lab);

    LabResult PatchLab(
        Guid labId,
        LabStatus newLabStatus,
        string newTitle,
        string newDescription,
        DateTimeOffset newDeadlineAt);
}
