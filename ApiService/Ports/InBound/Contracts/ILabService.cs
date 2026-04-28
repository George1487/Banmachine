using Domain.Labs;
using Domain.Users;
using Ports.InBound.DTO;

namespace Ports.InBound.Contracts;

public interface ILabService
{
    LabResult GetLab(Guid labId);
    
    LabsResult GetLabs(Guid userId, UserRole role);

    LabResult AddLab(Lab lab);
    
    LabResult PatchLab(
        Guid labId,
        Guid userId, 
        LabStatus newLabStatus,
        string newTitle, 
        string newDescription, 
        DateTimeOffset newDeadlineAt);
}